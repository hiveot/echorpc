package echoservice

import (
	"context"
	"net"
	"os"

	// use v3. By default it pulls in v2 (zombiezen.com/go/capnproto2)
	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"github.com/sirupsen/logrus"

	echocap "github.com/hiveot/echorpc/capnp/go"
	echocapnp "github.com/hiveot/echorpc/capnp/go"
)

// EchoBootstrapCapnpAdapter implements the generated echocap.EchoBootstrap_Server
// interface.
type EchoBootstrapCapnpAdapter struct {
	// svc *echocapnp.EchoBootstrap
}

func (adapter *EchoBootstrapCapnpAdapter) GetEcho(
	_ context.Context, call echocapnp.EchoBootstrap_getEcho) error {

	res, err := call.AllocResults()
	if err != nil {
		return err
	}
	// create the service main handler
	echoSvc := NewEchoService()
	main := echocap.EchoService_ServerToClient(&EchoServiceCapnpAdapter{
		svc: echoSvc,
	})

	err = res.SetService(main)
	return err
}

// EchoBootstrapCapnpAdapterStart start listening
//  address to list on: ":port", "host:port", "/tmp/path-to.socket"
//  isUDS set to true when address is a unix domain socket
func EchoBootstrapCapnpAdapterStart(address string, isUDS bool) {

	network := "tcp"
	if isUDS {
		_ = os.Remove(address)
		network = "unix"
	}
	listener, err := net.Listen(network, address)
	if err != nil {
		logrus.Fatalf("Failed open listener: %v", err)
	}

	// create the service main handler
	main := echocapnp.EchoBootstrap_ServerToClient(&EchoBootstrapCapnpAdapter{})

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
