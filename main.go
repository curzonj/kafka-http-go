package main

import (
	"github.com/curzonj/kafka-http/auth"
	"log"
	"net/http"
)

var config *Config

func NoContent(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header.Get("Content-Type"))

	w.WriteHeader(http.StatusNoContent)
}

func PublishByContentType(w http.ResponseWriter, r *http.Request) {
	requestMeter.Mark(1)

	if r.Method != "POST" {
		http.Error(w, "Method must be POST.", 400)
		return
	}

	authenticated := auth.AuthenticateHttpRequest(w, r)
	if authenticated == false {
		return
	}

	switch r.Header.Get("Content-Type") {
	case "application/logplex-1":
		HandleSyslog(w, r)
	case "application/json":
		PublishToKafka(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}

func main() {
	config = new(Config)

	connectToKafka()

	http.HandleFunc("/", NoContent)
	http.HandleFunc("/publish", PublishByContentType)

	err := http.ListenAndServe(":"+config.HttpPort(), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
