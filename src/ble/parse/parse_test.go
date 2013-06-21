package parse

import (
	"encoding/json"
	. "testing"
)

func TestParse(t *T) {
	case1 := []byte("( / (   + 12 (* 3   5 ) (- 8 7) 4))some leftover text")
	expr, left, err := parseExpr(case1)
	t.Log(expr)
	t.Log(string(left))
	t.Log(err)
	if expr != nil {
		json, err := json.Marshal(expr)
		t.Log(string(json))
		t.Log(err)
	}
}
