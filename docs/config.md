# Configuration

KeyAPI CLI stores local configuration in a JSON file and supports environment variable overrides.

## Config Location

By default, the config directory is based on the operating system user config directory:

```text
<user-config-dir>/keyapi/config.json
```

You can override it with:

```powershell
$env:KEYAPI_HOME='C:\path\to\keyapi-home'
```

When `KEYAPI_HOME` is set, the config file is:

```text
<KEYAPI_HOME>/config.json
```

## Config File Shape

```json
{
  "api_key": "sk_xxx",
  "base_url": "https://api.keyapi.ai",
  "auth_scheme": "Bearer",
  "language": "en"
}
```

Notes:

- `api_key` is the KeyAPI token.
- `base_url` defaults to `https://api.keyapi.ai`.
- `auth_scheme` defaults to `Bearer`.
- `language` defaults to `en`.

## Environment Variables

| Variable | Description | Priority |
| --- | --- | --- |
| `KEYAPI_API_KEY` | KeyAPI token. Overrides the token stored in config. | Highest for API token |
| `KEYAPI_HOME` | Config directory. | Highest for config location |



## Managing Token

Save a token:

```powershell
keyapi config set token sk_xxx
```

Show the masked token:

```powershell
keyapi config get token
```

Remove the stored token:

```powershell
keyapi config clear token
```

For temporary use, prefer the environment variable:

```powershell
$env:KEYAPI_API_KEY='sk_xxx'
```

## Managing Docs Language

```powershell
keyapi config set language en
keyapi config set language zh
```

| Value | Meaning |
| --- | --- |
| `en` | Use English OpenAPI docs. Default. |
| `zh` | Use Chinese OpenAPI docs. |
