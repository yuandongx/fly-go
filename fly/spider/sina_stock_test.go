package spider_test

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	Name string
	ID   string
}

func TestGetFundInfo(t *testing.T) {
	a := &A{"TestA", "a1"}
	rt := reflect.TypeOf(a)
	rv := reflect.ValueOf(a)
	fmt.Println(rt.Kind() == reflect.Pointer)
	fmt.Println(rv.Kind())
	fmt.Println(rt.Elem().Kind())
	fmt.Println(rv.Elem().Kind())
}
