# CloudSpec 本地开发工作流

本页用于 CloudSpec 项目的本地环境配置、验证命令和交付前检查。所有命令都在存在 `main.cspec` 的项目根目录执行。

## 1. 本地环境检查

```bash
pwd
test -f main.cspec
command -v aliyun
aliyun version
aliyun configure list
aliyun cspec --help
aliyun cspec baseinfo
```

检查点：

- `main.cspec` 存在，否则先通过 `cloudspec-amp-workflow` 初始化、建分支、clone cspec 仓库。
- `aliyun configure list` 能看到用户本机 profile 或凭证配置。不要打印或复述 AK/SK 明文。
- `aliyun cspec --help` 能执行，说明 cspec 插件已安装。
- `aliyun cspec baseinfo` 能读取当前项目，输出 `namespace`、`apiStyle`、`isInnerApi`、资源和操作列表。

如缺少 cspec 插件：

```bash
aliyun plugin install --names cspec --source-base https://cli.aliyun-inc.com/registry_id/2/env/pre/plugins
```

## 2. 编辑前检查

```bash
git branch --show-current
cat .amp/context.yaml 2>/dev/null
aliyun cspec baseinfo
```

- 不在 `master` / `main` 上直接编辑。
- 先识别 `apiStyle`：RPC 或 ROA。
- 先识别 `isInnerApi`：inner API 只需要 build，不强制跑规范 check。
- 修改前先读同项目同类资源、操作和测试用例，保持注解结构、命名、后端配置风格一致。

## 3. build / check / test

每次修改 `.cspec` 后先运行：

```bash
aliyun cspec build
```

如果改动了 `resources/*.cspec`，且不是 inner API：

```bash
aliyun cspec check --name <ResourceName>
```

如果改动了 operation 或用户要求规范校验，且不是 inner API：

```bash
aliyun cspec check --name <OperationName>
```

如果改动或生成了资源测试：

```bash
aliyun cspec test run -n <MainTestName>
```

资源测试同一时间只运行一个入口用例；失败时保留完整命令和关键日志，再判断是环境、元数据、operation 映射还是测试数据问题。

## 4. 本地转 yaml

不同版本 CLI 的 yaml 子命令可能不同，先确认当前版本支持的参数：

```bash
aliyun cspec --help
aliyun cspec yaml --help
```

若当前版本支持 `aliyun cspec yaml`，按 help 输出执行，例如：

```bash
aliyun cspec yaml -o <output-dir>
```

不要在没有 help 验证的情况下猜测 yaml 参数。若 `aliyun cspec yaml --help` 不存在，记录 CLI 版本和 help 输出，再向用户说明当前环境不支持该命令。

## 5. 资源测试生成前检查

资源测试语法、生成和运行走 `cloudspec-resource-test`，生成前必须额外确认：

```bash
command -v cloudspec-agent
cloudspec-agent --help
```

- cloudspec-agent 已安装。
- 安装后 MCP 服务已按提示配置完成。
- 用户明确确认目标资源相关后端代码已经发布到 POP 网关，目标环境可由 POP 网关正常调用。
- 资源测试执行账号优先使用 `execConfigUuid` / account uuid；不要把普通云账号 UID/UserId 当成执行配置入口。
- 当前 CLI 默认可能允许本地 AK/SK 覆盖 `execConfigUuid` 的账号配置；需要强制使用 uuid 时，在测试 runtime 中设置 `accountOverrideSupport: false`，或在结果中明确记录覆盖风险。
- 资源测试失败涉及 API 行为或后端配合时，输出诊断报告，不直接修改 API 原始定义。

未满足以上条件时，不运行 `cloudspec-agent --auto-login-mcp cover`。

## 6. 提交前清单

- `git diff` 只包含本次需求相关文件。
- `aliyun cspec build` 已通过，或明确说明未运行原因。
- 需要规范检查的组件已跑 `aliyun cspec check --name <Component>`。
- 需要测试的入口已跑 `aliyun cspec test run -n <MainTestName>`。
- 无法安全修复的规范问题已列出原因和需用户确认项。
- 没有提交 AK/SK、token、`.amp/local.yaml` 等本地敏感配置。
