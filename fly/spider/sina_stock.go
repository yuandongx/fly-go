package spider

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var Nodes = map[string]string{"sgt_sz": "sz",
	"hgt_sh": "sh",
	"hs_bjs": "bj",
	"kcb":    "kc",
	"cyb":    "cy"}

func GetStockInfo() ([]bson.M, error) {
	var result []bson.M
	var msg error
	query := QueryParams{
		"page":   1,
		"num":    40,
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
	for nodeKey, nodeValue := range Nodes {
		count := 0
		req := NewRequest(SinaStockCountURL+nodeKey, "GET", headers, nil, 10)

		if data, err := req.Get(); err != nil {
			fmt.Printf("%s 查询数量失败！", nodeValue)
			msg = err
			continue
		} else {
			if total, err := CovertInt(data); err != nil {
				msg = err
				continue
			} else {
				count = total
				fmt.Printf("%s 应有 %d 条数据 \n", nodeValue, total)
			}
		}
		for page := 1; page <= (count/50)+1; page++ {
			query["page"] = page
			query["node"] = nodeKey
			req.SetURL(SinaStockListURL + "?" + query.String())
			rows, err := req.Get()
			if err != nil {
				msg = err
				continue
			}
			var tmp []bson.M
			err = json.Unmarshal(rows, &tmp)
			if err != nil {
				msg = err
				continue
			}
			result = append(result, tmp...)
		}
	}
	fmt.Println("Total stocks:", len(result))
	return result, msg
}
