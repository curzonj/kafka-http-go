package main

import (
	"encoding/json"
	kafka "github.com/Shopify/sarama"
	"github.com/curzonj/kafka-http-go/auth"
	"net/http"
)

// Now this needs to handle more flexible request paths,
// topic metadata for the base and consuming for the sub paths
func HandleTopicConsumers(w http.ResponseWriter, r *http.Request) {
	authenticated := auth.AuthenticateHttpRequest(w, r)
	if authenticated == false {
		return
	}

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
