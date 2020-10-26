//
//  Asynchronous client-to-server (DEALER to ROUTER).
//
//  While this example runs in a single process, that is just to make
//  it easier to start and stop the example. Each task has its own
//  context and conceptually acts as a separate process.

package efogiotedgehub

import (
	zmq "github.com/pebbe/zmq4"

	"time"
)

// BackendBindDefaultEndpoint default value for backend bind endpoint
const BackendBindDefaultEndpoint = "tcp://*:5571"
// BackendConnectDefaultEndpoint default value for backend connect endpoint
const BackendConnectDefaultEndpoint = "tcp://localhost:5571"
// FrontendBindDefaultEndpoint default value for frontend bind endpoint
const FrontendBindDefaultEndpoint = "tcp://*:5570"
// FrontendConnectDefaultEndpoint default value for frontend connect endpoint
const FrontendConnectDefaultEndpoint = "tcp://localhost:5570"

// Server instance structure which binds frontend with backend workers.
type Server struct {
	BackendBindEndpoint     string
	BackendConnectEndpoint  string
	FrontendBindEndpoint    string
	FrontendConnectEndpoint string
}

// NewServer Instantiates a new Hub Server
func NewServer(backendBindEndpoint *string, backendConnectEndpoint *string, frontendBindEndpoint *string, frontendConnectEndpoint *string) *Server {
	server := new(Server)

	if backendBindEndpoint != nil {
		server.BackendBindEndpoint = *backendBindEndpoint
	} else {
		server.BackendBindEndpoint = BackendBindDefaultEndpoint
	}
	if backendConnectEndpoint != nil {
		server.BackendConnectEndpoint = *backendConnectEndpoint
	} else {
		server.BackendConnectEndpoint = BackendConnectDefaultEndpoint
	}

	if frontendBindEndpoint != nil {
		server.FrontendBindEndpoint = *frontendBindEndpoint
	} else {
		server.FrontendBindEndpoint = FrontendBindDefaultEndpoint
	}
	if frontendConnectEndpoint != nil {
		server.FrontendConnectEndpoint = *frontendConnectEndpoint
	} else {
		server.FrontendConnectEndpoint = FrontendConnectDefaultEndpoint
	}

	return server
}

// Run binds the frontend and backend endpoints with the proxy and telemetry counter module
func (server *Server) Run() {
	telemetryWorker := NewTelemetryWorker()
	go telemetryWorker.Start()
	time.Sleep(100 * time.Millisecond)
	subscriber, _ := zmq.NewSocket(zmq.XSUB)
	subscriber.Connect(server.FrontendConnectEndpoint)
	publisher, _ := zmq.NewSocket(zmq.XPUB)
	publisher.Bind(server.BackendBindEndpoint)
	listener, _ := zmq.NewSocket(zmq.PAIR)
	listener.Connect("inproc://pipe")
	zmq.Proxy(subscriber, publisher, listener)
}
