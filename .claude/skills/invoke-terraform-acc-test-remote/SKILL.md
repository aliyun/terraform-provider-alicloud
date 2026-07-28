---
name: invoke-terraform-acc-test-remote
description: Use when a Jarvis Terraform workflow needs to run terraform-provider-alicloud acceptance tests remotely via ACube/FC, especially for PR review, provider resource development, long-running TF_ACC tests, cross-account resources, or requests like "执行 AccTest", "远程测试", "跑集成测试", "不要占本地资源".
---

# Terraform AccTest 远程执行

把本地 `terraform-provider-alicloud` 代码打包上传到 ACube AccTest，由远程 FC 执行 `TF_ACC=1 go test ./alicloud`，再下载 `run.log` 和 `tf-debug.log` 做诊断。

核心原则：真实 AccTest 优先走远程 ACube/FC，避免占用本机；本地只做 `go test -run '^$'`、小单测、lint、示例 `terraform validate` 这类轻量校验。

## 边界

- 只支持本地代码上传，不走 GitLab 分支提交。
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
   ├─ YES  → use that path, jump to §2 提交并执行 (CLI will hard-fail if invalid)
   └─ NO   → continue to step 2

2. Auto-detect by inspecting the current working directory:
       test -f go.mod \
         && grep -q "terraform-provider-alicloud" go.mod \
         && test -d alicloud
   ├─ PASS → silently use "." as --dir, jump to §2 提交并执行
   └─ FAIL → continue to step 3

3. Resolve from the workspace registry, then common parents:
     - "$(bash <jarvis仓>/bootstrap/workspace.sh dir terraform_provider)"   # 登记真源(CLAUDE.md 工作纪律 #4)
     - ./terraform-provider-alicloud
     - ../terraform-provider-alicloud
   ├─ PASS → silently use the resolved path, jump to §2 提交并执行
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
| `success` | `Succeeded` | runner 未 FAIL；真 PASS/SKIP 需读 run.log 才能区分（见「SKIP ≠ PASS」小节） |
| `failed` | `Succeeded` | FC 正常结束，但测试失败；看 `tf-debug.log` |
| `failed` | `Failed` / `Stopped` | FC 实例异常或被取消；看 `fcInvocationErrorMessage` |
| `failed` | `Expired` | FC 排队过期；检查配额或重试 |

### SKIP ≠ PASS —— runner "success" 判定的陷阱

远程 runner 把 `--- SKIP: X (0.00s)` 也判成 `success`（"非 FAIL" 即绿），JSON 里的 `最终测试状态: passed` 和 `testStatus: "passed"` 都是 **"没跑失败"** 的意思，不是 **"验证过"**。看到 `status=success` 就以为绿了 → 用例可能连一个 API 包都没发出去。

**每次 upload-run 拿到终态后必做**：从 run.log 读真实的 case 结果，别只看顶层 status：

```bash
grep -E "^\s*--- (PASS|FAIL|SKIP):" acctest_logs/<taskId>_run.log
# 或找 "统计: PASS=x FAIL=y SKIP=z" 那行
```

- `PASS` → 真验证过
- `SKIP` → **没验证**，别当绿；进入下方排查链
- `FAIL` → 挂了

#### SKIP → PASS 的排查链

SKIP 通常来自 `testAccPreCheckWithXxx(t)` 门闩——读某 env var 不到就 `t.Skipf`。**修法不是删门闩**（gate 本身是设计出来兜"账号没这能力就别跑"的），而是**让门闩不必要**：

1. `grep -n "func testAccPreCheckWith" alicloud/provider_test.go` 找到 gate 函数，看它读什么 env var / 检查什么条件；
2. 全文 `grep -rn testAccPreCheckWithXxx alicloud/` 看**谁在用这个 gate**——常常只有 1-2 个 test；
3. 反过来找**不调这个 gate 的兄弟资源 test**（如 `resource_alicloud_xxx_test.go`），看它的 test config HCL 怎么让所需功能真跑起来——**多半是资源 schema 里某个 `config` / `enable_xxx` 字段在 create 时就把能力开出来**；
4. 把兄弟的 test config 模式抄到你的 test 里，preCheck 换成 `testAccPreCheckWithRegions` 之类通用门，原 gate 就不用调了。

