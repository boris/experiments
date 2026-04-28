package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/boris/experiments/go/kraken-market-scanner/internal/alerts"
	"github.com/boris/experiments/go/kraken-market-scanner/internal/kraken"
	"github.com/boris/experiments/go/kraken-market-scanner/internal/scanner"
)

type Config struct {
	APIBaseURL       string
	HTTPTimeout      time.Duration
	AlertThreshold   float64
	RequireAuthToken string
	Scanner          scanner.Config
	Alerts           alerts.Config
}

type Service struct {
	scanner        *scanner.Scanner
	alertSender    alerts.Sender
	alertThreshold float64
	authToken      string
}

type Response struct {
	Sent       bool               `json:"sent"`
	Message    string             `json:"message"`
	ScanResult scanner.ScanResult `json:"scan_result"`
}

func New(config Config) (*Service, error) {
	httpClient := &http.Client{Timeout: config.HTTPTimeout}
	if config.HTTPTimeout <= 0 {
		httpClient.Timeout = 15 * time.Second
	}
	marketScanner := scanner.New(kraken.NewClient(httpClient, config.APIBaseURL), config.Scanner)

	var sender alerts.Sender
	var err error
	if hasAlertConfig(config.Alerts) {
		sender, err = alerts.BuildSenders(config.Alerts)
		if err != nil {
			return nil, err
		}
	}
	return &Service{scanner: marketScanner, alertSender: sender, alertThreshold: config.AlertThreshold, authToken: config.RequireAuthToken}, nil
}

func (s *Service) Run(ctx context.Context, send bool) (Response, error) {
	result, err := s.scanner.Scan(ctx)
	if err != nil {
		return Response{}, err
	}
	message := scanner.FormatAlert(result)
	response := Response{Message: message, ScanResult: result}
	if send && s.alertSender != nil && shouldSend(result, s.alertThreshold) {
		if err := s.alertSender.Send(ctx, message); err != nil {
			return Response{}, err
		}
		response.Sent = true
	}
	return response, nil
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		if s.authToken != "" && !authorized(r, s.authToken) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		send := r.URL.Query().Get("send") == "1" || strings.EqualFold(r.URL.Query().Get("send"), "true")
		response, err := s.Run(r.Context(), send)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, response)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func shouldSend(result scanner.ScanResult, threshold float64) bool {
	if len(result.Candidates) == 0 {
		return false
	}
	if threshold <= 0 {
		return true
	}
	return result.Candidates[0].Score >= threshold
}

func authorized(r *http.Request, token string) bool {
	if r.Header.Get("X-Scanner-Token") == token {
		return true
	}
	authHeader := r.Header.Get("Authorization")
	return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer ")) == token && strings.HasPrefix(authHeader, "Bearer ")
}

func hasAlertConfig(config alerts.Config) bool {
	return config.TelegramBotToken != "" || config.TelegramChatID != "" || config.WebhookURL != ""
}
