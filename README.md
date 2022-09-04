# rpcecho
This is a crude microbenchmark that compares 'HTTP', gRPC and Cap'n proto RPCs. Using a simple echo service. The service itself is very simple (echo) in order to measure the performance of the RPC itself.   

The server receives requests for each RPC method using an adapter for that RPC - EchoServiceCapnAdapter, EchoServiceGRPCAdapter and EchoServiceHTTPAdapter - which converts the arguments and response from the native service into the format used by the protocol. The HTTP adapter converts to JSON, the gRPC adapter creates a  protobuf message, and the Capnp adapter .. gets/sets the parameters using methods.

The clients create a connection and invokes the methods. Creating the connection supports both unix domain sockets and TCP connections. TCP connections carry more overhead and are expected (correctly) to be slower.

To build, run make all
To run the test, run make testgrpc and make testcapnp


# Performance Comparisons
Microbenchmark on an i5-4570S CPU @ 2.90GHz

## 2022-09-04- added http-1.1 and a second call to retrieve stats

Added a get stats method that returns the latest echo and nr of calls. Each call below is therefore two calls, one to post echo and another to retrieve stats.

Unexpectedly, http (-1.1, no tls) is significantly faster for small payloads. The RPC's seem to have more overhead in invoking the method compared to marshalling the payload.


```
--- test with Hello World payload ---
Invoking echo directly
Invoking echo directly
1000 calls using direct call: 5 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 160 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 180 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 180 millisec
1000 calls using http  on :8990: 170 millisec
1000 calls using gRPC  on :8991: 220 millisec
1000 calls using Capnp on :8992: 250 millisec
--- test with 1K payload ---
Invoking echo directly
1000 calls using direct call: 1 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 180 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 170 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 190 millisec
1000 calls using http  on :8990: 220 millisec
1000 calls using gRPC  on :8991: 180 millisec
1000 calls using Capnp on :8992: 270 millisec
--- test with 10K payload ---
Invoking echo directly
1000 calls using direct call: 2 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 670 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 250 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 230 millisec
1000 calls using http  on :8990: 760 millisec
1000 calls using gRPC  on :8991: 310 millisec
1000 calls using Capnp on :8992: 310 millisec
--- test with 100K payload ---
Invoking echo directly
1000 calls using direct call: 2 microsec
1000 calls using http  on /tmp/echoservice-http.socket: 4340 millisec
1000 calls using gRPC  on unix:///tmp/echoservice-grpc.socket: 780 millisec
1000 calls using Capnp on /tmp/echoservice-capnp.socket: 580 millisec
1000 calls using http  on :8990: 4660 millisec
1000 calls using gRPC  on :8991: 830 millisec
1000 calls using Capnp on :8992: 660 millisec
```

## 2022-08-31 - initial comparison of 'echo' using gRPC and Capnp 

This test runs 'echo' <payload> using gRPC and Capnp using various payload sizes.
Capnp shows its marshalling advantage at higher payload but seems to have more overhead that hurts performance at smaller payloads.

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
