package builder

import (
	"fmt"
	"testing"

	"github.com/antlabs/deepcopy"
)

func TestCopy(t *testing.T) {
	q1 := &Query{}
	q2 := &Query{
		Fields: []string{"1", "2", "3"},
	}
	deepcopy.Copy(q1, q2).Do()
	fmt.Println(q1.Fields)
	if len(q1.Fields) != 3 {
		t.Error("copy is not valid")
	}
}
