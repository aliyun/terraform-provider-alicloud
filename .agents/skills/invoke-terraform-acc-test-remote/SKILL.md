---
name: invoke-terraform-acc-test-remote
description: Use when a Jarvis Terraform workflow needs to run terraform-provider-alicloud acceptance tests remotely via ACube/FC, especially for PR review, provider resource development, long-running TF_ACC tests, cross-account resources, or requests like "执行 AccTest", "远程测试", "跑集成测试", "不要占本地资源".
---

# Terraform AccTest 远程执行

把本地 `terraform-provider-alicloud` 代码打包上传到 ACube AccTest，由远程 FC 执行 `TF_ACC=1 go test ./alicloud`，再下载 `run.log` 和 `tf-debug.log` 做诊断。

核心原则：真实 AccTest 优先走远程 ACube/FC，避免占用本机；本地只做 `go test -run '^$'`、小单测、lint、示例 `terraform validate` 这类轻量校验。

## 边界

- 只支持本地代码上传；旧的 GitLab 分支提交路径不要再用。
- 远程执行环境负责注入 AK/SK 和测试变量；命令、示例、PR 评论、Aone 评论和日志摘要里都不要写真实 AK/SK。
- 跨账号资源只描述需要哪些环境变量，例如 `ALICLOUD_ACCESS_KEY_1` / `ALICLOUD_SECRET_KEY_1` / `ALICLOUD_ACCESS_KEY_2` / `ALICLOUD_SECRET_KEY_2`，不要展示值。
- `cancel` 会停止远程 FC 异步任务；只有用户明确同意取消时才能调用。
- API 报云产品业务错误时，以 `tf-debug.log` 的 RequestId 和请求参数为准，不凭本地复现猜测。

## 1. 定位 provider 目录

优先自动定位，只有失败时才问用户。合法目录必须同时满足：

- `<dir>/go.mod` 包含 `terraform-provider-alicloud`
- `<dir>/alicloud/` 存在并包含 provider 测试文件

```
1. Did the user explicitly give a path?
   ├─ YES  → use that path, jump to Step 1 (CLI will hard-fail if invalid)
   └─ NO   → continue to step 2

2. Auto-detect by inspecting the current working directory:
       test -f go.mod \
         && grep -q "terraform-provider-alicloud" go.mod \
         && test -d alicloud
   ├─ PASS → silently use "." as --dir, jump to Step 1
   └─ FAIL → continue to step 3

3. Walk a couple of common parents (only if the agent has obvious hints from
   the conversation, e.g., user mentioned a sibling repo). Examples worth trying:
     - ./terraform-provider-alicloud
     - ../terraform-provider-alicloud
   ├─ PASS → silently use the resolved path, jump to Step 1
   └─ FAIL → continue to step 4

4. Ask the user — only now is asking justified:
   > 当前目录看起来不是 `terraform-provider-alicloud` 仓库
   > (`go.mod` 不含 `terraform-provider-alicloud`,或者 `alicloud/` 目录缺失)。
   > 请提供本地 provider 仓库的绝对路径(例:`/Users/xxx/terraform-provider-alicloud`),
   > 或先 cd 进去再让我运行。
```

CLI 会在上传前再次校验目录；目录错误会直接失败，不会发起网络请求。

## 2. 提交并执行

以下命令里的 `{SKILL_DIR}` 是本 skill 目录的绝对路径；不要假设当前工作目录一定是 Jarvis 仓。

### 推荐：单用例

已知测试函数名时传 `--test-case`。这样更快，名称拼错也会快速失败。

```bash
python3 {SKILL_DIR}/scripts/acctest.py upload-run \
  --namespace ResourceManager \
  --resource Handshake \
  --dir "$(bootstrap/workspace.sh dir terraform_provider)" \
  --test-case TestAccAliCloudResourceManagerHandshake_basic \
  --download-dir ./acctest_logs
```

### 推荐：按 Terraform 资源名查映射

