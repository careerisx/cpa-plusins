# cpa-plusins

CLIProxyAPI v8 自定义插件商店索引。两个插件的源码、版本和 Release 在独立仓库中维护，本仓库只负责登记和使用说明。

索引仓库：<https://github.com/careerisx/cpa-plusins>。

## 插件列表

| 插件 | 插件 ID | ProviderID | 展示兜底版本 | 源码与 Release 仓库 |
|---|---|---|---|---|
| CommandCode Provider | `commandcode-cliproxyapi-provider` | `commandcode` | `0.2.0` | [careerisx/commandcode-cliproxyapi-provider](https://github.com/careerisx/commandcode-cliproxyapi-provider) |
| New API Provider | `newapi-provider` | `newapi` | `0.2.0` | [careerisx/newapi-cliproxy-provider](https://github.com/careerisx/newapi-cliproxy-provider) |

插件 ID 对应宿主注册名称和商店索引中的 `id`，不是源码目录名，也不是认证和模型归属使用的 ProviderID。索引的 `author` 与现有插件注册元信息保持一致，不代表仓库所有者；许可及版权声明见各仓库 LICENSE 文件。

New API 的认证配置与管理入口见其 [README](https://github.com/careerisx/newapi-cliproxy-provider/blob/main/README.md)。

## 添加到 CLIProxyAPI 插件商店

将以下内容合并到宿主现有配置的 `plugins` 节中，不要重复定义同名 YAML 根节点：

```yaml
plugins:
  enabled: true
  dir: "./plugins"
  store-sources:
    - "https://raw.githubusercontent.com/careerisx/cpa-plusins/main/registry.json"
```

索引文件需要先推送到本仓库的 `main` 分支，并确保宿主能够访问上述 URL。

`store-sources` 追加第三方商店，不替换官方商店。这只是在使用者的宿主中添加本仓库，不会自动将插件提交到官方商店。

**登记索引不等于已经发布插件。** 安装需要对应插件仓库中已发布的 Release、匹配平台的 ZIP 和校验和；当前工作流仅构建 `linux/amd64`。宿主本身也必须支持原生插件，不能使用禁用插件的构建。

## 索引与更新来源

根目录 [registry.json](./registry.json) 使用 `schema_version: 1`，在 `plugins` 数组中登记两个插件。每个条目的 `repository` 和对应插件源码中的 `GitHubRepository` 均指向该插件的独立 GitHub 仓库，而不是本索引仓库。

| 字段 | 作用 |
|---|---|
| `id` | 插件稳定标识，保持与插件注册名称一致 |
| `name`、`description` | 商店展示名称及说明 |
| `author` | 与插件当前注册元信息保持一致 |
| `repository` | 对应插件的 GitHub 源码及 Release 仓库 |
| `version` | 展示兜底版本，不决定实际安装版本 |
| `license`、`tags` | 许可及分类信息 |

宿主分别读取两个插件仓库的 latest release，获取实际安装版本及匹配的资产。发布新版本不需要将固定下载 URL 写回索引，也不需要跨仓提交；可以按需手动更新索引中的展示兜底版本。

## 通过 tag 自动发布

两个插件仓库各自提供 `.github/workflows/release.yml`，互不依赖，不需要统一版本号。

| 项目 | 发布行为 |
|---|---|
| 触发条件 | 向对应插件仓库推送正式 `vX.Y.Z` tag，例如 `v0.1.3`；不支持预发布或构建后缀 |
| Go 版本 | 读取各自 `go.mod`，当前为 `1.26.7` |
| 构建环境 | `ubuntu-22.04`，`CGO_ENABLED=1`，`GOOS=linux`，`GOARCH=amd64`，`GOWORK=off` |
| 构建方式 | `go build -mod=readonly -trimpath -buildmode=c-shared`，不启用 debug tag |
| 版本来源 | 从 tag 去掉前导 `v`，通过 `-ldflags -X` 注入 `internal/plugin.pluginVersion` |
| 发布权限 | 使用仓库内置 `GITHUB_TOKEN` 的 `contents: write`，无需另配 PAT |
| 发布内容 | 当前插件 ZIP 和 `checksums.txt` |

构建时使用 Go 模块路径注入版本，不使用 GitHub 仓库名推导包路径。CommandCode 的模块名、插件 ID 和发布产物名称统一为 `commandcode-cliproxyapi-provider`。

CommandCode 重命名后，宿主会将其识别为新的插件 ID。已有安装需要迁移 `plugins.configs` 下的配置 key，改用新插件的管理页面地址；旧页面使用的浏览器存储不会自动迁移。认证和模型归属的 ProviderID 仍为 `commandcode`，不需要改名。重命名后的源码应以新 tag 发布匹配新名称的 ZIP，再同步线上索引；已有 tag 和 Release 不会因源码改名而自动更新，不要覆盖旧版本。

未指定版本注入的本地构建仍使用源码默认版本：CommandCode `0.1.3`、New API `0.1.0`。正式 Release 的插件注册版本则始终来自本次 tag。

### 发布顺序

1. 在对应插件仓库提交并推送源码及工作流；本索引的地址改动也需要提交并推送。
2. 为包含工作流的目标提交创建正式版本 tag，并将 tag 推送到对应插件仓库。
3. 在该仓库的 Actions 页面查看构建和发布结果。
4. 确认 Release 中存在该版本的 ZIP 及 `checksums.txt`，再通过 CLIProxyAPI 商店安装或更新。

已有 tag 在工作流加入前就已推送，不会因后来新增工作流而自动触发。工作流不覆盖已有同名 Release；已发布版本应保持不变，修复后使用新 tag 发布。

### Release 资产布局

ZIP 命名格式：

```text
<pluginID>_<version>_linux_amd64.zip
```

首次按当前默认版本发布时，两个独立仓库分别生成：

| 插件仓库 | ZIP | ZIP 根目录中的文件 |
|---|---|---|
| CommandCode | `commandcode-cliproxyapi-provider_0.1.3_linux_amd64.zip` | `commandcode-cliproxyapi-provider.so` |
| New API | `newapi-provider_0.1.0_linux_amd64.zip` | `newapi-provider.so` |

ZIP 内只包含该插件共享库，不含外层目录或构建时生成的 C 头文件。每个 Release 单独附带 `checksums.txt`，内容由实际 ZIP 的 SHA-256 摘要生成，工作流在发布前执行校验。

插件使用 CGO 动态库，目标宿主环境需要兼容构建产物的系统库。当前工作流没有承诺 Alpine/musl 或比构建环境更旧的 glibc 兼容性。

## 验证建议

本次配置未在本机执行构建、测试或远程发布。可以按需手动检查索引 JSON（需要已有 Node.js）：

```bash
node -e 'const fs = require("node:fs"); const r = JSON.parse(fs.readFileSync("registry.json", "utf8")); console.log("schema_version:", r.schema_version); console.table(r.plugins.map(({id, repository, version}) => ({id, repository, version})));'
```

发布后，下载对应 Release 的 ZIP 和 `checksums.txt` 到同一目录，可以在 Linux 中执行：

```bash
sha256sum -c checksums.txt
```

最终应在支持原生插件的 Linux/amd64 CLIProxyAPI 宿主中确认商店读取、资产下载、插件加载和注册版本。索引解析成功或工作流文件已创建不代表在线安装已经验证。

## 参考资料

- [插件开发文档：商店索引与发布要求](https://cliproxyapi.gh.miniasp.com/cn/plugin/development.html)
- [官方插件商店索引](https://raw.githubusercontent.com/router-for-me/CLIProxyAPI-Plugins-Store/main/registry.json)
- [CLIProxyAPI 配置说明](https://help.router-for.me/cn/configuration/options)
