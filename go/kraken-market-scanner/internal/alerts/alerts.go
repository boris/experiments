package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Sender interface {
	Send(ctx context.Context, message string) error
}

type Config struct {
	TelegramBotToken string
	TelegramChatID   string
	WebhookURL       string
	HTTPTimeout      time.Duration
}

type MultiSender struct {
	Senders []Sender
}

func (m MultiSender) Send(ctx context.Context, message string) error {
	for _, sender := range m.Senders {
		if err := sender.Send(ctx, message); err != nil {
			return err
		}
	}
	return nil
}

type TelegramSender struct {
	BotToken string
	ChatID   string
	BaseURL  string
	Client   *http.Client
}

func (s TelegramSender) Send(ctx context.Context, message string) error {
	if s.BotToken == "" || s.ChatID == "" {
		return errors.New("telegram sender is missing bot token or chat id")
	}
	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	baseURL := strings.TrimRight(s.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.telegram.org"
	}
	form := url.Values{}
	form.Set("chat_id", s.ChatID)
	form.Set("text", message)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/bot%s/sendMessage", baseURL, s.BotToken), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram sender failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

type WebhookSender struct {
	Endpoint string
	Client   *http.Client
}

func (s WebhookSender) Send(ctx context.Context, message string) error {
	if s.Endpoint == "" {
		return errors.New("webhook sender is missing endpoint")
	}
	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	payload, err := json.Marshal(map[string]string{"text": message})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook sender failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func BuildSenders(config Config) (Sender, error) {
	timeout := config.HTTPTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	senders := make([]Sender, 0, 2)
	if config.TelegramBotToken != "" || config.TelegramChatID != "" {
		if config.TelegramBotToken == "" || config.TelegramChatID == "" {
			return nil, errors.New("telegram delivery requires both bot token and chat id")
		}
		senders = append(senders, TelegramSender{BotToken: config.TelegramBotToken, ChatID: config.TelegramChatID, Client: client})
	}
	if config.WebhookURL != "" {
		senders = append(senders, WebhookSender{Endpoint: config.WebhookURL, Client: client})
	}
	if len(senders) == 0 {
		return nil, errors.New("no alert destinations configured")
	}
	if len(senders) == 1 {
		return senders[0], nil
	}
	return MultiSender{Senders: senders}, nil
}
