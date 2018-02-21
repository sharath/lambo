package CMC

import (
	"net/http"
	"encoding/json"
)

type GlobalData struct {
	TotalMarketCapUsd            int64   `json:"total_market_cap_usd"`
	Total24HVolumeUsd            int64   `json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies             int     `json:"active_currencies"`
	ActiveAssets                 int     `json:"active_assets"`
	ActiveMarkets                int     `json:"active_markets"`
	LastUpdated                  int     `json:"last_updated"`
}

func (g *GlobalData) FetchStats() {
	resp, _ := http.Get("https://api.coinmarketcap.com/v1/global/")
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&g)
}