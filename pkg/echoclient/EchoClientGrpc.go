package echoclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hiveot/echorpc/grpc/go"
)

// InvokeEchoGrpc Invokes the echo service using grpc
func InvokeEchoGrpc(address string, isUDS bool, text string, count int) {
	// Set up a connection to the server. Max 200 second test run
	// fmt.Println("Invoking echo over grpc")
	cred := insecure.NewCredentials()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	if isUDS {
		address = "unix://" + address
	}
	conn, err := grpc.DialContext(ctx,
		address,
		grpc.WithTransportCredentials(cred),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)
	t1 := time.Now()
	for i := 0; i < count; i++ {
		response, err := client.Echo(ctx, &pb.TextParam{Text: text})
		_, _ = client.Stats(ctx, &emptypb.Empty{})
		if err != nil {
			log.Fatalf("error echo response: %s", err)
		}
		_ = response

		// fmt.Println("Response:", response)
	}
	d1 := time.Since(t1)
	msec := d1.Milliseconds() / 10 * 10
	fmt.Printf("%d calls using gRPC  on %s: %d millisec\n", count, address, msec)
	cancel()

	// return response.Text, err
}
