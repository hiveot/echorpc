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

// InvokeEchoCapnp Invoke the service using capnp rpc
func InvokeEchoCapnp(address string, isUDS bool, text string, count int) {
	// Set up a connection to the server. Max 200 second test run
	// fmt.Println("Invoking echo over capnp")

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

	// Q: in grpc this is NewEchoServiceClient. Why not do the same? - lower barrier to entry
	//    why the need to provide a bootstrap? Can the generated code create this?
	echoClient := cb.EchoService(rpcConn.Bootstrap(ctx))

	t1 := time.Now()
	for i := 0; i < count; i++ {

		// create the client capability (eg request instance)
		// Q: Why the callback and the 'funny' params type?
		//    Lower barrier to entry by not having these. Eg:
		// A: avoid having to expose EchoServiceCap_echo_Params...
		//    method, release := echoClient.Echo(ctx)
		//    method.Params.SetText(text)
		// B: the grpc way would look familiar to people used to gRPC
		//    params := EchoServiceCap_echo_Params{...}
		//    method, release := echoClient.Echo(ctx, params)
		//
		// get results with:
		//   result = method.GetResult()
		resp, release := echoClient.Echo(ctx,
			func(params cb.EchoService_echo_Params) error {
				err = params.SetText(text)
				return err
			})
		// invoke the request by asking for a result
		// Q: why resp.Struct() instead of something like resp.Get()? Isn't the root always a struct?
		result, err := resp.Struct()
		if err != nil {
			log.Fatalf("error getting response struct: %v", err)
		}
		echoText, err := result.EchoText()
		if err != nil {
			log.Fatalf("error getting echo result text: %s", err)
		}
		_ = echoText
		release()

		// second call to get stats
		resp2, release2 := echoClient.Stats(ctx,
			func(params cb.EchoService_stats_Params) error {
				return nil
			})
		result2, err := resp2.Struct()
		s2, err := result2.Stats()
		//fmt.Printf("Result of stats: %v\n", s2)
		_ = s2
		release2()

		// fmt.Println("Response:", echoText)
	}
	d1 := time.Since(t1)
	// round result
	msec := d1.Milliseconds() / 10 * 10
	fmt.Printf("%d calls using Capnp on %s: %d millisec\n", count, address, msec)
	cancel()

	// return response.Text, err
}
