package models

import (
	"bytes"
	"fmt"
	"strconv"
)

// OrderedMap is ordered map
type OrderedMap []struct {
	Key string
	Val interface{}
}

// MarshalJSON is json marshall
func (omap OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, kv := range omap {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Quote(kv.Key) + ":" + strconv.Quote(fmt.Sprintf("%v", kv.Val)))
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}
