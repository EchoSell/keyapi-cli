# Command Reference

This document lists the current KeyAPI CLI commands, arguments, flags, and JSON output fields.

## Root

| Command | Description | Output |
| --- | --- | --- |
| `keyapi --help` | Show root help and core workflow. | Plain text |
| `keyapi help` | Alias for root help. | Plain text |
| `keyapi --version` | Print CLI version. | Plain text |
| `keyapi version` | Alias for version. | Plain text |

## Top-Level Config Commands

Config commands keep their existing human-readable behavior. `keyapi show` returns JSON; single-key operations return plain text.

| Command | Description | Required Arguments | Output |
| --- | --- | --- | --- |
| `keyapi show` | Show current config. The token is masked. | None | JSON |
| `keyapi get` | Show current config. Same as `show` when no key is provided. | None | JSON |
| `keyapi get token` | Show masked token. | `token` | Plain text |
| `keyapi get api-key` | Show masked token. Alias for `token`. | `api-key` | Plain text |
| `keyapi get auth-scheme` | Show authorization scheme. | `auth-scheme` | Plain text |
| `keyapi get language` | Show OpenAPI docs language. | `language` | Plain text |
| `keyapi get lang` | Show OpenAPI docs language. Alias for `language`. | `lang` | Plain text |
| `keyapi set token <value>` | Save KeyAPI token locally. | `token`, `value` | `saved` |
| `keyapi set api-key <value>` | Save KeyAPI token locally. Alias for `token`. | `api-key`, `value` | `saved` |
| `keyapi set auth-scheme <value>` | Save authorization scheme. Default is `Bearer`. | `auth-scheme`, `value` | `saved` |
| `keyapi set language <en\|zh>` | Save OpenAPI docs language. | `language`, `en` or `zh` | `saved` |
| `keyapi set lang <en\|zh>` | Save OpenAPI docs language. Alias for `language`. | `lang`, `en` or `zh` | `saved` |
| `keyapi clear token` | Remove stored token. | `token` | `cleared` |
| `keyapi clear api-key` | Remove stored token. Alias for `token`. | `api-key` | `cleared` |

## `keyapi spec`

`spec` commands return JSON by default.

The CLI reads the latest KeyAPI OpenAPI documentation each time and returns real-time API details.

| Command | Description | Required Arguments | Flags | Output |
| --- | --- | --- | --- | --- |
| `keyapi spec --help` | Show spec help. | None | None | Plain text |
| `keyapi spec list-platforms` | List supported platforms. | None | None | JSON |
| `keyapi spec list-endpoints <platform...>` | List endpoint descriptions for one or more platforms. | One or more platform names | None | JSON |
| `keyapi spec list-endpoints all` | List endpoint descriptions for all platforms. | `all` | None | JSON |
| `keyapi spec get-endpoint <endpoint-or-path>` | Get one endpoint's full details. | `endpoint-or-path` | `--format json\|pretty` | JSON by default |

### `spec list-platforms`

Usage:

```powershell
keyapi spec list-platforms
```

Output:

```json
["google", "tiktok"]
```

| Data | Meaning |
| --- | --- |
| Array item | Platform name that currently has OpenAPI specs. |

### `spec list-endpoints`

Usage:

```powershell
keyapi spec list-endpoints tiktok
keyapi spec list-endpoints tiktok instagram
keyapi spec list-endpoints all
```

Output shape:

```json
{
  "tiktok": [
    {
      "method": "GET",
      "path": "/v1/tiktok/influencer/detail",
      "summary": "User detail",
      "description": "Get TikTok user detail."
    }
  ],
  "instagram": []
}
```

| Field | Meaning |
| --- | --- |
| Top-level platform key | Platform name, such as `tiktok` or `instagram`. |
| `method` | HTTP request method. |
| `path` | API request path under `https://api.keyapi.ai`. |
| `summary` | Short endpoint title. |
| `description` | Longer endpoint description. Use this field to decide which endpoint matches the task. |

### `spec get-endpoint`

Usage:

```powershell
keyapi spec get-endpoint google/search
keyapi spec get-endpoint /v1/google/search
keyapi spec get-endpoint google/search --format pretty
```

Default JSON output shape:

```json
{
  "platform": "google",
  "spec_name": "google.openapi",
  "method": "GET",
  "path": "/v1/google/search",
  "summary": "Search",
  "description": "Search Google web results.",
  "parameters": [
    {
      "name": "q",
      "in": "query",
      "type": "string",
      "required": true,
      "description": "Search keyword.",
      "schema": {
        "type": "string",
        "example": "openai"
      }
    }
  ],
  "body_fields": []
}
```

