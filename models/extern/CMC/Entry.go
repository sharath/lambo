package CMC

import (
	"net/http"
	"strconv"
	"encoding/json"
)

type Entry struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

func FetchEntries(lim int) []*Entry {
	c := make([]*Entry, lim)
	resp, _ := http.Get("https://api.coinmarketcap.com/v1/ticker/?limit=" + strconv.Itoa(lim))
	json.NewDecoder(resp.Body).Decode(&c)
	return c
}
