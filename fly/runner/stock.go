// Package runner, task runner
// 任务的执行单元
package runner

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