优先传 `--terraform-resource`，脚本会调用 Acube `getTerraformResourceSpec` 读取 `namespace` / `resourceTypeCode`，不要在本 skill 中维护固定映射表。映射查证规则引用 `provider-resource-dev` skill 的「Terraform 资源名解析」小节；接口查不到映射时，改用显式 `--namespace` / `--resource`。

```bash
python3 {SKILL_DIR}/scripts/acctest.py upload-run \
  --terraform-resource alicloud_schedulerx_job \
  --dir "$(bootstrap/workspace.sh dir terraform_provider)" \
  --download-dir ./acctest_logs
```

### 匹配一类用例

```bash
python3 {SKILL_DIR}/scripts/acctest.py upload-run \
  --namespace VPC --resource VSwitch \
  --dir "$(bootstrap/workspace.sh dir terraform_provider)" \
  --download-dir ./acctest_logs
```

不传 --test-case（即不设置该参数）时，远程 runner 按 `TestAccAliCloud{Namespace}{ResourceTypeCode}*` 匹配并去重执行，也就是一个任务跑该资源全部 TestAcc 用例。传 `--test-case` 时只跑指定精确函数名：单用例走 `testCaseName=<name>`，多用例（逗号分隔）经 `testCaseNames` 重复 query 参数提交（`testCaseNames=A&testCaseNames=B`）。服务端把每个值当作字面函数名，自行做 `^...$` 锚定并在 `alicloud/` 中校验存在性；客户端不再拼正则。

### 参数

| Parameter | Description | Example |
|-----------|-------------|---------|
| `--namespace` | Product namespace | `VPC`, `ECS`, `SLB`, `RDS` |
| `--resource` | Resource type code | `VSwitch`, `Instance`, `Listener` |
| `--terraform-resource` | Terraform resource name; resolves namespace/resource via Acube mapping | `alicloud_schedulerx_job` |
| `--dir` | Path to local `terraform-provider-alicloud` directory | `./terraform-provider-alicloud` |
| `--test-case` | Optional: exact test function name(s). 单用例 → `testCaseName`；逗号分隔多用例 → `testCaseNames` 重复参数（服务端自行锚定精确函数名） | `TestA,TestB` |
| `--download-dir` | Local dir for downloaded logs | `./acctest_logs` |
| `--insecure` | Skip TLS certificate verification for internal ACube endpoints when Python cert verification fails | flag |

`upload-run` 会：

1. 只打包 `go.mod`、`go.sum`、`alicloud/`
2. zip 根目录固定为 `terraform-provider-alicloud`，不使用本地 worktree 目录名
3. 上传到 ACube，服务端注入 `config.json`
4. 触发 FC 异步执行
5. 每 60 秒轮询状态，并增量打印 `run.log`
6. 终态后下载 `run.log` / `tf-debug.log`
7. 输出包含 taskId、状态、日志路径的 JSON

## 3. 长时间无日志

轮询输出 `STALL_DETECTED` 时，说明 5 分钟没有新日志。可能是资源 Apply/Destroy 正常耗时，也可能卡死。

必须先询问用户：

```
任务 <taskId> 已经 5 分钟没有新日志了，当前 status=<status>, fcStatus=<fcStatus>。
最近一条日志是：<last log line>
继续等待，还是取消远程任务？
```

只有用户明确说“取消”后才执行：

```bash
python3 {SKILL_DIR}/scripts/acctest.py cancel --task-id <taskId>
```

用户说继续等时，不要重新 `upload-run`，用 `status` 继续查同一个 task。

## 4. 分步命令

```bash
# Upload only (returns taskId)
python3 {SKILL_DIR}/scripts/acctest.py upload \
  --namespace VPC --resource VSwitch \
  --dir "$(bootstrap/workspace.sh dir terraform_provider)" \
  [--test-case TestAccAliCloudVPCVSwitch_basic]

# Query status (response includes recent log tail + FC lifecycle fields)
python3 {SKILL_DIR}/scripts/acctest.py status --task-id <taskId>

# Cancel only after user confirmation
python3 {SKILL_DIR}/scripts/acctest.py cancel --task-id <taskId>

# After terminal: get presigned download URLs / download logs
python3 {SKILL_DIR}/scripts/acctest.py logs --task-id <taskId> --download-dir ./acctest_logs
```

