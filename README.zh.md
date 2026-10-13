# KeyAPI CLI

[中文](https://github.com/EchoSell/keyapi-cli/blob/main/README.zh.md) | [English](https://github.com/EchoSell/keyapi-cli/blob/main/README.md)

KeyAPI CLI 是用于 [KeyAPI](https://keyapi.ai) 的轻量级命令行工具。它可以配置认证信息、发现支持的平台、查看接口 schema，并直接从终端调用 API。

## 功能

- 配置命令：`keyapi set`、`keyapi get`、`keyapi show`、`keyapi clear`
- 接口发现：`keyapi spec list-platforms`、`keyapi spec list-endpoints`
- schema 查看：`keyapi schema`
- 请求执行：`keyapi call`
- 默认 JSON 输出，适合脚本和 Agent 使用
- 支持 `--dry-run` 预览请求
- 支持 `--output "path"` 保存大响应

## 安装

需要 Node.js 18 或更高版本：

```bash
npm install -g keyapi-cli
keyapi --help
keyapi --version
```

npm 包会安装一个 JavaScript 启动器，并自动安装当前系统对应的 KeyAPI 原生二进制包。如果你已经有自定义 KeyAPI 二进制文件，可以在运行 `keyapi` 前设置 `KEYAPI_CLI_BINARY` 指向该文件路径。

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

调用前查看接口 schema：

```powershell
keyapi schema google/search
```

预览请求：

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

清除 token：

```powershell
keyapi clear token
```

## 常用命令

| 命令 | 作用 |
| --- | --- |
| `keyapi --help` | 查看帮助 |
| `keyapi set token <value>` | 保存 API token |
| `keyapi set language <en|zh>` | 设置 OpenAPI 文档语言 |
| `keyapi show` | 查看当前配置 |
| `keyapi spec list-platforms` | 列出支持的平台 |
| `keyapi spec list-endpoints <platform...|all>` | 列出平台接口 |
| `keyapi schema <endpoint-or-path>` | 查看接口 schema |
| `keyapi call <endpoint-or-path>` | 调用接口 |

每个命令都支持 `--help`：

```powershell
keyapi --help
keyapi spec --help
keyapi schema --help
keyapi call --help
```

## 使用建议

推荐流程：

1. `keyapi set token sk_xxx`
2. `keyapi spec list-platforms`
3. `keyapi spec list-endpoints <platform>`
4. `keyapi schema <endpoint-or-path>`
5. `keyapi call <endpoint-or-path> --dry-run`
6. `keyapi call <endpoint-or-path>`

对于较大的 JSON 响应，使用 `--output "path"` 保存到文件。