#### 实例：alikafka_sasl_acls datasource test（工单 83992418）

- Gate `testAccPreCheckWithAlikafkaAclEnable` 要 `ALICLOUD_ALIKAFKA_ACL_ENABLE=true`，数年 SKIP，没人真跑过；
- 兄弟 `resource_alicloud_alikafka_sasl_acl_test.go` 从不调它——它建实例时直接 `config = "{\"enable.acl\":\"true\"}"` 把 ACL 开出来，precheck 只 pin `cn-hangzhou`；
- 抄兄弟的 config（serverless + `enable.acl` JSON）+ preCheck（region pin + basic）→ 测试真跑通（680s：create instance → sasl_user → sasl_acl → query datasource → destroy）；
- **副产品**：真跑之后暴露了断言与 API 返回大小写不一致（`Topic`/`Write` 断言 vs API 返回 `TOPIC`/`WRITE`）——这类**只有真跑才能发现的差异**，是 SKIP 一直掩盖的老坑。

#### 反模式：改 gate function 强行跑过 SKIP

把 `testAccPreCheckWithXxx` 改成 no-op（`_ = os.Getenv(...)` 之类）只能**假装绕开 Go 侧 gate**，阿里云 API 侧真门（如 `BIZ_ACL_NOT_ENABLED`）会立刻把测试挂掉，而且引入的空调用还得清理。**正确路径永远是改 test config 让功能真能开出来**，不是改 gate。

### 既有用例全绿 ≠ 修复验证过 —— 修 bug 必用报告者的确切输入复现

跑既有 `TestAcc*` 全 PASS 只证明**没回归**，**不证明你的修复命中了报告的失败路径**。既有用例的 config 常常**恰好绕开触发条件**——bug 才一直没被它们暴露。所以修 bug 收口前，**必须新增一条用报告者确切输入的 repro 用例**（客户/工单里的那段 HCL、那个属性值），先让它在**修复前 FAIL、修复后 PASS**，才算验证闭环。repro 用例是临时验证件，不进 PR（跑完删）；namespace/实例名之类可随机化避免撞名，但**触发字段（key/值/组合）必须原样照抄**。

#### 实例：esa_kv 创建失败（工单 84090290）

- 工单描述判定根因是「EdgeKV 最终一致性 → 写后读 404」，建议「加重试」；据此改的第一版把 `NotFound` 当可重试。
- 既有 `TestAccAliCloudESAKvKv_test1/test2` 全 PASS（且 tf-debug.log 里 `InvalidKey.NotFound` **0 次**）——看着像修好了，其实**重试分支一次都没被执行**，只证明了 happy path。
- 用工单里客户**确切 HCL** 包成 repro 用例（key = `test:resource-managed`，**含冒号**原样保留）真跑：**FAIL 302.54s**（空转到 Create 超时），debug trace 显示反复查一个**错的 key**。
- 真根因是 id = `namespace:key` 被 `strings.Split(id, ":")` 按每个冒号拆，含冒号的 key 被拆错 → 查错 key → **确定性** 404，与时序无关。正解是 `strings.SplitN(id, ":", 2)`。修复后同一 repro **PASS 4.02s**，test1/test2 无回归（`PASS=3 FAIL=0 SKIP=0`）。
- **教训**：既有用例 key 不含冒号，永远触不到 bug；不跑客户确切输入，就会把「最终一致性」这个误诊一路带到 PR，还把「秒失败」改成「5 分钟空转」的负优化。与「SKIP≠PASS」同族——**acc 绿有多种假绿，真信心只来自读 run.log 里真实 case 结果 + 确认触发条件真被跑到**。

## 6. FAIL 定性:四分类(mandatory triage)

