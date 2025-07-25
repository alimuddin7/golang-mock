package model

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
}
