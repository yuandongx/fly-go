// Package fly, the spider
package fly

import (
	"fly-go/fly/spider"
)

type Stock struct {
	TaskBase
}

func NewStock() *Stock {
	return &Stock{}
}
func (s Stock) Run() ([]BM, error) {
	return spider.GetStockInfo()
}

func (s Stock) Stop() error {
	return nil
}
