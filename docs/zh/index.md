# KeyAPI CLI 文档

KeyAPI CLI 是一个小型命令行工具，用于配置 KeyAPI token、发现 API 接口、查看接口参数并执行请求。

## 命令组

| 命令 | 作用 |
| --- | --- |
| `keyapi config` | 管理本地 token 和 CLI 设置。 |
| `keyapi spec list-platforms` | 列出当前可用的 KeyAPI 平台。 |
| `keyapi spec list-endpoints <platform...\|all>` | 列出一个或多个平台的接口描述。 |
| `keyapi spec get-endpoint <endpoint-or-path>` | 获取单个接口的方法、路径、描述、参数和 body 字段。 |
| `keyapi call <endpoint-or-path>` | 执行 KeyAPI 请求。 |

## 文档

| 文件 | 作用 |
| --- | --- |
| [commands.md](commands.md) | 完整命令参考，包含参数、flags、返回值和 JSON 字段含义。 |
| [config.md](config.md) | 配置文件和环境变量说明。 |
| [../reference/cli-spec.yaml](../reference/cli-spec.yaml) | 轻量机器可读命令契约，仅维护英文版。 |

## 推荐流程

```powershell
keyapi config set token sk_xxx
keyapi spec list-platforms
keyapi spec list-endpoints google
keyapi spec get-endpoint google/search
keyapi call google/search --query q=openai
```

## 接口选择说明

选择接口时，先列出目标平台的接口描述，并比较：

- `method`
- `path`
- `summary`
- `description`

选择接口后，执行 `get-endpoint`，再根据返回的 `parameters` 和 `body_fields` 构造请求。

请求构造规则：

- query 参数使用 `--query`。
- JSON 请求体使用 `--data` 或 `--body-json`。
- 不确定请求是否正确时，先使用 `--dry-run`。
