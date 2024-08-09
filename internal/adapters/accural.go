package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/lionslon/yap-gophermart/internal/config"
	"github.com/lionslon/yap-gophermart/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Accrual struct {
	httpClient *retryablehttp.Client
	host       string
}

var ErrOrderNotRegistered = errors.New("the order is not registered in the payment system")

func NewAccrualClient(cfg *config.Config) *Accrual {
	retryClient := retryablehttp.NewClient()
	retryClient.CheckRetry = CheckRetry
	retryClient.Backoff = Backoff

	return &Accrual{
		host:       cfg.Accrual,
		httpClient: retryClient,
	}
}

func CheckRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if resp == nil {
		return false, err
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	return false, err
}

func Backoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	const defaultTimeout = 120

	if resp.StatusCode != http.StatusTooManyRequests {
		return defaultTimeout * time.Second
	}

	retryAfterString := resp.Header.Get("Retry-After")
	retryAfter, err := strconv.ParseInt(retryAfterString, 10, 64)
	if err != nil {
		return defaultTimeout * time.Second
	}

	return time.Duration(retryAfter) * time.Second
}

func (a *Accrual) GetOrderAccrual(ctx context.Context, order *models.Order) (*models.OrderAccrual, error) {
	url, err := url.JoinPath(a.host, "/api/orders/", order.Number)
	if err != nil {
		return nil, fmt.Errorf("failed build url err: %w", err)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed build accrual request err: %w", err)
	}
	req.Close = true

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed exec accrual request err: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, ErrOrderNotRegistered
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body err: %w", err)
	}

	var oa models.OrderAccrual
	if err := json.Unmarshal(res, &oa); err != nil {
		return nil, fmt.Errorf("failed unmarshal response body err: %w", err)
	}

	return &oa, nil
}
