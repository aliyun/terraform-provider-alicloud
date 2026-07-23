---
name: reproduce-terraform-provider-issue
description: 分析客户 Terraform Provider 问题并生成可复现、可审计的中英文报告交付包；使用最小化本地 HCL、Provider 调试日志、OpenAPI 请求标识、刷新漂移证据、受控的现场保留或清理，以及 AutomationAgent 静态 HTML 在线预览。适用于用户要求分析客户问题、实际复现 Provider 缺陷、根据给定 Terraform 创建资源、排查创建与读取无法回环或意外替换、保留云上现场、生成完整 API 时间线或生成复现报告的场景；普通验收测试优先使用 invoke-terraform-acc-test-remote，只有问题依赖用户原始 HCL 和真实 API 响应时才使用本技能。
---

# Terraform Provider 用户问题现场排查复现

把 Provider 问题转化为可重复、可审计的真实资源复现。精确保留用户输入，分离云产品 API 行为与 Provider 状态行为，并最终留下已验证的云上现场或已验证的干净账号。

中文会话默认输出中文报告；用户明确指定其他语言时，报告正文和 HTML `lang` 使用该语言。

## 路由边界

- 直接 Terraform 模板、创建/读取回环、刷新漂移、意外 ForceNew 替换、用户原始输入、RequestId 时间线或现场交接，使用本技能。
- 常规 `TestAcc*` 执行、PR 回归覆盖或大范围远程验证，使用 `invoke-terraform-acc-test-remote`。
- 现有 AccTest 即使通过，只要没有覆盖用户触发问题的精确字段，就只能作为回归证据，不能作为复现结论。
- 存在 Aone ID 时，先加载 `aone-triage` 并核对标题与描述。上传报告不代表获得评论、改状态或关闭工单的授权。

## 双宿主兼容

保持流程和内置脚本与宿主无关：

- 从当前宿主的项目技能根目录加载 `SKILL.md`。
- `agents/openai.yaml` 仅用于可选的界面展示，运行逻辑不得依赖它。
- 将 `skill_dir` 解析为当前已加载 `SKILL.md` 所在目录的绝对路径，不得硬编码另一个宿主的技能根目录。
- 两个宿主统一使用相同的 `scripts/`、`references/` 和 `assets/` 相对路径。

## 强制安全规则

1. 仅当用户明确要求创建、apply 或复现，或者既有授权流程明确包含真实资源创建时，才创建云资源。
2. apply 前记录现场处置方式：`preserve` 或 `destroy-after-evidence`。仅在用户要求或下游排查确实需要时保留现场。
3. 禁止执行漂移或替换计划。出现替换的 plan 是证据，不是执行产物。
4. 禁止打印、持久化、提交或上传 AK/SK、Bearer Token、签名、密码或 Provider 配置值。
5. 将 `.tfplan` 文件和原始 `TF_LOG` 日志视为敏感数据。先提取白名单时间线，再删除原文件。
6. 禁止销毁已声明保留的现场。对于临时现场，只销毁该复现 state 中的精确资源，并验证资源已删除。
7. 不得静默替换用户提供但当前账号无权访问的 VPC/VSwitch ID。要么停止，要么在获得认可后使用隔离的等价网络，并在报告醒目位置说明偏差。

## 执行流程

### 1. 建立复现用例

记录：

- 用户 HCL 和精确 Provider 版本；
- 预期行为与实际症状；
- Region、Zone、资源类型和生命周期阶段；
- 可用时记录 Aone ID；
- 用户指定的工作目录和现场处置方式。

默认在 `~/workspace/troubleshoot/<slug>-<aone-id-or-date>` 下创建隔离目录。不得覆盖其他复现目录或复用无关 Terraform state。

### 2. 前置检查身份与依赖

在不暴露凭证的前提下检查：

```bash
terraform version
aliyun sts GetCallerIdentity
```

