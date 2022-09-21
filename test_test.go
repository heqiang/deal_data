package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTrans(t *testing.T) {
	testStr := `["aa","bb","vv"]`
	var res []string
	_ = json.Unmarshal([]byte(testStr), &res)
	fmt.Println(res)
}