`run.log` 出现 `--- FAIL: TestAccXxx` 后,不要只把日志摘要甩给用户——**必须**把这次失败归到下面四类之一,证据不足就继续挖,不要用「大概」「疑似」下结论。

- **A. 后端云产品 API 问题**:API 行为/文档/合同不合理或不正确,不符合 Terraform 使用要求。
- **B. 测试用例问题**:代码本身正确,测试 HCL 或前置条件写法不合理(如写死不存在的依赖资源、precheck 门闩缺配套等)。
- **C. Resource / DataSource 代码 或 CloudSpec 资源定义 bug**:provider 侧代码,或其上游 cspec 定义,有真实缺陷(C1=手写代码 bug,C2=cspec 定义 bug)。
- **D. 生成器 bug**:cspec 定义正确、API 契约正确,但生成器把「对的 cspec」翻成了「错的 Go 代码」——问题在生成器本身。仅适用于**生成器产出**的资源。

**判定路径**(先做 6.1 取证,再按下图收敛):

```
FAIL
 ├─ API 契约有问题?           → A
 ├─ 用例 HCL / precheck 有问题? → B
 └─ 都不是 → 看资源是手写还是生成:
      ├─ 手写(alicloud/resource_*.go 无生成器标记)
      │    └─ 代码与 API 不对齐?     → C1
      └─ 生成器产出(有 `// Code generated` 类标记或对得上生成器 golden)
           ├─ cspec 定义与 API 不一致? → C2(修 cspec)
           └─ cspec 定义正确、生成产物错? → D(修生成器)
