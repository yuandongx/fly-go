package spider

import (
	"encoding/json"
	"fly-go/internal/models"
	"fmt"
)

func GetFundInfo(data interface{}) error {
	query := QueryParams{
		"page":   1,
		"num":    50,
		"sort":   "symbol",
		"asc":    1,
		"symbol": "",
		"node":   "",
		"_s_r_a": "init",
	}
	headers := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
		"User-Agent":   "insomnia/9.3.0-beta.6",
		"Connection":   "keep-alive",
	}
	// var rows []models.
	var count = 0
	for nodeKey, _ := range Nodes {
		req := NewRequest(SinaStockCountUrl+nodeKey, "GET", headers, nil, 10)
		data, err := req.Get()
		if err != nil {
			return err
		}
		if total, err := CovertInt(data); err != nil {
			return err
		} else {
			count = total
		}

		for page := 1; page <= (count/50)+1; page++ {
			query["page"] = page
			req.SetUrl(SinaStockListUrl + "?" + query.String())
			rows, err := req.Get()
			if err != nil {
				continue
			}
			rowList := []models.Stock{}
			err = json.Unmarshal(rows, &rowList)
			if err != nil {
				continue
			}
		}
	}
	fmt.Println("Total funds:", count)
	return nil
}
