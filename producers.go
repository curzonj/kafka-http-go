package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/lpx"
	"github.com/curzonj/kafka-http-go/auth"
	"log"
	"net/http"
	"net/url"
)

import kafka "github.com/Shopify/sarama"

type JsonPublishRequest struct {
	Topic    string
	Messages []interface{}
}

type HttpError struct {
	StatusMessage string
	StatusCode    int
	cause         error
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%v: %v", e.StatusMessage, e.StatusCode)
}

func HandlePublishByContentType(w http.ResponseWriter, r *http.Request) {
	requestMeter.Mark(1)

	if r.Method != "POST" {
		http.Error(w, "Method must be POST.", 400)
		return
	}

	authenticated := auth.AuthenticateHttpRequest(w, r)
	if authenticated == false {
		return
	}

	var handler func(*http.Request) error

	switch r.Header.Get("Content-Type") {
	case "application/logplex-1":
		handler = SyslogRequestToKafka
	case "application/json":
		handler = JSONRequestToKafka
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if err := handler(r); err != nil {
		switch err := err.(type) {
		case HttpError:
			http.Error(w, err.StatusMessage, err.StatusCode)
			if err.cause != nil {
				log.Println(err.cause)
			}
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func SyslogRequestToKafka(r *http.Request) error {
	params := r.URL.Query()
	topic := params.Get("topic")
	if topic == "" {
		return &HttpError{"topic query parameter required", http.StatusBadRequest, nil}
	}

	bufReader := bufio.NewReader(r.Body)
	reader := lpx.NewReader(bufReader)
	messageHandler := getMessageHandler(r.URL)

	for reader.Next() {
		line := reader.Bytes()
		line = bytes.Trim(line, "\n")

		err := messageHandler(topic, nil, kafka.ByteEncoder(line))
		if err != nil {
			return err
		}
	}

	return nil
}

func JSONRequestToKafka(r *http.Request) error {
	dec := json.NewDecoder(r.Body)
	messageHandler := getMessageHandler(r.URL)

	var v JsonPublishRequest
	if err := dec.Decode(&v); err != nil {
		return &HttpError{"Invalid json request body", http.StatusBadRequest, err}
	}

	params := r.URL.Query()
	// query param will override json content to present
	// a consistent api
	topic := params.Get("topic")
	if topic != "" {
		v.Topic = topic
	}

	for _, message := range v.Messages {
		switch vv := message.(type) {
		case string:
			err := messageHandler(v.Topic, nil, kafka.StringEncoder(vv))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getMessageHandler(url *url.URL) kafkaMessageSender {
	params := url.Query()
	if params.Get("ack") == "true" {
		return sendMessage
	} else {
		return queueMessage
	}
}