```

### 6.1 通用取证(判定前先做)

按顺序做,每一步都要指到具体行号 / API 名 / RequestId,不留模糊描述:

1. `grep -B2 -A30 "^--- FAIL:" run.log` → 记下:失败函数、失败的 Terraform 生命周期步骤(Plan/Apply/Refresh/Destroy)、错误短语。
2. 在 `tf-debug.log` 里定位对应的 API 调用:

   ```bash
   grep -B2 -A5 "<ErrorCode 或错误短语>" tf-debug.log | head -60
   ```

   拿到:API name / 请求 body / RequestId / 返回 payload。
3. 与 API 契约对齐:调 `amp-resource-metadata` skill 或 OpenAPI Explorer 拉该 API 的官方 spec(参数必填性、返回结构、错误码含义、幂等/重试语义)。
4. 与 provider 代码对齐:

   ```bash
   grep -n "<关键字段名或 API 名>" alicloud/resource_alicloud_<name>.go
   ```

5. 若怀疑上游 cspec 定义:去 `amp-resource-metadata` 或 `cloudspec-model` 仓拉当前 cspec resource definition 比对。

### 6.2 A - 后端云产品 API 问题

**判定信号(任一命中)**:

- API 文档说一种行为,实际返回另一种(如文档说 List 返回数组,实际返回单对象)
- 必填字段在 API 侧无稳定支持,Terraform 侧无法回填
- 返回结构缺少 Terraform Read 需要的字段,且 API 无替代读取通道
- 错误码语义混乱(如实际不存在返回 InternalError 而非 NotFound),CRUD 无法幂等
- 分页 / 最终一致性 / 幂等语义与 Terraform 生命周期契约不兼容

**报告格式**:

```
定性:A - 后端 API 问题
API:<Product>::<Action>  version <YYYY-MM-DD>
RequestId:<xxx>
契约 vs 实际差异:<一句话>
证据:tf-debug.log line NN,请求 <params>,期望 <expected>,实际 <actual>
建议:向 API owner 报 bug / 走 CloudSpec 加白 / 推动上游修合同(必要时按 aone-triage 建关联单)
```

### 6.3 B - 测试用例问题

**判定信号**:

- 用例 HCL 引用了当前 region/账号下不存在的固定 ID(如写死 `vsw-xxx` / OSS bucket / VPC / RAM role)
- 硬编码了区域或账号敏感值(别人的 UID / ARN)
- 前置条件(precheck)与用例真实需要的能力不匹配(照 §5 SKIP→PASS 里兄弟资源模式对照)
- 引用了已删除 / 重命名的 provider 字段(`Unsupported argument "X"`)
- 用例互相污染 / 未清理导致 `QuotaExceeded`(见「常见失败」表 sweep 处置)

**报告格式**:

```
定性:B - 测试用例问题
文件:alicloud/<file>_test.go
函数:TestAccAliCloudXxx
证据行:line NN,<HCL 片段>
问题:<写死 vsw-xxx / 缺 precheck / 引用已删字段>
修复:<改成 data source 动态取 / 加 testAccPreCheckWithRegions / 换替代字段>
```

### 6.4 C - Resource/DataSource 代码 或 CloudSpec 资源定义 bug

若不属于 A / B,就是 C。C 再细分两小类,**报告必须具体到「行」或「定义字段」**——用户拿到就能改。

#### C1. Provider 代码 bug

**判定信号**:

- API 契约正确、用例合理,但 provider Create/Read/Update/Delete 与 API 不对齐(参数拼装错、返回解析错、状态判定错、重试策略缺失)
- Refresh 后有 spurious diff → Read 未把 API 返回字段完整映射回 schema
- import/id 拆分错误(参照 §5 esa_kv `strings.Split` vs `SplitN` 实例)
- 错误码分类错(该 NotFound 被当作 Error retry,或反之)

**报告格式**:

```
定性:C1 - Resource 代码 bug
文件:alicloud/resource_alicloud_<name>.go   (或 data_source_alicloud_<name>.go)
函数:<resourceAlicloudXxxRead / Create / Update / Delete>
出错行:line NN(具体到语句)
根因:<一句话>—例如 strings.Split(id, ":") 对含冒号的 key 拆错
修复:<改成 strings.SplitN(id, ":", 2)>
验证:新增 repro 用例(报告者原样输入),必须先 FAIL 后 PASS(照 §5「既有用例全绿 ≠ 修复验证过」)
```

#### C2. CloudSpec 资源定义 bug

**判定信号**:

- `attributeMappings` 的 `responsePath` / `requestPath` 错
- `identifyDefinition.uniqueKeyFields` 缺主键
- `resourceTypeOperationMapping` 少 `list` / `gets` / `deletes` / `updates`
- 属性 `readOnly` / `required` / `enum` / `pattern` 与 API 实际不一致
- `rootMapping.responsePath` 与 API 返回结构不匹配

**报告格式**:

```
定性:C2 - CloudSpec 资源定义 bug
资源:<Namespace>::<ResourceName>
cspec 位置:cloudspec-model 分支 <feature/xxx>,文件 resources/<Name>.cspec
定义错误:<例如 attributeMappings 中 $.FileSystemId.responsePath 错;或 uniqueKeyFields 缺 $.FileSystemId>
修复路径:去 cspec 分支改定义 → `aliyun cspec build` → `aliyun cspec check --name <Res>` → `amp publish pre` → 生成器重跑
关联 skill:`cloudspec-resource-edit` / `cloudspec-operation-edit` / `cloudspec-norm-check-fix`
回写与修复:返回 RD，按 aone-triage 分支 E 和
`terraform-provider-release/references/cloudspec-pre-resource-loop.md` 在**原主单**内修复
CloudSpec feature 分支，完成 build/check、pre dry-run/发布、从 pre 重新生成与本 AccTest 复验。
`requested_external_actions: []`，不得为资源定义或文档源头问题另建 Aone；能力失败返回
`missing_capability` / `blocked`。pre 成功后 `release/idle`，不得 finish；prod/online 与
master/main merge/push 仍是人工硬门。
```

### 6.5 D - 生成器 bug

当 A/B/C 三层全部核对通过——**API 契约正确、用例合理、上游 cspec 资源定义正确**——但 provider 里的 resource/datasource 代码仍然错,而这个资源又是**通过生成器**(acube 生成链路 / codegen 工具)自动产出的,那问题就在**生成器**本身。

**判定信号(合取,四条都要成立)**:

- 已按 6.2 排除 A(API 契约核对通过)
- 已按 6.3 排除 B(用例合理)
- 已按 6.4 C2 排除 cspec 定义问题(cspec 与 API 一致、mapping 正确、`aliyun cspec check` 通过)
- provider 里的 resource/datasource 代码是**生成器产出**——判定方式:文件头部一般带 `// Code generated ...` / `// DO NOT EDIT` 类标记;或对照生成器 golden template 能对上。不确定时查 `provider-resource-dev` skill 的「生成 vs 手写判定」小节。

