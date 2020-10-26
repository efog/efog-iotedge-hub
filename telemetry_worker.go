package efogiotedgehub

import (
	"log"
	"sync"

	zmq "github.com/pebbe/zmq4"
)

// TelemetryWorker Telemetry Worker structure
type TelemetryWorker struct {
	bindEndpoint *string
}

// NewTelemetryWorker Instantiates a new telemetry worker
func NewTelemetryWorker(bindEndpoint string) *TelemetryWorker {
	worker := new(TelemetryWorker)
	worker.bindEndpoint = &bindEndpoint
	return worker
}

// Start Starts a telemetry worker
func (worker *TelemetryWorker) Start() {
	var m sync.Mutex
	pipe, _ := zmq.NewSocket(zmq.PAIR)
	defer pipe.Close()
	pipe.Bind(*worker.bindEndpoint)
	//  Print everything that arrives on pipe
	for {
		m.Lock()
		msg, err := pipe.RecvMessage(0)
		if err != nil {
			break //  Interrupted
		}
		log.Printf("Listened %q", msg)
		m.Unlock()
	}
}
