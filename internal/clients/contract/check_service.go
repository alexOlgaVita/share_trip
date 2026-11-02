package contractclient

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	contractgen "job4j.ru/share-trip/internal/clients/contract/gen"
	"net/http"
)

func (c *clientImpl) CheckService(ctx context.Context, companyID string, serviceCode string) (CheckResult, error) {

	companyIDUuid, err := uuid.Parse(companyID)
	if err != nil {
		return CheckResult{}, fmt.Errorf("invalid company ID: %w", err)
	}
	resp, err := c.generated.CheckServiceAvailabilityWithResponse(ctx, contractgen.CheckServiceAvailabilityJSONRequestBody{
		CompanyId:   companyIDUuid,
		ServiceCode: serviceCode,
	})
	if err != nil {
		return CheckResult{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return CheckResult{}, fmt.Errorf("api error: status code %d", resp.StatusCode())
	}

	if resp.JSON200 == nil {
		return CheckResult{}, fmt.Errorf("api returned 200 but body is empty")
	}

	result := CheckResult{
		Allowed: resp.JSON200.Allowed,
	}

	if resp.JSON200.Reason != nil {
		result.Reason = string(*resp.JSON200.Reason)
	}

	return result, nil
}
