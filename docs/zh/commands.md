# 命令参考

本文档列出当前 KeyAPI CLI 的命令、参数、flags 和 JSON 返回字段含义。

## 顶层命令

| 命令 | 说明 | 输出 |
| --- | --- | --- |
| `keyapi --help` | 显示顶层帮助和核心流程。 | 普通文本 |
| `keyapi help` | 顶层帮助别名。 | 普通文本 |
| `keyapi --version` | 输出 CLI 版本号。 | 普通文本 |
| `keyapi version` | 版本命令别名。 | 普通文本 |

## `keyapi config`

`config` 命令主要用于本地配置管理。`config show` 返回 JSON，单个配置项操作返回普通文本。

| 命令 | 说明 | 必填参数 | 输出 |
| --- | --- | --- | --- |
| `keyapi config --help` | 显示配置命令帮助。 | 无 | 普通文本 |
| `keyapi config show` | 显示当前配置，token 会被打码。 | 无 | JSON |
| `keyapi config get` | 不指定 key 时等同于 `show`。 | 无 | JSON |
| `keyapi config get token` | 显示打码后的 token。 | `token` | 普通文本 |
| `keyapi config get api-key` | 显示打码后的 token，`token` 的别名。 | `api-key` | 普通文本 |
| `keyapi config get auth-scheme` | 显示认证方案。 | `auth-scheme` | 普通文本 |
| `keyapi config get language` | 显示 OpenAPI 文档语言。 | `language` | 普通文本 |
| `keyapi config get lang` | 显示 OpenAPI 文档语言，`language` 的别名。 | `lang` | 普通文本 |
| `keyapi config set token <value>` | 保存 KeyAPI token。 | `token`, `value` | `saved` |
| `keyapi config set api-key <value>` | 保存 KeyAPI token，`token` 的别名。 | `api-key`, `value` | `saved` |
| `keyapi config set auth-scheme <value>` | 保存认证方案，默认值为 `Bearer`。 | `auth-scheme`, `value` | `saved` |
| `keyapi config set language <en\|zh>` | 保存 OpenAPI 文档语言。 | `language`, `en` 或 `zh` | `saved` |
| `keyapi config set lang <en\|zh>` | 保存 OpenAPI 文档语言，`language` 的别名。 | `lang`, `en` 或 `zh` | `saved` |
| `keyapi config clear token` | 删除已保存的 token。 | `token` | `cleared` |
| `keyapi config clear api-key` | 删除已保存的 token，`token` 的别名。 | `api-key` | `cleared` |

## `keyapi spec`

`spec` 命令默认返回 JSON。

CLI 每次都会读取最新的 KeyAPI OpenAPI 文档，并解析返回实时的 API 详情。

| 命令 | 说明 | 必填参数 | Flags | 输出 |
| --- | --- | --- | --- | --- |
| `keyapi spec --help` | 显示 spec 帮助。 | 无 | 无 | 普通文本 |
| `keyapi spec list-platforms` | 列出当前支持的平台。 | 无 | 无 | JSON |
| `keyapi spec list-endpoints <platform...>` | 列出一个或多个平台的接口描述。 | 一个或多个平台名 | 无 | JSON |
| `keyapi spec list-endpoints all` | 列出全部平台的接口描述。 | `all` | 无 | JSON |
| `keyapi spec get-endpoint <endpoint-or-path>` | 获取单个接口详情。 | `endpoint-or-path` | `--format json\|pretty` | 默认 JSON |

### `spec list-platforms`

用法：

```powershell
keyapi spec list-platforms
```

输出：

```json
["google", "tiktok"]
```

| 数据 | 含义 |
| --- | --- |
| 数组元素 | 当前拥有 OpenAPI spec 的平台名。 |

### `spec list-endpoints`

用法：

```powershell
keyapi spec list-endpoints tiktok
keyapi spec list-endpoints tiktok instagram
keyapi spec list-endpoints all
```

返回结构：

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

| 字段 | 含义 |
| --- | --- |
| 顶层平台 key | 平台名，例如 `tiktok` 或 `instagram`。 |
| `method` | HTTP 请求方法。 |
| `path` | `https://api.keyapi.ai` 下的 API 请求路径。 |
| `summary` | 接口简短标题。 |
| `description` | 接口详细描述，用于判断接口是否符合需求。 |

