package echoservice

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hiveot/echorpc/grpc/go"
)

// EchoServiceGrpcAdapter gRPC Adapter for echo service
type EchoServiceGrpcAdapter struct {
	pb.UnimplementedEchoServiceServer

	svc *EchoService
}

func (adapter *EchoServiceGrpcAdapter) Echo(_ context.Context, args *pb.TextParam) (
	*pb.TextParam, error) {

	if args == nil {
		return nil, fmt.Errorf("missing args")
	}
	res, err := adapter.svc.Echo(args.Text)
	return &pb.TextParam{Text: res}, err
}

func (adapter *EchoServiceGrpcAdapter) Latest(_ context.Context, _ *emptypb.Empty) (
	*pb.TextParam, error) {

	latestText, err := adapter.svc.Latest()
	response := &pb.TextParam{Text: latestText}
	return response, err
}

func (adapter *EchoServiceGrpcAdapter) Stats(_ context.Context, _ *emptypb.Empty) (
	*pb.EchoStats, error) {

	latest, count := adapter.svc.Stats()
	response := &pb.EchoStats{
		Latest: latest,
		Count:  int32(count),
	}
	return response, nil
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
