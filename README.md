# cpa-plusins

CLIProxyAPI 自定义插件商店，提供 CommandCode 和 New API 两个提供商插件。添加本商店后，可在 CLIProxyAPI 管理页面安装插件，再通过各插件的管理页面配置账号和上游服务。

## 可用插件

| 插件 | 主要功能 | 使用说明 | 下载 |
|---|---|---|---|
| CommandCode Provider | 接入 CommandCode，提供认证、模型发现、请求转发和额度查询 | [CommandCode 使用说明](https://github.com/careerisx/commandcode-cliproxyapi-provider#readme) | [Releases](https://github.com/careerisx/commandcode-cliproxyapi-provider/releases) |
| New API Provider | 按认证配置上游站点，适配 Chat Completions、Messages、Responses，查看模型和额度 | [New API 使用说明](https://github.com/careerisx/newapi-cliproxy-provider#readme) | [Releases](https://github.com/careerisx/newapi-cliproxy-provider/releases) |

两个插件可按需安装，不要求同时启用。

## 使用条件

- 使用支持原生动态插件的 CLIProxyAPI v8，不能使用禁用插件功能的宿主构建。
- 当前提供的预编译插件面向 **Linux x86_64（`linux/amd64`）**，不是 Windows、macOS 或 ARM64 安装包。
- 运行环境需要兼容插件依赖的 glibc；不要直接将其当作 Alpine/musl 插件使用。
- 宿主能够访问商店索引、GitHub Release 下载地址以及对应插件的上游服务。
- 准备好 CLIProxyAPI 管理页面的访问凭据；上游账号或 API key 按各插件说明配置。

## 方式一：通过插件商店安装

### 1. 添加商店源

在 CLIProxyAPI 配置文件中合并以下内容。如果已有 `plugins` 节，请在原配置中补充，避免重复定义；已有 `store-sources` 时追加本商店地址。

```yaml
plugins:
  enabled: true
  dir: "./plugins"
  store-sources:
    - "https://raw.githubusercontent.com/careerisx/cpa-plusins/main/registry.json"
```

`dir` 是宿主加载插件的目录，可按实际部署位置调整。添加第三方商店不会替换官方商店。

### 2. 安装并配置

1. 保存配置并重启 CLIProxyAPI，使商店源配置生效。
2. 打开 CLIProxyAPI 管理页面中的插件商店，找到 CommandCode Provider 或 New API Provider。
3. 安装所需插件。如果宿主提示需要重启，请按提示操作。
4. 按对应插件的使用说明启用插件，并在其管理页面添加认证。
5. 在宿主中确认模型已可用，再在客户端选择对应模型。

本仓库只提供插件商店索引，不需要作为单独服务部署，也不需要把本仓库放到宿主的插件目录。

## 方式二：手动安装

不通过插件商店时，可从各插件的 Releases 页面下载对应版本的 `linux_amd64.zip` 和同一 Release 的 `checksums.txt`。

| 插件 | ZIP 文件名格式 | 解压后的插件文件 |
|---|---|---|
| CommandCode | `commandcode-cliproxyapi-provider_<版本>_linux_amd64.zip` | `commandcode-cliproxyapi-provider.so` |
| New API | `newapi-provider_<版本>_linux_amd64.zip` | `newapi-provider.so` |

操作步骤：

1. 下载所需插件的 ZIP 和校验和文件到同一目录。
2. 可执行 `sha256sum -c checksums.txt` 检查下载文件是否完整。
3. 停止 CLIProxyAPI，将 ZIP 中的 `.so` 解压到宿主配置的 `plugins.dir` 目录，不要只复制 ZIP。
4. 确认 `plugins.enabled: true`，并按对应插件说明配置启用项。
5. 启动 CLIProxyAPI，检查插件是否加载，再配置认证。

容器部署时，插件文件应位于容器内的 `plugins.dir` 对应路径，并挂载到持久化目录，避免重建容器后丢失。宿主进程需要有读取插件文件的权限；通过商店安装和更新时还需要目录写权限。

## 更新与卸载

| 操作 | 建议方式 |
|---|---|
| 商店更新 | 在插件商店检查更新并安装；提示需重启时重启宿主 |
| 手动更新 | 停止宿主，备份原插件和配置，使用新版本 ZIP 中的 `.so` 替换旧文件，再启动宿主 |
| 暂时停用 | 按插件文档将对应配置中的 `enabled` 设为 `false`，使配置生效 |
| 卸载 | 使用宿主提供的卸载操作；手动卸载时先停止宿主，再移除对应 `.so` 和不再使用的启用配置 |

更新或卸载前建议备份宿主配置和认证文件。不要为了清理插件而删除整个认证目录。

## 常见问题

| 问题 | 排查方式 |
|---|---|
| 商店里找不到插件 | 检查商店源 URL、配置是否生效，以及宿主能否访问索引地址 |
| 能看到插件但下载失败 | 检查对应仓库是否有已发布的匹配平台安装包，并检查 GitHub 下载访问情况 |
| 显示的版本与预期不一致 | 区分商店版本和已安装版本；无法获取 Release 信息时，商店可能显示缓存或索引兜底版本，刷新页面不等于更新已安装插件 |
| 插件文件存在但没有加载 | 检查宿主是否支持原生插件、系统架构、文件路径、读取权限及启动日志中的动态库加载错误 |
| 插件加载后没有模型 | 按插件说明添加有效认证，确认上游地址和凭据可用，并检查模型发现结果 |
| 容器重建后插件消失 | 检查插件目录是否已持久化挂载，以及配置路径是否指向容器内实际目录 |
| CommandCode 改名后配置未生效 | 当前插件 ID 为 `commandcode-cliproxyapi-provider`；旧安装需迁移配置 key，并从新插件菜单重新打开管理页面 |

CommandCode 的认证和模型归属标识仍为 `commandcode`；New API 的认证和模型归属标识为 `newapi`。它们与插件 ID 不同，不要全局替换。

请妥善保管宿主管理密码、API key 和认证文件，不要将其提交到公开仓库，也不要把插件管理页面直接暴露给不可信访问者。

## 相关链接

- [CLIProxyAPI 官方配置说明](https://help.router-for.me/cn/configuration/options)
- [CommandCode 插件](https://github.com/careerisx/commandcode-cliproxyapi-provider)
- [New API 插件](https://github.com/careerisx/newapi-cliproxy-provider)
