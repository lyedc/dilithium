package westworld2

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Config struct {
	seqRandom            bool
	txPortalStartSz      int
	txPortalMinSz        int
	txPortalMaxSz        int
	txPortalIncreaseCt   int
	txPortalIncreaseFrac float64
	txPortalDupAckCt     int
	txPortalDupAckFrac   float64
	txPortalRetxCt       int
	txPortalRetxFrac     float64
	retxStartMs          int
	retxScale            float64
	retxAddMs            int
	rttProbeMs           int
	rttProbeAvgCt        int
	maxSegmentSz         int
	poolBufferSz         int
	rxBufferSz           int
	txBufferSz           int
	treeLen              int
	readsQLen            int
	listenerRxQLen       int
	acceptQLen           int
	i                    Instrument
}

func NewDefaultConfig() *Config {
	// san_fran_2
	return &Config{
		seqRandom:            true,
		txPortalStartSz:      16384,
		txPortalMinSz:        16384,
		txPortalMaxSz:        4096000,
		txPortalIncreaseCt:   224,
		txPortalIncreaseFrac: 1.0,
		txPortalDupAckCt:     64,
		txPortalDupAckFrac:   0.9,
		txPortalRetxCt:       64,
		txPortalRetxFrac:     0.75,
		retxStartMs:          200,
		retxScale:            2.0,
		retxAddMs:            100,
		rttProbeMs:           50,
		rttProbeAvgCt:        8,
		maxSegmentSz:         1420,
		poolBufferSz:         64 * 1024,
		rxBufferSz:           16000000,
		txBufferSz:           16000000,
		treeLen:              1024,
		readsQLen:            1024,
		listenerRxQLen:       1024,
		acceptQLen:           1024,
	}
}

func (self *Config) Load(data map[interface{}]interface{}) error {
	if v, found := data["seq_random"]; found {
		if b, ok := v.(bool); ok {
			self.seqRandom = b
		} else {
			return errors.New("invalid 'seq_random' value")
		}
	}
	if v, found := data["tx_portal_start_sz"]; found {
		if i, ok := v.(int); ok {
			self.txPortalStartSz = i
		} else {
			return errors.New("invalid 'tx_portal_start_sz' value")
		}
	}
	if v, found := data["tx_portal_min_sz"]; found {
		if i, ok := v.(int); ok {
			self.txPortalMinSz = i
		} else {
			return errors.New("invalid 'tx_portal_min_sz' value")
		}
	}
	if v, found := data["tx_portal_max_sz"]; found {
		if i, ok := v.(int); ok {
			self.txPortalMaxSz = i
		} else {
			return errors.New("invalid 'tx_portal_max_sz' value")
		}
	}
	if v, found := data["tx_portal_increase_ct"]; found {
		if i, ok := v.(int); ok {
			self.txPortalIncreaseCt = i
		} else {
			return errors.New("invalid 'tx_portal_increase_ct' value")
		}
	}
	if v, found := data["tx_portal_increase_frac"]; found {
		if f, ok := v.(float64); ok {
			self.txPortalIncreaseFrac = f
		} else {
			return errors.New("invalid 'tx_portal_increase_frac' value")
		}
	}
	if v, found := data["tx_portal_dup_ack_ct"]; found {
		if i, ok := v.(int); ok {
			self.txPortalDupAckCt = i
		} else {
			return errors.New("invalid 'tx_portal_dup_ack_ct' value")
		}
	}
	if v, found := data["tx_portal_dup_ack_frac"]; found {
		if f, ok := v.(float64); ok {
			self.txPortalDupAckFrac = f
		} else {
			return errors.New("invalid 'tx_portal_dup_ack_frac' value")
		}
	}
	if v, found := data["tx_portal_retx_ct"]; found {
		if i, ok := v.(int); ok {
			self.txPortalRetxCt = i
		} else {
			return errors.New("invalid 'tx_portal_retx_ct' value")
		}
	}
	if v, found := data["tx_portal_retx_frac"]; found {
		if f, ok := v.(float64); ok {
			self.txPortalRetxFrac = f
		} else {
			return errors.New("invalid 'tx_portal_retx_frac' value")
		}
	}
	if v, found := data["retx_start_ms"]; found {
		if i, ok := v.(int); ok {
			self.retxStartMs = i
		} else {
			return errors.New("invalid 'retx_start_ms' value")
		}
	}
	if v, found := data["retx_scale"]; found {
		if f, ok := v.(float64); ok {
			self.retxScale = f
		} else {
			return errors.New("invalid 'retx_scale' value")
		}
	}
	if v, found := data["retx_add_ms"]; found {
		if i, ok := v.(int); ok {
			self.retxAddMs = i
		} else {
			return errors.New("invalid 'retx_add_ms' value")
		}
	}
	if v, found := data["rtt_probe_ms"]; found {
		if i, ok := v.(int); ok {
			self.rttProbeMs = i
		} else {
			return errors.New("invalid 'rtt_probe_ms' value")
		}
	}
	if v, found := data["rtt_probe_avg_ct"]; found {
		if i, ok := v.(int); ok {
			self.rttProbeAvgCt = i
		} else {
			return errors.New("invalid 'rtt_probe_avg_ct' value")
		}
	}
	if v, found := data["max_segment_sz"]; found {
		if i, ok := v.(int); ok {
			self.maxSegmentSz = i
		} else {
			return errors.New("invalid 'max_segment_sz' value")
		}
	}
	if v, found := data["pool_buffer_sz"]; found {
		if i, ok := v.(int); ok {
			self.poolBufferSz = i
		} else {
			return errors.New("invalid 'pool_buffer_sz' value")
		}
	}
	if v, found := data["rx_buffer_sz"]; found {
		if i, ok := v.(int); ok {
			self.rxBufferSz = i
		} else {
			return errors.New("invalid 'rx_buffer_sz' value")
		}
	}
	if v, found := data["tx_buffer_sz"]; found {
		if i, ok := v.(int); ok {
			self.txBufferSz = i
		} else {
			return errors.New("invalid 'tx_buffer_sz' value")
		}
	}
	if v, found := data["tree_len"]; found {
		if i, ok := v.(int); ok {
			self.treeLen = i
		} else {
			return errors.New("invalid 'tree_len' value")
		}
	}
	if v, found := data["reads_q_len"]; found {
		if i, ok := v.(int); ok {
			self.readsQLen = i
		} else {
			return errors.New("invalid 'reads_q_len' value")
		}
	}
	if v, found := data["listener_rx_q_len"]; found {
		if i, ok := v.(int); ok {
			self.listenerRxQLen = i
		} else {
			return errors.New("invalid 'listener_rx_q_len' value")
		}
	}
	if v, found := data["accept_q_len"]; found {
		if i, ok := v.(int); ok {
			self.acceptQLen = i
		} else {
			return errors.New("invalid 'accept_q_len' value")
		}
	}
	if v, found := data["instrument"]; found {
		submap, oks := v.(map[string]interface{})
		if !oks {
			if subi, oki := v.(map[interface{}]interface{}); oki {
				submap = make(map[string]interface{})
				oks = true
				for k, v := range subi {
					if s, ok := k.(string); ok {
						submap[s] = v
					} else {
						oks = false
					}
				}
			}
		}
		if oks {
			if v, found := submap["name"]; found {
				if name, ok := v.(string); ok {
					i, err := NewInstrument(name, submap)
					if err != nil {
						return errors.Wrap(err, "error creating instrument")
					}
					self.i = i
				} else {
					return errors.New("invalid 'instrument/name' value")
				}
			} else {
				return errors.New("missing 'instrument/name'")
			}
		} else {
			return errors.Errorf("invalid 'instrument' value [%v]", reflect.TypeOf(v))
		}
	}
	return nil
}

