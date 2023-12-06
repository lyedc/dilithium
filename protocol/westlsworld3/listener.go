package westlsworld3

import (
	"gitee.com/zhaochuninhefei/gmgo/gmtls"
	"github.com/openziti/dilithium/protocol/westworld3"
	"net"
)

type listener struct {
	w3Listener net.Listener
	tlsConfig  *gmtls.Config
}

func Listen(addr *net.UDPAddr, tlsConfig *gmtls.Config, profileId byte) (net.Listener, error) {
	w3Listener, err := westworld3.Listen(addr, profileId)
	if err != nil {
		return nil, err
	}
	return &listener{
		w3Listener: w3Listener,
		tlsConfig:  tlsConfig,
	}, nil
}

func (self *listener) Accept() (net.Conn, error) {
	w3Conn, err := self.w3Listener.Accept()
	if err != nil {
		return nil, err
	}
	return gmtls.Server(w3Conn, self.tlsConfig), nil
}

func (self *listener) Close() error {
	return self.w3Listener.Close()
}

func (self *listener) Addr() net.Addr {
	return self.w3Listener.Addr()
}
