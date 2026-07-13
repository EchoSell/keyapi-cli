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

The CLI version is defined in `internal/cli/help.go`. Update the `version` constant there when releasing a new version.

## Quick Start

```powershell
.\keyapi.exe set token sk_xxx
.\keyapi.exe spec list-platforms
.\keyapi.exe spec list-endpoints google
.\keyapi.exe spec get-endpoint google/search
.\keyapi.exe call "/v1/google/search?q=openai&gl=us&hl=en"
```

## Core Commands

| Command | Purpose |
| --- | --- |
| `keyapi set <key> <value>` | Save local configuration, including token and docs language. |
| `keyapi get <key>` | Read one local configuration value. |
| `keyapi show` | Show current local configuration. |
| `keyapi clear <key>` | Clear one local configuration value. |
| `keyapi spec list-platforms` | List platforms that currently have KeyAPI endpoints. |
| `keyapi spec list-endpoints <platform...\|all>` | List endpoint descriptions for one or more platforms. |
| `keyapi spec get-endpoint <endpoint-or-path>` | Get method, path, description, parameters, and body fields for one endpoint. |
| `keyapi call <endpoint-or-path>` | Execute a KeyAPI request. |

Every command supports help:

```powershell
keyapi --help
keyapi set --help
keyapi get --help
keyapi show --help
keyapi clear --help
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
keyapi call "/v1/tiktok/influencer/detail?unique_id=mrbeast"
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
keyapi call "/v1/google/search?q=openai&gl=us"
keyapi call "https://api.keyapi.ai/v1/google/search?q=openai&gl=us"
```

Query parameters:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us"
keyapi call "/v1/tiktok/shop/product-search?keyword=hot%20sales&region=US"
keyapi call /v1/google/search --params "q=openai" --params "gl=us"
```

For GET endpoints, prefer putting query parameters directly in the endpoint path. Use `--params "k=v"` or `--query "k=v"` when parameters are easier to append programmatically.

Raw text body:

```powershell
keyapi call some/text/endpoint --body-raw 'plain text body'
```

`--body-raw` sends the body exactly as provided with `Content-Type: text/plain; charset=utf-8`.

JSON body:

```powershell
keyapi call tiktok/shop/product-search --body-json "{\"keyword\":\"dress\",\"region\":\"US\"}"
keyapi call tiktok/shop/product-search --body-json '{"timeRange":7,"countryCode":"US","page":1,"limit":100}'
```

`--body-json` must be valid JSON and is sent with `Content-Type: application/json`. Invalid JSON fails before the request is sent.

In PowerShell, wrap JSON with single quotes:

```powershell
keyapi call tiktok/shop/product-search --body-json '{"page":1,"verified":true}'
```

If PowerShell strips JSON double quotes before they reach the executable, the CLI normalizes common object forms such as `{page:1,verified:true}` back to valid JSON.

Dry run:

```powershell
keyapi call "/v1/google/search?q=openai" --dry-run
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

For HTTP status handling, `2xx` responses exit with code `0`. HTTP non-`2xx` responses still print the KeyAPI response, then exit with a non-zero code so scripts and CI can detect failures.

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
keyapi call "/v1/tiktok/influencer/detail?unique_id=mrbeast" --dry-run
keyapi call "/v1/tiktok/influencer/detail?unique_id=mrbeast"
```

Workflow rules:

- Choose endpoints by reading `method`, `path`, `summary`, and `description`.
- Inspect `parameters` and `body_fields` before constructing a request.
- For GET endpoints, put query parameters in the endpoint path, for example `"/v1/google/search?q=openai&gl=us"`.
- Use `--params "k=v"` or `--query "k=v"` when query parameters need to be appended separately.
- Use `--body-json` for JSON request bodies.
- Use `--body-raw` for raw text request bodies.
- Use `--dry-run` before execution when request construction is uncertain.
- Use `docs/reference/cli-spec.yaml` as the lightweight machine-readable command contract when integrating the CLI into other tools.

## Documentation

| File | Purpose                                        |
| --- |------------------------------------------------|
| [docs/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/index.md) | Documentation entry point.                     |
| [docs/commands.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/commands.md) | Full command and parameter reference.          |
| [docs/config.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/config.md) | Config file and environment variables.         |
| [docs/reference/cli-spec.yaml](https://github.com/EchoSell/keyapi-cli/blob/main/docs/reference/cli-spec.yaml) | Lightweight machine-readable CLI contract.     |
| [README.zh.md](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | 中文 README。    |
| [docs/zh/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/index.md) | 中文文档入口。       |
| [docs/commands.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/commands.md) | 中文完整命令和参数参考。  |
| [docs/config.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/config.md) | 中文配置文件和环境变量。    |
