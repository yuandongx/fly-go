package handlers

import (
	"fly-go/database"
)

// StockResource 股票资源处理器
type StockResource struct {
	*BaseHandler
	collection string
}

func NewStockResource(db *database.MongoDB) *StockResource {
	return &StockResource{
		BaseHandler: NewBaseHandler(db),
		collection:  "stock",
	}
}

func (r *StockResource) Name() string {
	return "stock"
}

func (r *StockResource) Routes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "", r.List(r.collection), "获取股票列表"},
	}
}

func (r *StockResource) CustomRoutes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "/:id", r.GetByID(r.collection), "获取股票详情"},
	}
}