func (self *Config) Dump() string {
	out := "westworld2.Config{\n"
	out += fmt.Sprintf("\t%-30s %t\n", "seq_random", self.seqRandom)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_portal_start_sz", self.txPortalStartSz)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_portal_min_sz", self.txPortalMinSz)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_portal_max_sz", self.txPortalMaxSz)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_portal_incrase_ct", self.txPortalIncreaseCt)
	out += fmt.Sprintf("\t%-30s %.4f\n", "tx_portal_increase_frac", self.txPortalIncreaseFrac)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_portal_dup_ack_ct", self.txPortalDupAckCt)
	out += fmt.Sprintf("\t%-30s %.4f\n", "tx_portal_dup_ack_frac", self.txPortalDupAckFrac)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_portal_retx_ct", self.txPortalRetxCt)
	out += fmt.Sprintf("\t%-30s %.4f\n", "tx_portal_retx_frac", self.txPortalRetxFrac)
	out += fmt.Sprintf("\t%-30s %d\n", "retx_start_ms", self.retxStartMs)
	out += fmt.Sprintf("\t%-30s %.4f\n", "retx_scale", self.retxScale)
	out += fmt.Sprintf("\t%-30s %d\n", "retx_add_ms", self.retxAddMs)
	out += fmt.Sprintf("\t%-30s %d\n", "rtt_probe_ms", self.rttProbeMs)
	out += fmt.Sprintf("\t%-30s %d\n", "rtt_probe_avg_ct", self.rttProbeAvgCt)
	out += fmt.Sprintf("\t%-30s %d\n", "max_segment_sz", self.maxSegmentSz)
	out += fmt.Sprintf("\t%-30s %d\n", "pool_buffer_sz", self.poolBufferSz)
	out += fmt.Sprintf("\t%-30s %d\n", "rx_buffer_sz", self.rxBufferSz)
	out += fmt.Sprintf("\t%-30s %d\n", "tx_buffer_sz", self.txBufferSz)
	out += fmt.Sprintf("\t%-30s %d\n", "tree_len", self.treeLen)
	out += fmt.Sprintf("\t%-30s %d\n", "reads_q_len", self.readsQLen)
	out += fmt.Sprintf("\t%-30s %d\n", "listener_rx_q_len", self.listenerRxQLen)
	out += fmt.Sprintf("\t%-30s %d\n", "accept_q_len", self.acceptQLen)
	out += fmt.Sprintf("\t%-30s %v\n", "instrument", reflect.TypeOf(self.i))
	out += "}"
	return out
}
