package main

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
)

var (
	requestMeter         metrics.Meter
	kafkaSendTimer       metrics.Timer
	kafkaQueueTimer      metrics.Timer
	kafkaMessageSizeHist metrics.Histogram
)

var request_id_chan = make(chan int, 10)

/*
* This stuff is all isolated here so that we can stub out
* the metrics with nil providers without touching all the
* code everywhere.
 */
func init() {
	go generateRequestIDs()
	microSecondLogger()

	requestMeter = metrics.NewRegisteredMeter("requests", nil)
	kafkaSendTimer = metrics.NewRegisteredTimer("kafka.sendMessage", nil)
	kafkaQueueTimer = metrics.NewRegisteredTimer("kafka.queueMessage", nil)

	s := metrics.NewExpDecaySample(1028, 0.015)
	kafkaMessageSizeHist = metrics.NewRegisteredHistogram("kafka.messageSize", nil, s)

	// TODO build our own metrics printer based on l2met format
	go metrics.Log(metrics.DefaultRegistry, 60e9, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
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
