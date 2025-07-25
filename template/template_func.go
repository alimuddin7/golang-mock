package template

import (
	"encoding/json"

	"golang-mock/model"
)

// ToJSONPretty ...
func ToJSONPretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

// HTTPMethods ...
func HTTPMethods() []string {
	return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
}

// ConfigLen ...
func ConfigLen(arr []model.MockConfig) int {
	return len(arr)
}
