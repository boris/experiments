package scanner

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/boris/experiments/go/kraken-market-scanner/internal/kraken"
)

type MarketData interface {
	AssetPairs(ctx context.Context) ([]kraken.AssetPair, error)
	Tickers(ctx context.Context, pairs []string) (map[string]kraken.Ticker, error)
	OHLC(ctx context.Context, pair string, interval int) ([]kraken.Candle, error)
}

type Config struct {
	TopN                int
	MinQuoteVolume      float64
	MaxSpreadPercent    float64
	OHLCIntervalMinutes int
	MaxOHLCPairs        int
	AllowedQuotes       []string
}

type Scanner struct {
	marketData MarketData
	config     Config
}

type Candidate struct {
	Symbol             string    `json:"symbol"`
	Score              float64   `json:"score"`
	LastPrice          float64   `json:"last_price"`
	SpreadPercent      float64   `json:"spread_percent"`
	DayChangePercent   float64   `json:"day_change_percent"`
	ShortChangePercent float64   `json:"short_change_percent"`
	Volume24h          float64   `json:"volume_24h"`
	VWAP24h            float64   `json:"vwap_24h"`
	Reasons            []string  `json:"reasons"`
	GeneratedAt        time.Time `json:"generated_at"`
}

type ScanResult struct {
	GeneratedAt time.Time   `json:"generated_at"`
	Candidates  []Candidate `json:"candidates"`
}

func New(marketData MarketData, config Config) *Scanner {
	if config.TopN <= 0 {
		config.TopN = 5
	}
	if config.MinQuoteVolume <= 0 {
		config.MinQuoteVolume = 500000
	}
	if config.MaxSpreadPercent <= 0 {
		config.MaxSpreadPercent = 0.4
	}
	if config.OHLCIntervalMinutes <= 0 {
		config.OHLCIntervalMinutes = 5
	}
	if config.MaxOHLCPairs <= 0 {
		config.MaxOHLCPairs = 25
	}
	if len(config.AllowedQuotes) == 0 {
		config.AllowedQuotes = []string{"USD", "USDT", "EUR"}
	}
	return &Scanner{marketData: marketData, config: config}
}

func (s *Scanner) Scan(ctx context.Context) (ScanResult, error) {
	pairs, err := s.marketData.AssetPairs(ctx)
	if err != nil {
		return ScanResult{}, err
	}

	filteredPairs := filterPairs(pairs, s.config.AllowedQuotes)
	symbols := make([]string, 0, len(filteredPairs))
	for _, pair := range filteredPairs {
		symbols = append(symbols, pair.Symbol)
	}

	tickers, err := s.marketData.Tickers(ctx, symbols)
	if err != nil {
		return ScanResult{}, err
	}

	preselected := make([]preCandidate, 0, len(filteredPairs))
	for _, pair := range filteredPairs {
		ticker, ok := tickers[pair.Symbol]
		if !ok {
			continue
		}
		spreadPct := spreadPercent(ticker.Bid, ticker.Ask)
		if ticker.Volume24h < s.config.MinQuoteVolume || spreadPct > s.config.MaxSpreadPercent || ticker.Bid <= 0 || ticker.Ask <= 0 {
			continue
		}
		preselected = append(preselected, preCandidate{pair: pair, ticker: ticker, spreadPercent: spreadPct})
	}

	sort.Slice(preselected, func(i, j int) bool {
		return preselected[i].ticker.Volume24h > preselected[j].ticker.Volume24h
	})
	if len(preselected) > s.config.MaxOHLCPairs {
		preselected = preselected[:s.config.MaxOHLCPairs]
	}

	candidates := make([]Candidate, 0, len(preselected))
	for _, item := range preselected {
		candles, err := s.marketData.OHLC(ctx, item.pair.Symbol, s.config.OHLCIntervalMinutes)
		if err != nil || len(candles) < 2 {
			continue
		}
		candidate := scoreCandidate(item.pair.Symbol, item.ticker, item.spreadPercent, candles, s.config.OHLCIntervalMinutes)
		if candidate.Score <= 0 {
			continue
		}
		candidates = append(candidates, candidate)
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].Volume24h > candidates[j].Volume24h
		}
		return candidates[i].Score > candidates[j].Score
	})
	if len(candidates) > s.config.TopN {
		candidates = candidates[:s.config.TopN]
	}
	return ScanResult{GeneratedAt: time.Now().UTC(), Candidates: candidates}, nil
}

type preCandidate struct {
	pair          kraken.AssetPair
	ticker        kraken.Ticker
	spreadPercent float64
}

func filterPairs(pairs []kraken.AssetPair, allowedQuotes []string) []kraken.AssetPair {
	allowed := make(map[string]struct{}, len(allowedQuotes))
	for _, quote := range allowedQuotes {
		allowed[strings.ToUpper(quote)] = struct{}{}
	}
	filtered := make([]kraken.AssetPair, 0, len(pairs))
	for _, pair := range pairs {
		if pair.Status != "online" {
			continue
		}
		if _, ok := allowed[strings.ToUpper(pair.Quote)]; !ok {
			continue
		}
		filtered = append(filtered, pair)
	}
	return filtered
}

