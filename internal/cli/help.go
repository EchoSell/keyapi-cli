package cli

import "fmt"

const version = "1.0.0"

func printRootHelp() {
	fmt.Println(`KeyAPI CLI

Usage:
  keyapi set token <value>
  keyapi show
  keyapi spec list-platforms
  keyapi spec list-endpoints <platform...|all>
  keyapi spec get-endpoint <endpoint-or-path>
  keyapi call <endpoint-or-path> [flags]

Core workflow:
  1. Configure token:    keyapi set token sk_xxx
  2. List platforms:     keyapi spec list-platforms
  3. List endpoints:     keyapi spec list-endpoints tiktok instagram
  4. Inspect endpoint:   keyapi spec get-endpoint tiktok/influencer/detail
  5. Execute request:    keyapi call "/v1/<platform>/<endpoint>?k=v"

Environment:
  KEYAPI_API_KEY     API key, overrides config file
  KEYAPI_HOME        Config directory, defaults to the user config dir`)
}

func printConfigHelp() {
	fmt.Println(`Manage local KeyAPI CLI configuration.

When to use:
  Use before authenticated API requests, or to set the OpenAPI docs language.

Usage:
  keyapi set token <value>
  keyapi set language <en|zh>
  keyapi get token
  keyapi show
  keyapi clear token

Output:
  JSON for show; plain text for single-key get.`)
}

func printSpecHelp() {
	fmt.Println(`Discover and inspect KeyAPI OpenAPI specs.

When to use:
  Use before calling APIs. First list supported platforms, then list endpoint descriptions, then inspect one endpoint's parameters.

Usage:
  keyapi spec list-platforms
  keyapi spec list-endpoints <platform...|all>
  keyapi spec get-endpoint <endpoint-or-path> [--format json|pretty]

Output:
  JSON by default.

Search guidance:
  There is no keyword search command. To choose an API, list endpoint descriptions for the relevant platform and select by method, path, summary, and description.`)
}

func printSpecPlatformsHelp() {
	fmt.Println(`List supported KeyAPI platforms.

When to use:
  Use when you need to know which platforms currently have KeyAPI endpoints.

Usage:
  keyapi spec list-platforms

Output:
  JSON array:
  ["google", "tiktok"]`)
}

func printSpecListHelp() {
	fmt.Println(`List endpoint descriptions for one or more platforms.

When to use:
  Use this before selecting an endpoint. Read method, path, summary, and description to decide which endpoint matches the task.

Usage:
  keyapi spec list-endpoints <platform...|all>

Examples:
  keyapi spec list-endpoints tiktok
  keyapi spec list-endpoints tiktok instagram
  keyapi spec list-endpoints all

Output:
  JSON object keyed by platform:
  {
    "tiktok": [
      {
        "method": "GET",
        "path": "/v1/tiktok/influencer/detail",
        "summary": "User detail",
        "description": "Get TikTok user detail."
      }
    ]
  }`)
}

func printSpecShowHelp() {
	fmt.Println(`Get details for one endpoint.

When to use:
  Use after selecting an endpoint from list-endpoints and before calling it.

Usage:
  keyapi spec get-endpoint <endpoint-or-path> [--format json|pretty]

Parameters:
  endpoint-or-path: google/search, /v1/google/search, or a full KeyAPI URL.
  --format: json or pretty. Default is json.

Output:
  JSON by default. Fields include:
  - platform: platform name
  - method: HTTP request method
  - path: request path
  - summary: short endpoint title
  - description: detailed endpoint description
  - parameters: OpenAPI-style request parameters
  - parameters[].schema: original OpenAPI schema, including example/default/enum when present
  - body_fields: JSON body fields extracted from OpenAPI requestBody`)
}

func printCallHelp() {
	fmt.Println(`Call a KeyAPI endpoint by path convention.

When to use:
  Use after endpoint details are known from spec get-endpoint.

Usage:
  keyapi call <endpoint-or-path-with-query> [--params "k=v"] [--body-raw TEXT] [--method GET] [--body-json JSON]

Parameters:
  endpoint-or-path-with-query: google/search, /v1/google/search?q=openai, or a full KeyAPI URL.
  --query: Add or append a URL query parameter as "key=value". Repeatable.
  --params: Add or append a URL query parameter as "key=value". Repeatable.
  --body-raw: Raw text request body. Sent as text/plain without JSON validation.
  --body-json: JSON request body. Must be valid JSON and is sent as application/json. PowerShell-style {q:openai} input is normalized when possible.
  --method: HTTP method override.
  --dry-run: Print request JSON without sending.
  --format: json, pretty, or raw. Default is json.

Output:
  Original KeyAPI response JSON by default, formatted JSON with --format pretty, raw response with --format raw.
  HTTP non-2xx responses are printed and return a non-zero exit code.`)
}
