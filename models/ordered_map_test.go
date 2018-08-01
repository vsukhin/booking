package models

import (
	"encoding/json"
	"testing"
)

func Test_OrderedMap_MarshalJSON_Success(t *testing.T) {
	var orderedMap = OrderedMap{
		{
			"testKey1",
			"testValue1",
		},
		{
			"testKey2",
			"testValue2",
		},
	}

	text, err := json.Marshal(orderedMap)
	if err != nil {
		t.Error("Expected marshal ordered map successfully")
	}
	if string(text) != "{\"testKey1\":\"testValue1\",\"testKey2\":\"testValue2\"}" {
		t.Error("Expected to have matching json")
	}
}
