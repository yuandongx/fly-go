package spider_test

import (
	"fly-go/fly/spider"
	"fmt"
	"testing"
)

func TestGetFundInfo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test GetFundInfo",
			data: &[]map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Start test " + tt.name)
			gotErr := spider.GetFundInfo(tt.data)
			fmt.Printf("Result: %+v\n", tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetFundInfo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetFundInfo() succeeded unexpectedly")
			}
		})
	}
}
