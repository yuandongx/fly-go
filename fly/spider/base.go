// Package spider, the base
package spider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	SinaBaseUrl       = "https://vip.stock.finance.sina.com.cn"
	SinaStockCountUrl = SinaBaseUrl + "/quotes_service/api/json_v2.php/Market_Center.getHQNodeStockCount?node="
	SinaStockListUrl  = SinaBaseUrl + "/quotes_service/api/json_v2.php/Market_Center.getHQNodeData"
)

var Nodes = map[string]string{"sgt_sz": "sz",
	"hgt_sh": "sh",
	"hs_bjs": "bj",
	"kcb":    "kc",
	"cyb":    "cy"}

type QueryParams map[string]interface{}
type Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Args    map[string]string
	Timeout int
	Client  http.Client
}

func NewRequest(url string, method string, headers map[string]string, args map[string]string, timeout int) *Request {
	return &Request{
		URL:     url,
		Method:  method,
		Headers: headers,
		Args:    args,
		Timeout: timeout,
		Client:  http.Client{},
	}
}
func (r *Request) SetUrl(url string) {
	r.URL = url
}
func (r *Request) Get() ([]byte, error) {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		return nil, err
	}
	// for k, v := range r.Headers {
	// 	req.Header.Set(k, v)
	// }
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 从body读取数据判断是不是json,如果是则解析为data
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *Request) GetJSON(v interface{}) error {
	body, err := r.Get()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func (q QueryParams) String() string {
	rtn := ""
	for k, v := range q {
		rtn += fmt.Sprintf("%s=%s&", k, v)
	}
	return rtn
}

func CovertInt(data []byte) (int, error) {
	var value int
	v := strings.Trim(string(data), "\\\n\r\" ")
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	value = int(i)
	return value, nil
}