plan 前使用只读 API 验证目标 Region/Zone 以及引用的 VPC/VSwitch。报告只保留内部排查必要的身份标签或 ID，禁止展示凭证值。

精确锁定用户报告的 Provider 版本。通过环境变量或已有 CLI Profile 传递凭证，不得把明文密钥写入 `.tfvars`。

### 3. 构造最小复现

保留全部触发字段及其精确值和依赖关系，只移除无关资源。若原网络不可访问且已获准使用隔离等价网络，保持 Region、Zone、资源形态和依赖图一致。

执行：

```bash
terraform fmt -recursive
terraform init -backend=false -input=false -no-color
terraform validate -no-color
```

### 4. 审核首次 plan

生成保存的创建计划并检查 JSON：

```bash
terraform plan -input=false -no-color -out=create.tfplan
terraform show -json create.tfplan | jq '{changes: [.resource_changes[] | {address, actions: .change.actions}]}'
```

仅当计划中只有预期创建动作，且没有更新、删除或替换时才继续。否则停止并排查配置、身份、Region 或陈旧 state。

### 5. 携带追踪证据执行 apply

记录本地时区和 UTC 的开始、结束时间。只对已审核的 apply 开启 Provider 调试日志：

```bash
DEBUG=terraform \
TF_LOG=DEBUG \
TF_LOG_PROVIDER=DEBUG \
TF_LOG_PATH="$evidence_dir/terraform-apply.raw.log" \
terraform apply -input=false -auto-approve -no-color create.tfplan
```

立即记录 Terraform 资源 ID。禁止发布原始日志。

### 6. 获取 API 真实结果

使用已创建的资源 ID 调用云产品只读 API。每次调用记录：

- 时间戳和 API 名；
- RequestId；
- 实例或资源 ID；
- 相关状态和字段；
- 能证明触发字段已下发的创建请求字段。

对于状态无法回环的问题，至少比较：

1. 创建请求参数和创建响应；
2. 主查询或读取响应；
3. 可能保存缺失关系的关联资源 API；
4. Provider Read 后的 Terraform state。

对于 NAS VPC/VSwitch 问题，比较 `CreateFileSystem`、`DescribeFileSystems` 和 `DescribeMountTargets`，并明确区分 `VpcId=""`、`QuorumVswId=null` 和响应中完全缺失 `QuorumVswId`。

### 7. 提取脱敏 apply 时间线

使用内置解析器，按当前资源的 API 契约显式传入白名单字段。解析器不内置任何产品或资源字段；未传参时只输出时间、API、RequestId 和通用状态：

```bash
skill_dir="<当前已加载 SKILL.md 所在目录的绝对路径>"
python3 "$skill_dir/scripts/extract-api-timeline.py" \
  "$evidence_dir/terraform-apply.raw.log" \
  --request-field <trigger-field> \
  --target-field <resource-id-field> \
  --observe-field <read-api>:<disputed-response-field> \
  --format markdown > "$evidence_dir/apply-api-timeline.md"
```

每个参数可重复传入。`--observe-field` 优先使用 `API:字段` 限定只在该 API 响应中判断字段存在性，避免把其他 API 本就不返回的字段误报为缺失。只选择本次诊断所需字段；解析器会拒绝密钥、Token、密码、签名、Cookie、UserData 和私钥类字段。

分享证据或编写最终报告前，阅读 [references/evidence-contract.md](references/evidence-contract.md)。

### 8. 证明刷新漂移

使用独立原始日志执行正常 refresh/plan：

```bash
set +e
DEBUG=terraform \
TF_LOG=DEBUG \
TF_LOG_PROVIDER=DEBUG \
TF_LOG_PATH="$evidence_dir/terraform-drift-plan.raw.log" \
terraform plan -input=false -no-color -detailed-exitcode -out=drift.tfplan
plan_rc=$?
set -e
```

