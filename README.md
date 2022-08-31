# rpcecho
Comparison of gRPC and capnproto RPC with an echo service

The service has 3 methods, echo, upper and reverse. 

The server for each is defined in capnserver and grpcserver, which contains an adapter that converts the rpc API to the service. 

The client is defined in the testcase itself, using the generated client code.

The testcases uses a unix domain socket for communication between client and server.

To build, run make all
To run the test, run make testgrpc and make testcapnp


# Performance Comparison v1
Initial crude results on an i5-4570S CPU @ 2.90GHz:
 (10% variation on repeated calls)

$ go run pkg/main.go
```
--- test with Hello World payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 880 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 910 millisec
10000 calls using gRPC  on :8991: 1160 millisec
10000 calls using Capnp on :8992: 1270 millisec

--- test with 1K payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 1040 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 1150 millisec
10000 calls using gRPC  on :8991: 1260 millisec
10000 calls using Capnp on :8992: 1690 millisec

--- test with 10K payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 2340 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 2260 millisec
10000 calls using gRPC  on :8991: 2320 millisec
10000 calls using Capnp on :8992: 2650 millisec

--- test with 100K payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 14130 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 12580 millisec
10000 calls using gRPC  on :8991: 14660 millisec
10000 calls using Capnp on :8992: 13130 millisec
```