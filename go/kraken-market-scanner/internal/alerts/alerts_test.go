package alerts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMultiSenderSendsToAllSinks(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sender := MultiSender{Senders: []Sender{
		WebhookSender{Endpoint: server.URL, Client: server.Client()},
		WebhookSender{Endpoint: server.URL, Client: server.Client()},
	}}

	if err := sender.Send(context.Background(), "hello"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if calls != 2 {
		t.Fatalf("expected 2 webhook calls, got %d", calls)
	}
}

func TestTelegramSenderUsesBotAPI(t *testing.T) {
	var gotPath string
	var gotValues url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		gotValues = r.Form
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	sender := TelegramSender{
		BotToken: "token",
		ChatID:   "123",
		BaseURL:  server.URL,
		Client:   server.Client(),
	}

	if err := sender.Send(context.Background(), "market alert"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if gotPath != "/bottoken/sendMessage" {
		t.Fatalf("unexpected path %q", gotPath)
	}
	if gotValues.Get("chat_id") != "123" {
		t.Fatalf("unexpected chat id %q", gotValues.Get("chat_id"))
	}
	if gotValues.Get("text") != "market alert" {
		t.Fatalf("unexpected text %q", gotValues.Get("text"))
	}
}

func TestWebhookSenderPostsJSONPayload(t *testing.T) {
	var payload map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	sender := WebhookSender{Endpoint: server.URL, Client: server.Client()}
	if err := sender.Send(context.Background(), "payload text"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if payload["text"] != "payload text" {
		t.Fatalf("unexpected payload %#v", payload)
	}
}

func TestBuildSendersRequiresConfiguration(t *testing.T) {
	_, err := BuildSenders(Config{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}
	if !strings.Contains(err.Error(), "no alert destinations configured") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestBuildSendersRejectsPartialTelegramConfiguration(t *testing.T) {
	_, err := BuildSenders(Config{TelegramBotToken: "bot-only"})
	if err == nil {
		t.Fatal("expected error for partial telegram config")
	}
	if !strings.Contains(err.Error(), "requires both bot token and chat id") {
		t.Fatalf("unexpected error %v", err)
	}
}
