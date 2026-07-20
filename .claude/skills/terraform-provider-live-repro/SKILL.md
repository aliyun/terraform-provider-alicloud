---
name: terraform-provider-live-repro
description: 使用最小化本地 HCL、Provider 调试日志、OpenAPI RequestId、刷新漂移证据、受控的现场保留或清理，以及可选的 AutomationAgent HTML 报告，在真实阿里云资源上复现 Terraform Provider 问题。适用于用户要求实际复现 Provider 缺陷、根据给定 Terraform 创建资源、排查创建与读取无法回环或意外替换、为产品团队保留云上现场，或生成完整 API 时间线的场景；普通验收测试优先使用 invoke-terraform-acc-test-remote，只有问题依赖用户原始 HCL 和真实 API 响应时才使用本技能。
---

# Terraform Provider 真实资源现场复现

把 Provider 问题转化为可重复、可审计的真实资源复现。精确保留用户输入，分离云产品 API 行为与 Provider 状态行为，并最终留下已验证的云上现场或已验证的干净账号。

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

默认在 `~/Workspace/troubleshoot/<slug>-<aone-id-or-date>` 下创建隔离目录。不得覆盖其他复现目录或复用无关 Terraform state。

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

使用内置解析器，只输出请求和响应的白名单字段：

```bash
skill_dir="<当前已加载 SKILL.md 所在目录的绝对路径>"
python3 "$skill_dir/scripts/extract-api-timeline.py" \
  "$evidence_dir/terraform-apply.raw.log" \
  --format markdown > "$evidence_dir/apply-api-timeline.md"
```

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

生成不含 base64 图片的自包含 HTML：

```bash
skill_dir="<当前已加载 SKILL.md 所在目录的绝对路径>"
python3 "$skill_dir/scripts/render-report-html.py" report.md report.html
```

用户要求在线报告时，调用 `html-report-preview`，并使用已验证的 Aone ID 执行仓库脚本：

```bash
bash bootstrap/html-report-preview.sh upload <aone-id> report.html \
  --base-url https://agent.aliyun-inc.com --format jsonl
```

服务端 Token 缺失时，禁止使用浏览器 Cookie 或个人 BUC 会话。无凭证访问返回的 `/reports/aone/.../view` URL，验证 HTTP 200、`text/html`、报告标题、至少一个实例 ID、一个 API 名和一个关键 RequestId。

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
- 要求在线报告时，预览链接返回 200 且包含可辨识的证据标记。
