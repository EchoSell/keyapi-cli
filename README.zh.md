# KeyAPI CLI

---

[中文版](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | [English](https://github.com/EchoSell/keyapi-cli/blob/main/README.md)

KeyAPI CLI 是 [KeyAPI](https://keyapi.ai) 平台的轻量级命令行客户端。它可以帮助你配置认证信息、发现可用的 API 平台、查看接口描述和参数，并在不编写请求代码的情况下调用 API。

CLI 聚焦用户需要的核心流程：

- 一个配置命令组
- 一个 API 发现命令组
- 一个 API 调用命令
- `spec` 和 `call` 默认输出 JSON

CLI 每次都会读取最新的 KeyAPI OpenAPI 规范，并提供最新的接口详情，方便你在发起 API 请求前查看请求方法、参数和接口能力。

## 安装

### 下载 Release

从 GitHub Releases 页面下载适合你操作系统的二进制文件，然后运行：

```powershell
.\keyapi.exe --help
```

普通用户不需要安装 Go。编译好的二进制文件即可使用。

### 从源码构建

要求：

- Go 1.22 或更高版本

Windows：

```powershell
go build -buildvcs=false -o keyapi.exe .
```

Linux 或 macOS：

```bash
go build -buildvcs=false -o keyapi .
```

## 快速开始

```powershell
.\keyapi.exe config set token sk_xxx
.\keyapi.exe spec list-platforms
.\keyapi.exe spec list-endpoints google
.\keyapi.exe spec get-endpoint google/search
.\keyapi.exe call "/v1/google/search?q=openai&gl=us&hl=en"
```

## 核心命令

| 命令 | 作用 |
| --- | --- |
| `keyapi config` | 管理本地配置，包括 token 和文档语言。 |
| `keyapi spec list-platforms` | 列出当前拥有 KeyAPI 接口的平台。 |
| `keyapi spec list-endpoints <platform...\|all>` | 列出一个或多个平台的接口描述。 |
| `keyapi spec get-endpoint <endpoint-or-path>` | 获取单个接口的方法、路径、描述、参数和 body 字段。 |
| `keyapi call <endpoint-or-path>` | 执行 KeyAPI 请求。 |

每个命令都支持 help：

```powershell
keyapi --help
keyapi config --help
keyapi spec --help
keyapi spec list-platforms --help
keyapi spec list-endpoints --help
keyapi spec get-endpoint --help
keyapi call --help
```

## API 发现

接口选择建议：

```text
选择 API 接口时，先列出相关平台的接口描述。
阅读每个接口的 method、path、summary 和 description。
选择 description 最符合用户需求的接口。
然后调用 get-endpoint 查看参数，再执行请求。
```

示例：

```powershell
keyapi spec list-endpoints tiktok
keyapi spec get-endpoint tiktok/influencer/detail
keyapi call "/v1/tiktok/influencer/detail?unique_id=mrbeast"
```

## JSON 输出

`spec` 和 `call` 命令默认返回 JSON。

`keyapi spec list-platforms`：

```json
["google", "tiktok"]
```

`keyapi spec list-endpoints tiktok`：

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

字段含义：

| 字段 | 含义 |
| --- | --- |
| `method` | HTTP 请求方法，例如 `GET` 或 `POST`。 |
| `path` | `https://api.keyapi.ai` 下的 API 请求路径。 |
| `summary` | 来自文档或 OpenAPI 的接口简短标题。 |
| `description` | 更详细的接口描述，用于判断接口是否符合任务。 |

`keyapi spec get-endpoint google/search` 返回：

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
| `type` | 参数 schema 类型。 |
| `required` | 参数是否必填。 |
| `description` | 参数说明。 |
| `schema` | 原始 OpenAPI 参数 schema。会保留 spec 中提供的 `example`、`default`、`enum` 等元数据。 |

## 调用 API

查看接口详情后，使用 `keyapi call`。

```powershell
keyapi call <endpoint-or-path> [flags]
```

接口格式：

```powershell
keyapi call google/search
keyapi call "/v1/google/search?q=openai&gl=us"
keyapi call "https://api.keyapi.ai/v1/google/search?q=openai&gl=us"
```

Query 参数：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us"
keyapi call "/v1/tiktok/shop/product-search?keyword=hot%20sales&region=US"
keyapi call /v1/google/search --qparam "q=openai" --qparam "gl=us"
```

对于 GET 接口，推荐把 query 参数直接写在接口路径里。需要用程序追加参数时，可以使用 `--qparam "k=v"` 或 `--query "k=v"`。

原始文本 body：

```powershell
keyapi call some/text/endpoint --body-raw 'plain text body'
```

`--body-raw` 会原样发送 body，并使用 `Content-Type: text/plain; charset=utf-8`。

JSON body：

```powershell
keyapi call tiktok/shop/product-search --body-json "{\"keyword\":\"dress\",\"region\":\"US\"}"
keyapi call tiktok/shop/product-search --body-json '{"timeRange":7,"countryCode":"US","page":1,"limit":100}'
```

`--body-json` 必须是合法 JSON，并使用 `Content-Type: application/json`。JSON 非法时，CLI 会在发送请求前失败。

Dry run：

```powershell
keyapi call "/v1/google/search?q=openai" --dry-run
```

默认情况下，`keyapi call` 会输出 KeyAPI 接口返回的原始响应。CLI 不会包裹或改写 API 响应。

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

常见 KeyAPI 响应码：

| Code | Message | 含义 |
| --- | --- | --- |
| `0` | `success` | 请求成功。 |
| `400` | `Missing required parameter: <name>` | 缺少必填参数或参数不合法。 |
| `401` | `apiKey is null` | API key 缺失。 |
| `402` | `apiKey quota is not enough` | API key 额度不足。 |
| `429` | `qps limit exceeded` | 请求频率超过 QPS 限制。 |
| `500` | `Internal server error` | 服务端错误。 |

这些状态码和字段由 KeyAPI 接口定义。CLI 会保持它们不变。

对于 HTTP 状态处理，`2xx` 响应会以退出码 `0` 结束。HTTP 非 `2xx` 响应仍会输出 KeyAPI 响应，然后以非零退出码结束，方便脚本和 CI 判断失败。

## 使用流程

推荐流程：

```text
列出平台 -> 列出接口描述 -> 查看接口参数 -> dry-run -> 执行请求
```

示例：

```powershell
keyapi spec list-platforms
keyapi spec list-endpoints tiktok
keyapi spec get-endpoint tiktok/influencer/detail
keyapi call "/v1/tiktok/influencer/detail?unique_id=mrbeast" --dry-run
keyapi call "/v1/tiktok/influencer/detail?unique_id=mrbeast"
```

流程规则：

- 通过阅读 `method`、`path`、`summary` 和 `description` 选择接口。
- 构造请求前先查看 `parameters` 和 `body_fields`。
- 对于 GET 接口，把 query 参数写在接口路径里，例如 `"/v1/google/search?q=openai&gl=us"`。
- 需要单独追加 query 参数时，使用 `--qparam "k=v"` 或 `--query "k=v"`。
- JSON 请求体使用 `--body-json`。
- 原始文本请求体使用 `--body-raw`。
- 不确定请求构造是否正确时，先使用 `--dry-run`。
- 将 CLI 集成到其他工具时，可使用 `docs/reference/cli-spec.yaml` 作为轻量机器可读命令契约。

## 文档

| 文件 | 作用 |
| --- | --- |
| [docs/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/index.md) | 文档入口。 |
| [docs/commands.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/commands.md) | 完整命令和参数参考。 |
| [docs/config.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/config.md) | 配置文件和环境变量。 |
| [docs/reference/cli-spec.yaml](https://github.com/EchoSell/keyapi-cli/blob/main/docs/reference/cli-spec.yaml) | 轻量机器可读 CLI 契约。 |
| [README.zh.md](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | 中文 README。 |
| [docs/zh/index.md](https://github.com/EchoSell/keyapi-cli/blob/main/docs/zh/index.md) | 中文文档入口。 |
