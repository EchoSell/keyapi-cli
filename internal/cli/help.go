package cli

import "fmt"

const version = "1.0.0"

func printRootHelp() {
	fmt.Println(`KeyAPI CLI

KeyAPI CLI configures authentication, discovers KeyAPI endpoints, inspects
endpoint schemas, and executes API requests from one command-line tool.

Usage:
  keyapi set token <value>
  keyapi show
  keyapi spec list-platforms
  keyapi spec list-endpoints <platform...|all>
  keyapi spec get-endpoint <endpoint-or-path>
  keyapi schema <endpoint-or-path>
  keyapi call <endpoint-or-path> [flags]

Recommended order:
  1. Configure token:    keyapi set token sk_xxx
  2. List platforms:     keyapi spec list-platforms
  3. List endpoints:     keyapi spec list-endpoints tiktok instagram
  4. Inspect endpoint:   keyapi schema tiktok/influencer/detail
  5. Dry-run if unsure:  keyapi call "/v1/<platform>/<endpoint>?k=v" --dry-run
  6. Execute request:    keyapi call "/v1/<platform>/<endpoint>?k=v"

Commands:
  set/get/show/clear        Manage local configuration.
  spec list-platforms       List supported API platforms as JSON.
  spec list-endpoints       List endpoint summaries for platforms.
  spec get-endpoint         Inspect one endpoint's method, path, and parameters.
  schema                    Alias for spec get-endpoint.
  call                      Execute a KeyAPI API request.

More help:
  keyapi set --help
  keyapi show --help
  keyapi spec --help
  keyapi schema --help
  keyapi call --help

Request rules:
  Choose endpoints from spec list-endpoints, then inspect schema before calling.
  Prefer inline query strings for GET endpoints.
  Use --dry-run when request construction is uncertain.
  Use --output "path" for large responses.

Environment:
  KEYAPI_API_KEY     API key, overrides config file
  KEYAPI_HOME        Config directory, defaults to the user config dir`)
}

func printConfigHelp() {
	fmt.Println(`Manage local KeyAPI CLI configuration.

When to use:
  Use these commands before authenticated requests, when checking the current
  token, or when changing the OpenAPI documentation language used by spec output.

Usage:
  keyapi set token <value>
  keyapi set language <en|zh>
  keyapi get token
  keyapi show
  keyapi clear token

Output:
  keyapi show returns JSON.
  keyapi get returns the raw value for the requested key.
  keyapi set and keyapi clear return a short confirmation message.

Examples:
  keyapi set token sk_xxx
  keyapi show
  keyapi clear token`)
}

func printSpecHelp() {
	fmt.Println(`Discover and inspect KeyAPI OpenAPI specs.

When to use:
  Use spec commands before calling APIs. First list supported platforms, then
  list endpoint summaries for a platform, then inspect the selected endpoint.

Usage:
  keyapi spec list-platforms
  keyapi spec list-endpoints <platform...|all>
  keyapi spec get-endpoint <endpoint-or-path> [--format json|pretty]
  keyapi schema <endpoint-or-path> [--format json|pretty]

Output:
  JSON by default.

Endpoint selection:
  Read method, path, summary, and description from list-endpoints output.
  Use the endpoint whose description matches the task, then run schema or
  spec get-endpoint before constructing request parameters.

Examples:
  keyapi spec list-platforms
  keyapi spec list-endpoints google
  keyapi spec list-endpoints tiktok instagram
  keyapi schema google/search`)
}

func printSpecPlatformsHelp() {
	fmt.Println(`List supported KeyAPI platforms.

When to use:
  Use this as the first discovery command when you need to know which platforms
  currently have KeyAPI endpoints.

Usage:
  keyapi spec list-platforms

Output:
  JSON array of platform names:
  ["google", "tiktok"]`)
}

func printSpecListHelp() {
	fmt.Println(`List endpoint descriptions for one or more platforms.

When to use:
  Use this after list-platforms and before schema/call. Read method, path,
  summary, and description to decide which endpoint matches the task.

Usage:
  keyapi spec list-endpoints <platform...|all>

Arguments:
  <platform...|all>
    One platform, multiple platforms, or all.
    Examples: google, tiktok instagram, all.

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
  }

Fields:
  method: HTTP request method from OpenAPI.
  path: KeyAPI request path.
  summary: Short endpoint title.
  description: Endpoint description used to select the right API.`)
}