将退出码 `0` 解释为空计划，`2` 解释为存在变更，其他值解释为执行失败。提取替换路径：

```bash
terraform show -json drift.tfplan | jq '{changes: [
  .resource_changes[]
  | select(.change.actions == ["delete", "create"])
  | {address, replace_paths: .change.replace_paths}
]}'
```

使用同一解析器提取 refresh API 时间线。禁止执行 `terraform apply drift.tfplan`。

### 9. 保留或清理现场

选择 `preserve` 时：

- 保留云资源和 Terraform state；
- 验证 state 资源数量和云资源实时状态；
- 列出 ID、Region/Zone、挂载地址和明确的现场保留警告；
- 脱敏证据完成后删除敏感 plan 文件和原始日志；
- 提供只读验证命令和醒目的“禁止销毁”说明。

选择 `destroy-after-evidence` 时：

- 审核 destroy 计划只包含本次复现资源；
- 使用同一凭证来源执行 `terraform destroy`；
- 验证 state 资源数量为零，云产品查询返回 NotFound 或零条结果；
- 明确列出已删除资源，并说明资源无法恢复。

### 10. 生成并发布报告

复制 [assets/report-template.md](assets/report-template.md)，填写所有适用章节。分层输出结论：

1. 云产品 API 行为；
2. Provider Schema/Read 行为；
3. Terraform plan 后果；
4. 安全修复方向和需要产品团队确认的问题。

默认交付目录必须包含：

```text
REPORT.md
REPORT.html
template/main.tf
template/README.md
template/.terraform.lock.hcl  # init 生成时保留
```

以 `template/main.tf` 为唯一 HCL 源。把它的完整 UTF-8 字节（包括末尾换行）放入
`REPORT.md` 唯一的 `hcl` fence，再由 renderer 生成 `REPORT.html`；不得分别手写三个副本。
模板中的 `profile` 必须为 `string`、`default = null`、`nullable = true`，本地只用
`TF_VAR_profile` 或已有 CLI Profile 传凭据，禁止把凭据放入模板、报告或 URL。

在线预填链接固定使用：

```text
https://api.aliyun.com/terraform?spm=XToCode.TerraformAI.QA.0&activeTab=code&source=PlayGround&sourcePath=TerraformAI/<13位时间戳>::<同一时间戳>&params=<编码后的main.tf>
```

`params` 必须等同 Java `URLEncoder.encode(hcl, UTF_8).replace("+", "%20")` 且只编码一次。
Markdown、HTML 和 URL 解码后的 HCL 必须与 `template/main.tf` 字节一致。链接只预填代码，
不代表执行 plan/apply。

生成零脚本、无 base64/data URI 的静态 HTML：

```bash
skill_dir="<当前已加载 SKILL.md 所在目录的绝对路径>"
python3 "$skill_dir/scripts/render-report-html.py" \
  --lang zh-CN REPORT.md REPORT.html
```

不得加入 `<script>`、`<button>`、`on*=`, `javascript:` 或任何复制按钮/脚本。即使无网络行为的
inline script 也可能被 pre-agent WAF 拦截。Viewer 的复制能力归 AutomationAgent 平台；
当前报告侧状态固定为 `platform_blocked`，不得通过上传 HTML 绕过。

先运行统一 validator。它只执行 fmt/init/validate，绝不执行 plan/apply：

```bash
python3 "$skill_dir/scripts/validate-report-package.py" <package-dir> \
  --format json
```

退出码 `0` 表示包已验证，`2` 表示确定性格式/安全/配置错误，`3` 表示 init、预览或平台能力
真实阻塞。不得把 `3` 改写成成功或伪造验证记录。

validator 会扫描包括 `.txt` 在内的包内文本，拒绝原始 TF_LOG、真实凭据赋值和 HTML entity
解码后形成的可执行协议/事件属性；只提及 `accessKeyId` 等字段名不算凭据。Terraform
fmt/init/validate 使用隔离 HOME、TF_DATA_DIR 和白名单环境，不继承 `TF_CLI_ARGS`、
`TF_VAR_*` 或云凭据。

