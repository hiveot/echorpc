package echoclient

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"capnproto.org/go/capnp/v3/rpc"
	cb "github.com/hiveot/echorpc/capnp/go"
)

// Invoke the service using capnp rpc
func InvokeUpperCapnp(address string, isUDS bool, text string, count int) {
	// Set up a connection to the server. Max 200 second test run
	// fmt.Println("Invoking upper over capnp")

	// create a net connection with context and timeout
	// a socket must have been created by the service
	network := "tcp"
	if isUDS {
		network = "unix"
	}
	connection, err := net.Dial(network, address)
	if err != nil {
		log.Fatalf("error connecting to '%s': %s", address, err)
	}
	transport := rpc.NewStreamTransport(connection)
	rpcConn := rpc.NewConn(transport, nil)
	defer rpcConn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	echoClient := cb.EchoServiceCap(rpcConn.Bootstrap(ctx))

	t1 := time.Now()
	for i := 0; i < count; i++ {

		// create the client capability (eg message)
		resp, release := echoClient.Upper(ctx,
			func(params cb.EchoServiceCap_upper_Params) error {
				err = params.SetText(text)
				return err
			})

		result, err := resp.Struct()
		if err != nil {
			log.Fatalf("error getting response struct: %v", err)
		}
		upperText, err := result.UpperText()
		if err != nil {
			log.Fatalf("error getting upper result text: %s", err)
		}
		_ = upperText
		release()
		// fmt.Println("Response:", upperText)
	}
	d1 := time.Since(t1)
	// round result
	msec := d1.Milliseconds() / 10 * 10
	fmt.Printf("%d calls using Capnp on %s: %d millisec\n", count, address, msec)
	cancel()

	// return response.Text, err
}
