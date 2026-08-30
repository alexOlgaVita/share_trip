package configs

import (
	"errors"
	"strings"
)

type KafkaConf struct {
	BrokerList []string
}

func LoadKafkaConf() (KafkaConf, error) {
	rawBrokers := Env("NOTIFICATION_KAFKA_BROKERS_LIST", "localhost:9092,localhost:9093,127.0.0.1:9094")

	// 1. Разбиваем строку в слайс []string
	brokerList := strings.Split(rawBrokers, ",")

	// 2. Очищаем от лишних пробелов (на случай "host:9092, host:9093")
	for i := range brokerList {
		brokerList[i] = strings.TrimSpace(brokerList[i])
	}

	return KafkaConf{
		BrokerList: brokerList,
	}, nil
}

func (c KafkaConf) Validate() error {
	if len(c.BrokerList) == 0 {
		return errors.New("NOTIFICATION_KAFKA_BROKERS_LIST are required")
	}
	return nil
}

func GetBrokerList(c KafkaConf) []string {
	return c.BrokerList
}
