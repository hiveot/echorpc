# rpcecho
This is a crude microbenchmark that compares 'HTTP', gRPC and Cap'n proto RPCs. Using a simple echo service. The service itself is very simple (echo) in order to measure the performance of the RPC itself.   

The server receives requests for each RPC method using an adapter for that RPC - EchoServiceCapnAdapter, EchoServiceGRPCAdapter and EchoServiceHTTPAdapter - which converts the arguments and response from the native service into the format used by the protocol. The HTTP adapter converts to JSON, the gRPC adapter creates a  protobuf message, and the Capnp adapter .. gets/sets the parameters using methods.

The clients create a connection and invokes the methods. Creating the connection supports both unix domain sockets and TCP connections. TCP connections carry more overhead and are expected (correctly) to be slower.

To build, run make all
To run the test, run make testgrpc and make testcapnp


# Performance Comparisons

## 2022-09-04- added http and a second call to retrieve stats

Added a get stats method that returns the latest echo and nr of calls. Each call below is therefore two calls, one to post echo and another too retrieve stats.
Unexpectedly, http is fastest for small payloads. 

```
--- test with Hello World payload ---
Invoking echo directly
1000 calls using direct call: 2 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 110 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 150 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 190 millisec
1000 calls using http  on :8990: 180 millisec
1000 calls using gRPC  on :8991: 220 millisec
1000 calls using Capnp on :8992: 280 millisec
--- test with 1K payload ---
Invoking echo directly
1000 calls using direct call: 2 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 190 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 190 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 210 millisec
1000 calls using http  on :8990: 200 millisec
1000 calls using gRPC  on :8991: 230 millisec
1000 calls using Capnp on :8992: 320 millisec
--- test with 10K payload ---
Invoking echo directly
1000 calls using direct call: 2 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 570 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 250 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 250 millisec
1000 calls using http  on :8990: 730 millisec
1000 calls using gRPC  on :8991: 310 millisec
1000 calls using Capnp on :8992: 350 millisec
--- test with 100K payload ---
Invoking echo directly
1000 calls using direct call: 3 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 4410 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 820 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 590 millisec
1000 calls using http  on :8990: 4500 millisec
1000 calls using gRPC  on :8991: 840 millisec
1000 calls using Capnp on :8992: 710 millisec
```

## 2022-08-31
Initial crude results on an i5-4570S CPU @ 2.90GHz:
 (10% variation on repeated calls)

$ go run pkg/main.go
```
--- test with Hello World payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 950 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 950 millisec
10000 calls using gRPC  on :8991: 1150 millisec
10000 calls using Capnp on :8992: 1430 millisec
--- test with 1K payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 870 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 1140 millisec
10000 calls using gRPC  on :8991: 1260 millisec
10000 calls using Capnp on :8992: 1670 millisec
--- test with 10K payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 1410 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 1310 millisec
10000 calls using gRPC  on :8991: 1170 millisec
10000 calls using Capnp on :8992: 1690 millisec
--- test with 100K payload ---
10000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 4670 millisec
10000 calls using Capnp on /tmp/echoservice-capnp.socket: 3270 millisec
10000 calls using gRPC  on :8991: 5190 millisec
10000 calls using Capnp on :8992: 3850 millisec

```
