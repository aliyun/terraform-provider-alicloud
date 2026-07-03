# ticket-template —— probe 工单/草稿骨架

probe 发现的 provider 问题按此骨架落 `escalation/probe-drafts/<日期>-<资源或场景>-<code>.md`（draft 模式），
或未来毕业后按此建 Aone 需求单。

## 硬规则（写之前先记）

- **禁贴 AK/SK**：任何凭证、SecretKey、token 绝不出现在标题/正文/日志摘录里。日志摘录只留字段名、错误码、RequestId。
- **禁 AI 署名水印**：不写 `Co-Authored-By: Codex` / `🤖 Generated with Codex`；溯源只写「来源：jarvis tf-customer-probe」。
- **draft 文件头**加 `status: pending-review`（未审核信号；未跟踪文件在 `git status` 天然充当待审队列）。

## 标题

```
[probe][<resource>] <一句话症状>
```
例：`[probe][alicloud_oss_bucket] 开启 versioning 后 lifecycle expiration 每次 plan 永久 diff`

## 正文骨架

```markdown
---
status: pending-review
---

# [probe][<resource>] <一句话症状>

## 环境
- provider: aliyun/alicloud 1.284.0
- terraform: <verdict.terraform_version>
- region: <verdict.region>

## 最小复现 .tf
（贴 probes/scenarios/<id>/main.tf 的相关最小片段；含 pin 版本块）

## 复现步骤
1. terraform init / validate / plan / apply
2. <触发问题的具体步骤，如 apply 后再 plan / 改 tags 再 apply / state rm 后 import>

## 期望行为
（引官方文档 + OpenAPI 文档链接说明“本应如何”）
- 文档：<website/docs 链接>
- API：<next.api.aliyun.com 对应 OpenAPI 链接，如涉及>

## 实际行为
（日志摘录：错误码 / 关键 diff 字段 / RequestId；**不贴 AK/SK**）

## 危害评估与建议优先级
- severity: S<n>（依据 severity-rubric.md，注明升降级理由）
- 建议 Aone 优先级：紧急/高/中/低

## 溯源
- tier-1:场景 probes/scenarios/<id>/ + verdict runs/probe/<日期>-<id>.json
- tier-0:资源 <alicloud_xxx> + verdict runs/probe/<日期>-tier0.json(含 doc/source 位置、judgment_queue)
- 来源：jarvis tf-customer-probe
```

## 建单参数块（毕业后 mode=file 用；draft 阶段只记录不执行）

| 字段 | 值（来自 config/probe.json） |
|------|------|
| project | 528766（tf_provider 池，terraform-alicloud 内部研发） |
| category | req（需求） |
| assignee | WORKER_1782379562571 |
| tag | jarvis-probe |
| priority | 按 severity-rubric 映射（枚举值首次建单时用 a1 查证项目字段固化） |

> 客户来源的问题若需与客户主单双向关联，按 loops/aone-triage.md §4.2 的 tf_provider 池纪律处理。
> daily 新建上限 `config.limits.daily_new_tickets`（默认 3），超限的顺延下一轮。
