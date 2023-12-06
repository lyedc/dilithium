package dilithium

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"gitee.com/zhaochuninhefei/gmgo/gmtls"
	"math/big"
	"net"
	"time"
)

func generateTLSConfig() *gmtls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-24 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
		IsCA:         true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := gmtls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &gmtls.Config{
		Certificates:       []gmtls.Certificate{tlsCert},
		NextProtos:         []string{"dilithium"},
		InsecureSkipVerify: true,
	}
}
