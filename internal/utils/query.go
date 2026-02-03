package utils

type BaseQuery struct {
	Page    int    `json:"page" form:"page" validate:"required,min=1" default:"1"`
	Size    int    `json:"size" form:"size" validate:"required,min=1,max=100" default:"10"`
	Search  string `json:"search" form:"search" validate:"omitempty"`
	OrderBy string `json:"order_by" form:"order_by" validate:"omitempty"`
	Order   string `json:"order" form:"order" validate:"omitempty"`
}

type Query struct {
	BaseQuery `json:",inline"`
	ID        string `json:"id" form:"id" validate:"required"`
}
