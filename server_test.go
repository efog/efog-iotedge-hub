package efogiotedgehub

import (
	"fmt"
	"log"
	"math/rand"

	zmq "github.com/pebbe/zmq4"

	"testing"
	"time"
)

func TestServerInstantiationWithDefaults(t *testing.T) {
	server := NewServer(nil, nil, nil, nil)
	if server.BackendBindEndpoint != BackendBindDefaultEndpoint {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.BackendBindEndpoint, BackendBindDefaultEndpoint)
	}
	if server.FrontendBindEndpoint != FrontendBindDefaultEndpoint {
		t.Errorf("NewServer.FrontendBindEndpoint = %q, want %q", server.FrontendBindEndpoint, FrontendBindDefaultEndpoint)
	}
}

func TestServerInstantiationWithValues(t *testing.T) {
	wantFrontEndBind := "tcp://*:12345"
	wantBackEndBind := "tcp://*:56789"
	wantFrontEndConnect := "tcp://localhost:12345"
	wantBackEndConnect := "tcp://localhost:56789"
	server := NewServer(&wantBackEndBind, &wantBackEndConnect, &wantFrontEndBind, &wantFrontEndConnect)
	if server.BackendBindEndpoint != wantBackEndBind {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.BackendBindEndpoint, wantBackEndBind)
	}
	if server.BackendConnectEndpoint != wantBackEndConnect {
		t.Errorf("NewServer.BackendConnectEndpoint = %q, want %q", server.BackendConnectEndpoint, wantBackEndConnect)
	}
	if server.FrontendBindEndpoint != wantFrontEndBind {
		t.Errorf("NewServer.FrontendBindEndpoint = %q, want %q", server.FrontendBindEndpoint, wantFrontEndBind)
	}
	if server.FrontendConnectEndpoint != wantFrontEndConnect {
		t.Errorf("NewServer.FrontendConnectEndpoint = %q, want %q", server.FrontendConnectEndpoint, wantFrontEndConnect)
	}
}

func testSubscriberThread(endpoint *string) {
	//  Subscribe to "A" and "B"
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect(*endpoint)
	subscriber.SetSubscribe("A")
	subscriber.SetSubscribe("B")
	subscriber.SetSubscribe("C")
	defer subscriber.Close() // cancel subscribe
	for i := 0; i < 10; i++ {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			break //  Interrupted
		}
		log.Printf("Received %q", msg)
	}
}

func testPublisherThread(endpoint *string) {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	publisher.Bind(*endpoint)
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%c-%05d", rand.Intn(3)+'A', rand.Intn(100000))
		log.Printf("Sending %q", s)
		_, err := publisher.SendMessage(s)
		if err != nil {
			break //  Interrupted
		}
		time.Sleep(1000 * time.Millisecond) //  Wait for 1/10th second
	}
}

func TestServerCanRun(t *testing.T) {
	wantFrontEndBind := "tcp://*:12345"
	wantFrontEndConnect := "tcp://localhost:12345"
	wantBackEndBind := "tcp://*:56789"
	wantBackEndConnect := "tcp://localhost:56789"
	server := NewServer(&wantBackEndBind, &wantBackEndConnect, &wantFrontEndBind, &wantFrontEndConnect)
	go server.Run()
	go testPublisherThread(&wantFrontEndBind)
	go testSubscriberThread(&wantBackEndConnect)
	time.Sleep(11000 * time.Millisecond)
}
