//
//  Asynchronous client-to-server (DEALER to ROUTER).
//
//  While this example runs in a single process, that is just to make
//  it easier to start and stop the example. Each task has its own
//  context and conceptually acts as a separate process.

package efogiotedgehub

import (
	"log"

	zmq "github.com/pebbe/zmq4"
)

// SubscriberDefaultEndpoint default value for subscriber endpoint
const SubscriberDefaultEndpoint = "tcp://*:5571"

// PublisherDefaultEndpoint default value for publisher endpoint
const PublisherDefaultEndpoint = "tcp://localhost:5571"

// ListenerConnectDefaultEndpoint default value for listener connect endpoint
const ListenerConnectDefaultEndpoint = "inproc://listener"

// Server instance structure which binds frontend with backend workers.
type Server struct {
	SubscriberEndpoint string
	PublisherEndpoint  string
}

// NewServer Instantiates a new Hub Server
func NewServer(subscriberEndpoint *string, publisherEndpoint *string) *Server {
	server := new(Server)

	if subscriberEndpoint != nil {
		server.SubscriberEndpoint = *subscriberEndpoint
	} else {
		server.SubscriberEndpoint = SubscriberDefaultEndpoint
	}
	if publisherEndpoint != nil {
		server.PublisherEndpoint = *publisherEndpoint
	} else {
		server.PublisherEndpoint = PublisherDefaultEndpoint
	}

	return server
}

// Run binds the frontend and backend endpoints with the proxy and telemetry counter module
func (server *Server) Run() {

	telemetryWorker := NewTelemetryWorker(ListenerConnectDefaultEndpoint)
	go telemetryWorker.Start()

	log.Print("Starting xsub server")
	subscriber, _ := zmq.NewSocket(zmq.XSUB)
	defer subscriber.Close()
	subscriber.Bind(server.SubscriberEndpoint)
	log.Printf("Binded subscriber to endpoint %s", server.SubscriberEndpoint)

	log.Print("Starting xpub server")
	publisher, _ := zmq.NewSocket(zmq.XPUB)
	defer publisher.Close()
	publisher.Bind(server.PublisherEndpoint)
	log.Printf("Binded publisher to endpoint %s", server.PublisherEndpoint)

	log.Print("Starting paired listener")
	listener, _ := zmq.NewSocket(zmq.PAIR)
	defer listener.Close()
	listener.Connect("inproc://listener")
	log.Print("Started paired listener")

	log.Print("Starting up ZMQ Proxy")
	zmq.Proxy(subscriber, publisher, listener)
}
