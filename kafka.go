package main

import (
	"log"
	"time"
)

import kafka "github.com/Shopify/sarama"

var producer *kafka.Producer

// TODO here is where lots of configuration is needed
func connectToKafka() {
	client, err := kafka.NewClient("client_id", []string{"localhost:9092"}, &kafka.ClientConfig{MetadataRetries: 1, WaitForElection: 250 * time.Millisecond})
	if err != nil {
		panic(err)
	} else {
		log.Println("at=connectedToKafka")
	}

	producer, err = kafka.NewProducer(client, &kafka.ProducerConfig{RequiredAcks: kafka.WaitForLocal, MaxBufferTime: 2, MaxBufferedBytes: 16000})
	if err != nil {
		panic(err)
	}
}

func queueMessage(topic string, key, value kafka.Encoder) error {
	return producer.QueueMessage(topic, key, value)
}

func sendMessage(topic string, key, value kafka.Encoder) error {
	var err error

	kafkaMessageMeter.Mark(1)
	kafkaMessageTimer.Time(func() {
		err = producer.SendMessage(topic, key, value)
	})

	return err
}
