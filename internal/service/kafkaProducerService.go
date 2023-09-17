package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"golang-kafka/internal/entity"
	"os"
	"strings"
)

type Service struct {
	ProduceService
}

func NewService(service ProduceService) *Service {
	return &Service{
		ProduceService: service,
	}
}

type ProduceService interface {
	Produce(fio entity.Fio) error
	Shutdown()
}

type KafkaProduceService struct {
	topicName  string
	configFile string
	producer   *kafka.Producer
}

func NewKafkaProduceService(topicName, configFile string) (*KafkaProduceService, error) {
	kafkaService := &KafkaProduceService{
		topicName:  topicName,
		configFile: configFile,
		producer:   nil,
	}
	configMap, err := kafkaService.readConfig()
	if err != nil {
		return nil, err
	}
	kafkaService.producer, err = kafka.NewProducer(&configMap)
	if err != nil {
		return nil, err
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range kafkaService.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Value))
				}
			}
		}
	}()

	return kafkaService, nil
}

func (s *KafkaProduceService) readConfig() (kafka.ConfigMap, error) {
	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(s.configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			before, after, found := strings.Cut(line, "=")
			if found {
				parameter := strings.TrimSpace(before)
				value := strings.TrimSpace(after)
				m[parameter] = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *KafkaProduceService) Produce(fio entity.Fio) error {
	fioMarshalled, err := json.Marshal(fio)
	if err != nil {
		return err
	}

	logrus.Debugf("Sending %s", fio.String())
	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topicName, Partition: kafka.PartitionAny},
		Value:          fioMarshalled,
	}, nil)

	if err != nil {
		return err
	}
	return nil
}

func (s *KafkaProduceService) Shutdown() {
	logrus.Debug("Shutdown kafkaProducerService")
	s.producer.Flush(15 * 1000)
	s.producer.Close()
}
