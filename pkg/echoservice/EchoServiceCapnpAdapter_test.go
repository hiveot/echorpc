package echoservice_test

import (
	"context"
	"net"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/rpc"
	"github.com/stretchr/testify/assert"

	cb "github.com/hiveot/echorpc/capnp/go"
	"github.com/hiveot/echorpc/pkg/echoservice"
)

const testSocketAddr = "/tmp/echoservice-test.socket"

// Testcase for performance profiling
func TestUpper(t *testing.T) {
	const count = 1
	// launch the server
	go echoservice.EchoServiceCapnpStart(testSocketAddr, true)
	time.Sleep(time.Second)

	// create the network connection
	connection, err := net.Dial("unix", testSocketAddr)
	assert.NoError(t, err)
	transport := rpc.NewStreamTransport(connection)
	rpcConn := rpc.NewConn(transport, nil)
	ctx := context.Background()
	echoClient := cb.EchoServiceCap(rpcConn.Bootstrap(ctx))

	for i := 0; i < count; i++ {
		// Create the request. This does not yet invoke the method
		resp, release := echoClient.Upper(ctx,
			func(params cb.EchoServiceCap_upper_Params) error {
				err = params.SetText("Hello world")
				return err
			})
		// The remote method is invoked by asking for the result
		result, err := resp.Struct()
		assert.NoError(t, err)

		upperText, err := result.UpperText()
		_ = upperText
		assert.NoError(t, err)
		release()
	}
	//fmt.Println(upperText)

}
