package leq

import (
	"encoding/json"
)

type LeqConfig struct {
	URL    string          `json:"url"`
	Method string          `json:"method"`
	Header json.RawMessage `json:"header"`
	Body   json.RawMessage `json:"body"`
}

