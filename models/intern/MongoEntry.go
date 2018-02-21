package intern

type Token struct {
	ID               string `json:"id" bson:"id"`
	Name             string `json:"name" bson:"name"`
	Symbol           string `json:"symbol" bson:"symbol"`
	Rank             string `json:"rank" bson:"rank"`
	PriceUsd         string `json:"price_usd" bson:"price_usd"`
	PriceBtc         string `json:"price_btc" bson:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd" bson:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd" bson:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply" bson:"available_supply"`
	TotalSupply      string `json:"total_supply" bson:"total_supply"`
	MaxSupply        string `json:"max_supply" bson:"max_supply"`
	PercentChange1H  string `json:"percent_change_1h" bson:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h" bson:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d" bson:"percent_change_7d"`
	LastUpdated      string `json:"last_updated" bson:"last_updated"`
}

type Global struct {
	TotalMarketCapUsd            int64   `json:"total_market_cap_usd" bson:"total_market_cap_usd"`
	Total24HVolumeUsd            int64   `json:"total_24h_volume_usd" bson:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap" bson:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies             int     `json:"active_currencies" bson:"active_currencies"`
	ActiveAssets                 int     `json:"active_assets" bson:"active_assets"`
	ActiveMarkets                int     `json:"active_markets" bson:"active_markets"`
	LastUpdated                  int     `json:"last_updated" bson:"last_updated"`
}

type MongoEntry struct {
	Tokens []*Token `json:"token_data" bson:"token_data"`
	Global *Global  `json:"global_data" bson:"global_data"`
}
