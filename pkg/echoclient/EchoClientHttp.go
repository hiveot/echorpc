package echoclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// InvokeEchoHttp Invokes the echo service using http
func InvokeEchoHttp(address string, isUDS bool, text string, count int) {
	// Set up a connection to the server.
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	prefix := ""
	network := "tcp"
	if isUDS {
		//address = "unix://" + address
		prefix = "http://unix"
		network = "unix"
	} else {
		//address = "http://" + address
		prefix = "http://" + address
		network = "tcp"
	}
	// use net.Dial to support both network sockets and unix domain sockets
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(network, address)
			},
		},
	}

	t1 := time.Now()
	for i := 0; i < count; i++ {
		// post an echo
		echoAddr := prefix + "/echo"
		payload, _ := json.Marshal(text)
		response, err := client.Post(echoAddr, "", bytes.NewBuffer(payload))
		if err != nil {
			log.Fatalf("Error posting to '%s': %v", echoAddr, err)
		}
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			log.Fatalf("Error reading response from '%s': %v", echoAddr, err)
		}
		var respText string
		err = json.Unmarshal(body, &respText)
		//fmt.Printf("response: %s\n", respText)

		// get stats -> map[]
		statsAddr := prefix + "/stats"
		response, err = client.Get(statsAddr)
		if err != nil {
			log.Fatalf("Error get from '%s': %v", statsAddr, err)
		}
		body, err = ioutil.ReadAll(response.Body)
		response.Body.Close()

		var respMsg map[string]interface{}
		err = json.Unmarshal(body, &respMsg)
		//fmt.Printf("response: %v\n", respMsg)
		//_ = respMsg

	}
	d1 := time.Since(t1)
	msec := d1.Milliseconds() / 10 * 10
	fmt.Printf("%d calls using http  on %s: %d millisec\n", count, address, msec)

	// return response.Text, err
}
