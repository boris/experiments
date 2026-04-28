# Kraken market scanner

This experiment scans Kraken spot markets, ranks the best bearish/short candidates, and optionally pushes alerts to Telegram and/or a generic webhook.

It is a scanner, not an auto-trader. The goal is to surface liquid markets with tight spreads and short-term weakness so you can review them quickly.

## What it does

- pulls Kraken asset pairs, ticker data, and short-interval OHLC candles
- filters out thin or wide-spread markets
- ranks candidates with a fee-aware bearish score
- emits a compact alert digest
- can run once from the CLI or expose an HTTP endpoint for scheduled triggers
- supports Telegram and generic webhook delivery

## Project layout

```text
.
├── cmd/kraken-market-scanner
│   └── main.go
├── deploy
│   └── cloudflare-worker.mjs
├── internal
│   ├── alerts
│   ├── app
│   ├── kraken
│   └── scanner
└── go.mod
```

## Local usage

Build it:

```bash
go build ./cmd/kraken-market-scanner
```

Run one scan locally:

```bash
./kraken-market-scanner \
  -mode once \
  -top 5 \
  -min-volume 1000000 \
  -max-spread 0.5
```

Send alerts to Telegram:

```bash
export TELEGRAM_BOT_TOKEN="..."
export TELEGRAM_CHAT_ID="..."
./kraken-market-scanner -mode once -send
```

Send alerts to a generic webhook instead:

```bash
export ALERT_WEBHOOK_URL="https://example.com/hooks/market-alerts"
./kraken-market-scanner -mode once -send
```

Run as a small HTTP service:

```bash
export SCANNER_API_TOKEN="replace-me"
./kraken-market-scanner -mode serve -addr :8080
```

Then trigger a scan:

```bash
curl -sS \
  -H "Authorization: Bearer $SCANNER_API_TOKEN" \
  "http://localhost:8080/scan?send=1"
```

## Configuration

All flags also support environment variables.

| Flag | Env | Default | Purpose |
|---|---|---:|---|
| `-mode` | `SCANNER_MODE` | `once` | Run once or serve HTTP |
| `-addr` | `SCANNER_ADDR` | `:8080` | Listen address in serve mode |
| `-api-base-url` | `KRAKEN_API_BASE_URL` | `https://api.kraken.com` | Kraken API base URL |
| `-http-timeout` | `HTTP_TIMEOUT` | `15s` | Outbound HTTP timeout |
| `-top` | `TOP_N` | `5` | Number of candidates to keep |
| `-min-volume` | `MIN_QUOTE_VOLUME` | `1000000` | Minimum 24h quote volume |
| `-max-spread` | `MAX_SPREAD_PERCENT` | `0.5` | Maximum spread percentage |
| `-ohlc-interval` | `OHLC_INTERVAL_MINUTES` | `5` | OHLC interval |
| `-max-ohlc-pairs` | `MAX_OHLC_PAIRS` | `25` | OHLC fetch fan-out cap |
| `-allowed-quotes` | `ALLOWED_QUOTES` | `USD,USDT,EUR` | Quotes to scan |
| `-alert-threshold` | `ALERT_THRESHOLD` | `0` | Minimum top score required to send |
| `-send` | `SEND_ALERTS` | `false` | Send alerts on run |
| `-telegram-bot-token` | `TELEGRAM_BOT_TOKEN` | | Telegram bot token |
| `-telegram-chat-id` | `TELEGRAM_CHAT_ID` | | Telegram chat ID |
| `-webhook-url` | `ALERT_WEBHOOK_URL` | | Generic JSON webhook |
| `-scanner-token` | `SCANNER_API_TOKEN` | | Shared token for `/scan` |

## How the ranking works

Each candidate is prefiltered and scored on:

- negative short-interval momentum
- negative recent trend
- trading below the 24h open or VWAP
- strong 24h liquidity
- acceptable spread
- active daily range

The score is intentionally simple and explainable so it is easy to tune. The alert includes the reasons behind each rank.

## Tests

Run:

```bash
go test ./...
```

## Deployment

### Recommended Go-native serverless path: Cloud Run

Cloud Run is the cleanest serverless target for this code because it runs Go directly.

1. Build and push a container image.
2. Deploy the binary in `serve` mode.
3. Protect `/scan` with `SCANNER_API_TOKEN`.
4. Use Cloud Scheduler or an external trigger to call `/scan?send=1`.
5. Keep Telegram/webhook credentials in Secret Manager or Cloud Run env vars.

A minimal Dockerfile is included in this directory. Build and deploy from `go/kraken-market-scanner` so the build context matches the module root:

```bash
cd go/kraken-market-scanner
docker build -t kraken-market-scanner .
```

The included Dockerfile is:

```dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/kraken-market-scanner ./cmd/kraken-market-scanner

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/kraken-market-scanner /kraken-market-scanner
ENTRYPOINT ["/kraken-market-scanner"]
CMD ["-mode", "serve", "-addr", ":8080"]
```

Example deploy command:

```bash
cd go/kraken-market-scanner
gcloud run deploy kraken-market-scanner \
  --source . \
  --region us-central1 \
  --set-env-vars SCANNER_MODE=serve,SCANNER_ADDR=:8080,SCANNER_API_TOKEN=replace-me,TELEGRAM_BOT_TOKEN=your-bot-token,TELEGRAM_CHAT_ID=your-chat-id \
  --allow-unauthenticated=false
```

### Optional Cloudflare Worker scheduler relay

Cloudflare Workers do not run this Go service directly, but they are a good edge scheduler/trigger.

Use the included `deploy/cloudflare-worker.mjs` as a cron-triggered relay that calls your deployed scanner endpoint with the shared token.

Set these Worker secrets:

- `SCANNER_URL`
- `SCANNER_API_TOKEN`
- `WORKER_TRIGGER_TOKEN`

Then add a cron trigger such as `*/5 * * * *` to scan every 5 minutes.
If you keep the optional `/trigger` route enabled, call it with `Authorization: Bearer $WORKER_TRIGGER_TOKEN`.

This gives you:

- serverless scheduling at the edge
- a Go-native scanner deployment
- fast webhook-style triggering
- Telegram delivery without polling

## Notes

- Kraken fees, spread, and slippage still dominate short-horizon profitability. Use this as a discovery and alerting tool first.
- If you want duplicate-alert suppression or paper-trading state, add a persistent store next. Cloudflare KV, D1, Redis, or Firestore would all work.
