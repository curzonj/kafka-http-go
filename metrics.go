package main

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
)

var (
	requestMeter      metrics.Meter
	kafkaMessageMeter metrics.Meter
	kafkaMessageTimer metrics.Timer
)

var request_id_chan = make(chan int, 10)

func init() {
	go generateRequestIDs()
	microSecondLogger()

	requestMeter = metrics.NewRegisteredMeter("requests", nil)
	kafkaMessageMeter = metrics.NewRegisteredMeter("kafka.messages", nil)
	kafkaMessageTimer = metrics.NewRegisteredTimer("kafka.sendMessage", nil)

	// TODO build our own metrics printer based on l2met format
	go metrics.Log(metrics.DefaultRegistry, 10e9, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
}

func generateRequestIDs() {
	for i := 0; true; i++ {
		request_id_chan <- i
	}
}

func microSecondLogger() {
	i := log.Flags()
	log.SetFlags(i | log.Lmicroseconds)
}
