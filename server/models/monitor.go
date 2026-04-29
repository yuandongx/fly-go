package models

type NoticeConfig struct {
	NoticeType  string `bson:"notice_type,omitempty" json:"notice_type,omitempty"`
	NoticeUnit  string `bson:"notice_unit,omitempty" json:"notice_unit,omitempty"`
	NoticeValue int    `bson:"notice_value,omitempty" json:"notice_value,omitempty"`
}

type Monitor struct {
	BaseModel     `bson:",inline" json:",inline"`
	Code          string         `bson:"code,omitempty" json:"code,omitempty"`
	Id            string         `bson:"id,omitempty" json:"id,omitempty"`
	Name          string         `bson:"name,omitempty" json:"name,omitempty"`
	StartDate     string         `bson:"start_date,omitempty" json:"start_date,omitempty"`
	StartPrice    float64        `bson:"start_price,omitempty" json:"start_price,omitempty"`
	NoticeConfigs []NoticeConfig `bson:"notice_configs,omitempty" json:"notice_configs,omitempty"`
}
