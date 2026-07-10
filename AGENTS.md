# KeyAPI CLI Command Guide

This file describes the CLI usage workflow and command contracts for tool-based invocation.

## Workflow

Use this sequence:

```text
list platforms -> list endpoint descriptions -> inspect endpoint details -> dry-run -> execute
```

## Commands

List supported platforms:

```powershell
keyapi spec list-platforms
```

Returns:

```json
["google", "tiktok"]
```

List endpoint descriptions for one or more platforms:

```powershell
keyapi spec list-endpoints tiktok
keyapi spec list-endpoints tiktok instagram
keyapi spec list-endpoints all
```

Returns an object keyed by platform. Each endpoint item contains:

- `method`: HTTP request method.
- `path`: KeyAPI request path.
- `summary`: short endpoint title.
- `description`: detailed endpoint description used to choose the endpoint.

Inspect one endpoint:

```powershell
keyapi spec get-endpoint tiktok/influencer/detail
```

Returns endpoint details:

- `platform`: platform name.
- `spec_name`: source OpenAPI spec name.
- `method`: HTTP request method.
- `path`: request path.
- `summary`: short endpoint title.
- `description`: detailed endpoint description.
- `parameters`: OpenAPI-style request parameters.
- `body_fields`: JSON body fields extracted from OpenAPI `requestBody`.

Each item in `parameters` preserves the original OpenAPI `schema` field when present, including `example`, `default`, and `enum`.

Preview a request:

```powershell
keyapi call tiktok/influencer/detail --query unique_id=mrbeast --dry-run
```

Execute a request:

```powershell
keyapi call tiktok/influencer/detail --query unique_id=mrbeast
```

## Selection Rules

- Choose endpoints by reading `method`, `path`, `summary`, and `description`.
- Inspect `parameters` and `body_fields` before constructing a request.
- Use `--query` for query parameters.
- Use `--data` or `--body-json` for JSON request bodies.
- Use `--dry-run` before execution when request construction is uncertain.
