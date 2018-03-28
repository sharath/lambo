package response

import (
	"encoding/json"
	"github.com/sharath/lambo/database"
	"gopkg.in/mgo.v2"
)

// Hist is the response format for the hist endpoint
type Hist struct {
	Name           string `json:"name" bson:"name"`
	Symbol         string `json:"symbol" bson:"symbol"`
	PriceUsd       string `json:"price_usd" bson:"price_usd"`
	PriceBtc       string `json:"price_btc" bson:"price_btc"`
	Two4HVolumeUsd string `json:"24h_volume_usd" bson:"24h_volume_usd"`
	MarketCapUsd   string `json:"market_cap_usd" bson:"market_cap_usd"`
	LastUpdated    string `json:"last_updated" bson:"last_updated"`
}

// NewHist returns a response for the hist endpoint
func NewHist(name string, entries *mgo.Collection) []Hist {
	var all []*database.MongoEntry
	var hist []Hist
	entries.Find(nil).All(&all)
	for _, e := range all {
		for _, t := range e.Tokens {
			if t.Name == name {
				var temp Hist
				enc, _ := json.Marshal(t)
				json.Unmarshal(enc, &temp)
				hist = append(hist, temp)
			}
		}
	}
	return hist
}
