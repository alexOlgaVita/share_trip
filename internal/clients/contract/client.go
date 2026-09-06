package contractclient

import (
	"context"
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

type checkServiceRequest struct {
	CompanyID   string `json:"CompanyID"`
	ServiceCode string `json:"ServiceCode"`
}

func mapHTTPError(resp *resty.Response) error {
	// Пример реализации, чтобы не возвращать пустой nil при ошибке
	return fmt.Errorf("http error: status code %d", resp.StatusCode())
}
