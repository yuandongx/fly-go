// Package models provides the data models for the application.
package models

type Fund struct {
	BaseModel           `bson:",inline" json:",inline"`
	Code                string  `bson:"code,omitempty" json:"code,omitempty"`
	Name                string  `bson:"name,omitempty" json:"name,omitempty"`
	UnitNetValue        float64 `bson:"unit_net_value,omitempty" json:"unit_net_value,omitempty"`
	Date                string  `bson:"date,omitempty" json:"date,omitempty"`
	DayGrowthRate       float64 `bson:"day_growth_rate,omitempty" json:"day_growth_rate,omitempty"`
	WeekGrowthRate      float64 `bson:"week_growth_rate,omitempty" json:"week_growth_rate,omitempty"`
	MonthGrowthRate     float64 `bson:"month_growth_rate,omitempty" json:"month_growth_rate,omitempty"`
	QuarterGrowthRate   float64 `bson:"quarter_growth_rate,omitempty" json:"quarter_growth_rate,omitempty"`
	HalfYearGrowthRate  float64 `bson:"half_year_growth_rate,omitempty" json:"half_year_growth_rate,omitempty"`
	YearGrowthRate      float64 `bson:"year_growth_rate,omitempty" json:"year_growth_rate,omitempty"`
	TwoYearGrowthRate   float64 `bson:"two_year_growth_rate,omitempty" json:"two_year_growth_rate,omitempty"`
	ThreeYearGrowthRate float64 `bson:"three_year_growth_rate,omitempty" json:"three_year_growth_rate,omitempty"`
	ThisYearGrowthRate  float64 `bson:"this_year_growth_rate,omitempty" json:"this_year_growth_rate,omitempty"`
	TotalGrowthRate     float64 `bson:"total_growth_rate,omitempty" json:"total_growth_rate,omitempty"`
	Fee                 string  `bson:"fee,omitempty" json:"fee,omitempty"`
	MinPurchaseAmount   float64 `bson:"min_purchase_amount,omitempty" json:"min_purchase_amount,omitempty"`
	TrackIndex          int     `bson:"track_index,omitempty" json:"track_index,omitempty"`
	TrackType           int     `bson:"track_type,omitempty" json:"track_type,omitempty"`
	Category            int     `bson:"category,omitempty" json:"category,omitempty"`
	UpdateTime          string  `bson:"update_time,omitempty" json:"update_time,omitempty"`
}

type FundFollow struct {
	Fund      `bson:",inline" json:",inline"`
	StartTime string `bson:"start_time,omitempty" json:"start_time,omitempty"`
	Active    bool   `bson:"active,omitempty" json:"active,omitempty"`
}
