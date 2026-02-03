package models

type Stock struct {
	BaseModel     `bson:",inline" json:",inline"`
	Symbol        string  `bson:"symbol,omitempty" json:"symbol,omitempty"`
	Code          string  `bson:"code,omitempty" json:"code,omitempty"`
	Name          string  `bson:"name,omitempty" json:"name,omitempty"`
	Trade         string  `bson:"trade,omitempty" json:"trade,omitempty"`
	Pricechange   float64 `bson:"pricechange,omitempty" json:"pricechange,omitempty"`
	Changepercent float64 `bson:"changepercent,omitempty" json:"changepercent,omitempty"`
	Buy           string  `bson:"buy,omitempty" json:"buy,omitempty"`
	Sell          string  `bson:"sell,omitempty" json:"sell,omitempty"`
	Settlement    string  `bson:"settlement,omitempty" json:"settlement,omitempty"`
	Open          string  `bson:"open,omitempty" json:"open,omitempty"`
	High          string  `bson:"high,omitempty" json:"high,omitempty"`
	Low           string  `bson:"low,omitempty" json:"low,omitempty"`
	Volume        int64   `bson:"volume,omitempty" json:"volume,omitempty"`
	Amount        int64   `bson:"amount,omitempty" json:"amount,omitempty"`
	Ticktime      string  `bson:"ticktime,omitempty" json:"ticktime,omitempty"`
	Per           float64 `bson:"per,omitempty" json:"per,omitempty"`
	Pb            float64 `bson:"pb,omitempty" json:"pb,omitempty"`
	Mktcap        float64 `bson:"mktcap,omitempty" json:"mktcap,omitempty"`
	Nmc           float64 `bson:"nmc,omitempty" json:"nmc,omitempty"`
	Turnoverratio float64 `bson:"turnoverratio,omitempty" json:"turnoverratio,omitempty"`
	Date          string  `bson:"date,omitempty" json:"date,omitempty"`
	Area          string  `bson:"area,omitempty" json:"area,omitempty"`
	Follow        int     `bson:"follow,omitempty" json:"follow,omitempty"`
}

type StockFollow struct {
	Stock     `bson:",inline" json:",inline"`
	StartTime string `bson:"start_time,omitempty" json:"start_time,omitempty"`
	Active    bool   `bson:"active,omitempty" json:"active,omitempty"`
}
