package template

import (
	"encoding/json"
	"html/template"

	"golang-mock/model"
)

// ToJSONPretty ...
func ToJSONPretty(v interface{}) template.HTML {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return template.HTML(b)
}

// HTTPMethods ...
func HTTPMethods() []string {
	return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
}

// ConfigLen ...
func ConfigLen(arr []model.MockConfig) int {
	return len(arr)
}
