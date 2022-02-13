package integratedkafka

import (
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/transport"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var kafkaProducerCreate, kafkaProducerUpdate, kafkaProducerDelete *kafka.Producer

// KafkaInit - KafkaInit
func KafkaInit(kafkaURL string) {
	//UseKafka = true
	producer, err := transport.InitProducer(kafkaURL)
	if err != nil {
		panic(err)
	}
	kafkaProducerCreate = producer

	producer1, err := transport.InitProducer(kafkaURL)
	if err != nil {
		panic(err)
	}
	kafkaProducerUpdate = producer1

	producer2, err := transport.InitProducer(kafkaURL)
	if err != nil {
		panic(err)
	}
	kafkaProducerDelete = producer2
}
