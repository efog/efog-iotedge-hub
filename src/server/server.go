//
//  Asynchronous client-to-server (DEALER to ROUTER).
//
//  While this example runs in a single process, that is just to make
//  it easier to start and stop the example. Each task has its own
//  context and conceptually acts as a separate process.

package efogIotEdgeHubServer

import (
	zmq "github.com/pebbe/zmq4"

	"log"
	"fmt"
	"math/rand"
	"time"
)

const BackendBindDefaultEndpoint = "tcp://*:5571"
const BackendConnectDefaultEndpoint = "tcp://localhost:5571"
const FrontendBindDefaultEndpoint = "tcp://*:5570"
const FrontendConnectDefaultEndpoint = "tcp://localhost:5570"

// Server instance structure which binds frontend with backend workers.
type Server struct {
	BackendBindEndpoint  string
	BackendConnectEndpoint  string
	FrontendBindEndpoint string
	FrontendConnectEndpoint string
}

// Instantiates a new Hub Server
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

func listener_thread() {
	pipe, _ := zmq.NewSocket(zmq.PAIR)
	pipe.Bind("inproc://pipe")

	//  Print everything that arrives on pipe
	for {
		msg, err := pipe.RecvMessage(0)
		if err != nil {
			break //  Interrupted
		}
		log.Printf("received message %q", msg)
	}
}

func subscriber_thread(endpoint *string) {
	//  Subscribe to "A" and "B"
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect(*endpoint)
	subscriber.SetSubscribe("A")
	subscriber.SetSubscribe("B")
	defer subscriber.Close() // cancel subscribe

	for count := 0; count < 5; count++ {
		_, err := subscriber.RecvMessage(0)
		if err != nil {
			break //  Interrupted
		}
	}
}

func publisher_thread(endpoint *string) {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	publisher.Bind(*endpoint)
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%c-%05d", rand.Intn(10)+'A', rand.Intn(100000))
		log.Printf("Sending %q", s)
		_, err := publisher.SendMessage(s)
		if err != nil {
			break //  Interrupted
		}
		time.Sleep(100 * time.Millisecond) //  Wait for 1/10th second
	}
}

// Binds the frontend and backend endpoints with the proxy and telemetry counter module
func (server *Server) Run() {
	// Start the telemetry worker listening on the Pub/Sub
	// telemetryWorker := NewTelemetryWorker()
	// go telemetryWorker.Start()

	// time.Sleep(100 * time.Millisecond)

	// frontend, _ := zmq.NewSocket(zmq.XSUB)
	// defer frontend.Close()
	// feBindErr := frontend.Bind(server.FrontendBindEndpoint)
	// if feBindErr != nil {
	// 	log.Fatalf("Failed to bind frontend: %q", feBindErr)
	// }

	// backend, _ := zmq.NewSocket(zmq.XPUB)
	// defer backend.Close()
	// beBindErr := backend.Connect(server.BackendBindEndpoint)
	// if beBindErr != nil {
	// 	log.Fatalf("Failed to bind backend: %q", beBindErr)
	// }

	// listener, _ := zmq.NewSocket(zmq.PAIR)
	// defer listener.Close()
	// listener.Connect("inproc://pipe")
	// zmq.Proxy(frontend, backend, listener)

	go publisher_thread(&server.FrontendBindEndpoint)
	go subscriber_thread(&server.BackendConnectEndpoint)
	go listener_thread()

	time.Sleep(100 * time.Millisecond)

	subscriber, _ := zmq.NewSocket(zmq.XSUB)
	subscriber.Connect(server.FrontendConnectEndpoint)

	publisher, _ := zmq.NewSocket(zmq.XPUB)
	publisher.Bind(server.BackendBindEndpoint)

	listener, _ := zmq.NewSocket(zmq.PAIR)
	listener.Connect("inproc://pipe")

	zmq.Proxy(subscriber, publisher, listener)

}
