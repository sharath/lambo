package intern

type MongoEntry struct {
	Tokens []struct {
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
	} `json:"token_data",bson:"token_data"`
	Global struct {
		TotalMarketCapUsd            int64   `json:"total_market_cap_usd"`
		Total24HVolumeUsd            int64   `json:"total_24h_volume_usd"`
		BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
		ActiveCurrencies             int     `json:"active_currencies"`
		ActiveAssets                 int     `json:"active_assets"`
		ActiveMarkets                int     `json:"active_markets"`
		LastUpdated                  int     `json:"last_updated"`
	} `json:"global_data",bson:"global_data"`
}
