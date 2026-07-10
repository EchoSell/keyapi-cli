package cli

type apiResult struct {
	OK       bool        `json:"ok"`
	Endpoint string      `json:"endpoint,omitempty"`
	Method   string      `json:"method,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type dryRunRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type configView struct {
	APIKey     string `json:"api_key"`
	AuthScheme string `json:"auth_scheme"`
	Language   string `json:"language"`
}

type endpointSummary struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}
