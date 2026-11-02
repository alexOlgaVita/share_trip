package contractclient

import (
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"job4j.ru/share-trip/configs"
	contractgen "job4j.ru/share-trip/internal/clients/contract/gen"
)

type Client interface {
	CheckService(ctx context.Context, companyID string, serviceCode string) (CheckResult, error)
}

type CheckResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

type clientImpl struct {
	generated *contractgen.ClientWithResponses
}

func NewClient(contractCfg configs.Contract) (Client, error) {
	httpClient := &http.Client{
		Timeout: contractCfg.Timeout,
		Transport: &retryTransport{
			base:       http.DefaultTransport,
			retryCount: contractCfg.RetryCount,
			wait:       200 * time.Millisecond,
			maxWait:    time.Second,
		},
	}

	generated, err := contractgen.NewClientWithResponses(
		contractCfg.BaseUrl,
		contractgen.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	return &clientImpl{generated: generated}, nil
}

type retryTransport struct {
	base       http.RoundTripper
	retryCount int
	wait       time.Duration
	maxWait    time.Duration
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	var resp *http.Response
	var err error
	for attempt := 0; attempt <= t.retryCount; attempt++ {
		if attempt > 0 {
			if resp != nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				resp = nil
			}
			if waitErr := sleepCtx(req.Context(), backoff(t.wait, t.maxWait, attempt-1)); waitErr != nil {
				return nil, waitErr
			}
			if rewindErr := rewindBody(req); rewindErr != nil {
				return nil, rewindErr
			}
		}

		resp, err = base.RoundTrip(req)
		if !shouldRetry(resp, err) || attempt == t.retryCount {
			return resp, err
		}
	}

	return resp, err
}

func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func rewindBody(req *http.Request) error {
	if req.Body == nil || req.Body == http.NoBody {
		return nil
	}
	if req.GetBody == nil {
		return http.ErrNotSupported
	}
	body, err := req.GetBody()
	if err != nil {
		return err
	}
	req.Body = body
	return nil
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func backoff(min, max time.Duration, attempt int) time.Duration {
	wait := time.Duration(float64(min) * math.Exp2(float64(attempt)))
	if max > 0 && wait > max {
		return max
	}
	return wait
}
