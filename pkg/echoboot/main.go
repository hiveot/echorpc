// Package main with a demonstration on how to obtain additional capabilities with capnp
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"capnproto.org/go/capnp/v3/rpc"

	bph "github.com/hiveot/echorpc/capnp/go"
	"github.com/hiveot/echorpc/pkg/echoservice"
)

const socketAddrEchoCapnp = "/tmp/echocapabilities-capnp.socket"

// Return the capability to echo
func GetEchoCapability(rpcConn *rpc.Conn) (*bph.EchoService, error) {

	// step 2: create a client instance of the bootstrap service
	ctx := context.Background()
	echoBootClient := bph.EchoBootstrap(rpcConn.Bootstrap(ctx))

	// step 3: call the get echo method
	resp, release := echoBootClient.GetEcho(ctx,
		func(bph.EchoBootstrap_getEcho_Params) error {
			return nil
		})
	//defer release()
	_ = release

	// step 4: get the result
	result, err := resp.Struct()
	if err != nil {
		return nil, fmt.Errorf("EchoBootstrap has no result")
	}

	// step 5: get the capability from the result
	echoSvc := result.Service()
	// if echoSvc == nil {
	// 	return nil, fmt.Errorf("EchoBootstrap result has no echoservice")
	// }
	return &echoSvc, err
}

// Invoke the echo function
func CallEcho(echoCapability *bph.EchoService, text string) (echoText string, err error) {
	// Invoke echo
	ctx := context.Background()
	echoResp, release := echoCapability.Echo(ctx,
		func(params bph.EchoService_echo_Params) error {
			params.SetText(text)
			return nil
		})
	defer release()

	echoRespData, err := echoResp.Struct()
	if err != nil {
		return "", err
	}
	echoText, err = echoRespData.EchoText()
	return echoText, err
}

// Invoke the echo service after getting it as a capability
func main() {
	// Start the echo and capability services
	go echoservice.EchoBootstrapCapnpAdapterStart(socketAddrEchoCapnp, true)
	time.Sleep(time.Second)

	// Get capability
	connection, err := net.Dial("unix", socketAddrEchoCapnp)
	if err != nil {
		fmt.Println("Error retrieving Echo capability:", err)
		os.Exit(1)
	}
	transport := rpc.NewStreamTransport(connection)
	rpcConn := rpc.NewConn(transport, nil)

	echoCap, err := GetEchoCapability(rpcConn)
	if err != nil {
		fmt.Println("Error retrieving Echo capability:", err)
		os.Exit(1)
	}
	// Invoke echo
	echoText, err := CallEcho(echoCap, "Hello world")
	if err != nil {
		fmt.Println("Error invoking Echo capability:", err)
		os.Exit(1)
	}

	fmt.Println("Hoora, received echo response: ", echoText)

	// Do it again
	echoText2, err := CallEcho(echoCap, "Hello world 2")
	if err != nil {
		fmt.Println("but failed not a second time:", err)
		os.Exit(1)
	}
	fmt.Println("Second result: ", echoText2)

}
