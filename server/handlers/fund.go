package handlers

import (
	"fly-go/database"
)

// FundResource 基金资源处理器
type FundResource struct {
	*BaseHandler
	collection string
}

func NewFundResource(db *database.MongoDB) *FundResource {
	return &FundResource{
		BaseHandler: NewBaseHandler(db),
		collection:  "fund",
	}
}

func (r *FundResource) Name() string {
	return "fund"
}

func (r *FundResource) Routes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "", r.List(r.collection), "获取基金列表"},
	}
}

func (r *FundResource) CustomRoutes() []RouteConfig {
	return []RouteConfig{
		{MethodGet, "/:id", r.GetByID(r.collection), "获取基金详情"},
	}
}
