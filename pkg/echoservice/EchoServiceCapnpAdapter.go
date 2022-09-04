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
)

// EchoServiceCapnpAdapter cap'n'proto adapter for echo service
// EchoServiceCapnpAdapter implements the generated echocap.EchoService_Server
// interface. Copy the interface method when building the adapter.
type EchoServiceCapnpAdapter struct {
	svc *EchoService
}

func (adapter *EchoServiceCapnpAdapter) Echo(
	_ context.Context, call echocap.EchoService_echo) error {
	// Ack is optional for simple functions like these.
	call.Ack()
	text, _ := call.Args().Text()
	echoText, err := adapter.svc.Echo(text)
	res, _ := call.AllocResults()
	err = res.SetEchoText(echoText)
	return err
}

// Latest returns the latest echo
func (adapter *EchoServiceCapnpAdapter) Latest(
	_ context.Context, call echocap.EchoService_latest) error {
	// Ack is optional for simple functions like these.
	call.Ack()
	latestText, err := adapter.svc.Latest()
	res, _ := call.AllocResults()
	err = res.SetEchoText(latestText)
	return err
}

func (adapter *EchoServiceCapnpAdapter) Stats(
	_ context.Context, call echocap.EchoService_stats) error {
	latestText, echoCount := adapter.svc.Stats()

	arena := capnp.SingleSegment(nil)
	_, seg, err := capnp.NewMessage(arena)
	stats, err := echocap.NewRootEchoStats(seg)

	err = stats.SetLatest(latestText)
	stats.SetCount(uint32(echoCount))
	res, _ := call.AllocResults()
	err = res.SetStats(stats)
	return err
}

// EchoServiceCapnpStart start listening
//  address to list on: ":port", "host:port", "/tmp/path-to.socket"
//  isUDS set to true when address is a unix domain socket
func EchoServiceCapnpStart(address string, isUDS bool) {
	// fmt.Println("EchoServiceCapnpStart starting echo-service on capnp address", address)

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
	echoSvc := NewEchoService()
	main := echocap.EchoService_ServerToClient(&EchoServiceCapnpAdapter{
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
