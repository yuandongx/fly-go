// Package spider, the base
package spider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SinaBaseURL       = "https://vip.stock.finance.sina.com.cn"
	SinaStockCountURL = SinaBaseURL + "/quotes_service/api/json_v2.php/Market_Center.getHQNodeStockCount?node="
	SinaStockListURL  = SinaBaseURL + "/quotes_service/api/json_v2.php/Market_Center.getHQNodeData"
)

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
		Client: http.Client{
			Transport: &http.Transport{},
			Timeout:   30 * time.Second,
		},
	}
}
func (r *Request) SetURL(url string) {
	r.URL = url
}
func (r *Request) Get() ([]byte, error) {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
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
		rtn += fmt.Sprintf("%s=%v&", k, v)
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
