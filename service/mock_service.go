package service

import (
	"fmt"
	"strings"
)

// RenderTemplateRecursive ...
func RenderTemplateRecursive(data interface{}, bodyMap, headerMap, queryMap, pathMap map[string]string) interface{} {
	switch v := data.(type) {
	case string:
		for k, val := range bodyMap {
			v = strings.ReplaceAll(v, "{{body."+k+"}}", val)
		}
		for k, val := range headerMap {
			v = strings.ReplaceAll(v, "{{header."+k+"}}", val)
		}
		for k, val := range queryMap {
			v = strings.ReplaceAll(v, "{{query."+k+"}}", val)
		}
		for k, val := range pathMap {
			v = strings.ReplaceAll(v, "{{path."+k+"}}", val)
		}
		return v
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = RenderTemplateRecursive(val, bodyMap, headerMap, queryMap, pathMap)
		}
		return result
	case []interface{}:
		for i, val := range v {
			v[i] = RenderTemplateRecursive(val, bodyMap, headerMap, queryMap, pathMap)
		}
	}
	return data
}

// MapToStringMap ...
func MapToStringMap(input map[string]interface{}) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = fmt.Sprintf("%v", v)
	}
	return output
}
