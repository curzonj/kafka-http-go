package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
    "strconv"
)

import . "github.com/Shopify/sarama"

func main() {
    client, err := NewClient("client_id", []string{"localhost:9092"}, &ClientConfig{MetadataRetries: 1, WaitForElection: 250 * time.Millisecond})
    if err != nil {
          panic(err)
    } else {
          fmt.Println("> connected")
    }
    // defer client.Close()

    producer, err := NewProducer(client, &ProducerConfig{RequiredAcks: WaitForLocal, MaxBufferTime: 1, MaxBufferedBytes: 1})
    if err != nil {
          panic(err)
    }
    // defer producer.Close()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if ("/kafka" == r.URL.Path) {
          t := time.Now()
          timestamp := t.Unix()
          ts_string := strconv.FormatInt(timestamp, 10)

          err = producer.SendMessage("test", nil, StringEncoder(ts_string))
          if err != nil {
                panic(err)
          }
        }
 
        w.WriteHeader(http.StatusNoContent)
    })

    err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)

    if err != nil {
      panic(err)
    } else {
      fmt.Println("listening...")
    }
}
