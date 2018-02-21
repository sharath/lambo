package CMC

import (
	"encoding/json"
	"github.com/sharath/lambo/util"
	"net/http"
)

// GlobalData is the storage for the JSON from the CMC api
type GlobalData struct {
	TotalMarketCapUsd            int64   `json:"total_market_cap_usd"`
	Total24HVolumeUsd            int64   `json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies             int     `json:"active_currencies"`
	ActiveAssets                 int     `json:"active_assets"`
	ActiveMarkets                int     `json:"active_markets"`
	LastUpdated                  int     `json:"last_updated"`
}

// FetchStats fetches the global data from the CMC api
func FetchStats() *GlobalData {
	g := new(GlobalData)
	resp, err := http.Get("https://api.coinmarketcap.com/v1/global/")
	if err != nil {
		util.HandleError(err, false)
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&g)
	return g
}
