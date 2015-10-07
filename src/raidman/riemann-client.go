package main

import (
	"flag"
	"os"
	"time"

	"github.com/amir/raidman"
)

func main() {

	var riemannHost string
	var riemannPort string
	var eventState string
	var eventMessage string
	var eventHost string
	var eventService string
	var eventMetric float64

	hostname, _ := os.Hostname()

	flag.StringVar(&riemannHost, "riemannHost", "localhost", "Riemann hostname")
	flag.StringVar(&riemannPort, "riemannPort", "5555", "Riemann listen port")
	flag.StringVar(&eventState, "eventState", "ok", "event state: ok or critical")
	flag.StringVar(&eventMessage, "eventMessage", "ping", "event description to send to Riemann")
	flag.StringVar(&eventHost, "eventHost", hostname, "hostname of the server that send the event")
	flag.StringVar(&eventService, "eventService", "cf-operational-tools", "name of the service")
	flag.Float64Var(&eventMetric, "eventMetric", 1.0, "metric for the event")
	flag.Parse()

	c, err := raidman.Dial("tcp", riemannHost+":"+riemannPort)
	if err != nil {
		panic(err)
	}

	var event = &raidman.Event{
		State:       eventState,
		Service:     eventService,
		Metric:      eventMetric,
		Description: eventMessage,
		Host:        eventHost,
		Time: 			 time.Now().Unix(),
	}

	err = c.Send(event)
	if err != nil {
		panic(err)
	}

	events, err := c.Query("service = \"" + eventService + "\"")
	if err != nil {
		panic(err)
	}

	if len(events) < 1 {
		panic("Submitted event not found")
	}

	c.Close()
}
