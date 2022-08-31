package echoservice

import (
	"context"
	"net"
	"os"

	// use v3. By default it pulls in v2 (zombiezen.com/go/capnproto2)
	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	echocap "github.com/hiveot/echorpc/capnp/go"
	"github.com/sirupsen/logrus"
)

// cap'n'proto adapter for echo service
// EchoServiceCapnpAdapter implements the generated echocap.EchoServiceCap_Server
// interface. Copy the interface method when building the adapter.
type EchoServiceCapnpAdapter struct {
	svc *EchoService
}

// The adapter provides the capnp parameters defined
func (adapter *EchoServiceCapnpAdapter) Echo(
	_ context.Context, call echocap.EchoServiceCap_echo) error {
	text, _ := call.Args().Text()
	echoText, err := adapter.svc.Echo(text)
	res, _ := call.AllocResults()
	res.SetEchoText(echoText)
	call.Ack()
	return err
}

func (adapter *EchoServiceCapnpAdapter) Reverse(
	_ context.Context, call echocap.EchoServiceCap_reverse) error {
	text, _ := call.Args().Text()
	revText, err := adapter.svc.Reverse(text)
	res, _ := call.AllocResults()
	res.SetReverseText(revText)
	call.Ack()
	return err
}

func (adapter *EchoServiceCapnpAdapter) Upper(
	_ context.Context, call echocap.EchoServiceCap_upper) error {
	text, _ := call.Args().Text()
	upText, err := adapter.svc.Upper(text)
	res, _ := call.AllocResults()
	res.SetUpperText(upText)
	call.Ack()
	return err
}

// EchoServiceCapnpStart start listening
//  address to list on: ":port", "host:port", "/tmp/path-to.socket"
//  isUDS set to true when address is a unix domain socket
func EchoServiceCapnpStart(address string, isUDS bool) {
	// fmt.Println("EchoServiceCapnpStart starting echo-service on capnp address", address)

	network := "tcp"
	if isUDS {
		os.Remove(address)
		network = "unix"
	}
	listener, err := net.Listen(network, address)
	if err != nil {
		logrus.Fatalf("Failed open listener: %v", err)
	}

	// create the service main handler
	echoSvc := NewEchoService()
	main := echocap.EchoServiceCap_ServerToClient(&EchoServiceCapnpAdapter{
		svc: echoSvc,
	})

	// listen for calls
	for {
		rwc, _ := listener.Accept()
		go func() error {
			transport := rpc.NewStreamTransport(rwc)
			conn := rpc.NewConn(transport, &rpc.Options{
				BootstrapClient: capnp.Client(main.AddRef()),
			})
			defer conn.Close()
			// Wait for connection to abort.
			select {
			case <-conn.Done():
				return nil
				// case <-ctx.Done():
				// 	return conn.Close()
			}
		}()
	}
}
