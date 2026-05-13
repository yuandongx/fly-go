package spider

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var Nodes = map[string]string{
	"sgt_sz": "sz",
	"hgt_sh": "sh",
	"hs_bjs": "bj",
	"kcb":    "kc",
	"cyb":    "cy",
}

// GetStockInfo 获取股票信息
func GetStockInfo() ([]bson.M, error) {
	var result []bson.M
	var lastErr error

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
		"Accept":     "*/*",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Connection": "keep-alive",
	}

	for nodeKey, nodeValue := range Nodes {
		// 获取总数
		countURL := SinaStockCountURL + nodeKey
		req := NewRequest(countURL, "GET", headers, nil, 10)
		data, err := req.Get()
		if err != nil {
			fmt.Printf("%s: Failed to get count\n", nodeValue)
			lastErr = err
			continue
		}

		count, err := CovertInt(data)
		if err != nil {
			fmt.Printf("%s: Failed to parse count\n", nodeValue)
			lastErr = err
			continue
		}

		fmt.Printf("%s: Expecting %d records\n", nodeValue, count)

		// 分页获取数据
		for page := 1; page <= (count/50)+1; page++ {
			query["page"] = page
			query["node"] = nodeKey
			req.SetURL(SinaStockListURL + "?" + query.String())

			rows, err := req.Get()
			if err != nil {
				lastErr = err
				continue
			}

			var tmp []bson.M
			if err := json.Unmarshal(rows, &tmp); err != nil {
				lastErr = err
				continue
			}
			result = append(result, tmp...)
		}
	}

	fmt.Printf("Total stocks fetched: %d\n", len(result))
	return result, lastErr
}
