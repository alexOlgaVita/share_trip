package contractclient

import (
	"context"
	"encoding/json"
	"fmt"
)

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
