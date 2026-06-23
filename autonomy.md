# Jarvis 决策权策略（Autonomy Policy）

> Aone 为研发事项唯一真源，进展必同步（bootstrap/wrap.sh）。

## 运行模式

### supervised（默认）
- 每轮先输出行动计划，**逐项等待用户授权**后才执行 Aone 写操作。
- 适合日常使用；切换方式：显式指令 `--mode unattended`。

### unattended
- 在置信度高且操作可逆的前提下全自动执行；**仅 escalate 触发时才通知人**。
- 需显式开启，不得默认激活。

---

## 自动执行项（授权后 / unattended 模式下高置信+可逆）

| 操作 | 说明 |
|------|------|
| `reply` | 回复工单 / 工作项评论 |
| `create_req` | 建需求（Cloudspec 需求池） |
| `tag` | 打标签 / 更新标签 |
| `create_cr` | 建变更 / CR |
| `worktree` | worktree 开发（本地分支） |
| `prestage` | 预发部署 |
| `adhoc_aone` | ad-hoc 建/补单（loops/adhoc-intake.md，PR 默认落 tf_provider） |
| `pr_review` | 只读 PR 评审（不写不合并） |
| `wrap_sync` | 中途回填 Aone 进展评论（wrap.sh sync，不改状态） |
| `wrap_done` | 收尾回填 Aone：评论+run_done+可选改状态（wrap.sh done） |

---

## 永久停止项（Always Stop）

| 操作 | 说明 |
|------|------|
| `release_prod` | **正式发布**——无论任何模式，必须人工确认后才能执行 |

---

## 置信度判定

- **高置信（high_conf）**：OpenAPI 查证 + provider 源码两层一致，规则命中明确。
- **低置信（low_conf）**，触发 escalate：
  - OpenAPI 与源码结果冲突
  - 缺少源码映射（missing source）
  - 路由规则未命中（routing miss）

---

## Escalate 触发条件

| 触发器 | 说明 |
|--------|------|
| `low_conf` | 置信度低，无法自动决策 |
| `verify_fail` | 验证步骤失败（查证返回矛盾结果） |
| `redline` | 红线操作：推送 master / 零差异 CR / 正式发布 |
| `missing_capability` | 缺工作区/工具/池映射（config/workspaces.json 未登记） |

Escalate 行为：暂停执行，输出摘要，通知用户决策。

---

## 机读策略块

```json
{"mode":"supervised","auto":["reply","create_req","tag","create_cr","worktree","prestage","adhoc_aone","pr_review","wrap"],"stop":["release_prod"],"escalate_if":["low_conf","verify_fail","redline","missing_capability"]}
```
