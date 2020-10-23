package efogIotEdgeHubServer

import (
	"log"

	zmq "github.com/pebbe/zmq4"
)

// Telemetry Worker structure
type TelemetryWorker struct {
}

// Instantiates a new telemetry worker
func NewTelemetryWorker() *TelemetryWorker {
	worker := new(TelemetryWorker)
	return worker
}

// Starts a telemetry worker
func (worker *TelemetryWorker) Start() {
	pipe, _ := zmq.NewSocket(zmq.PAIR)
	defer pipe.Close()
	pipe.Bind("inproc://pipe")
	for {
		msg, err := pipe.RecvMessage(0)
		if err != nil {
			log.Fatalf("err listening to pipe %q", err)
			break //  Interrupted
		}
		log.Printf("%q\n", msg)
	}
}
