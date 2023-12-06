package westlsworld3

import (
	"gitee.com/zhaochuninhefei/gmgo/gmtls"
	"github.com/openziti/dilithium/protocol/westworld3"
	"net"
)

func Dial(addr *net.UDPAddr, tlsConfig *gmtls.Config, profileId byte) (net.Conn, error) {
	w3Conn, err := westworld3.Dial(addr, profileId)
	if err != nil {
		return nil, err
	}
	return gmtls.Client(w3Conn, tlsConfig), nil
}
