# Jarvis 决策权策略（Autonomy Policy）

> Aone 为研发事项唯一真源，进展必同步（bootstrap/wrap.sh）。

## 运行模式

### supervised（默认）
- 每轮先输出行动计划，**逐项等待用户授权**后才执行 Aone 写操作。
- 适合日常使用。切换 unattended 分两层：会话层=仓库主人显式指令；脚本链路（bridge/serve 的 `plan.sh` 机械闸门）读本文末机读块的 `mode` 字段，需一并改 JSON 才放行。

### unattended
- 在置信度高且操作可逆的前提下全自动执行；**仅 escalate 触发时才通知人**。
- 需显式开启，不得默认激活。

### headless（Tata 委派 / bridge dispatch）
- Jarvis 由 bridge 后台 spawn，无终端交互，始终持有一个 Aone 工单。
- **自主权高**：自动执行项（auto 列表）全部免授权直接执行。
- **身份约束**：编排层默认 jarvis 身份；数字人子代理按职责用角色身份（terraform-pd/terraform-rd/terraform-qa，未登录回退 jarvis）；需使用 chenyi/guozai/linjun 等个人身份时，必须在 Aone 工单评论中 @对应人并获得明确授权回复后方可使用。
- **遇阻挂起**：遇到必须人类确认/决策的点时，在 Aone 工单评论中 @对应人，输出 `[[SUSPEND:...]]` 哨兵信号后退出进程。bridge 的 WaitWatcher 轮询评论，检测到回复后用 `--resume` 唤醒 Jarvis 继续。
- **外化契约（多机安全）**：SUSPEND 挂起或 release 释放**之前**必须先把上下文与代码外化到远端，否则换一台机器无法续跑——依次执行 `wrap.sh sync`（进展/上下文入 Aone 评论）+ `github-identity.sh push`（代码入远端分支）+ `coord.sh checkpoint <aid> <stage> <wt> <branch> <repo> <pushed_branch>`（把已 push 的远端分支写进 checkpoint）。缺任一即视为 `unexternalized`，`JARVIS_REQUIRE_PUSH=1` 时 wrap-check 会阻断收尾。
- **超时**：单轮执行上限 12 小时（`JARVIS_DISPATCH_TIMEOUT`）；挂起等待上限 14 天。

---

## 自动执行项（授权后 / unattended 模式下高置信+可逆）

| 操作 | 说明 |
|------|------|
| `reply` | 回复工单 / 工作项评论 |
| `create_req` | 建需求（按 `config/pools.json` 路由落池） |
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
| `unexternalized` | headless 收尾前外化契约未完成（wrap sync + push + checkpoint，见 headless 节；`JARVIS_REQUIRE_PUSH=1` 时 wrap-check 阻断） |

Escalate 行为：暂停执行，输出摘要，通知用户决策。

---

## 机读策略块

```json
{"mode":"supervised","modes":{"supervised":"逐项授权","unattended":"高置信自动","headless":"bridge委派,高自主+挂起唤醒"},"auto":["reply","create_req","tag","create_cr","worktree","prestage","adhoc_aone","pr_review","wrap_sync","wrap_done"],"stop":["release_prod"],"escalate_if":["low_conf","verify_fail","redline","missing_capability","unexternalized"],"headless":{"dispatch_timeout":43200,"suspend_expire":1209600,"suspend_signal":"[[SUSPEND:{...}]]"}}
```
