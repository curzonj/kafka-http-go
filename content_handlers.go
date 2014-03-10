package main

import (
	"bufio"
	"encoding/json"
	"github.com/bmizerany/lpx"
	"log"
	"net/http"
)

import kafka "github.com/Shopify/sarama"

type JsonPublishRequest struct {
	Topic    string
	Messages []interface{}
}

func HandleSyslog(w http.ResponseWriter, r *http.Request) {
	bufReader := bufio.NewReader(r.Body)
	reader := lpx.NewReader(bufReader)
	for reader.Next() {
		bytes := reader.Bytes()
		sz := len(bytes)
		newLine := bytes[sz-1]
		if newLine == 10 {
			bytes = bytes[:sz-1]
		}

		// TODO get the topic from the query string
		err := sendMessage("test", nil, kafka.ByteEncoder(bytes))
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func PublishToKafka(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	var v JsonPublishRequest
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
		return
	}

	for _, message := range v.Messages {
		switch vv := message.(type) {
		case string:
			err := sendMessage(v.Topic, nil, kafka.StringEncoder(vv))
			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TODO handle nested JSON messages
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
