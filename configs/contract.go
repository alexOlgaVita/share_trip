package configs

import (
	"errors"
	"strconv"
	"time"
)

type Contract struct {
	BaseUrl    string
	Timeout    time.Duration
	RetryCount int
}

func LoadContract() (Contract, error) {
	baseUrl := Env("CONTRACT_BASE_URL", "http://localhost:8082")

	timeoutStr := Env("CONTRACT_TIMEOUT", "2s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return Contract{}, errors.New("contract-service.config: timeout has bad value")
	}

	retryStr := Env("CONTRACT_RETRY_COUNT", "2")
	if retryStr == "" {
		retryStr = "2"
	}
	retryCount, err := strconv.Atoi(retryStr)
	if err != nil {
		return Contract{}, errors.New("contract-service.config: retryCount has bad value")
	}

	return Contract{
		BaseUrl:    baseUrl,
		Timeout:    timeout,
		RetryCount: retryCount,
	}, nil
}

func (c Contract) Validate() error {
	if c.BaseUrl == "" {
		return errors.New("CONTRACT_BASE_URL, CONTRACT_TIMEOUT and CONTRACT_RETRY_COUNT are required")
	}
	return nil
}