func printSpecShowHelp() {
	fmt.Println(`Get details for one endpoint.

When to use:
  Use after selecting an endpoint from list-endpoints and before calling it.
  This command shows the exact request method, path, query parameters, and
  request body fields defined by the KeyAPI OpenAPI document.

Usage:
  keyapi spec get-endpoint <endpoint-or-path> [--format json|pretty]

Arguments:
  endpoint-or-path
    Endpoint alias, API path, or full KeyAPI URL.
    Examples: google/search, /v1/google/search, https://api.keyapi.ai/v1/google/search

Flags:
  --format json|pretty
    json returns machine-readable JSON. pretty returns a human-readable view.
    Default: json.

Output:
  JSON by default. Fields include:
  - platform: platform name
  - method: HTTP request method
  - path: request path
  - summary: short endpoint title
  - description: detailed endpoint description
  - parameters: OpenAPI-style request parameters
  - parameters[].schema: original OpenAPI schema, including example/default/enum when present
  - body_fields: JSON body fields extracted from OpenAPI requestBody

Examples:
  keyapi spec get-endpoint google/search
  keyapi spec get-endpoint /v1/google/search --format pretty`)
}

func printSchemaHelp() {
	fmt.Println(`Get details for one endpoint.

When to use:
  Use after selecting an endpoint from list-endpoints and before calling it.
  This is a short top-level alias for keyapi spec get-endpoint.

Usage:
  keyapi schema <endpoint-or-path> [--format json|pretty]

Arguments:
  endpoint-or-path
    Endpoint alias, API path, or full KeyAPI URL.
    Examples: google/search, /v1/google/search, https://api.keyapi.ai/v1/google/search

Flags:
  --format json|pretty
    json returns machine-readable JSON. pretty returns a human-readable view.
    Default: json.

Output:
  Same as keyapi spec get-endpoint. JSON by default. Fields include method,
  path, summary, description, parameters, parameters[].schema, and body_fields.

Examples:
  keyapi schema google/search
  keyapi schema /v1/google/search --format pretty`)
}

func printCallHelp() {
	fmt.Println(`Call a KeyAPI endpoint.

When to use:
  Use after endpoint details are known from schema or spec get-endpoint. For
  uncertain requests, run the same command with --dry-run before sending it.

Usage:
  keyapi call <endpoint-or-path-with-query> [flags]

Endpoint formats:
  google/search
  /v1/google/search?q=openai&gl=us
  https://api.keyapi.ai/v1/google/search?q=openai&gl=us

Flags:
  --params "k=v"
    Add or append a URL query parameter. Repeatable.
    Quote the full key=value pair, especially when values contain spaces.

  --query "k=v"
    Alias for --params. Repeatable.

  --method METHOD
    Override the HTTP method, for example GET or POST.

  --body-json JSON
    JSON request body. Must be valid JSON and is sent as application/json.

  --body-raw TEXT
    Raw text request body. Sent as text/plain. The value is not parsed as JSON.

  --dry-run
    Print request JSON without sending the HTTP request.

  --format json|pretty|raw
    json returns the original KeyAPI response JSON. pretty formats JSON for
    reading. raw prints the raw response body. Default: json.

  --output "PATH"
    Write the response body to a file. Quote the path. Relative paths are
    resolved from the current working directory.

Output:
  Original KeyAPI response body by default. KeyAPI status fields such as code,
  message, and data are preserved. With --output, the CLI prints a small JSON
  object containing the written file path and byte count.

Exit behavior:
  HTTP 2xx returns exit code 0.
  HTTP non-2xx prints or writes the response body and returns a non-zero exit code.

Examples:
  keyapi call "/v1/google/search?q=openai&gl=us&hl=en"
  keyapi call google/search --params "q=openai" --params "gl=us"
  keyapi call "/v1/example/post" --method POST --body-json '{"page":1}'
  keyapi call "/v1/example/text" --method POST --body-raw "plain text"
  keyapi call "/v1/google/search?q=openai" --output "outputs/google.json"
  keyapi call "/v1/google/search?q=openai" --dry-run`)
}
