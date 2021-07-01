package schema

import "bytes"

type Ticker struct {
	ID      string
	Convert string
	Start   string
	End     string
}

type CurrencySparkline struct {
	Currency   string   `json:"currency"`
	Timestamps []string `json:"timestamps"`
	Prices     []string `json:"prices"`
}

type Chart struct {
	Image *bytes.Buffer
}
