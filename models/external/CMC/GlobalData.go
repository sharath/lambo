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

func FetchGlobalStats() (*GlobalData, error) {
	g := new(GlobalData)
	resp, err := http.Get("https://api.coinmarketcap.com/v1/global/")
	if err != nil {
		return g, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&g)
	return g, err
}