**报告格式**:

```
定性:D - 生成器 bug
资源:alicloud_<name>(生成器产出)
生成器仓库/工具:<例如 acube 生成链路 / codegen 工具名>
生成器代码位置:<generator repo>/<path>/<file>.go line NN(具体到规则/模板)
生成缺陷:<一句话>—例如生成的 Read 把 List 型返回按单对象拆解;或复合字段展开丢了嵌套层
证据链:
  - cspec 定义正确(引用 attributeMappings 片段)
  - API 契约与 cspec 一致(RequestId <xxx>)
  - 生成产物 alicloud/resource_alicloud_<name>.go line NN 与 cspec/API 不一致
修复路径:review 并修正生成器代码/模板 → 生成器重新构建 → **重新生成该资源**(不要只手工 patch provider 产物,否则下次生成会覆盖)→ 复跑本 ACC 用例验证
关联 skill:`provider-resource-dev`(生成链路诊断)、`amp-resource-metadata`(cspec/API 契约核对)
```

**为什么不能只手工 patch 生成产物**:生成器下次跑会覆盖你的 patch,bug 复发;而且同一生成器产出的其他资源大概率带同类 bug——**修生成器 = 一次修一类**。

### 6.6 定性纪律

- **必须报主因**:一个失败可能横跨多类(如 API 不合理 + provider workaround 引入 bug),报告时**先列主因,再列并发问题**,不要「都可能」糊在一起。
- **证据不足禁止臆测**:若无法唯一定性,先补证据(读更多 debug、比对文档、比对兄弟资源、拉 cspec 与生成器模板对照),证据仍不足则明写「证据只能到 X 或 Y 二选一,还需要 Z」;不允许在没有 API 契约核对的情况下就断言「API 问题」,也不允许在没有确认「资源是生成器产出且 cspec 正确」的情况下就断言「生成器 bug」。
- **C2 vs D 的边界不能糊**:C2 修的是 cspec 定义,D 修的是生成器代码/模板。二者的证据要点完全不同——C2 要指出 cspec 里哪个字段错,D 要指出生成器规则/模板哪里错。混淆会让修复方向跑偏。
- **基础设施异常不算四分类**:FC OOM / 超时 / 网络 / 配额,归 §5 状态解读 和「常见失败」表处置,不进 A/B/C/D 分诊。
- 详细 pattern 匹配材料见 `references/error-patterns.md`,各条已标注归属分类。

## 7. 没有跑到用例

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
| 0 tests run | Name mismatch (see §7 没有跑到用例) | Fix test function names or pass exact `--test-case` |
| FC 解压后找不到 `terraform-provider-alicloud` | 本地打包根目录错误(zip root 未固定为 `terraform-provider-alicloud`) | 使用本 skill 的 `acctest.py`; zip root 必须固定为 `terraform-provider-alicloud` |
| `指定的测试用例 X 在 alicloud/ 中未找到` | `--test-case` value misspelled | grep local test file for actual names |
| 需要跑多个具体用例 | 把 `A,B` 当作单个函数名传入 | 用逗号分隔多个用例；脚本经 `testCaseNames` 重复参数提交，服务端锚定各自精确函数名;或不传 `--test-case` 跑该资源全部用例 |
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
