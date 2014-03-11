package main

import (
	"fmt"
	kafka "github.com/Shopify/sarama"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"time"
)

const (
	CLIENTID = "kafka-http"
)

var kafkaProducer *kafka.Producer
var kafkaClient *kafka.Client

// TODO here is where lots of configuration is needed
func connectToKafka() {
	var err error

	kafkaClient, err = kafka.NewClient(CLIENTID, config.KafkaBrokers, &kafka.ClientConfig{MetadataRetries: 1, WaitForElection: 250 * time.Millisecond})
	if err != nil {
		panic(err)
	} else {
		log.Println("at=connectedToKafka")
	}

	kafkaProducer, err = kafka.NewProducer(kafkaClient, &kafka.ProducerConfig{RequiredAcks: kafka.WaitForAll, MaxBufferTime: config.KafkaMaxBufferTime, MaxBufferedBytes: config.KafkaMaxBufferedBytes})
	if err != nil {
		panic(err)
	}

	go handleKafkaAsyncErrors()
}

func handleKafkaAsyncErrors() {
	for {
		// TODO it would be nice if we could track which request
		// caused these errors
		if err := <-kafkaProducer.Errors(); err != nil {
			fmt.Println(err)
		}
	}
}

type kafkaMessageSender func(topic string, key, value kafka.Encoder) error

func queueMessage(topic string, key, value kafka.Encoder) error {
	return handleMessage(topic, key, value, kafkaQueueTimer, kafkaProducer.QueueMessage)
}

func sendMessage(topic string, key, value kafka.Encoder) error {
	return handleMessage(topic, key, value, kafkaSendTimer, kafkaProducer.SendMessage)
}

func handleMessage(topic string, key, value kafka.Encoder, timer metrics.Timer, sender kafkaMessageSender) error {
	var err error

	var valBytes []byte
	if value != nil {
		if valBytes, err = value.Encode(); err != nil {
			return err
		}
	}

	kafkaMessageSizeHist.Update(int64(len(valBytes)))

	timer.Time(func() {
		err = sender(topic, key, kafka.ByteEncoder(valBytes))
	})

	return err
}
