package echoservice

import (
	"context"
	"fmt"
	"net"
	"os"

	pb "github.com/hiveot/echorpc/grpc/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// gRPC Adapter for echo service
type EchoServiceGrpcAdapter struct {
	pb.UnimplementedEchoServiceServer

	svc *EchoService
}

func (adapter *EchoServiceGrpcAdapter) Echo(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	if args == nil {
		return nil, fmt.Errorf("missing args")
	}
	res, err := adapter.svc.Echo(args.Text)
	return &pb.TextParam{Text: res}, err
}

func (adapter *EchoServiceGrpcAdapter) Upper(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	// fmt.Println("EchoServiceGrpcAdapter.Upper")
	if args == nil {
		return nil, fmt.Errorf("missing args")
	}
	upperText, err := adapter.svc.Upper(args.Text)
	response := pb.TextParam{Text: upperText}
	return &response, err
}

func (adapter *EchoServiceGrpcAdapter) Reverse(ctx context.Context, args *pb.TextParam) (*pb.TextParam, error) {
	if args == nil {
		return nil, fmt.Errorf("missing args")
	}
	reverseText, err := adapter.svc.Reverse(args.Text)
	response := pb.TextParam{Text: reverseText}
	return &response, err
}

// EchoServiceGrpcStart start listening
//  address to list on: ":port", "host:port", "/tmp/path-to.socket"
//  isUDS set to true when address is a unix domain socket
func EchoServiceGrpcStart(address string, isUDS bool) {
	// fmt.Println("EchoServiceGrpcStart starting echo-service on grpc address", socketAddrGrpc)

	var opts []grpc.ServerOption

	network := "tcp"
	if isUDS {
		os.Remove(address)
		network = "unix"
	}
	grpcServer := grpc.NewServer(opts...)
	echoSvc := NewEchoService()
	adapter := &EchoServiceGrpcAdapter{
		svc: echoSvc,
	}

	// register the gRPC service that can be invoked by any gRPC client
	pb.RegisterEchoServiceServer(grpcServer, adapter)
	if isUDS {
		os.Remove(address)
	}
	lis, err := net.Listen(network, address)
	if err != nil {
		logrus.Fatalf("Failed open listener: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
