package service

import (
	"testing"
)

func TestRenderTemplateRecursive(t *testing.T) {
	template := "Hello {{body.name}}, your token is {{header.X-Token}} in {{query.env}}"
	bodyMap := map[string]string{"name": "Alice"}
	headerMap := map[string]string{"X-Token": "secret123"}
	queryMap := map[string]string{"env": "prod"}
	pathMap := map[string]string{}

	result := RenderTemplateRecursive(template, bodyMap, headerMap, queryMap, pathMap)

	expected := "Hello Alice, your token is secret123 in prod"
	if result.(string) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestRenderTemplateRecursive_NoMatch(t *testing.T) {
	template := "Hello {{body.missing}}"
	bodyMap := map[string]string{"name": "Alice"}

	result := RenderTemplateRecursive(template, bodyMap, nil, nil, nil)

	expected := "Hello {{body.missing}}"
	if result.(string) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
