package main

import (
	"encoding/json"
	kafka "github.com/Shopify/sarama"
	"github.com/curzonj/kafka-http-go/auth"
	"net/http"
	"strconv"
	"strings"
)

// Now this needs to handle more flexible request paths,
// topic metadata for the base and consuming for the sub paths
func HandleTopicConsumers(w http.ResponseWriter, r *http.Request) {
	authenticated := auth.AuthenticateHttpRequest(w, r)
	if authenticated == false {
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	// We'll always have an empty string in the
	// first position and 'topics' in the 2nd
	switch len(parts) {
	default:
		w.WriteHeader(404)
		return
	case 3:
		HandleTopicMetadata(parts[2], w, r)
	case 4:
		partitionId, err := strconv.ParseInt(parts[3], 10, 32)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		HandlePartitionConsumer(parts[2], int32(partitionId), w, r)
	}
}

type batchedConsumerResponse struct {
	HighwaterMark int
	Messages      []*simpleBatchedMessage
}

type simpleBatchedMessage struct {
	Key    string
	Value  string
	Offset int64
}

func HandlePartitionConsumer(topic string, partition int32, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offsetParam := params.Get("offset")
	var offsetMethod kafka.OffsetMethod
	var offset int64

	switch offsetParam {
	case "oldest", "":
		offsetMethod = kafka.OffsetMethodOldest
	case "newest":
		offsetMethod = kafka.OffsetMethodNewest
	default:
		var err error
		offset, err = strconv.ParseInt(offsetParam, 10, 32)
		if err != nil {
			http.Error(w, "Invaled offset query parameter", 400)
			return
		}
	}

	countParam := params.Get("count")
	var messageLimit int = 10
	if countParam != "" {
		messageLimit64, err := strconv.ParseInt(countParam, 10, 32)
		if err != nil {
			http.Error(w, "Invaled count query parameter", 400)
			return
		} else {
			messageLimit = int(messageLimit64)
		}
	}

	consumer, err := kafka.NewConsumer(kafkaClient, topic, partition, "group-name",
		&kafka.ConsumerConfig{
			// TODO This should all be configurable from the params
			MaxWaitTime:  10,
			OffsetMethod: offsetMethod,
			OffsetValue:  offset,
		})

	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer consumer.Close()

	batchedResponse := new(batchedConsumerResponse)
	eventsChan := consumer.Events()

	// TODO switch to the lower level api with explicit fetches
	// so that we don't fetch more than we need to and so we
	// can return the highwatermark and so that we don't wait
	// on the channel if there are no messages currently available
	// TODO support binary keys and values, that causes the
	// json encoder to render them as base64 encoded strings
	for event := range eventsChan {
		batchedResponse.Messages = append(
			batchedResponse.Messages,
			&simpleBatchedMessage{
				Key:    string(event.Key),
				Value:  string(event.Value),
				Offset: event.Offset,
			},
		)

		if len(batchedResponse.Messages) >= messageLimit {
			break
		}
	}

	enc := json.NewEncoder(w)
	enc.Encode(batchedResponse)
}

func HandleTopicMetadata(topic string, w http.ResponseWriter, r *http.Request) {
	metadata, err := buildSimpleMetadata("test")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(metadata)
}

type simplePartitionMetadata struct {
	Id       int32
	Leader   int32
	Replicas []int32
	Isr      []int32
}

type simpleBrokerMetadata struct {
	Id      int32
	Address string
}

type simpleTopicMetadata struct {
	Brokers    []*simpleBrokerMetadata
	Partitions []*simplePartitionMetadata
}

// TODO deal with KError values in the metadata
func buildSimpleMetadata(topicName string) (*simpleTopicMetadata, error) {
	var metadata *simpleTopicMetadata

	// The leader of the first partition will give
	// us full metadata on all partitions in the topic
	leader, err := kafkaClient.Leader(topicName, 0)
	if err != nil {
		return metadata, err
	}

	request := &kafka.MetadataRequest{Topics: []string{topicName}}

	response, err := leader.GetMetadata(CLIENTID, request)
	if err != nil {
		return metadata, err
	}

	metadata = new(simpleTopicMetadata)

	for _, broker := range response.Brokers {
		brokerMetadata := &simpleBrokerMetadata{
			Id:      broker.ID(),
			Address: broker.Addr(),
		}
		metadata.Brokers = append(metadata.Brokers, brokerMetadata)
	}

	// We requested a single topic
	topic := response.Topics[0]

	for _, p := range topic.Partitions {
		partitionMetadata := &simplePartitionMetadata{
			Id:       p.ID,
			Leader:   p.Leader,
			Replicas: p.Replicas,
			Isr:      p.Isr,
		}
		metadata.Partitions = append(metadata.Partitions, partitionMetadata)
	}

	return metadata, nil
}