## 5. 读日志

下载后必须读日志再下结论：

- `run.log`: 结构化执行日志、最终状态、每个 case 的 PASS/FAIL/SKIP。
- `tf-debug.log`: Terraform DEBUG 日志，包含 API 请求、响应和 RequestId。引用给用户时只摘取参数名、错误码、RequestId，不贴 AK/SK 或敏感字段值。

### 状态解读

| ACube `status` | FC `fcStatus` | 含义 |
|---|---|---|
| `processing` | `Enqueued` / `Running` / `Retrying` | 仍在运行 |
| `success` | `Succeeded` | AccTest 通过 |
| `failed` | `Succeeded` | FC 正常结束，但测试失败；看 `tf-debug.log` |
| `failed` | `Failed` / `Stopped` | FC 实例异常或被取消；看 `fcInvocationErrorMessage` |
| `failed` | `Expired` | FC 排队过期；检查配额或重试 |

## 6. 没有跑到用例

`run.log` 显示 0 个 case 或 `no tests to run` 时：

- 未传 `--test-case`：`TestAccAliCloud{Namespace}{ResourceTypeCode}*` 没匹配到。常见原因：
  - Namespace casing: `Vpc` vs `VPC` (case-sensitive)
  - Resource code spelling: `VSwitch` vs `Vswitch`
  - Legacy naming: `TestAccAlicloudXxx` (lowercase 'c') vs `TestAccAliCloudXxx`
  - Underscore in name: `TestAcc_AliCloudXxx`
- 传了 `--test-case`：如果拼错，FC 会快速报 `指定的测试用例 X 在 alicloud/ 中未找到`。

本地查真实函数名：

```bash
grep -n "^func TestAcc" alicloud/resource_alicloud_<name>_test.go
```

修正后用准确的 `--test-case` 重跑。

## 常见失败

| Symptom | Likely Cause | Action |
|---------|-------------|--------|
| 0 tests run | Name mismatch (see Step 4) | Fix test function names or pass exact `--test-case` |
| FC 解压后找不到 `terraform-provider-alicloud` | 本地打包根目录错误或旧版脚本使用了 worktree 目录名 | 更新/使用本 skill 的 `acctest.py`; zip root 必须固定为 `terraform-provider-alicloud` |
| `指定的测试用例 X 在 alicloud/ 中未找到` | `--test-case` value misspelled | grep local test file for actual names |
| 需要跑多个具体用例 | 旧写法把 `A,B` 当作一个函数名 | 用逗号分隔多个用例；脚本经 `testCaseNames` 重复参数提交，服务端锚定各自精确函数名;或不传 `--test-case` 跑该资源全部用例 |
| `CERTIFICATE_VERIFY_FAILED` | 本机 Python 证书链不信任内网 ACube 证书 | 在 `acctest.py` 后、子命令前加 `--insecure` |
| `Unsupported argument "X"` | Test HCL references a deleted field | Update test config to use the replacement field |
| `daring resource` / destroy timeout | Subscription resource can't be destroyed | Treat as PASS if all Apply/Read steps succeeded |
| Cloud API error with RequestId | API-side issue | Check RequestId in `tf-debug.log`, investigate params |
| `fcStatus=Failed` + `fcInvocationErrorMessage` set | FC instance crash (OOM / panic / timeout) | Check FC console runtime logs |
| `QuotaExceeded.<Resource>` at Create | Test account quota in that region is full from leftover resources | Run `make sweep REGION=<region> RESOURCE=<resource_type>` locally to clean up, then re-submit (see error-patterns §7) |

详细诊断见 `references/error-patterns.md`。

## 下载单个文件

需要从预签名 OSS URL 单独下载日志时：

```bash
python3 {SKILL_DIR}/scripts/acctest.py download --url "<presigned_url>" --download-dir ./acctest_logs
```

## Tools

| Tool | Path | Description |
|------|------|-------------|
| AccTest CLI | `{SKILL_DIR}/scripts/acctest.py` | Submit / poll / cancel / download AccTest tasks (pure Python 3, no dependencies) |
