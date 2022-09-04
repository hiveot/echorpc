// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package echo

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EchoServiceClient is the client API for EchoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EchoServiceClient interface {
	// Return an copy of the text
	Echo(ctx context.Context, in *TextParam, opts ...grpc.CallOption) (*TextParam, error)
	// Return the upper case converted text
	Latest(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TextParam, error)
	// Return the echo statistics
	Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*EchoStats, error)
}

type echoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEchoServiceClient(cc grpc.ClientConnInterface) EchoServiceClient {
	return &echoServiceClient{cc}
}

func (c *echoServiceClient) Echo(ctx context.Context, in *TextParam, opts ...grpc.CallOption) (*TextParam, error) {
	out := new(TextParam)
	err := c.cc.Invoke(ctx, "/echo.EchoService/echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) Latest(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TextParam, error) {
	out := new(TextParam)
	err := c.cc.Invoke(ctx, "/echo.EchoService/latest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*EchoStats, error) {
	out := new(EchoStats)
	err := c.cc.Invoke(ctx, "/echo.EchoService/stats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServiceServer is the server API for EchoService service.
// All implementations must embed UnimplementedEchoServiceServer
// for forward compatibility
type EchoServiceServer interface {
	// Return an copy of the text
	Echo(context.Context, *TextParam) (*TextParam, error)
	// Return the upper case converted text
	Latest(context.Context, *emptypb.Empty) (*TextParam, error)
	// Return the echo statistics
	Stats(context.Context, *emptypb.Empty) (*EchoStats, error)
	mustEmbedUnimplementedEchoServiceServer()
}

// UnimplementedEchoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEchoServiceServer struct {
}

func (UnimplementedEchoServiceServer) Echo(context.Context, *TextParam) (*TextParam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedEchoServiceServer) Latest(context.Context, *emptypb.Empty) (*TextParam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Latest not implemented")
}
func (UnimplementedEchoServiceServer) Stats(context.Context, *emptypb.Empty) (*EchoStats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedEchoServiceServer) mustEmbedUnimplementedEchoServiceServer() {}

// UnsafeEchoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EchoServiceServer will
// result in compilation errors.
type UnsafeEchoServiceServer interface {
	mustEmbedUnimplementedEchoServiceServer()
}

func RegisterEchoServiceServer(s grpc.ServiceRegistrar, srv EchoServiceServer) {
	s.RegisterService(&EchoService_ServiceDesc, srv)
}

func _EchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.EchoService/echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Echo(ctx, req.(*TextParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoService_Latest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Latest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.EchoService/latest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Latest(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoService_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.EchoService/stats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Stats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// EchoService_ServiceDesc is the grpc.ServiceDesc for EchoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EchoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "echo.EchoService",
	HandlerType: (*EchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "echo",
			Handler:    _EchoService_Echo_Handler,
		},
		{
			MethodName: "latest",
			Handler:    _EchoService_Latest_Handler,
		},
		{
			MethodName: "stats",
			Handler:    _EchoService_Stats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "echo.proto",
}
