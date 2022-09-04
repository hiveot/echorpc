package echoservice

// EchoService demonstrates how to build a microservice for grpc, pub/sub and http invocation.
// The service itself is written independent of any RPC. Protocol adapters map from RPC to this service.
type EchoService struct {
	latestText string
	count      int
}

// Echo the text without modification, similar to a 'ping'
func (service *EchoService) Echo(text string) (string, error) {
	// todo: handle concurrency
	service.latestText = text
	service.count++
	return text, nil
}

// Latest returns the latest echo text
func (service *EchoService) Latest() (string, error) {
	return service.latestText, nil
}

// Stats returns echo statistics
func (service *EchoService) Stats() (latest string, count int) {
	return service.latestText, service.count
}

// NewEchoService creates and registers the service with gRPC interface
// onShutDown callback is used to handle stop request
func NewEchoService() *EchoService {
	service := &EchoService{}
	return service
}
