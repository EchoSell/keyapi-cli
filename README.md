# KeyAPI CLI

---

[Chinese](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | [English](https://github.com/EchoSell/keyapi-cli/blob/main/README.md)

KeyAPI CLI is a lightweight command line client for the [KeyAPI](https://keyapi.ai) platform. It helps you configure authentication, discover supported API platforms, inspect endpoint schemas, and call APIs without writing request code.

## Why KeyAPI CLI?

KeyAPI provides many APIs across platforms such as Google, TikTok, Instagram, and more. When a project uses many endpoints, manually writing and maintaining request code can become repetitive.

KeyAPI CLI provides one consistent entry point:

- configure your API token once
- list supported platforms
- list endpoints under each platform
- inspect endpoint methods, paths, descriptions, parameters, and request body fields
- execute API requests from the terminal
- save large responses to files

The CLI reads the latest KeyAPI OpenAPI documentation when endpoint details are requested, so endpoint descriptions and schemas stay aligned with the current KeyAPI API documentation.

## Features

- Simple configuration commands: `keyapi set`, `keyapi get`, `keyapi show`, `keyapi clear`
- API discovery commands: `keyapi spec list-platforms`, `keyapi spec list-endpoints`, `keyapi schema`
- One request command: `keyapi call`
- JSON output by default for `spec` and `call`
- Built-in `--help` output designed for humans, scripts, and agents
- `--dry-run` support before sending requests
- `--output "path"` support for large responses

## Installation

Requires Node.js 18 or later:

```bash
npm install -g keyapi-cli
keyapi --help
keyapi --version
```

The npm package installs a small JavaScript launcher plus the matching native
KeyAPI binary package for your operating system. If you already have a custom
KeyAPI binary, set `KEYAPI_CLI_BINARY` to its path before running `keyapi`.

## Quick Start

Configure your KeyAPI token:

```powershell
keyapi set token sk_xxx
```

Check the current configuration:

```powershell
keyapi show
```

List supported platforms:

```powershell
keyapi spec list-platforms
```

List endpoints for a platform:

```powershell
keyapi spec list-endpoints google
```

Inspect an endpoint before calling it:

```powershell
keyapi schema google/search
```

Dry-run a request:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&hl=en" --dry-run
```

Execute the request:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&hl=en"
```

Save a large response:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&hl=en" --output "google-search.json"
```

Clear the saved token when needed:

```powershell
keyapi clear token
```

## Commands

| Command | Purpose |
| --- | --- |
| `keyapi --help` | Show the main help, recommended order, and command help entry points. |
| `keyapi set token <value>` | Save the API token locally. |
| `keyapi set language <en\|zh>` | Set OpenAPI documentation language for endpoint descriptions. |
| `keyapi get <key>` | Read one local configuration value. |
| `keyapi show` | Show current local configuration. |
| `keyapi clear <key>` | Clear one local configuration value. |
| `keyapi spec list-platforms` | Return a JSON array of supported platforms. |
| `keyapi spec list-endpoints <platform...\|all>` | Return endpoint summaries for one or more platforms. |
| `keyapi spec get-endpoint <endpoint-or-path>` | Return method, path, description, parameters, and body fields for one endpoint. |
| `keyapi schema <endpoint-or-path>` | Short alias for `keyapi spec get-endpoint`. |
| `keyapi call <endpoint-or-path>` | Execute a KeyAPI request. |

Every command supports `--help`:

```powershell
keyapi --help
keyapi set --help
keyapi spec --help
keyapi spec list-platforms --help
keyapi spec list-endpoints --help
keyapi spec get-endpoint --help
keyapi schema --help
keyapi call --help
```

## Advanced Usage

### JSON Output

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

### Save Large Responses

Some endpoints return large JSON responses, such as product details, images, SKU lists, or description fields. Use `--output "path"` to write the full response body to a file instead of printing it in the terminal:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us" --output "outputs/google-search.json"
```

Relative output paths are resolved from the current working directory.

### Dry Run

Use `--dry-run` to print the request without sending it:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us" --dry-run
```

### Request Body

Use `--body-json` for JSON request bodies:

```powershell
keyapi call "/v1/example/post-endpoint" --method POST --body-json '{"page":1,"verified":true}'
```

Use `--body-raw` for raw text request bodies:

```powershell
keyapi call "/v1/example/text-endpoint" --method POST --body-raw "plain text"
```

For GET endpoints, prefer putting query parameters directly in the endpoint path:

```powershell
keyapi call "/v1/google/search?q=openai&gl=us"
```

Use `--params "k=v"` or `--query "k=v"` when query parameters are easier to append programmatically.

## Documentation

`docs/` contains usage documentation for installing, configuring, and using the CLI. The CLI help output is the primary built-in reference for command usage.

| File | Purpose |
| --- | --- |
| [docs/index.md](docs/index.md) | Documentation entry point. |
| [docs/commands.md](docs/commands.md) | Full command and parameter reference. |
| [docs/config.md](docs/config.md) | Config file and environment variables. |
| [README.zh.md](README.zh.md) | Chinese README. |
| [docs/zh/index.md](docs/zh/index.md) | Chinese documentation entry point. |
