# rpcecho
Comparison of gRPC and capnproto RPC with an echo service

The service has 3 methods, echo, upper and reverse. 

The server for each is defined in capnserver and grpcserver, which contains an adapter that converts the rpc API to the service. 

The client is defined in the testcase itself, using the generated client code.

The testcases uses a unix domain socket for communication between client and server.

To build, run make all
To run the test, run make testgrpc and make testcapnp
