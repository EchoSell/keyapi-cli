# KeyAPI CLI

---
[中文版](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | [English](https://github.com/EchoSell/keyapi-cli/blob/main/README.md)

KeyAPI CLI is a lightweight command line client for the [KeyAPI](https://keyapi.ai) platform. It helps configure authentication, discover available API platforms, inspect endpoint descriptions and parameters, and call APIs without writing request code.

The CLI focuses on the core workflow users need:

- one command group for configuration
- one command group for API discovery
- one command for API execution
- JSON output by default for `spec` and `call`

The CLI always fetches the latest KeyAPI OpenAPI specification and provides up-to-date endpoint details, allowing you to review request methods, parameters, and capabilities before making API calls.

## Installation

### Download a Release

Download the binary for your operating system from the GitHub Releases page, then run:

```powershell
.\keyapi.exe --help
```

End users do not need Go installed. A compiled binary is enough.

### Build from Source

Requirements:

- Go 1.22 or later

Windows:

```powershell
go build -buildvcs=false -o keyapi.exe .
```

Linux or macOS:

```bash
go build -buildvcs=false -o keyapi .
```

## Quick Start

```powershell
.\keyapi.exe config set token sk_xxx
.\keyapi.exe spec list-platforms
.\keyapi.exe spec list-endpoints google
.\keyapi.exe spec get-endpoint google/search
.\keyapi.exe call google/search --query q=openai --query gl=us --query hl=en
```

## Core Commands

| Command | Purpose |
| --- | --- |
| `keyapi config` | Manage local configuration, including token and docs language. |
| `keyapi spec list-platforms` | List platforms that currently have KeyAPI endpoints. |
| `keyapi spec list-endpoints <platform...\|all>` | List endpoint descriptions for one or more platforms. |
| `keyapi spec get-endpoint <endpoint-or-path>` | Get method, path, description, parameters, and body fields for one endpoint. |
| `keyapi call <endpoint-or-path>` | Execute a KeyAPI request. |

Every command supports help:

```powershell
keyapi --help
keyapi config --help
keyapi spec --help
keyapi spec list-platforms --help
keyapi spec list-endpoints --help
keyapi spec get-endpoint --help
keyapi call --help
```

## API Discovery

Endpoint selection guidance:

```text
To choose an API endpoint, first list endpoint descriptions for the relevant platform.
Read each endpoint's method, path, summary, and description.
Select the endpoint whose description best matches the user request.
Then call get-endpoint to inspect parameters before executing the request.
```

Example:

```powershell
keyapi spec list-endpoints tiktok
keyapi spec get-endpoint tiktok/influencer/detail
keyapi call tiktok/influencer/detail --query unique_id=mrbeast
```

## JSON Output

`spec` and `call` commands return JSON by default.

`keyapi spec list-platforms`:

```json
["google", "tiktok"]
```

`keyapi spec list-endpoints tiktok`:

```json
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
```

Field meanings:

| Field | Meaning |
| --- | --- |
| `method` | HTTP request method, such as `GET` or `POST`. |
| `path` | API request path under `https://api.keyapi.ai`. |
| `summary` | Short endpoint title from docs or OpenAPI. |
| `description` | Longer endpoint description. Use this to decide whether the endpoint matches the task. |

`keyapi spec get-endpoint google/search` returns:

| Field | Meaning |
| --- | --- |
| `platform` | Platform name. |
| `spec_name` | Source OpenAPI spec name. |
| `method` | HTTP request method. |
| `path` | API request path. |
| `summary` | Short endpoint title. |
| `description` | Detailed endpoint description. |
| `parameters` | OpenAPI-style request parameters. |
| `body_fields` | JSON body fields extracted from OpenAPI `requestBody`. |

Each item in `parameters` contains:

| Field | Meaning |
| --- | --- |
| `name` | Parameter name. |
| `in` | Parameter location, such as `query`, `path`, or `header`. |
| `type` | Parameter schema type. |
| `required` | Whether the parameter is required. |
| `description` | Parameter description. |
| `schema` | Original OpenAPI parameter schema. Preserves metadata such as `example`, `default`, and `enum` when provided by the spec. |

## Calling APIs

Use `keyapi call` after inspecting endpoint details.

```powershell
keyapi call <endpoint-or-path> [flags]
```

Endpoint formats:

```powershell
keyapi call google/search
keyapi call /v1/google/search
keyapi call https://api.keyapi.ai/v1/google/search
```

Query parameters:

```powershell
keyapi call google/search --query q=openai --query gl=us
```

JSON body fields:

```powershell
keyapi call tiktok/shop/product-search --data keyword=dress --data region=US
```

Raw JSON body:

```powershell
keyapi call tiktok/shop/product-search --body-json "{\"keyword\":\"dress\",\"region\":\"US\"}"
```

Dry run:

```powershell
keyapi call google/search --query q=openai --dry-run
```

By default, `keyapi call` prints the original response returned by the KeyAPI endpoint. The CLI does not wrap or rewrite the API response.

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

Common KeyAPI response codes:

| Code | Message | Meaning |
| --- | --- | --- |
| `0` | `success` | Request succeeded. |
| `400` | `Missing required parameter: <name>` | A required parameter is missing or invalid. |
| `401` | `apiKey is null` | API key is missing. |
| `402` | `apiKey quota is not enough` | API key quota is insufficient. |
| `429` | `qps limit exceeded` | Request rate exceeded the QPS limit. |
| `500` | `Internal server error` | Server-side error. |

These status codes and fields are defined by KeyAPI endpoints. The CLI keeps them unchanged.

## Usage Workflow

Recommended workflow:

```text
list platforms -> list endpoint descriptions -> inspect endpoint parameters -> dry-run -> execute
```

Example:

```powershell
keyapi spec list-platforms
keyapi spec list-endpoints tiktok
keyapi spec get-endpoint tiktok/influencer/detail
keyapi call tiktok/influencer/detail --query unique_id=mrbeast --dry-run
keyapi call tiktok/influencer/detail --query unique_id=mrbeast
```

Workflow rules:

- Choose endpoints by reading `method`, `path`, `summary`, and `description`.
- Inspect `parameters` and `body_fields` before constructing a request.
- Use `--query` for query parameters.
- Use `--data` or `--body-json` for JSON request bodies.
- Use `--dry-run` before execution when request construction is uncertain.
- Use `docs/reference/cli-spec.yaml` as the lightweight machine-readable command contract when integrating the CLI into other tools.

## Documentation

| File | Purpose |
| --- | --- |
| [docs/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/index.md) | Documentation entry point. |
| [docs/commands.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/commands.md) | Full command and parameter reference. |
| [docs/config.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/config.md) | Config file and environment variables. |
| [docs/reference/cli-spec.yaml](https://github.com/EchoSell/keyapi-cli/blob/main/docs/reference/cli-spec.yaml) | Lightweight machine-readable CLI contract. |
| [README.zh.md](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | Chinese README. |
| [docs/zh/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/index.md) | Chinese documentation entry point. |

## Development

Run tests:

```powershell
go test ./...
```

Build:

```powershell
go build -buildvcs=false -o keyapi.exe .
```
