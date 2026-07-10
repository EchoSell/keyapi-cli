# KeyAPI CLI 本地中文命令说明

这份文档用于说明当前 CLI 的命令、入参、返回数据，以及返回字段的含义。

## 命令总览

```text
keyapi config                  管理本地配置
keyapi spec list-platforms     查看当前支持哪些平台
keyapi spec list-endpoints     查看一个或多个平台下有哪些接口
keyapi spec get-endpoint       查看某个接口的详细描述和参数
keyapi call                    调用某个 KeyAPI 接口
```

## 推荐使用流程

```powershell
keyapi config set token sk_xxx
keyapi spec list-platforms
keyapi spec list-endpoints tiktok
keyapi spec get-endpoint tiktok/influencer/detail
keyapi call tiktok/influencer/detail --query unique_id=mrbeast --dry-run
keyapi call tiktok/influencer/detail --query unique_id=mrbeast
```

流程含义：

| 步骤 | 命令 | 作用 |
| --- | --- | --- |
| 1 | `config set token` | 保存 KeyAPI token |
| 2 | `spec list-platforms` | 查看有哪些平台可用 |
| 3 | `spec list-endpoints` | 查看平台下有哪些接口 |
| 4 | `spec get-endpoint` | 查看接口详情和参数 |
| 5 | `call --dry-run` | 预览请求，不真正发送 |
| 6 | `call` | 发送真实请求 |

## 顶层命令

| 命令 | 作用 | 输出 |
| --- | --- | --- |
| `keyapi --help` | 查看 CLI 总帮助 | 普通文本 |
| `keyapi help` | `--help` 的别名 | 普通文本 |
| `keyapi --version` | 查看 CLI 版本号 | 普通文本 |
| `keyapi version` | `--version` 的别名 | 普通文本 |

Windows 当前目录执行：

```powershell
.\keyapi.exe --help
```

## `keyapi config`：配置管理

`config` 用于管理本地配置。这个命令组不强制默认 JSON 输出。

| 命令 | 作用 | 入参 | 输出 |
| --- | --- | --- | --- |
| `keyapi config show` | 查看当前配置，token 会打码 | 无 | JSON |
| `keyapi config get token` | 查看已保存的 token，打码显示 | `token` | 普通文本 |
| `keyapi config get api-key` | `token` 的别名 | `api-key` | 普通文本 |
| `keyapi config get auth-scheme` | 查看认证方案 | `auth-scheme` | 普通文本 |
| `keyapi config get language` | 查看 OpenAPI 文档语言 | `language` | 普通文本 |
| `keyapi config set token <value>` | 保存 KeyAPI token | `value` 是 token | `saved` |
| `keyapi config set api-key <value>` | 保存 KeyAPI token，`token` 的别名 | `value` 是 token | `saved` |
| `keyapi config set auth-scheme <value>` | 设置认证方案，默认 `Bearer` | 例如 `Bearer` | `saved` |
| `keyapi config set language <en\|zh>` | 设置 OpenAPI 文档语言 | `en` 或 `zh` | `saved` |
| `keyapi config clear token` | 删除本地保存的 token | `token` | `cleared` |

示例：

```powershell
keyapi config set token sk_xxx
keyapi config set language en
keyapi config show
```

## `keyapi spec list-platforms`：查看支持平台

作用：查看当前 KeyAPI 文档中支持哪些平台。

```powershell
keyapi spec list-platforms
```

返回 JSON：

```json
["google", "tiktok"]
```

数据说明：

| 数据 | 含义 |
| --- | --- |
| 数组元素 | 当前支持的平台名称 |

## `keyapi spec list-endpoints`：查看平台接口列表

作用：查看一个或多个平台下有哪些接口。这个命令用于选择接口。

```powershell
keyapi spec list-endpoints tiktok
keyapi spec list-endpoints tiktok instagram
keyapi spec list-endpoints all
```

返回 JSON 结构：

```json
{
  "tiktok": [
    {
      "method": "GET",
      "path": "/v1/tiktok/influencer/detail",
      "summary": "Get Influencer Detail",
      "description": "Retrieve real-time profile information for a TikTok influencer by unique ID."
    }
  ]
}
```

字段说明：

| 字段 | 含义 |
| --- | --- |
| 顶层平台名 | 例如 `tiktok`、`instagram`，每个平台对应一个接口数组 |
| `method` | 请求方法，例如 `GET`、`POST` |
| `path` | 请求路径，例如 `/v1/tiktok/influencer/detail` |
| `summary` | 接口简短标题 |
| `description` | 接口详细描述，用于判断这个接口是否符合需求 |

