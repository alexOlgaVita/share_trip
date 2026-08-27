package contractclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"job4j.ru/share-trip/configs"
	"net/http"
	"time"
)

type Client interface {
	CheckService(ctx context.Context, companyID string, serviceCode string) (CheckResult, error)
}

type CheckResult struct {
	//Allowed bool
	//Reason  string
	Allowed bool   `json:"allowed"` // Обязательно добавь этот тег
	Reason  string `json:"reason"`  // И этот
}

type clientImpl struct {
	http *resty.Client
}

var ClientFactory = func() Client {
	return &clientImpl{}
}

func NewClient(contractCfg configs.Contract) Client {
	r := resty.New().
		SetBaseURL(contractCfg.BaseUrl).
		SetTimeout(contractCfg.Timeout).
		SetRetryCount(contractCfg.RetryCount).
		SetRetryWaitTime(200 * time.Millisecond).
		SetRetryMaxWaitTime(1 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				return true
			}
			return r.StatusCode() == http.StatusTooManyRequests ||
				r.StatusCode() == http.StatusBadGateway ||
				r.StatusCode() == http.StatusServiceUnavailable ||
				r.StatusCode() == http.StatusGatewayTimeout
		})

	return &clientImpl{
		http: r,
	}
}

func (c *clientImpl) CheckService(ctx context.Context, companyID string, serviceCode string) (CheckResult, error) {
	result := CheckResult{}

	if c.http == nil {
		return CheckResult{}, fmt.Errorf("resty client is not initialized")
	}
	resp, err := c.http.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(checkServiceRequest{
			CompanyID:   companyID,
			ServiceCode: serviceCode,
		}).
		//SetResult(&result)- работает не корректно, поэтому ниже делаем отдельно Unmarshal
		Post("/contracts/service-availability-checks")
	if err != nil {
		return CheckResult{}, err
	}

	if resp.IsError() {
		return CheckResult{}, mapHTTPError(resp)
	}

	// вместо SetResult в коде выше используем ручной разбор
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return CheckResult{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return result, nil
}

type checkServiceRequest struct {
	CompanyID   string `json:"CompanyID"`
	ServiceCode string `json:"ServiceCode"`
}

func mapHTTPError(resp *resty.Response) error {
	// Пример реализации, чтобы не возвращать пустой nil при ошибке
	return fmt.Errorf("http error: status code %d", resp.StatusCode())
}
