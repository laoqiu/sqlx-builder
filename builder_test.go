package builder

import (
	"fmt"
	"testing"
)

type City struct {
	name string
}
type Address struct {
	c *City
}

func TestStructToMap(t *testing.T) {
	addr := &Address{}
	fmt.Println(addr.c == nil)
	m := StructToMap(addr.c)
	if len(m) != 0 {
		t.Error("错误")
	}
}
