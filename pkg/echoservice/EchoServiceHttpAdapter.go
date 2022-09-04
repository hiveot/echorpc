package echoservice

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// EchoServiceHttpAdapter gRPC Adapter for echo service
type EchoServiceHttpAdapter struct {
	httpServer *http.Server
	svc        *EchoService
}

func (adapter *EchoServiceHttpAdapter) HandleEcho(resp http.ResponseWriter, req *http.Request) {
	var params interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = json.Unmarshal(body, params)
	}
	if err != nil {
		http.Error(resp, "bad payload", http.StatusBadRequest)
		return
	}
	echoText, err := adapter.svc.Echo(string(body))
	respMsg, _ := json.Marshal(echoText)
	resp.Write(respMsg)
}

func (adapter *EchoServiceHttpAdapter) HandleLatest(resp http.ResponseWriter, req *http.Request) {

	latestText, _ := adapter.svc.Latest()
	respMsg, _ := json.Marshal(latestText)
	resp.Write(respMsg)
}

func (adapter *EchoServiceHttpAdapter) HandleStats(resp http.ResponseWriter, req *http.Request) {
	params := make(map[string]interface{})
	latestText, count := adapter.svc.Stats()
	params["latest"] = latestText
	params["count"] = count
	respMsg, _ := json.Marshal(params)
	resp.Write(respMsg)
}

// EchoServiceHttpStart start listening
//  address to list on: ":port", "host:port", "/tmp/path-to.socket"
//  isUDS set to true when address is a unix domain socket
func EchoServiceHttpStart(address string, isUDS bool) {

	network := "tcp"
	if isUDS {
		os.Remove(address)
		network = "unix"
	}
	router := mux.NewRouter()
	httpServer := &http.Server{Handler: router}
	echoSvc := NewEchoService()
	adapter := &EchoServiceHttpAdapter{
		httpServer: httpServer,
		svc:        echoSvc,
	}
	router.HandleFunc("/echo", adapter.HandleEcho)
	router.HandleFunc("/stats", adapter.HandleStats)
	router.HandleFunc("/latest", adapter.HandleLatest)

	listener, err := net.Listen(network, address)
	if err != nil {
		logrus.Fatalf("Failed open listener: %v", err)
	}
	httpServer.Serve(listener)
}
