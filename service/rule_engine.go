package service

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang-mock/model"
)

// RequestContext holds the request data for rule evaluation
type RequestContext struct {
	Body       map[string]string
	Headers    map[string]string
	Query      map[string]string
	PathParams map[string]string
}

var regexCache sync.Map

// EvaluateRules evaluates all rules with the given operator (AND/OR)
func EvaluateRules(rules []model.Rule, operator string, ctx RequestContext) bool {
	if len(rules) == 0 {
		return true // No rules means always match
	}

	results := make([]bool, len(rules))
	for i, rule := range rules {
		results[i] = evaluateRule(rule, ctx)
	}

	if strings.ToUpper(operator) == "OR" {
		return containsTrue(results)
	}
	// Default to AND
	return allTrue(results)
}

// evaluateRule evaluates a single rule against the request context
func evaluateRule(rule model.Rule, ctx RequestContext) bool {
	var actualValue string

	// Get the actual value from the appropriate source
	switch strings.ToLower(rule.Target) {
	case "body":
		actualValue = ctx.Body[rule.Field]
	case "header":
		actualValue = ctx.Headers[rule.Field]
	case "query":
		actualValue = ctx.Query[rule.Field]
	case "path":
		actualValue = ctx.PathParams[rule.Field]
	default:
		return false
	}

	// Evaluate based on operator
	switch strings.ToLower(rule.Operator) {
	case "equals", "==":
		return actualValue == rule.Value
	case "not_equals", "!=":
		return actualValue != rule.Value
	case "contains", "~":
		return strings.Contains(actualValue, rule.Value)
	case "regex", ".*":
		var re *regexp.Regexp
		if v, ok := regexCache.Load(rule.Value); ok {
			re = v.(*regexp.Regexp)
		} else {
			var err error
			re, err = regexp.Compile(rule.Value)
			if err != nil {
				return false
			}
			regexCache.Store(rule.Value, re)
		}
		return re.MatchString(actualValue)
	case "exists", "?":
		// For exists, we check if the field has a non-empty value
		// The rule.Value can be "true" or "false"
		exists := actualValue != ""
		expectedExists := strings.ToLower(rule.Value) != "false"
		return exists == expectedExists
	case "gt", ">":
		// Greater than - try to parse as numbers
		actual, err1 := strconv.ParseFloat(actualValue, 64)
		expected, err2 := strconv.ParseFloat(rule.Value, 64)
		if err1 != nil || err2 != nil {
			// Fallback to string comparison
			return actualValue > rule.Value
		}
		return actual > expected
	case "lt", "<":
		// Less than - try to parse as numbers
		actual, err1 := strconv.ParseFloat(actualValue, 64)
		expected, err2 := strconv.ParseFloat(rule.Value, 64)
		if err1 != nil || err2 != nil {
			// Fallback to string comparison
			return actualValue < rule.Value
		}
		return actual < expected
	case "gte", ">=":
		// Greater than or equal
		actual, err1 := strconv.ParseFloat(actualValue, 64)
		expected, err2 := strconv.ParseFloat(rule.Value, 64)
		if err1 != nil || err2 != nil {
			return actualValue >= rule.Value
		}
		return actual >= expected
	case "lte", "<=":
		// Less than or equal
		actual, err1 := strconv.ParseFloat(actualValue, 64)
		expected, err2 := strconv.ParseFloat(rule.Value, 64)
		if err1 != nil || err2 != nil {
			return actualValue <= rule.Value
		}
		return actual <= expected
	default:
		return false
	}
}

// allTrue returns true if all values in the slice are true
func allTrue(values []bool) bool {
	for _, v := range values {
		if !v {
			return false
		}
	}
	return true
}

// containsTrue returns true if any value in the slice is true
func containsTrue(values []bool) bool {
	for _, v := range values {
		if v {
			return true
		}
	}
	return false
}
