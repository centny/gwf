package handler

import (
	"encoding/json"
	"testing"
)

func TestHandler(t *testing.T) {
	NewRC_S_Json_M_S(nil)
	NewChan_Json_S(nil)
	NewRC_S_V_M_S(nil, json.Marshal, json.Unmarshal)
}
