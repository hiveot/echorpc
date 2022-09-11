package main_test

import (
	"context"
	"net"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiveot/echorpc/pkg/echoservice"

	bph "github.com/hiveot/echorpc/capnp/go"
)

const testSocketAddr = "/tmp/echocapabilities-test.socket"

// Testcase for performance profiling
// TODO: test how many echo service instances there are if also invoked directly
func TestGetService(t *testing.T) {
	const text = "hello echo"
	// launch the server
	go echoservice.EchoBootstrapCapnpAdapterStart(testSocketAddr, true)
	time.Sleep(time.Second)

	// step 1: create the network connection
	connection, err := net.Dial("unix", testSocketAddr)
	assert.NoError(t, err)
	transport := rpc.NewStreamTransport(connection)
	rpcConn := rpc.NewConn(transport, nil)
	ctx := context.Background()

	// step 2: create a client instance of the bootstrap service
	echoBootClient := bph.EchoBootstrap(rpcConn.Bootstrap(ctx))

	// step 3: call the get echo method
	resp, release := echoBootClient.GetEcho(ctx,
		func(bph.EchoBootstrap_getEcho_Params) error {
			return nil
		})
	require.NotNil(t, resp)

	// step 4: get the result
	result, err := resp.Struct()
	require.NoError(t, err)

	// step 5: get the result again
	echoCapability := result.Service()
	assert.NotNil(t, echoCapability)

	// step 6: invoke echo
	//echoText, err := CallEcho(echoCap, text)
	echoResp, release := echoCapability.Echo(ctx,
		func(params bph.EchoService_echo_Params) error {
			params.SetText(text)
			return nil
		})
	defer release()
	assert.NotNil(t, echoResp)

	// step 7: get result
	echoRespData, err := echoResp.Struct()
	require.NoError(t, err)

	echoText, err := echoRespData.EchoText()
	require.NoError(t, err)
	assert.Equal(t, text, echoText)
	//echoSvc.Echo(ctx, )
	release()

}
