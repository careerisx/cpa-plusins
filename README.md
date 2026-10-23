# cpa-plusins

CLIProxyAPI v8 自定义插件仓库，集中维护 CommandCode 和 New API 两个原生提供商插件。

仓库地址：<https://github.com/careerisx/cpa-plusins>。

## 插件列表

| 插件 | 插件 ID | ProviderID | 当前源码版本 | 源码目录 |
|---|---|---|---|---|
| CommandCode Provider | `commandcode-go-cliproxyapi` | `commandcode` | `0.1.3` | [commandcode-cliproxyapi-provider](./commandcode-cliproxyapi-provider/) |
| New API Provider | `newapi-provider` | `newapi` | `0.1.0` | [newapi-cliproxy-provider](./newapi-cliproxy-provider/) |

插件 ID 对应宿主注册名称和商店索引中的 `id`，不是源码目录名，也不是认证和模型归属使用的 ProviderID。索引的 `author` 与现有插件注册元信息保持一致，不代表仓库所有者；许可及版权声明见现有 LICENSE 文件。

New API 的认证配置与管理入口见其 [README](./newapi-cliproxy-provider/README.md)。

## 仓库结构

```text
cpa-plusins/
├── registry.json                         # 自定义插件商店索引
├── README.md                             # 仓库与发布说明
├── commandcode-cliproxyapi-provider/       # CommandCode 插件独立 Go 模块
└── newapi-cliproxy-provider/               # New API 插件独立 Go 模块
```

## 添加到 CLIProxyAPI 插件商店

将以下内容合并到宿主现有配置的 `plugins` 节中，不要重复定义同名 YAML 根节点：

```yaml
plugins:
  enabled: true
  dir: "./plugins"
  store-sources:
    - "https://raw.githubusercontent.com/careerisx/cpa-plusins/main/registry.json"
```

上述 URL 以默认分支 `main` 为例；如果实际分支不同，需要替换分支名。文件推送后还需确保宿主能访问该 URL。

`store-sources` 追加第三方商店，不替换官方商店。这只是在使用者的宿主中添加本仓库，不会自动将插件提交到官方商店。

**当前仅配置了源码侧的商店索引和仓库元信息，未在本次操作中创建或发布 GitHub Release，也未验证在线安装。** 添加索引不等于插件已经可以下载；安装还需要匹配宿主平台的 Release 产物和校验和。

## 索引与更新来源

根目录 [registry.json](./registry.json) 使用 `schema_version: 1`，在 `plugins` 数组中登记两个插件。两者的 `repository` 和源码注册元信息中的 `GitHubRepository` 均指向本仓库。

| 字段 | 作用 |
|---|---|
| `id` | 插件稳定标识，保持与插件注册名称一致 |
| `name`、`description` | 商店展示名称及说明 |
| `author` | 与插件当前注册元信息保持一致 |
| `repository` | GitHub 源码及 Release 仓库地址 |
| `version` | 展示兜底版本，不决定实际安装版本 |
| `license`、`tags` | 许可及分类信息 |

宿主实际安装版本来自 GitHub latest release tag，而不是索引中的 `version`。修改插件源码版本时应同步维护索引，避免展示信息过期。

## Release 发布要求

根据插件开发文档，每个平台的 ZIP 资产命名为：

```text
<pluginID>_<version>_<goos>_<goarch>.zip
```

ZIP 根目录直接包含该插件动态库，不要再套一层目录。动态库扩展名按目标平台使用 `.so`、`.dll` 或 `.dylib`；实际平台能否使用插件还取决于宿主构建是否支持原生插件。

发布时还需附带 `checksums.txt`，使用标准 SHA-256 校验和格式，每个 ZIP 对应一行：64 位十六进制摘要、两个空格、ZIP 文件名。摘要必须基于实际发布文件计算，不能使用占位内容。

### 同仓库多插件的版本限制

两个插件当前源码版本分别为 `0.1.3` 和 `0.1.0`，但它们共用同一 GitHub 仓库。宿主依据仓库的 latest release tag 获取版本，不能据此假定每个插件拥有独立的 latest release。

在加入自动发布前，需结合目标宿主实现确认同仓库多插件的资产选择、版本匹配和更新行为，再确定统一发布版本等策略。当前未修改两个插件的源码版本，也未配置自动构建或发布流程。

## 本地验证建议

本次修改未执行构建或测试。可以按需手动检查 JSON 是否能够解析：

```bash
node -e 'const fs = require("node:fs"); const r = JSON.parse(fs.readFileSync("registry.json", "utf8")); console.log("schema_version:", r.schema_version); console.table(r.plugins.map(({id, repository, version}) => ({id, repository, version})));'
```

该命令需要本机已有 Node.js，仅检查 JSON 解析并显示关键字段，不验证宿主安装兼容性。发布后仍应在支持原生插件的 CLIProxyAPI 宿主中验证商店读取、资产下载及插件加载。

## 参考资料

- [插件开发文档：商店索引与发布要求](https://cliproxyapi.gh.miniasp.com/cn/plugin/development.html)
- [官方插件商店索引](https://raw.githubusercontent.com/router-for-me/CLIProxyAPI-Plugins-Store/main/registry.json)
- [CLIProxyAPI 配置说明](https://help.router-for.me/cn/configuration/options)
