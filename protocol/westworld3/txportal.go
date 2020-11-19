package westworld3

import (
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
	"github.com/openziti/dilithium/util"
	"github.com/pkg/errors"
	"math"
	"net"
	"sync"
	"time"
)

type txPortal struct {
	lock         *sync.Mutex
	tree         *btree.Tree
	capacity     int
	ready        *sync.Cond
	txPortalSz   int
	rxPortalSz   int
	successCt    int
	successAccum int
	dupAckCt     int
	retxCt       int
	lastRttProbe time.Time
	monitor      *retxMonitor
	closeWaitSeq int32
	closed       bool
	conn         *net.UDPConn
	peer         *net.UDPAddr
	pool         *pool
	profile      *Profile
	ii           InstrumentInstance
}

func newTxPortal(conn *net.UDPConn, peer *net.UDPAddr, profile *Profile, pool *pool, ii InstrumentInstance) *txPortal {
	p := &txPortal{
		lock:         new(sync.Mutex),
		tree:         btree.NewWith(profile.TxPortalTreeLen, utils.Int32Comparator),
		capacity:     profile.TxPortalStartSz,
		rxPortalSz:   -1,
		closeWaitSeq: -1,
		closed:       false,
		conn:         conn,
		peer:         peer,
		pool:         pool,
		profile:      profile,
		ii:           ii,
	}
	p.ready = sync.NewCond(p.lock)
	p.monitor = newRetxMonitor(profile, conn, peer, p.lock, p.ii)
	p.monitor.setRetxF(p.retx)
	return p
}

func (self *txPortal) tx(p []byte, seq *util.Sequence) (n int, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.closeWaitSeq != -1 || self.closed {
		return 0, errors.New("closed")
	}

	remaining := len(p)
	n = 0
	for remaining > 0 {
		segmentSz := int(math.Min(float64(remaining), float64(self.profile.MaxSegmentSz)))

		var rtt *uint16
		if time.Since(self.lastRttProbe).Milliseconds() > int64(self.profile.RttProbeMs) {
			now := time.Now()
			rtt = new(uint16)
			*rtt = uint16(now.UnixNano() / int64(time.Millisecond))
			segmentSz -= 2
			self.lastRttProbe = now
		}

		for math.Min(float64(self.capacity-(self.txPortalSz+segmentSz)), float64(self.capacity-self.rxPortalSz)) < 0 {
			self.ready.Wait()
		}

		wm, err := newData(seq.Next(), rtt, p[n:n+segmentSz], self.pool)
		if err != nil {
			return 0, errors.Wrap(err, "new data")
		}
		self.tree.Put(wm.seq, wm)
		self.txPortalSz += segmentSz
		self.ii.TxPortalSzChanged(self.peer, self.txPortalSz)

		if err := writeWireMessage(wm, self.conn, self.peer); err != nil {
			return 0, errors.Wrap(err, "tx")
		}
		self.ii.WireMessageTx(self.peer, wm)

		self.monitor.add(wm)

		n += segmentSz
		remaining -= segmentSz
	}

	return n, nil
}

func (self *txPortal) ack(acks []ack) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	lastTxPortalSz := self.txPortalSz
	for _, ack := range acks {
		for seq := ack.start; seq <= ack.end; seq++ {
			if v, found := self.tree.Get(seq); found {
				wm := v.(*wireMessage)
				self.monitor.remove(wm)
				self.tree.Remove(seq)
				sz, err := wm.asDataSize()
				if err != nil {
					return errors.Wrap(err, "internal tree error")
				}
				self.txPortalSz -= int(sz)
				self.ii.TxPortalSzChanged(self.peer, self.txPortalSz)
				wm.buffer.unref()

				self.successfulAck(int(sz))

				if wm.seq == self.closeWaitSeq {
					self.closed = true
				}
			} else {
				self.duplicateAck(seq)
			}
		}
	}

	if self.txPortalSz != lastTxPortalSz {
		self.ii.TxPortalSzChanged(self.peer, self.txPortalSz)
	}

	self.ready.Broadcast()
	return nil
}

func (self *txPortal) updateRxPortalSz(rxPortalSz int) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.rxPortalSz = rxPortalSz
	self.ready.Broadcast()
	self.ii.TxPortalRxSzChanged(self.peer, rxPortalSz)
}

func (self *txPortal) rtt(probeTs uint16) {
	now := time.Now().UnixNano()
	self.lock.Lock()
	clockTs := uint16(now / int64(time.Millisecond))
	rttMs := clockTs - probeTs
	self.monitor.updateRttMs(rttMs)
	self.lock.Unlock()
}

func (self *txPortal) close(seq *util.Sequence) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if !self.closed {
		wm, err := newClose(seq.Next(), self.pool)
		if err != nil {
			return errors.Wrap(err, "close")
		}
		self.closeWaitSeq = wm.seq
		self.tree.Put(wm.seq, wm)
		self.monitor.add(wm)

		if err := writeWireMessage(wm, self.conn, self.peer); err != nil {
			return errors.Wrap(err, "tx close")
		}
		self.ii.WireMessageTx(self.peer, wm)
	}

	return nil
}

func (self *txPortal) successfulAck(sz int) {
	self.successCt++
	self.successAccum += sz
	if self.successCt == self.profile.TxPortalIncreaseThresh {
		newCapacity := self.capacity + int(float64(self.successAccum)*self.profile.TxPortalIncreaseScale)
		self.updatePortalCapacity(newCapacity)
		self.successCt = 0
		self.successAccum = 0
	}
}

func (self *txPortal) duplicateAck(seq int32) {
	self.dupAckCt++
	self.successCt = 0
	if self.dupAckCt >= self.profile.TxPortalDupAckThresh {
		newCapacity := int(float64(self.capacity) * self.profile.TxPortalDupAckCapacityScale)
		self.profile.RetxAddMs += 10 // #93: Self-Adjusting retxMs
		self.updatePortalCapacity(newCapacity)
		self.dupAckCt = 0
		self.successAccum = int(float64(self.successAccum) * self.profile.TxPortalDupAckSuccessScale)
	}
	self.ii.DuplicateAck(self.peer, seq)
}

func (self *txPortal) retx() {
	self.retxCt++
	self.successCt = 0
	if self.retxCt >= self.profile.TxPortalRetxThresh {
		newCapacity := int(float64(self.capacity) * self.profile.TxPortalRetxCapacityScale)
		self.updatePortalCapacity(newCapacity)
		self.retxCt = 0
		self.successAccum = int(float64(self.successAccum) * self.profile.TxPortalRetxSuccessScale)
	}
}

func (self *txPortal) updatePortalCapacity(newCapacity int) {
	oldCapacity := self.capacity
	self.capacity = newCapacity
	if self.capacity < self.profile.TxPortalMinSz {
		self.capacity = self.profile.TxPortalMinSz
	}
	if self.capacity > self.profile.TxPortalMaxSz {
		self.capacity = self.profile.TxPortalMaxSz
	}
	if self.capacity != oldCapacity {
		self.ii.TxPortalCapacityChanged(self.peer, self.capacity)
	}
}
