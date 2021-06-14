package tunnel

import (
	"github.com/openziti/dilithium/cmd/dilithium/dilithium"
	"github.com/openziti/dilithium/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"net"
	"runtime"
)

func init() {
	tunnelCmd.AddCommand(tunnelClientCmd)
}

var tunnelClientCmd = &cobra.Command{
	Use:   "client <serverAddress> <listenTcpAddress>",
	Short: "Start tunnel client",
	Args:  cobra.ExactArgs(2),
	Run:   tunnelClient,
}

func tunnelClient(_ *cobra.Command, args []string) {
	cl, err := util.GetCtrlListener(".", "tunnel")
	if err != nil {
		logrus.Fatalf("error getting ctrl listener (%v)", err)
	}
	cl.AddCallback("stacks", func(_ string, conn net.Conn) (int64, error) {
		buf := make([]byte, 64 * 1024)
		n := runtime.Stack(buf, true)
		var err error
		n, err = conn.Write(buf[:n])
		return int64(n), err
	})
	cl.Start()

	serverAddress := args[0]
	listenAddress, err := net.ResolveTCPAddr("tcp", args[1])
	if err != nil {
		logrus.Fatalf("error resolving listen address [%s] (%v)", args[1], err)
	}

	initiatorListener, err := net.ListenTCP("tcp", listenAddress)
	if err != nil {
		logrus.Infof("error creating initiator listener at [%s] (%v)", listenAddress, err)
	}
	logrus.Infof("created initiator listener at [%s]", initiatorListener.Addr())

	for {
		initiator, err := initiatorListener.Accept()
		if err != nil {
			logrus.Errorf("error accepting initiator (%v)", err)
			continue
		}
		go handleTunnelInitiator(initiator, serverAddress)
	}
}

func handleTunnelInitiator(initiator net.Conn, serverAddress string) {
	defer func() { _ = initiator.Close() }()

	logrus.Infof("tunneling for initiator at [%s]", initiator.RemoteAddr())
	defer logrus.Warnf("end tunnel for initiator at [%s]", initiator.RemoteAddr())

	protocol, err := dilithium.ProtocolFor(dilithium.SelectedProtocol)
	if err != nil {
		logrus.Fatalf("error selecting protocol (%v)", err)
	}

	tunnel, err := protocol.Dial(serverAddress)
	if err != nil {
		logrus.Errorf("error dialing tunnel server at [%s] (%v)", serverAddress, err)
		return
	}
	go handleTunnelInitiatorReader(initiator, tunnel)
	defer func() { _ = tunnel.Close() }()
	logrus.Infof("tunnel established to [%s]", serverAddress)

	buffer := make([]byte, bufferSize)
	for {
		n, err := initiator.Read(buffer)
		if err != nil {
			logrus.Errorf("error reading from initiator (%v)", err)
			return
		}
		n, err = tunnel.Write(buffer[:n])
		if err != nil {
			logrus.Errorf("error writing to tunnel (%v)", err)
			return
		}
	}
}

func handleTunnelInitiatorReader(initiator net.Conn, tunnel net.Conn) {
	defer logrus.Warnf("closing")
	defer func() {
		_ = tunnel.Close()
		_ = initiator.Close()
	}()

	buffer := make([]byte, bufferSize)
	for {
		n, err := tunnel.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			logrus.Errorf("error reading from tunnel (%v)", err)
		}
		n, err = initiator.Write(buffer[:n])
		if err != nil {
			logrus.Errorf("error writing to initiator (%v)", err)
			return
		}
	}
}
