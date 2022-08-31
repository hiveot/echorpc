// Package main with running echo service using grpc or capnp RPC
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hiveot/echorpc/pkg/echoclient"
	"github.com/hiveot/echorpc/pkg/echoservice"
)

const socketAddrGrpc = "/tmp/echoservice-grpc.socket"
const socketAddrCapnp = "/tmp/echoservice-capnp.socket"
const portGrpc = ":8991"
const portCapnp = ":8992"
const payload1kFile = "test/payload-1K.txt"
const payload10kFile = "test/payload-10K.txt"
const payload100kFile = "test/payload-100K.txt"

// Run test with using echo
func main() {
	count := 1000
	payloadHello := "Hello world"
	// echoclient.InvokeUpperDirect(text, count)
	payload1k, err := ioutil.ReadFile(payload1kFile)
	if err != nil {
		log.Fatalf("Failed loading payload: %s", err)
	}
	payload10k, err := ioutil.ReadFile(payload10kFile)
	if err != nil {
		log.Fatalf("Failed loading payload: %s", err)
	}
	payload100k, err := ioutil.ReadFile(payload100kFile)
	if err != nil {
		log.Fatalf("Failed loading payload: %s", err)
	}

	// start all the servers
	go echoservice.EchoServiceGrpcStart(socketAddrGrpc, true)
	go echoservice.EchoServiceCapnpStart(socketAddrCapnp, true)
	go echoservice.EchoServiceGrpcStart(portGrpc, false)
	go echoservice.EchoServiceCapnpStart(portCapnp, false)

	time.Sleep(time.Second)

	fmt.Println("--- test with Hello World payload ---")
	// compare GRPC and Capnproto using unix domain sockets
	echoclient.InvokeUpperGrpc(socketAddrGrpc, true, payloadHello, count)
	echoclient.InvokeUpperCapnp(socketAddrCapnp, true, payloadHello, count)

	// compare GRPC and Capnproto using tcp sockets
	echoclient.InvokeUpperGrpc(portGrpc, false, payloadHello, count)
	echoclient.InvokeUpperCapnp(portCapnp, false, payloadHello, count)

	fmt.Println("--- test with 1K payload ---")
	// compare GRPC and Capnproto using unix domain sockets
	echoclient.InvokeUpperGrpc(socketAddrGrpc, true, string(payload1k), count)
	echoclient.InvokeUpperCapnp(socketAddrCapnp, true, string(payload1k), count)

	// compare GRPC and Capnproto using tcp sockets
	echoclient.InvokeUpperGrpc(portGrpc, false, string(payload1k), count)
	echoclient.InvokeUpperCapnp(portCapnp, false, string(payload1k), count)

	fmt.Println("--- test with 10K payload ---")
	// compare GRPC and Capnproto using unix domain sockets
	echoclient.InvokeUpperGrpc(socketAddrGrpc, true, string(payload10k), count)
	echoclient.InvokeUpperCapnp(socketAddrCapnp, true, string(payload10k), count)

	// compare GRPC and Capnproto using tcp sockets
	echoclient.InvokeUpperGrpc(portGrpc, false, string(payload10k), count)
	echoclient.InvokeUpperCapnp(portCapnp, false, string(payload10k), count)

	fmt.Println("--- test with 100K payload ---")
	// compare GRPC and Capnproto using unix domain sockets
	echoclient.InvokeUpperGrpc(socketAddrGrpc, true, string(payload100k), count)
	echoclient.InvokeUpperCapnp(socketAddrCapnp, true, string(payload100k), count)

	// compare GRPC and Capnproto using tcp sockets
	echoclient.InvokeUpperGrpc(portGrpc, false, string(payload100k), count)
	echoclient.InvokeUpperCapnp(portCapnp, false, string(payload100k), count)
}
