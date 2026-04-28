package scanner

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/boris/experiments/go/kraken-market-scanner/internal/kraken"
)

type fakeMarketData struct {
	pairs   []kraken.AssetPair
	tickers map[string]kraken.Ticker
	ohlc    map[string][]kraken.Candle
	err     error
}

func (f fakeMarketData) AssetPairs(ctx context.Context) ([]kraken.AssetPair, error) {
	return f.pairs, f.err
}

func (f fakeMarketData) Tickers(ctx context.Context, pairs []string) (map[string]kraken.Ticker, error) {
	return f.tickers, f.err
}

func (f fakeMarketData) OHLC(ctx context.Context, pair string, interval int) ([]kraken.Candle, error) {
	if f.err != nil {
		return nil, f.err
	}
	candles, ok := f.ohlc[pair]
	if !ok {
		return nil, errors.New("missing candles")
	}
	return candles, nil
}

func TestScanRanksBearishLiquidCandidates(t *testing.T) {
	market := fakeMarketData{
		pairs: []kraken.AssetPair{
			{Symbol: "WEAKUSD", Quote: "USD", Status: "online", OrderMinimum: 10},
			{Symbol: "FASTUSD", Quote: "USD", Status: "online", OrderMinimum: 10},
			{Symbol: "WIDEUSD", Quote: "USD", Status: "online", OrderMinimum: 10},
		},
		tickers: map[string]kraken.Ticker{
			"WEAKUSD": {Ask: 100.2, Bid: 100.0, Volume24h: 2_200_000, VWAP24h: 104, Low24h: 96, High24h: 111},
			"FASTUSD": {Ask: 52.1, Bid: 52.0, Volume24h: 3_100_000, VWAP24h: 53, Low24h: 49, High24h: 54},
			"WIDEUSD": {Ask: 80.6, Bid: 80.0, Volume24h: 8_000_000, VWAP24h: 92, Low24h: 77, High24h: 96},
		},
		ohlc: map[string][]kraken.Candle{
			"WEAKUSD": {
				{Close: 108},
				{Close: 102},
				{Close: 100},
			},
			"FASTUSD": {
				{Close: 54},
				{Close: 53},
				{Close: 52},
			},
			"WIDEUSD": {
				{Close: 93},
				{Close: 87},
				{Close: 80.4},
			},
		},
	}

	s := New(market, Config{TopN: 2, MinQuoteVolume: 1_000_000, MaxSpreadPercent: 0.5, OHLCIntervalMinutes: 5})

	result, err := s.Scan(context.Background())
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(result.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(result.Candidates))
	}

	if result.Candidates[0].Symbol != "WEAKUSD" {
		t.Fatalf("expected top candidate WEAKUSD, got %s", result.Candidates[0].Symbol)
	}

	if result.Candidates[1].Symbol != "FASTUSD" {
		t.Fatalf("expected second candidate FASTUSD, got %s", result.Candidates[1].Symbol)
	}

	if got := result.Candidates[0].SpreadPercent; got > 0.5 {
		t.Fatalf("expected filtered spread <= 0.5, got %.4f", got)
	}
}

func TestScanSkipsPairsWithoutRequiredLiquidity(t *testing.T) {
	market := fakeMarketData{
		pairs: []kraken.AssetPair{{Symbol: "THINUSD", Quote: "USD", Status: "online", OrderMinimum: 10}},
		tickers: map[string]kraken.Ticker{
			"THINUSD": {Ask: 10.1, Bid: 10.0, Volume24h: 50_000, VWAP24h: 10.5, Low24h: 9.9, High24h: 11},
		},
		ohlc: map[string][]kraken.Candle{
			"THINUSD": {{Close: 10.6}, {Close: 10.4}, {Close: 10.0}},
		},
	}

	s := New(market, Config{TopN: 5, MinQuoteVolume: 100_000, MaxSpreadPercent: 0.5, OHLCIntervalMinutes: 5})

	result, err := s.Scan(context.Background())
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(result.Candidates) != 0 {
		t.Fatalf("expected no candidates, got %d", len(result.Candidates))
	}
}

func TestFormatAlertIncludesWhyLines(t *testing.T) {
	text := FormatAlert(ScanResult{Candidates: []Candidate{{
		Symbol:             "WEAKUSD",
		Score:              8.42,
		LastPrice:          100,
		SpreadPercent:      0.2,
		DayChangePercent:   -3.8,
		ShortChangePercent: -6.1,
		Volume24h:          2_200_000,
		Reasons:            []string{"5m momentum is negative", "trading below 24h VWAP"},
	}}})

	if text == "" {
		t.Fatal("expected non-empty alert")
	}

	for _, fragment := range []string{"WEAKUSD", "Score", "5m momentum is negative", "trading below 24h VWAP"} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("expected alert to contain %q, got %q", fragment, text)
		}
	}
}
