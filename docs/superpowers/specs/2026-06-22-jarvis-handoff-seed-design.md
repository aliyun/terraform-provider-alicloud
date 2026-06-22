# jarvis：可被 Claude 完全接替的工作起点

> 状态：设计 v0 · 2026-06-22 · 待用户评审

## 一句话目标

把"我"的日常工作收敛成一个**自包含、可克隆、可无人值守运行**的起点仓。
任何 Claude `git clone jarvis + 注入凭证`，即可接替我的工作，且这是所有后续项目的母版。

## 三个已锁定的决策

1. **接替程度 = 无人值守自治**：定时/主动跑，只在不可逆或低置信时找人。
2. **起步范围 = 只做 Aone triage**：最高频 loop 吃透，立为模板，再复制。
3. **自治边界 = 全链到预发/CR**：read → verify → reply → 建需求 → 建 CR → worktree 开发 → 预发，**唯一硬门 = 正式发布**。

## 定位

jarvis 是 `my-day` / `aone-triage` / `a1` / `cloudspec` 之上的**交接层 + 母版**，
但技能**收进项目**（vendored 到 `skills/`，不依赖 `~/.claude`）。
克隆即自包含，新项目从 jarvis fork。它存的是工具之外的四样东西：你是谁、活有哪些、敢干到哪、错了怎么回滚。

## 架构：六个件

| 件 | 路径 | 作用 |
|---|---|---|
| 自举 | `IDENTITY.md` | 进来先读：你是谁、对谁负责、开局动作。补 README 空缺 |
| 清单 | `loops/aone-triage.md` | 触发/输入/工具链/决策点/done/仅人工步。第一条吃透 |
| 决策权 | `autonomy.md` | 硬门=正式发布；预发/CR 以下全自动；低置信降级为"起草不发出+入队" |
| 真源/记忆 | `.my-day/` + `memory/` | Aone=真源、可重建；memory 存跨会话连续性 |
| 验收 | `loops/aone-triage.md` 内 done 门 | 预发即验证门 + benchmark，过不了即停 |
| 审计回滚 | `runs/` 日志 + `escalation/` 队列 | 每步可追、CR 未合即可逆、停的活一键续 |

技能 vendored 到 `skills/aone-triage` 等；凭证由 `bootstrap/` 注入（容器内 a1 凭证、gh token、其余用户给）。

## 触发链

`定时 scan → triage → 自交付到预发/CR → wrap`，**停在正式发布前**。
低置信 / 验收不过 → 降级起草并入 `escalation/` 队列，待人一键续。

## 路径

- **P0 自举**：IDENTITY + README，写清"任何 Claude 怎么开局"
- **P1 盘点**：把 aone-triage 这条 loop 写成清单（仅 triage）
- **P2 单点跑通**：确认门下全自动，技能收敛进项目，鉴权 smoke test
- **P3 拆门放权**：信任建立后减少确认，全链到预发
- **P4 自治**：定时跑，只升级才找人

## 待闭合缺口（按阻塞排序）

1. **headless 鉴权**（最硬）：容器内 a1 凭证 / gh token / 用户兜底 → 必须 smoke test 跑通
2. **escalation 队列**：停下的活+原因+一键续
3. **审计回滚**：全量 run 日志、可撤销
4. **验收门**：预发前过测/benchmark

## 成功判据

容器里克隆 jarvis、注凭证、定时一夜，第二天：triage 工作项被处理到预发/CR、停在正式发布、escalation 队列可读、全程可审。

## YAGNI

不做：全量日常盘点（先 triage 一条）、自动正式发布、重写技能、Web UI。
