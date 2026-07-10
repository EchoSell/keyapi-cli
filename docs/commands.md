# Command Reference

This document lists the current KeyAPI CLI commands, arguments, flags, and JSON output fields.

## Root

| Command | Description | Output |
| --- | --- | --- |
| `keyapi --help` | Show root help and core workflow. | Plain text |
| `keyapi help` | Alias for root help. | Plain text |
| `keyapi --version` | Print CLI version. | Plain text |
| `keyapi version` | Alias for version. | Plain text |

## `keyapi config`

`config` commands keep their existing human-readable behavior. `config show` returns JSON; single-key operations return plain text.

| Command | Description | Required Arguments | Output |
| --- | --- | --- | --- |
| `keyapi config --help` | Show config help. | None | Plain text |
| `keyapi config show` | Show current config. The token is masked. | None | JSON |
| `keyapi config get` | Show current config. Same as `show` when no key is provided. | None | JSON |
| `keyapi config get token` | Show masked token. | `token` | Plain text |
| `keyapi config get api-key` | Show masked token. Alias for `token`. | `api-key` | Plain text |
| `keyapi config get auth-scheme` | Show authorization scheme. | `auth-scheme` | Plain text |
| `keyapi config get language` | Show OpenAPI docs language. | `language` | Plain text |
| `keyapi config get lang` | Show OpenAPI docs language. Alias for `language`. | `lang` | Plain text |
| `keyapi config set token <value>` | Save KeyAPI token locally. | `token`, `value` | `saved` |
| `keyapi config set api-key <value>` | Save KeyAPI token locally. Alias for `token`. | `api-key`, `value` | `saved` |
| `keyapi config set auth-scheme <value>` | Save authorization scheme. Default is `Bearer`. | `auth-scheme`, `value` | `saved` |
| `keyapi config set language <en\|zh>` | Save OpenAPI docs language. | `language`, `en` or `zh` | `saved` |
| `keyapi config set lang <en\|zh>` | Save OpenAPI docs language. Alias for `language`. | `lang`, `en` or `zh` | `saved` |
| `keyapi config clear token` | Remove stored token. | `token` | `cleared` |
| `keyapi config clear api-key` | Remove stored token. Alias for `token`. | `api-key` | `cleared` |

## `keyapi spec`

`spec` commands return JSON by default.

The CLI always fetches the latest `https://docs.keyapi.ai/llms.txt`, parses OpenAPI links from its `## OpenAPI Specs` section, and fetches the selected OpenAPI JSON directly.

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
/v1/google/search
https://api.keyapi.ai/v1/google/search
```

### `call` Flags

| Flag | Required | Repeatable | Description |
| --- | --- | --- | --- |
| `--query k=v` | No | Yes | Add a URL query parameter. |
| `--qparam k=v` | No | Yes | Alias for `--query`. |
| `--data k=v` | No | Yes | Add a JSON body field. Values are sent as strings. |
| `--body-json <json>` | No | No | Send a raw JSON request body. Use this for nested objects or arrays. |
| `--method <method>` | No | No | Override HTTP method. If omitted, `POST` is used when body data exists; otherwise `GET` is used. |
| `--api-key <value>` | No | No | Use a token for this command only. Overrides config and environment. |
| `--token <value>` | No | No | Alias for `--api-key`. |
| `--dry-run` | No | No | Print the final request as JSON without sending it. |
| `--format json\|pretty\|raw` | No | No | Select response output format. Default is `json`. |

Default JSON response envelope:

```json
{
  "ok": true,
  "endpoint": "/v1/google/search",
  "method": "GET",
  "data": {}
}
```

| Field | Meaning |
| --- | --- |
| `ok` | Whether the CLI request succeeded. |
| `endpoint` | Resolved KeyAPI endpoint path. |
| `method` | HTTP method used by the request. |
| `data` | Parsed API response body, or raw text when the response is not JSON. |

## Endpoint Selection Guidance

Use this process to choose an endpoint:

```text
List endpoint descriptions for the relevant platform.
Compare method, path, summary, and description.
Choose the endpoint whose description best matches the user request.
Inspect that endpoint with get-endpoint.
Use parameters and body_fields to construct the call command.
```