| Field | Meaning |
| --- | --- |
| `platform` | Platform name. |
| `spec_name` | Source OpenAPI spec name. |
| `method` | HTTP request method. |
| `path` | API request path. |
| `summary` | Short endpoint title. |
| `description` | Detailed endpoint description. |
| `parameters` | OpenAPI-style request parameters. |
| `body_fields` | JSON request body fields extracted from OpenAPI `requestBody`. |

Parameter fields:

| Field | Meaning |
| --- | --- |
| `name` | Parameter name. |
| `in` | Parameter location: `query`, `path`, or `header`. |
| `type` | Schema type. |
| `required` | Whether the parameter is required. |
| `description` | Parameter description. |
| `schema` | Original OpenAPI parameter schema. Preserves metadata such as `example`, `default`, and `enum`. |

## `keyapi call`

`call` executes a KeyAPI endpoint. Default output is JSON.

| Command | Description | Required Arguments | Flags | Output |
| --- | --- | --- | --- | --- |
| `keyapi call --help` | Show call help. | None | None | Plain text |
| `keyapi call <endpoint-or-path>` | Execute a KeyAPI endpoint request. | `endpoint-or-path` | See below | JSON by default |

Endpoint formats:

```text
google/search
/v1/google/search?q=openai&gl=us
https://api.keyapi.ai/v1/google/search?q=openai&gl=us
```

### `call` Flags

| Flag | Required | Repeatable | Description |
| --- | --- | --- | --- |
| `--query "k=v"` | No | Yes | Add or append a URL query parameter. Always quote the full `k=v` value. |
| `--params "k=v"` | No | Yes | Add or append a URL query parameter. Useful when query parameters are appended separately from the endpoint path. Always quote the full `k=v` value. |
| `--body-raw <text>` | No | No | Send a raw text request body with `Content-Type: text/plain; charset=utf-8`. |
| `--body-json <json>` | No | No | Send a JSON request body with `Content-Type: application/json`. The value must be valid JSON. |
| `--method <method>` | No | No | Override HTTP method. If omitted, `POST` is used when body data exists; otherwise `GET` is used. |
| `--api-key <value>` | No | No | Use a token for this command only. Overrides config and environment. |
| `--token <value>` | No | No | Alias for `--api-key`. |
| `--dry-run` | No | No | Print the final request as JSON without sending it. |
| `--format json\|pretty\|raw` | No | No | Select response output format. Default is `json`. |
| `--output "path"` | No | No | Write the response body to a file. Quote the path. Relative paths are resolved from the current working directory. Parent directories are created automatically. |

Default output:

`keyapi call` prints the original response returned by the KeyAPI endpoint. The CLI does not add an extra wrapper around the API response.

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| Field | Meaning |
| --- | --- |
| `code` | KeyAPI endpoint status code. `0` means success. |
| `message` | KeyAPI endpoint status message. |
| `data` | KeyAPI endpoint response data. |

Common KeyAPI response codes:

| Code | Message | Meaning |
| --- | --- | --- |
| `0` | `success` | Request succeeded. |
| `400` | `Missing required parameter: <name>` | A required parameter is missing or invalid. |
| `401` | `apiKey is null` | API key is missing. |
| `402` | `apiKey quota is not enough` | API key quota is insufficient. |
| `429` | `qps limit exceeded` | Request rate exceeded the QPS limit. |
| `500` | `Internal server error` | Server-side error. |

These fields and status codes are returned by KeyAPI itself. The CLI keeps them consistent with the API response.

File output:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us" --output "result/google.json"
```

When `--output` is used, the response body is written to the file and the terminal prints a small JSON object containing the resolved output path and byte count. HTTP non-`2xx` responses are still written before the command exits with a non-zero code.

For GET endpoints, prefer putting query parameters directly in the endpoint path:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&page=1"
```

`--params` and `--query` are also supported and can be mixed with inline query parameters:

```powershell
keyapi call "/v1/google/search?q=openai" --params "gl=us" --params "page=1"
```

HTTP status handling:

- HTTP `2xx` responses print the KeyAPI response and exit with code `0`.
- HTTP non-`2xx` responses still print the KeyAPI response, then exit with a non-zero code so scripts and CI can detect failures.

## Endpoint Selection Guidance

Use this process to choose an endpoint:

```text
List endpoint descriptions for the relevant platform.
Compare method, path, summary, and description.
Choose the endpoint whose description best matches the user request.
Inspect that endpoint with get-endpoint.
Use parameters and body_fields to construct the call command.
```

Selection rules:

- Choose endpoints by reading `method`, `path`, `summary`, and `description`.
- Inspect `parameters` and `body_fields` before constructing a request.
- For GET endpoints, put query parameters in the endpoint path.
- Use `--params "k=v"` or `--query "k=v"` when query parameters need to be appended separately.
- Use `--body-json` for JSON request bodies.
- Use `--body-raw` for raw text request bodies.
- Use `--output "path"` for large responses or saved artifacts.
- Use `--dry-run` before execution when request construction is uncertain.
