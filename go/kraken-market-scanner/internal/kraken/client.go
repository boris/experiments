package kraken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.kraken.com"

type AssetPair struct {
	Symbol       string
	AltName      string
	WSName       string
	Base         string
	Quote        string
	Status       string
	OrderMinimum float64
}

type Ticker struct {
	Ask       float64
	Bid       float64
	Volume24h float64
	VWAP24h   float64
	Low24h    float64
	High24h   float64
	Open24h   float64
}

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	VWAP   float64
	Volume float64
	Count  int
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), httpClient: httpClient}
}

func (c *Client) AssetPairs(ctx context.Context) ([]AssetPair, error) {
	endpoint := c.baseURL + "/0/public/AssetPairs"
	var response struct {
		Error  []string                     `json:"error"`
		Result map[string]assetPairResponse `json:"result"`
	}
	if err := c.getJSON(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	if len(response.Error) > 0 {
		return nil, fmt.Errorf("kraken asset pairs: %s", strings.Join(response.Error, ", "))
	}

	pairs := make([]AssetPair, 0, len(response.Result))
	for symbol, pair := range response.Result {
		pairs = append(pairs, AssetPair{
			Symbol:       symbol,
			AltName:      pair.AltName,
			WSName:       pair.WSName,
			Base:         normalizeAssetCode(pair.Base),
			Quote:        normalizeAssetCode(pair.Quote),
			Status:       pair.Status,
			OrderMinimum: parseFloat(pair.OrderMin),
		})
	}
	return pairs, nil
}

func (c *Client) Tickers(ctx context.Context, pairs []string) (map[string]Ticker, error) {
	if len(pairs) == 0 {
		return map[string]Ticker{}, nil
	}
	const batchSize = 40
	result := make(map[string]Ticker, len(pairs))
	for start := 0; start < len(pairs); start += batchSize {
		end := start + batchSize
		if end > len(pairs) {
			end = len(pairs)
		}
		batchResult, err := c.fetchTickerBatch(ctx, pairs[start:end])
		if err != nil {
			return nil, err
		}
		for symbol, ticker := range batchResult {
			result[symbol] = ticker
		}
	}
	return result, nil
}

func (c *Client) fetchTickerBatch(ctx context.Context, pairs []string) (map[string]Ticker, error) {
	endpoint := c.baseURL + "/0/public/Ticker?pair=" + url.QueryEscape(strings.Join(pairs, ","))
	var response struct {
		Error  []string                  `json:"error"`
		Result map[string]tickerResponse `json:"result"`
	}
	if err := c.getJSON(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	if len(response.Error) > 0 {
		return nil, fmt.Errorf("kraken ticker: %s", strings.Join(response.Error, ", "))
	}

	result := make(map[string]Ticker, len(response.Result))
	for symbol, ticker := range response.Result {
		result[symbol] = Ticker{
			Ask:       firstFloat(ticker.Ask),
			Bid:       firstFloat(ticker.Bid),
			Volume24h: secondFloat(ticker.Volume),
			VWAP24h:   secondFloat(ticker.VWAP),
			Low24h:    secondFloat(ticker.Low),
			High24h:   secondFloat(ticker.High),
			Open24h:   parseFloat(ticker.Open),
		}
	}
	return result, nil
}

func (c *Client) OHLC(ctx context.Context, pair string, interval int) ([]Candle, error) {
	if interval <= 0 {
		interval = 5
	}
	endpoint := fmt.Sprintf("%s/0/public/OHLC?pair=%s&interval=%d", c.baseURL, url.QueryEscape(pair), interval)
	var response struct {
		Error  []string                   `json:"error"`
		Result map[string]json.RawMessage `json:"result"`
	}
	if err := c.getJSON(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	if len(response.Error) > 0 {
		return nil, fmt.Errorf("kraken ohlc: %s", strings.Join(response.Error, ", "))
	}

	for key, raw := range response.Result {
		if key == "last" {
			continue
		}
		var rows [][]any
		if err := json.Unmarshal(raw, &rows); err != nil {
			return nil, err
		}
		candles := make([]Candle, 0, len(rows))
		for _, row := range rows {
			if len(row) < 8 {
				continue
			}
			candles = append(candles, Candle{
				Time:   time.Unix(int64(numberFromAny(row[0])), 0).UTC(),
				Open:   parseAnyFloat(row[1]),
				High:   parseAnyFloat(row[2]),
				Low:    parseAnyFloat(row[3]),
				Close:  parseAnyFloat(row[4]),
				VWAP:   parseAnyFloat(row[5]),
				Volume: parseAnyFloat(row[6]),
				Count:  int(numberFromAny(row[7])),
			})
		}
		return candles, nil
	}

	return nil, fmt.Errorf("kraken ohlc: no candle data for %s", pair)
}

type assetPairResponse struct {
	AltName  string `json:"altname"`
	WSName   string `json:"wsname"`
	Base     string `json:"base"`
	Quote    string `json:"quote"`
	Status   string `json:"status"`
	OrderMin string `json:"ordermin"`
}

type tickerResponse struct {
	Ask    []string `json:"a"`
	Bid    []string `json:"b"`
	Volume []string `json:"v"`
	VWAP   []string `json:"p"`
	Low    []string `json:"l"`
	High   []string `json:"h"`
	Open   string   `json:"o"`
}

func (c *Client) getJSON(ctx context.Context, endpoint string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kraken-market-scanner/0.1")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("kraken request failed: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func firstFloat(values []string) float64 {
	if len(values) == 0 {
		return 0
	}
	return parseFloat(values[0])
}

func secondFloat(values []string) float64 {
	if len(values) < 2 {
		return 0
	}
	return parseFloat(values[1])
}

func parseFloat(value string) float64 {
	parsed, _ := strconv.ParseFloat(value, 64)
	return parsed
}

func parseAnyFloat(value any) float64 {
	switch typed := value.(type) {
	case string:
		return parseFloat(typed)
	case float64:
		return typed
	case json.Number:
		parsed, _ := typed.Float64()
		return parsed
	default:
		return 0
	}
}

func numberFromAny(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case json.Number:
		parsed, _ := typed.Float64()
		return parsed
	case string:
		return parseFloat(typed)
	default:
		return 0
	}
}

func normalizeAssetCode(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if normalized, ok := normalizedAssetCodes[code]; ok {
		return normalized
	}
	return code
}

var normalizedAssetCodes = map[string]string{
	"XXBT":  "XBT",
	"XETH":  "ETH",
	"XLTC":  "LTC",
	"XETC":  "ETC",
	"XREP":  "REP",
	"XMLN":  "MLN",
	"XXMR":  "XMR",
	"XXDG":  "DOGE",
	"XZEC":  "ZEC",
	"ZUSD":  "USD",
	"ZEUR":  "EUR",
	"ZGBP":  "GBP",
	"ZJPY":  "JPY",
	"ZCAD":  "CAD",
	"ZAUD":  "AUD",
	"ZCHF":  "CHF",
	"ZUSDT": "USDT",
}
