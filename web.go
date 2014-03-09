package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/lpx"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

import kafka "github.com/Shopify/sarama"

type JsonPublishRequest struct {
	Topic    string
	Messages []interface{}
}

var producer *kafka.Producer
var request_id_chan = make(chan int, 10)

func GenerateRequestIDs() {
	for i := 0; true; i++ {
		request_id_chan <- i
	}
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func ConnectToKafka() {
	client, err := kafka.NewClient("client_id", []string{"localhost:9092"}, &kafka.ClientConfig{MetadataRetries: 1, WaitForElection: 250 * time.Millisecond})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("> connected")
	}
	// defer client.Close()

	producer, err = kafka.NewProducer(client, &kafka.ProducerConfig{RequiredAcks: kafka.WaitForLocal, MaxBufferTime: 2, MaxBufferedBytes: 16000})
	if err != nil {
		panic(err)
	}
	// defer producer.Close()
}

func HandleSyslog(w http.ResponseWriter, r *http.Request) {
	numGoroutines := runtime.NumGoroutine()
	request_id := <-request_id_chan
	start_request := time.Now()
	log.Printf("at=received_connection request_id=%d goroutines=%d", request_id, numGoroutines)

	bufReader := bufio.NewReader(r.Body)
	reader := lpx.NewReader(bufReader)
	for reader.Next() {
		bytes := reader.Bytes()
		sz := len(bytes)
		newLine := bytes[sz-1]
		if newLine == 10 {
			bytes = bytes[:sz-1]
		}

		start_message := time.Now()
		log.Printf("at=start_SendMessage request_id=%d", request_id)
		// TODO get the topic from the query string
		err := producer.SendMessage("test", nil, kafka.ByteEncoder(bytes))
		time_since_message := time.Since(start_message).Nanoseconds()
		log.Printf("at=end_SendMessage request_id=%d duration=%d", request_id, time_since_message)

		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	time_since_request := time.Since(start_request).Nanoseconds()
	log.Printf("at=end_connection request_id=%d duration=%d", request_id, time_since_request)
}

func PublishToKafka(w http.ResponseWriter, r *http.Request) {
	numGoroutines := runtime.NumGoroutine()
	request_id := <-request_id_chan
	start_request := time.Now()
	log.Printf("at=received_connection request_id=%d goroutines=%d", request_id, numGoroutines)
	dec := json.NewDecoder(r.Body)

	var v JsonPublishRequest
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
		return
	}

	for _, message := range v.Messages {
		switch vv := message.(type) {
		case string:
			start_message := time.Now()
			log.Printf("at=start_SendMessage request_id=%d", request_id)
			err := producer.SendMessage(v.Topic, nil, kafka.StringEncoder(vv))
			time_since_message := time.Since(start_message).Nanoseconds()
			log.Printf("at=end_SendMessage request_id=%d duration=%d", request_id, time_since_message)

			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO handle nested JSON messages
		}
	}

	w.WriteHeader(http.StatusNoContent)
	time_since_request := time.Since(start_request).Nanoseconds()
	log.Printf("at=end_connection request_id=%d duration=%d", request_id, time_since_request)
}

func MicroSecondLogger() {
	i := log.Flags()
	log.SetFlags(i | log.Lmicroseconds)
}

func main() {
	go GenerateRequestIDs()
	MicroSecondLogger()
	ConnectToKafka()

	http.HandleFunc("/", NoContent)
	http.HandleFunc("/publish", PublishToKafka)
	http.HandleFunc("/syslog", HandleSyslog)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if err != nil {
		panic(err)
	} else {
		log.Println("listening...")
	}
}
