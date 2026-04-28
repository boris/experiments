package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boris/experiments/go/kraken-market-scanner/internal/alerts"
	"github.com/boris/experiments/go/kraken-market-scanner/internal/kraken"
	"github.com/boris/experiments/go/kraken-market-scanner/internal/scanner"
)

type stubMarketData struct{}

func (stubMarketData) AssetPairs(ctx context.Context) ([]kraken.AssetPair, error) {
	return []kraken.AssetPair{{Symbol: "WEAKUSD", Quote: "USD", Status: "online"}}, nil
}

func (stubMarketData) Tickers(ctx context.Context, pairs []string) (map[string]kraken.Ticker, error) {
	return map[string]kraken.Ticker{
		"WEAKUSD": {Ask: 100.2, Bid: 100.0, Volume24h: 2_200_000, Open24h: 110, Low24h: 95, High24h: 111},
	}, nil
}

func (stubMarketData) OHLC(ctx context.Context, pair string, interval int) ([]kraken.Candle, error) {
	return []kraken.Candle{{Close: 108}, {Close: 102}, {Close: 100}}, nil
}

type stubSender struct{ calls int }

func (s *stubSender) Send(ctx context.Context, message string) error {
	s.calls++
	if message == "" {
		return errors.New("empty message")
	}
	return nil
}

func TestHandlerRequiresAuthToken(t *testing.T) {
	service := &Service{
		scanner:     scanner.New(stubMarketData{}, scanner.Config{TopN: 1, MinQuoteVolume: 1_000_000, MaxSpreadPercent: 0.5, OHLCIntervalMinutes: 15}),
		authToken:   "secret",
		alertSender: &stubSender{},
	}

	request := httptest.NewRequest(http.MethodGet, "/scan?send=1", nil)
	response := httptest.NewRecorder()
	service.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestHandlerAllowsAuthorizedScanAndSend(t *testing.T) {
	sender := &stubSender{}
	service := &Service{
		scanner:        scanner.New(stubMarketData{}, scanner.Config{TopN: 1, MinQuoteVolume: 1_000_000, MaxSpreadPercent: 0.5, OHLCIntervalMinutes: 15}),
		authToken:      "secret",
		alertSender:    sender,
		alertThreshold: 1,
	}

	request := httptest.NewRequest(http.MethodGet, "/scan?send=1", nil)
	request.Header.Set("Authorization", "Bearer secret")
	response := httptest.NewRecorder()
	service.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if sender.calls != 1 {
		t.Fatalf("expected 1 alert send, got %d", sender.calls)
	}
}

func TestHasAlertConfig(t *testing.T) {
	if hasAlertConfig(alerts.Config{}) {
		t.Fatal("expected empty config to be false")
	}
	if !hasAlertConfig(alerts.Config{WebhookURL: "https://example.com"}) {
		t.Fatal("expected webhook config to be true")
	}
}

func TestNewBuildsScannerWithDefaultTimeout(t *testing.T) {
	service, err := New(Config{HTTPTimeout: 0})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if service.scanner == nil {
		t.Fatal("expected scanner to be initialized")
	}
	if service.alertSender != nil {
		t.Fatal("expected no sender without alert config")
	}
	if service.alertThreshold != 0 {
		t.Fatalf("expected zero alert threshold, got %v", service.alertThreshold)
	}
}
