package main

import (
	"log"
	"net/http"
)

var config *Config

func init() {
	config = new(Config)
	config.Configure()
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	connectToKafka()

	// TODO handle these funcs via a wrapper
	// so that not every handler has to deal
	// with authentication and request metrics

	http.HandleFunc("/health", HealthCheck)
	http.HandleFunc("/publish", HandlePublishByContentType)
	http.HandleFunc("/topics/", HandleTopicConsumers)

	err := http.ListenAndServe(":"+config.HttpPort, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
