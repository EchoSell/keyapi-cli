package cli

type dryRunRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type outputFileResult struct {
	Output string `json:"output"`
	Bytes  int    `json:"bytes"`
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