func scoreCandidate(symbol string, ticker kraken.Ticker, spreadPct float64, candles []kraken.Candle, intervalMinutes int) Candidate {
	last := candles[len(candles)-1].Close
	previous := candles[len(candles)-2].Close
	first := candles[0].Close
	if last <= 0 || previous <= 0 || first <= 0 {
		return Candidate{}
	}

	shortChange := percentChange(previous, last)
	trendChange := percentChange(first, last)
	dayChange := 0.0
	if ticker.Open24h > 0 {
		dayChange = percentChange(ticker.Open24h, last)
	} else if ticker.VWAP24h > 0 {
		dayChange = percentChange(ticker.VWAP24h, last)
	}

	volatilityRange := 0.0
	if ticker.Low24h > 0 {
		volatilityRange = ((ticker.High24h - ticker.Low24h) / ticker.Low24h) * 100
	}

	reasons := make([]string, 0, 4)
	score := 0.0
	if shortChange < 0 {
		score += math.Min(math.Abs(shortChange)*0.8, 4)
		reasons = append(reasons, fmt.Sprintf("%dm momentum is negative (%.2f%%)", intervalMinutes, shortChange))
	}
	if trendChange < 0 {
		score += math.Min(math.Abs(trendChange)*0.6, 3)
		reasons = append(reasons, fmt.Sprintf("recent trend is down (%.2f%%)", trendChange))
	}
	if dayChange < 0 {
		score += math.Min(math.Abs(dayChange)*0.3, 2)
		reasons = append(reasons, fmt.Sprintf("trading below 24h VWAP/open (%.2f%%)", dayChange))
	}
	if ticker.Volume24h >= 2_000_000 {
		score += 1.5
		reasons = append(reasons, fmt.Sprintf("24h liquidity is strong (%.0f)", ticker.Volume24h))
	} else if ticker.Volume24h >= 1_000_000 {
		score += 1
		reasons = append(reasons, fmt.Sprintf("24h liquidity is healthy (%.0f)", ticker.Volume24h))
	}
	if spreadPct <= 0.2 {
		score += 1
		reasons = append(reasons, fmt.Sprintf("spread is tight (%.3f%%)", spreadPct))
	} else {
		score += math.Max(0.2, 0.8-(spreadPct*0.8))
		reasons = append(reasons, fmt.Sprintf("spread is still tradable (%.3f%%)", spreadPct))
	}
	if volatilityRange >= 5 {
		score += math.Min(volatilityRange/10, 1.5)
		reasons = append(reasons, fmt.Sprintf("24h range is active (%.2f%%)", volatilityRange))
	}

	return Candidate{
		Symbol:             symbol,
		Score:              math.Round(score*100) / 100,
		LastPrice:          last,
		SpreadPercent:      math.Round(spreadPct*1000) / 1000,
		DayChangePercent:   math.Round(dayChange*100) / 100,
		ShortChangePercent: math.Round(shortChange*100) / 100,
		Volume24h:          ticker.Volume24h,
		VWAP24h:            ticker.VWAP24h,
		Reasons:            reasons,
		GeneratedAt:        time.Now().UTC(),
	}
}

func percentChange(start, end float64) float64 {
	if start == 0 {
		return 0
	}
	return ((end - start) / start) * 100
}

func spreadPercent(bid, ask float64) float64 {
	if bid <= 0 || ask <= 0 {
		return math.MaxFloat64
	}
	mid := (bid + ask) / 2
	if mid == 0 {
		return math.MaxFloat64
	}
	return ((ask - bid) / mid) * 100
}

func FormatAlert(result ScanResult) string {
	if len(result.Candidates) == 0 {
		return fmt.Sprintf("Kraken short scanner found no candidates at %s", result.GeneratedAt.Format(time.RFC3339))
	}

	var builder strings.Builder
	builder.WriteString("Kraken short scanner\n")
	builder.WriteString(result.GeneratedAt.Format(time.RFC3339))
	builder.WriteString("\n\n")
	for index, candidate := range result.Candidates {
		builder.WriteString(fmt.Sprintf("%d. %s — Score %.2f\n", index+1, candidate.Symbol, candidate.Score))
		builder.WriteString(fmt.Sprintf("   Last %.4f | Spread %.3f%% | 24h %.2f%% | Short %.2f%% | Volume %.0f\n", candidate.LastPrice, candidate.SpreadPercent, candidate.DayChangePercent, candidate.ShortChangePercent, candidate.Volume24h))
		for _, reason := range candidate.Reasons {
			builder.WriteString("   - ")
			builder.WriteString(reason)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	return strings.TrimSpace(builder.String())
}
