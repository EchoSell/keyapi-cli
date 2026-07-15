# KeyAPI CLI

---

[中文版](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | [English](https://github.com/EchoSell/keyapi-cli/blob/main/README.md)

KeyAPI CLI 是 [KeyAPI](https://keyapi.ai) 平台的轻量级命令行客户端。它可以帮助你配置认证信息、发现支持的 API 平台、查看接口 schema，并在不编写请求代码的情况下调用 API。

## 为什么选择 KeyAPI CLI？

KeyAPI 提供了 Google、TikTok、Instagram 等多个平台的接口。当项目中需要使用大量接口时，手写和维护请求代码会变得重复且容易出错。

KeyAPI CLI 提供一个统一的命令行入口：

- 配置一次 API token
- 查看支持的平台
- 查看每个平台下的接口
- 查看接口方法、路径、描述、参数和请求体字段
- 从终端直接执行 API 请求
- 将大响应保存到文件

CLI 会在请求接口详情时读取最新的 KeyAPI OpenAPI 文档，使接口描述和 schema 与当前 KeyAPI API 文档保持一致。

## 功能

- 简单的配置命令：`keyapi set`、`keyapi get`、`keyapi show`、`keyapi clear`
- API 查看命令：`keyapi spec list-platforms`、`keyapi spec list-endpoints`、`keyapi schema`
- 一个请求命令：`keyapi call`
- `spec` 和 `call` 默认返回 JSON
- 内置 `--help` 输出，适合人类用户、脚本和 Agent 阅读
- 支持 `--dry-run`，发送请求前先查看请求内容
- 支持 `--output "path"`，保存大响应

## 安装

### 从 Release 安装

1. 打开 [GitHub Releases](https://github.com/EchoSell/keyapi-cli/releases) 页面。
2. 下载适合当前系统的二进制文件。
```powershell
cd ~/keyapi-cli
keyapi --help
```

如果你希望 AI 助手帮你安装 release 二进制文件，可以将下面的提示词复制给 Claude Code、Codex、Cursor、Trae 等工具：

```text
请从 https://github.com/EchoSell/keyapi-cli/releases 下载并安装适用于当前系统的最新 KeyAPI CLI，完成可执行文件配置，使其可通过 `keyapi` 命令直接调用，并执行 `keyapi --version` 验证安装是否成功。
```

### 从源码构建

需要 Go `v1.22` 或更高版本。

克隆仓库并构建可执行文件：

```powershell
git clone https://github.com/EchoSell/keyapi-cli.git
cd keyapi-cli
go build -buildvcs=false -o keyapi.exe .
```

Linux 或 macOS 使用：

```bash
git clone https://github.com/EchoSell/keyapi-cli.git
cd keyapi-cli
go build -buildvcs=false -o keyapi .
```


如果你希望 AI 助手帮你从源码构建，可以将下面的提示词复制给 Claude Code、Codex、Cursor、Trae 等工具：

```text
请从源码安装 KeyAPI CLI：https://github.com/EchoSell/keyapi-cli.git，编译生成适用于当前系统的可执行文件，并配置为可通过 `keyapi` 命令直接调用，最后执行 `keyapi --version` 验证安装结果。
```

构建完成后，继续查看下面的 PATH 配置说明。

### 将 keyapi 加入 PATH

安装或构建完成后，可以将二进制文件放到已加入 `PATH` 的目录中。

Windows 示例：

```powershell
# 示例：keyapi.exe 存放在 C:\Tools\keyapi
$env:Path += ";C:\Tools\keyapi"
keyapi --help
```

上面的命令只会影响当前 PowerShell 会话。长期使用时，请将该目录添加到系统或用户的 `PATH` 环境变量中。

Linux 或 macOS 示例：

```bash
chmod +x keyapi
sudo mv keyapi /usr/local/bin/keyapi
keyapi --help
```

## 快速开始

配置 KeyAPI token：

```powershell
keyapi set token sk_xxx
```

查看当前配置：

```powershell
keyapi show
```

列出支持的平台：

```powershell
keyapi spec list-platforms
```

列出某个平台的接口：

```powershell
keyapi spec list-endpoints google
```

调用前查看接口详情：

```powershell
keyapi schema google/search
```

先 dry-run 查看请求：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&hl=en" --dry-run
```

执行请求：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&hl=en"
```

保存大响应：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us&hl=en" --output "google-search.json"
```

需要时清除已保存的 token：

```powershell
keyapi clear token
```

## 命令

| 命令 | 作用 |
| --- | --- |
| `keyapi --help` | 查看主帮助、推荐顺序和子命令帮助入口。 |
| `keyapi set token <value>` | 保存 API token 到本地。 |
| `keyapi set language <en\|zh>` | 设置接口描述使用的 OpenAPI 文档语言。 |
| `keyapi get <key>` | 读取一个本地配置值。 |
| `keyapi show` | 查看当前本地配置。 |
| `keyapi clear <key>` | 清除一个本地配置值。 |
| `keyapi spec list-platforms` | 返回支持的平台 JSON 数组。 |
| `keyapi spec list-endpoints <platform...\|all>` | 返回一个或多个平台的接口摘要。 |
| `keyapi spec get-endpoint <endpoint-or-path>` | 返回单个接口的方法、路径、描述、参数和 body 字段。 |
| `keyapi schema <endpoint-or-path>` | `keyapi spec get-endpoint` 的简短别名。 |
| `keyapi call <endpoint-or-path>` | 执行 KeyAPI 请求。 |

每个命令都支持 `--help`：

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

## 进阶用法

### JSON 输出

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

### 保存大响应

部分接口会返回较大的 JSON，例如商品详情、图片列表、SKU 列表或描述字段。建议使用 `--output "path"` 将完整响应体写入文件，避免直接在终端中输出大量内容：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us" --output "outputs/google-search.json"
```

相对路径会基于当前工作目录解析。

### Dry Run

使用 `--dry-run` 打印请求但不发送：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us" --dry-run
```

### 请求体

JSON 请求体使用 `--body-json`：

```powershell
keyapi call "/v1/example/post-endpoint" --method POST --body-json '{"page":1,"verified":true}'
```

原始文本请求体使用 `--body-raw`：

```powershell
keyapi call "/v1/example/text-endpoint" --method POST --body-raw "plain text"
```

对于 GET 接口，推荐直接把 query 参数写在接口路径中：

```powershell
keyapi call "/v1/google/search?q=openai&gl=us"
```

需要用程序追加 query 参数时，可以使用 `--params "k=v"` 或 `--query "k=v"`。

## 文档

`docs/` 包含 CLI 的安装、配置和使用文档。CLI 的 help 输出是内置的命令使用参考。

| 文件 | 作用 |
| --- | --- |
| [docs/index.md](docs/index.md) | 文档入口。 |
| [docs/commands.md](docs/commands.md) | 完整命令和参数参考。 |
| [docs/config.md](docs/config.md) | 配置文件和环境变量说明。 |
| [README.md](README.md) | 英文 README。 |
| [docs/zh/index.md](docs/zh/index.md) | 中文文档入口。 |
