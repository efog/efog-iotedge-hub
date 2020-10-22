package efogIotEdgeHubServer

import ( "testing" )

func TestServerInstantiationWithDefaults(t *testing.T) {
	wantFrontEnd := "tcp://*:5570"
	wantBackEnd := "inproc://backend"
	server := NewServer(nil, nil)
	if server.BackendBindEndpoint != wantBackEnd {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.BackendBindEndpoint, wantBackEnd)
	}
	if server.FrontendBindEndpoint != wantFrontEnd {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.FrontendBindEndpoint, wantFrontEnd)
	}
}
func TestServerInstantiationWithValues(t *testing.T) {
	wantFrontEnd := "tcp://*:12345"
	wantBackEnd := "tcp://*:56789"
	server := NewServer(&wantBackEnd, &wantFrontEnd)
	if server.BackendBindEndpoint != wantBackEnd {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.BackendBindEndpoint, wantBackEnd)
	}
	if server.FrontendBindEndpoint != wantFrontEnd {
		t.Errorf("NewServer.BackendBindEndpoint = %q, want %q", server.FrontendBindEndpoint, wantFrontEnd)
	}
}

func TestServerCanBind(t *testing.T) {
	wantFrontEnd := "tcp://*:12345"
	wantBackEnd := "tcp://*:56789"
	server := NewServer(&wantBackEnd, &wantFrontEnd)
	server.Bind()
}