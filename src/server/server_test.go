package efogIotEdgeHubServer

import (
	"testing"
)

func TestServerInstantiationWithDefaults(t *testing.T) {
	wantFrontEnd := "tcp://*:5570"
	wantBackEnd := "inproc://backend"
	server := NewServer(nil, nil, nil, nil)
	if server.BackendBindEndpoint != wantBackEnd {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.BackendBindEndpoint, wantBackEnd)
	}
	if server.FrontendBindEndpoint != wantFrontEnd {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.FrontendBindEndpoint, wantFrontEnd)
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

func TestServerCanRun(t *testing.T) {
	// go subscriber_thread("tcp://localhost:56789")
	// go publisher_thread("tcp://localhost:12345")

	wantFrontEndBind := "tcp://*:12345"
	wantFrontEndConnect := "tcp://localhost:12345"
	wantBackEndBind := "tcp://*:56789"
	wantBackEndConnect := "tcp://localhost:56789"
	server := NewServer(&wantBackEndBind, &wantBackEndConnect, &wantFrontEndBind, &wantFrontEndConnect)
	server.Run()
}
