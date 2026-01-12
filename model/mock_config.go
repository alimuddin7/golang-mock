package model

// Rule represents a condition to evaluate
type Rule struct {
	Target   string `json:"target"`   // "body", "header", "query", "path"
	Field    string `json:"field"`    // field name to check
	Operator string `json:"operator"` // "equals", "contains", "regex", "exists", "gt", "lt"
	Value    string `json:"value"`    // value to compare against
}

// Response represents a mock response configuration
type Response struct {
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"`
	StatusCode int                    `json:"statusCode"`
	Timeout    int                    `json:"timeout"`
}

// ConditionalResponse represents a response with conditions
type ConditionalResponse struct {
	Name         string   `json:"name"`
	Rules        []Rule   `json:"rules"`
	RuleOperator string   `json:"ruleOperator"` // "AND" or "OR"
	Response     Response `json:"response"`
}

// MockConfig ...
type MockConfig struct {
	Name            string                 `json:"name"`
	Method          string                 `json:"method"`
	Path            string                 `json:"path"`
	RequestHeaders  map[string]string      `json:"requestHeaders"`
	RequestBody     map[string]interface{} `json:"requestBody"`
	ResponseHeaders map[string]string      `json:"responseHeaders"`
	ResponseBody    map[string]interface{} `json:"responseBody"`
	StatusCode      int                    `json:"statusCode"`
	Timeout         int                    `json:"timeout"`

	// New fields for conditional logic
	Responses       []ConditionalResponse `json:"responses,omitempty"`
	DefaultResponse *Response             `json:"defaultResponse,omitempty"`
}
