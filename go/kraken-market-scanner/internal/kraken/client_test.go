package kraken

import "testing"

func TestNormalizeAssetCode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "kraken wrapped bitcoin", in: "XXBT", want: "XBT"},
		{name: "kraken wrapped usd", in: "ZUSD", want: "USD"},
		{name: "kraken wrapped tether", in: "ZUSDT", want: "USDT"},
		{name: "real x-prefixed asset preserved", in: "XCN", want: "XCN"},
		{name: "real z-prefixed asset preserved", in: "ZETA", want: "ZETA"},
		{name: "already clean asset preserved", in: "ETH", want: "ETH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeAssetCode(tt.in); got != tt.want {
				t.Fatalf("normalizeAssetCode(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
