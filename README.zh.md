# KeyAPI CLI

---

[中文版](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | [English](https://github.com/EchoSell/keyapi-cli/blob/main/README.md)

KeyAPI CLI 是一个轻量级命令行客户端，用于配置认证信息、发现可用 API 平台、查看接口描述和参数，并在不编写请求代码的情况下调用 API。

CLI 聚焦用户最常用的核心流程：

- 一个配置命令组
- 一个接口发现命令组
- 一个接口调用命令
- `spec` 和 `call` 默认输出 JSON

CLI 每次都会读取最新的 KeyAPI OpenAPI 文档，并解析返回实时的 API 详情，方便用户在发起请求前查看最新的接口能力和参数说明。

## 安装

### 下载 Release

从 GitHub Releases 页面下载适合你操作系统的二进制文件，然后运行：

```powershell
.\keyapi.exe --help
```

普通用户不需要安装 Go。编译好的二进制文件即可直接使用。

### 从源码构建

```powershell
go build -buildvcs=false -o keyapi.exe .
```

## 快速开始

```powershell
keyapi config set --token YOUR_KEYAPI_TOKEN
keyapi spec list-platforms
keyapi spec list-endpoints google
keyapi spec get-endpoint google/search
keyapi call google/search --query q=openai --query gl=us --query hl=en
```

## 核心命令

| 命令 | 作用 |
| --- | --- |
| `keyapi config set --token <token>` | 保存 KeyAPI token。 |
| `keyapi config get` | 查看当前配置。 |
| `keyapi config unset` | 删除已保存的 token。 |
| `keyapi spec list-platforms` | 返回当前支持的平台列表。 |
| `keyapi spec list-endpoints <platform...|all>` | 返回一个或多个平台的接口列表。 |
| `keyapi spec get-endpoint <endpoint-or-path>` | 查看某个接口的方法、路径、描述和参数。 |
| `keyapi call <endpoint-or-path>` | 调用 KeyAPI 接口。 |

## API 发现

### 查看支持的平台

```powershell
keyapi spec list-platforms
```

默认返回 JSON 数组：

```json
["google", "tiktok"]
```

### 查看平台接口

```powershell
keyapi spec list-endpoints tiktok instagram
```

也可以查看全部平台：

```powershell
keyapi spec list-endpoints all
```

返回结构：

```json
{
  "tiktok": [
    {
      "method": "GET",
      "path": "/v1/tiktok/user/detail",
      "summary": "Get TikTok user detail",
      "description": "Get public TikTok user profile data."
    }
  ]
}
```

字段含义：

| 字段 | 含义 |
| --- | --- |
| `method` | HTTP 请求方法，例如 `GET` 或 `POST`。 |
| `path` | `https://api.keyapi.ai` 下的 API 请求路径。 |
| `summary` | 来自文档或 OpenAPI 的接口简短标题。 |
| `description` | 接口详细描述，用于判断接口是否符合任务。 |

### 查看接口详情

```powershell
keyapi spec get-endpoint google/search
```

也可以传入完整路径：

```powershell
keyapi spec get-endpoint /v1/google/search
```

返回结构：

```json
{
  "platform": "google",
  "spec_name": "google",
  "method": "GET",
  "path": "/v1/google/search",
  "summary": "Google Search",
  "description": "Search Google results.",
  "parameters": [
    {
      "name": "q",
      "in": "query",
      "required": true,
      "description": "Search query.",
      "schema": {
        "type": "string",
        "example": "openai"
      }
    }
  ],
  "body_fields": []
}
```

字段含义：

| 字段 | 含义 |
| --- | --- |
| `platform` | 平台名称。 |
| `spec_name` | 来源 OpenAPI spec 名称。 |
| `method` | HTTP 请求方法。 |
| `path` | API 请求路径。 |
| `summary` | 接口简短标题。 |
| `description` | 接口详细描述。 |
| `parameters` | OpenAPI 风格的请求参数。 |
| `body_fields` | 从 OpenAPI `requestBody` 中解析出的 JSON body 字段。 |

`parameters` 中的每一项包含：

| 字段 | 含义 |
| --- | --- |
| `name` | 参数名称。 |
| `in` | 参数位置，例如 `query`、`path` 或 `header`。 |
| `required` | 参数是否必填。 |
| `description` | 参数说明。 |
| `schema` | 原始 OpenAPI 参数 schema，会保留 `example`、`default`、`enum` 等元数据。 |

## JSON 输出

`spec` 和 `call` 命令默认输出 JSON，便于脚本、自动化工具和其他程序解析。

```powershell
keyapi spec list-endpoints google
keyapi spec get-endpoint google/search
keyapi call google/search --query q=openai
```

`config` 命令默认输出面向人类阅读的文本，因为它主要用于本地配置管理。

## 调用 API

查看接口详情后，使用 `keyapi call` 调用接口：

```powershell
keyapi call google/search --query q=openai --query gl=us --query hl=en
```

### Query 参数

```powershell
keyapi call google/search --query q=openai --query gl=us
```

### JSON 请求体

```powershell
keyapi call some/post/endpoint --method POST --body-json '{"keyword":"openai"}'
```

也可以使用多个 `--data`：

```powershell
keyapi call some/post/endpoint --method POST --data keyword=openai --data page=1
```

### Dry Run

```powershell
keyapi call google/search --query q=openai --dry-run
```

`--dry-run` 会打印将要发送的请求信息，但不会真正发起 API 请求。

### API 返回状态

`keyapi call` 会直接输出 KeyAPI 接口返回的原始 JSON，不会再包一层 CLI 自定义结构。

成功示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

常见状态：

| `code` | `message` | 含义 |
| --- | --- | --- |
| `0` | `success` | 请求成功。 |
| `400` | `Missing required parameter: handle` | 缺少必填参数，具体参数名以接口返回为准。 |
| `401` | `apiKey is null` | 未提供 API token，或 token 配置不正确。 |
| `402` | `apiKey quota is not enough` | API 额度不足。 |
| `429` | `qps limit exceeded` | 请求频率超过 QPS 限制。 |
| `500` | `Internal server error` | 服务端错误。 |

这些状态码和字段由 KeyAPI 接口定义。CLI 会保持它们不变。

## 使用流程

推荐流程：

```text
查看平台 -> 查看接口列表 -> 查看接口参数 -> dry-run -> 执行请求
```

示例：

```powershell
keyapi spec list-platforms
keyapi spec list-endpoints tiktok
keyapi spec get-endpoint tiktok/user/detail
keyapi call tiktok/user/detail --query unique_id=mrbeast --dry-run
keyapi call tiktok/user/detail --query unique_id=mrbeast
```

流程建议：

- 通过 `method`、`path`、`summary`、`description` 选择接口。
- 构造请求前先查看 `parameters` 和 `body_fields`。
- query 参数使用 `--query`。
- JSON 请求体使用 `--data` 或 `--body-json`。
- 不确定请求构造是否正确时，先使用 `--dry-run`。
- 将 CLI 集成到其他工具时，可使用 `docs/reference/cli-spec.yaml` 作为轻量机器可读命令契约。

## 文档

| 文件 | 作用 |
| --- | --- |
| [docs/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/index.md) | 文档入口。 |
| [docs/commands.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/commands.md) | 完整命令参考。 |
| [docs/config.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/config.md) | 配置文件和环境变量说明。 |
| [docs/reference/cli-spec.yaml](https://github.com/EchoSell/keyapi-cli/blob/main/docs/reference/cli-spec.yaml) | 轻量机器可读 CLI 契约。 |
| [README.zh.md](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | 中文 README。 |
| [docs/zh/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/index.md) | 中文文档入口。 |
