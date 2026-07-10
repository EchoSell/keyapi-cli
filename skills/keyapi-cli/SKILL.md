# KeyAPI CLI Skill

Use this skill to discover, inspect, and call KeyAPI APIs from the command line.

## Executable

```powershell
keyapi
```

On Windows in the project directory:

```powershell
.\keyapi.exe
```

## Workflow

```text
list platforms -> list endpoint descriptions -> inspect endpoint details -> dry-run -> execute
```

## Command Details

### List Platforms

```powershell
keyapi spec list-platforms
```

Returns JSON:

- array of supported platform names.

### List Endpoints

```powershell
keyapi spec list-endpoints <platform...|all>
```

Examples:

```powershell
keyapi spec list-endpoints tiktok
keyapi spec list-endpoints tiktok instagram
keyapi spec list-endpoints all
```

Returns an object keyed by platform. Each endpoint contains:

- `method`: HTTP request method.
- `path`: request path.
- `summary`: short endpoint title.
- `description`: detailed endpoint description.

### Get Endpoint

```powershell
keyapi spec get-endpoint <endpoint-or-path>
```

Returns:

- `platform`: platform name.
- `spec_name`: source OpenAPI spec name.
- `method`: HTTP request method.
- `path`: request path.
- `summary`: short endpoint title.
- `description`: detailed endpoint description.
- `parameters`: OpenAPI-style request parameters.
- `body_fields`: JSON body fields extracted from OpenAPI `requestBody`.

Each item in `parameters` preserves the original OpenAPI `schema` field when present, including `example`, `default`, and `enum`.

### Call Endpoint

```powershell
keyapi call <endpoint-or-path> [flags]
```

Common flags:

- `--query k=v`: URL query parameter. Repeatable.
- `--data k=v`: JSON body field. Repeatable.
- `--body-json <json>`: raw JSON request body.
- `--method <method>`: HTTP method override.
- `--dry-run`: print request JSON without sending.
- `--format json|pretty|raw`: output format. Default is `json`.

## Endpoint Selection

Select an endpoint by reading endpoint descriptions from `list-endpoints`, then inspect the selected endpoint with `get-endpoint`.