接口选择方式：

```text
先查看相关平台的接口列表。
阅读每个接口的 method、path、summary、description。
根据 summary 和 description 判断哪个接口最符合需求。
选中接口后，再使用 get-endpoint 查看参数。
```

## `keyapi spec get-endpoint`：查看单个接口详情

作用：查看某个接口的请求方式、路径、描述和参数。

```powershell
keyapi spec get-endpoint google/search
keyapi spec get-endpoint /v1/google/search
keyapi spec get-endpoint google/search --format pretty
```

默认返回 JSON：

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

字段说明：

| 字段 | 含义 |
| --- | --- |
| `platform` | 平台名称 |
| `spec_name` | 来源 OpenAPI spec 名称 |
| `method` | 请求方法 |
| `path` | 请求路径 |
| `summary` | 接口简短标题 |
| `description` | 接口详细描述 |
| `parameters` | OpenAPI 风格的请求参数 |
| `body_fields` | 从 OpenAPI `requestBody` 中解析出的 JSON body 字段 |

`parameters` 字段说明：

| 字段 | 含义 |
| --- | --- |
| `name` | 参数名 |
| `in` | 参数位置，例如 `query`、`path`、`header` |
| `type` | 参数类型 |
| `required` | 是否必填 |
| `description` | 参数说明 |
| `schema` | OpenAPI 原始参数 schema，会保留 `example`、`default`、`enum` 等字段 |

`body_fields` 字段说明：

| 字段 | 含义 |
| --- | --- |
| `name` | body 字段名 |
| `type` | 字段类型 |
| `required` | 是否必填 |
| `description` | 字段说明 |

## `keyapi call`：调用接口

作用：发送真实的 KeyAPI 请求。

```powershell
keyapi call <endpoint-or-path> [flags]
```

`<endpoint-or-path>` 支持：

```powershell
keyapi call google/search
keyapi call /v1/google/search
keyapi call https://api.keyapi.ai/v1/google/search
```

参数说明：

| 参数 | 是否可重复 | 作用 |
| --- | --- | --- |
| `--query k=v` | 是 | 添加 URL query 参数 |
| `--qparam k=v` | 是 | `--query` 的别名 |
| `--data k=v` | 是 | 添加 JSON body 字段，值按字符串发送 |
| `--body-json <json>` | 否 | 直接传完整 JSON body |
| `--method <method>` | 否 | 手动指定 HTTP 方法 |
| `--api-key <value>` | 否 | 当前命令临时使用 token |
| `--token <value>` | 否 | `--api-key` 的别名 |
| `--dry-run` | 否 | 只打印请求，不真正发送 |
| `--format json\|pretty\|raw` | 否 | 输出格式，默认 `json` |

示例：

```powershell
keyapi call google/search --query q=openai --query gl=us
keyapi call google/search --query q=openai --dry-run
keyapi call tiktok/shop/product-search --data keyword=dress --data region=US
keyapi call tiktok/shop/product-search --body-json "{\"keyword\":\"dress\",\"region\":\"US\"}"
```

`--dry-run` 返回结构：

```json
{
  "method": "GET",
  "url": "https://api.keyapi.ai/v1/google/search?q=openai",
  "headers": {
    "Accept": "application/json",
    "Authorization": "Bearer ********",
    "Content-Type": ""
  },
  "body": null
}
```

字段说明：

| 字段 | 含义 |
| --- | --- |
| `method` | 最终请求方法 |
| `url` | 最终请求 URL |
| `headers` | 最终请求头，认证信息会打码 |
| `body` | 最终请求体 |

真实请求默认返回 JSON 包装：

```json
{
  "ok": true,
  "endpoint": "/v1/google/search",
  "method": "GET",
  "data": {}
}
```

字段说明：

| 字段 | 含义 |
| --- | --- |
| `ok` | CLI 请求是否成功 |
| `endpoint` | 最终请求的 KeyAPI 路径 |
| `method` | 使用的 HTTP 方法 |
| `data` | API 返回的数据 |

## OpenAPI 获取规则

执行 `spec` 命令时，CLI 会：

```text
1. 请求 https://docs.keyapi.ai/llms.txt
2. 解析其中 ## OpenAPI Specs 部分
3. 找到当前语言对应的平台 OpenAPI JSON 链接
4. 请求 OpenAPI JSON
5. 解析并输出接口信息
```
