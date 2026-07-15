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

### Install from Release

1. Open the [GitHub Releases](https://github.com/EchoSell/keyapi-cli/releases) page.
2. Download the binary for your operating system.
```powershell
cd ~/keyapi-cli
keyapi --help
```

Copy this prompt into an AI assistant such as Claude Code, Codex, Cursor, or Trae if you want it to install the release binary for you:

```text
Help me install KeyAPI CLI from https://github.com/EchoSell/keyapi-cli/releases. Download the binary for my operating system, make it available as keyapi on PATH, and verify it with keyapi --help.
```

### Build from Source

Requires Go `v1.22` or later.

Clone the repository and build the executable:

```powershell
git clone https://github.com/EchoSell/keyapi-cli.git
cd keyapi-cli
go build -buildvcs=false -o keyapi.exe .
```

On Linux or macOS, build with:

```bash
git clone https://github.com/EchoSell/keyapi-cli.git
cd keyapi-cli
go build -buildvcs=false -o keyapi .
```


Copy this prompt into an AI assistant such as Claude Code, Codex, Cursor, or Trae if you want it to build the project from source:

```text
Help me install this project: https://github.com/EchoSell/keyapi-cli.git. After installation, compile it into an executable file and make the keyapi command available on PATH.
```

After building, continue with the PATH setup section below.

### Make keyapi available on PATH

After installing or building the binary, place it in a directory that is available on `PATH`.

Windows example:

```powershell
# Example: keyapi.exe is stored in C:\Tools\keyapi
$env:Path += ";C:\Tools\keyapi"
keyapi --help
```

The command above only updates the current PowerShell session. For permanent use, add the directory to the system or user `PATH` environment variable.

Linux or macOS example:

```bash
chmod +x keyapi
sudo mv keyapi /usr/local/bin/keyapi
keyapi --help
```

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

Use `--output "path"` to write a response body to a file:

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
| [docs/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/index.md) | Documentation entry point. |
| [docs/commands.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/commands.md) | Full command and parameter reference. |
| [docs/config.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/config.md) | Config file and environment variables. |
| [README.zh.md](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | Chinese README. |
| [docs/zh/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/index.md) | Chinese documentation entry point. |
