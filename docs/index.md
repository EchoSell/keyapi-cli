# KeyAPI CLI Documentation

KeyAPI CLI is a small command line tool for configuring a KeyAPI token, discovering API endpoints, inspecting endpoint parameters, and executing requests.

## Command Groups

| Command | Purpose |
| --- | --- |
| `keyapi config` | Manage local token and CLI settings. |
| `keyapi spec list-platforms` | List platforms with available KeyAPI endpoints. |
| `keyapi spec list-endpoints <platform...\|all>` | List endpoint descriptions for one or more platforms. |
| `keyapi spec get-endpoint <endpoint-or-path>` | Get one endpoint's method, path, description, parameters, and body fields. |
| `keyapi call <endpoint-or-path>` | Execute a KeyAPI request. |

## Documents

| File | Purpose |
| --- | --- |
| [commands.md](commands.md) | Full command reference with arguments, flags, return values, and JSON field meanings. |
| [config.md](config.md) | Configuration file and environment variables. |
| [commands.zh.md](commands.zh.md) | Chinese command guide. |
| [reference/cli-spec.yaml](reference/cli-spec.yaml) | Lightweight machine-readable command contract. |

## Recommended Workflow

```powershell
keyapi config set token sk_xxx
keyapi spec list-platforms
keyapi spec list-endpoints google
keyapi spec get-endpoint google/search
keyapi call "/v1/google/search?q=openai"
```

## Endpoint Selection Guidance

To choose an endpoint, list endpoint descriptions for the target platform and compare:

- `method`
- `path`
- `summary`
- `description`

After selecting an endpoint, run `get-endpoint` and use its `parameters` and `body_fields` to construct the request.

Request construction rules:

- For GET endpoints, put query parameters in the endpoint path.
- Use `--qparam "k=v"` or `--query "k=v"` when query parameters need to be appended separately.
- Use `--body-json` for JSON request bodies.
- Use `--body-raw` for raw text request bodies.
- Use `--dry-run` before execution when request construction is uncertain.
