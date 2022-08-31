package echoservice

import (
	"strings"
)

// EchoService demonstrates how to build a microservice for grpc, pub/sub and http invocation.
// The service itself is written independent of any RPC. Protocol adapters map from RPC to this service.
type EchoService struct {
}

// Echo the text without modification, similar to a 'ping'
func (service *EchoService) Echo(text string) (string, error) {
	return text, nil
}

// Upper converts the text to upper case
func (service *EchoService) Upper(text string) (string, error) {
	upper := strings.ToUpper(text)
	return upper, nil
}

// Reverse the text
func (service *EchoService) Reverse(text string) (string, error) {
	rns := []rune(text)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns), nil
}

// NewEchoService creates and registers the service with gRPC interface
// onShutDown callback is used to handle stop request
func NewEchoService() *EchoService {
	service := &EchoService{}
	return service
}
