# 配置

KeyAPI CLI 会把本地配置保存为 JSON 文件，并支持通过环境变量覆盖部分配置。

## 配置位置

默认情况下，配置目录基于操作系统的用户配置目录：

```text
<user-config-dir>/keyapi/config.json
```

可以通过环境变量覆盖配置目录：

```powershell
$env:KEYAPI_HOME='C:\path\to\keyapi-home'
```

设置 `KEYAPI_HOME` 后，配置文件路径为：

```text
<KEYAPI_HOME>/config.json
```

## 配置文件结构

```json
{
  "api_key": "sk_xxx",
  "base_url": "https://api.keyapi.ai",
  "auth_scheme": "Bearer",
  "language": "en"
}
```

说明：

- `api_key` 是 KeyAPI token。
- `base_url` 是 KeyAPI API 地址，默认值为 `https://api.keyapi.ai`。
- `auth_scheme` 是认证方案，默认值为 `Bearer`。
- `language` 是 OpenAPI 文档语言，默认值为 `en`。

## 环境变量

| 变量 | 说明 | 优先级 |
| --- | --- | --- |
| `KEYAPI_API_KEY` | KeyAPI token，会覆盖配置文件中保存的 token。 | API token 最高优先级 |
| `KEYAPI_HOME` | 配置目录。 | 配置位置最高优先级 |


## 管理 Token

保存 token：

```powershell
keyapi set token sk_xxx
```

查看打码后的 token：

```powershell
keyapi get token
```

删除已保存的 token：

```powershell
keyapi clear token
```

临时使用 token 时，推荐使用环境变量：

```powershell
$env:KEYAPI_API_KEY='sk_xxx'
```

## 管理文档语言

```powershell
keyapi set language en
keyapi set language zh
```

| 值 | 含义 |
| --- | --- |
| `en` | 使用英文 OpenAPI 文档，默认值。 |
| `zh` | 使用中文 OpenAPI 文档。 |