用户要求在线报告时，先调用 `html-report-preview`，使用已验证的 Aone ID 执行仓库 helper，
并保存其单行 JSON：

```bash
bash bootstrap/html-report-preview.sh upload <aone-id> REPORT.html \
  --format jsonl > preview.json
```

helper 默认上传到 `https://pre-agent.aliyun-inc.com`。缺少 `JARVIS_HTML_REPORT_TOKEN` 时，它会
在任何 curl/Aone 调用前以 `missing_token`、退出码 `3` 阻塞；禁止改用浏览器 Cookie 或个人
BUC 会话。

使用同一个 validator 验证上传记录和匿名 GET，不发送 Authorization/Cookie：

```bash
python3 "$skill_dir/scripts/validate-report-package.py" <package-dir> \
  --require-preview --preview-json preview.json \
  --preview-marker '<报告专属 marker>' --format json
```

匿名 GET 必须为 HTTP 200、`text/html`，标题、HCL 和用户指定 marker 均一致；`viewUrl` 必须是
`/reports/aone/.../view` 只读路由，绝对 `url` 必须使用预期 origin 且原始 path 与 `viewUrl`
完全相同，不得带 query、fragment、百分号编码或发生重定向。`--require-preview` 必须至少传入
一个非空 marker。helper 使用显式 `--base-url` 时，validator 同步传
`--preview-origin <同一 origin>`；HTTP origin 仅允许 loopback 前向测试。只有用户把 Viewer
复制能力列为硬验收项时才加
`--require-viewer-copy`；在平台修复前它会如实退出 `3`。

### 11. 删除敏感产物

脱敏报告完成后，只定位并删除以下精确生成文件：

- `create.tfplan`；
- `drift.tfplan`；
- apply/refresh 原始调试日志。

保留现场时，保留 `.tf` 文件、lockfile、脱敏证据、报告和 state。交付前扫描可分享产物，确认不包含真实凭证值。

## 完成门槛

所有适用条件满足前，不得宣称完成：

- 首次 plan 已审核且 apply 成功；
- 已记录创建资源 ID 和 API RequestId；
- 已比较直接 API 响应与 Terraform state；
- 已记录漂移 plan 退出码和替换路径；
- 已验证现场保留或资源清理结果；
- 已删除原始敏感产物；
- 报告包含输入偏差和精确时间戳；
- 默认四文件交付完整，必要 lockfile 已保留；
- `main.tf`、MD/HTML HCL 和在线 `params` 字节一致，`profile` 契约通过；
- `terraform fmt -check -recursive`、隔离 `init -backend=false` 和 `validate` 真实通过；
- 统一 validator 退出 `0`，报告静态 HTML 不含脚本、事件属性、base64/data URI 或凭据；
- 要求在线报告时，helper 记录为 `success=true,status=uploaded`，匿名预览返回 200、
  `text/html` 且包含标题、完整 HCL 和可辨识 marker；
- Viewer 复制未修复时明确记录 `platform_blocked`，不得宣称完成该平台能力。

## 技能维护与前向测试

- 修改后同步 `agents/openai.yaml`，确保 `default_prompt` 显式包含
  `$reproduce-terraform-provider-issue`，并对 `.agents`、`.claude` 两端运行官方
  `quick_validate.py`。
- 用 fake terraform、fake curl 和本地匿名 HTTP server 前向测试 validator/helper；必须覆盖
  字节不一致、双重编码、JSON 凭据、`.txt` 原始 TF_LOG、HTML entity、环境净化、重定向、
  空 marker、错误 view URL 绑定、服务端反射脱敏和缺 Token 零网络调用。
- 前向测试不得执行 terraform plan/apply、创建云资源、真实上传或操作 Aone。