### `spec get-endpoint`

用法：

```powershell
keyapi spec get-endpoint google/search
keyapi spec get-endpoint /v1/google/search
keyapi spec get-endpoint google/search --format pretty
```

默认 JSON 返回结构：

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

| 字段 | 含义 |
| --- | --- |
| `platform` | 平台名称。 |
| `spec_name` | 来源 OpenAPI spec 名称。 |
| `method` | HTTP 请求方法。 |
| `path` | API 请求路径。 |
| `summary` | 接口简短标题。 |
| `description` | 接口详细描述。 |
| `parameters` | OpenAPI 风格的请求参数。 |
| `body_fields` | 从 OpenAPI `requestBody` 中解析出的 JSON 请求体字段。 |

参数字段：

| 字段 | 含义 |
| --- | --- |
| `name` | 参数名称。 |
| `in` | 参数位置，例如 `query`、`path` 或 `header`。 |
| `type` | schema 类型。 |
| `required` | 参数是否必填。 |
| `description` | 参数说明。 |
| `schema` | 原始 OpenAPI 参数 schema，会保留 `example`、`default`、`enum` 等元数据。 |

## `keyapi call`

`call` 用于执行 KeyAPI 接口请求，默认输出接口原始 JSON 响应。

| 命令 | 说明 | 必填参数 | Flags | 输出 |
| --- | --- | --- | --- | --- |
| `keyapi call --help` | 显示 call 帮助。 | 无 | 无 | 普通文本 |
| `keyapi call <endpoint-or-path>` | 执行 KeyAPI 接口请求。 | `endpoint-or-path` | 见下方 | 默认输出 KeyAPI 原始响应 |

接口路径格式：

```text
google/search
/v1/google/search
https://api.keyapi.ai/v1/google/search
```

### `call` Flags

| Flag | 必填 | 可重复 | 说明 |
| --- | --- | --- | --- |
| `--query k=v` | 否 | 是 | 添加 URL query 参数。 |
| `--qparam k=v` | 否 | 是 | `--query` 的别名。 |
| `--data k=v` | 否 | 是 | 添加 JSON body 字段，值会按字符串发送。 |
| `--body-json <json>` | 否 | 否 | 发送完整 JSON 请求体，适合嵌套对象或数组。 |
| `--method <method>` | 否 | 否 | 覆盖 HTTP 方法。未指定时，有 body 默认为 `POST`，否则为 `GET`。 |
| `--api-key <value>` | 否 | 否 | 仅当前命令使用的 API key，会覆盖配置和环境变量。 |
| `--token <value>` | 否 | 否 | `--api-key` 的别名。 |
| `--dry-run` | 否 | 否 | 打印最终请求 JSON，不发送请求。 |
| `--format json\|pretty\|raw` | 否 | 否 | 响应输出格式，默认 `json`。 |

`keyapi call` 会直接输出 KeyAPI 接口返回的原始响应，保持与 KeyAPI 接口一致。

成功响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

常见 KeyAPI 响应状态：

| code | message | 含义 |
| --- | --- | --- |
| `0` | `success` | 请求成功。 |
| `400` | `Missing required parameter: <name>` | 缺少必填参数或参数不合法。 |
| `401` | `apiKey is null` | API key 缺失。 |
| `402` | `apiKey quota is not enough` | API key 额度不足。 |
| `429` | `qps limit exceeded` | 请求频率超过 QPS 限制。 |
| `500` | `Internal server error` | 服务端错误。 |

这些字段和状态码由 KeyAPI 接口定义，CLI 会保持原样输出。

## 接口选择说明

使用以下流程选择接口：

```text
列出相关平台的接口描述。
比较 method、path、summary、description。
选择 description 最符合需求的接口。
使用 get-endpoint 查看接口详情。
根据 parameters 和 body_fields 构造 call 命令。
```

选择和请求构造规则：

- 通过 `method`、`path`、`summary`、`description` 选择接口。
- 构造请求前先查看 `parameters` 和 `body_fields`。
- query 参数使用 `--query`。
- JSON body 使用 `--data` 或 `--body-json`。
- 不确定请求是否正确时，先使用 `--dry-run`。
