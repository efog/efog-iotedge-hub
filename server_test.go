package efogiotedgehub

import (
	"testing"
)

func TestServerInstantiationWithDefaults(t *testing.T) {
	server := NewServer(nil, nil)
	if server.SubscriberEndpoint != SubscriberDefaultEndpoint {
		t.Errorf("NewServer.SubscriberDefaultEndpoint = %q, want %q", server.SubscriberEndpoint, SubscriberDefaultEndpoint)
	}
	if server.PublisherEndpoint != PublisherDefaultEndpoint {
		t.Errorf("NewServer.PublisherEndpoint = %q, want %q", server.PublisherEndpoint, PublisherDefaultEndpoint)
	}
}

func TestServerInstantiationWithValues(t *testing.T) {
	wantPublisherEndpoint := "tcp://*:56789"
	wantSubscriberEndpoint := "tcp://*:56789"
	server := NewServer(&wantSubscriberEndpoint, &wantPublisherEndpoint)
	if server.SubscriberEndpoint != wantPublisherEndpoint {
		t.Errorf("NewServer.SubscriberEndpoint = %q, want %q", server.SubscriberEndpoint, wantPublisherEndpoint)
	}
	if server.PublisherEndpoint != wantPublisherEndpoint {
		t.Errorf("NewServer.PublisherEndpoint = %q, want %q", server.PublisherEndpoint, wantPublisherEndpoint)
	}
}
