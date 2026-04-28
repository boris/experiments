FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/kraken-market-scanner ./cmd/kraken-market-scanner

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/kraken-market-scanner /kraken-market-scanner
ENTRYPOINT ["/kraken-market-scanner"]
CMD ["-mode", "serve", "-addr", ":8080"]
