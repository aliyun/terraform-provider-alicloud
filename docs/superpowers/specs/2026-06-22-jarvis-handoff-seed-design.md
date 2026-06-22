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
| 自举 | `CLAUDE.md`（自动入上下文）`@import` 身份/清单/决策权 | 进来必读：你是谁、对谁负责、开局动作。任意 IDENTITY.md 不会被自动读，故用 CLAUDE.md 作主入口 |
| 清单 | `loops/aone-triage.md` | 触发/输入/工具链/决策点/done/仅人工步。第一条吃透 |
| 决策权 | `autonomy.md` | 硬门=正式发布；预发/CR 以下全自动；低置信降级为"起草不发出+入队" |
| 真源/记忆 | `.my-day/` + `memory/` | Aone=真源、可重建；memory 存跨会话连续性 |
| 验收 | `loops/aone-triage.md` 内 done 门 | 预发即验证门 + benchmark，过不了即停 |
| 审计回滚 | `runs/` 日志 + `escalation/` 队列 | 每步可追、CR 未合即可逆、停的活一键续 |

技能 vendored 到 `skills/aone-triage` 等；凭证由 `bootstrap/` 注入（容器内 a1 凭证、gh token、其余用户给）。

## 触发链

`定时 scan → triage → 自交付到预发/CR → wrap`，**停在正式发布前**。
低置信 / 验收不过 → 降级起草并入 `escalation/` 队列，待人一键续。

## 前置依赖（全显式、可自动装）

原则：`bootstrap/install.sh` 幂等；`bootstrap/verify.sh` **每个依赖一条独立 check，单独 PASS/FAIL，任一 FAIL 整体退非零，绝不聚合成"全 ok"**。逐项独立验意味着缺哪个一眼定位。新容器 = 克隆 + 注凭证 + 一键装 + 全绿即可干活。

| 类 | 清单 | 装法 | 验法 |
|---|---|---|---|
| CLI | a1 / gh / git / aliyun / terraform | 锁版本，脚本装 | `--version` |
| 技能 | aone-triage 等 vendored 到 `skills/` | 入仓即随克隆 | 能被 Skill 加载 |
| MCP | claude.ai / 其它必需 | `bootstrap/mcp.json` 声明 | 列工具成功 |
| 凭证 | a1(容器内) / gh token / aliyun key | `.env.example` 模板，用户注入 | `verify.sh` 各调一下 |

`bootstrap/deps.lock` 记全部版本。每个凭证也各自单独验（gh token 调一次、a1 调一次、aliyun 调一次），不互相代替。`verify.sh` 输出形如 `PASS a1 / PASS gh / FAIL aliyun`，整体退非零。

## 路径

- **P0 自举**：`CLAUDE.md`（`@import` 身份/`loops/`/`autonomy.md`）+ README，写清"任何 Claude 怎么开局"
- **P1 盘点**：把 aone-triage 这条 loop 写成清单（仅 triage）
- **P2 单点跑通**：确认门下全自动，技能收敛进项目，鉴权 smoke test
- **P3 拆门放权**：信任建立后减少确认，全链到预发
- **P4 自治**：定时跑，只升级才找人

## 待闭合缺口（按阻塞排序）

1. **依赖+鉴权**（最硬）：CLI/技能/MCP/凭证全显式，`install.sh` 装、`verify.sh` smoke test 跑通
2. **escalation 队列**：停下的活+原因+一键续
3. **审计回滚**：全量 run 日志、可撤销
4. **验收门**：预发前过测/benchmark

## 成功判据

容器里克隆 jarvis、注凭证、定时一夜，第二天：triage 工作项被处理到预发/CR、停在正式发布、escalation 队列可读、全程可审。

## YAGNI

不做：全量日常盘点（先 triage 一条）、自动正式发布、重写技能、Web UI、独立 IDENTITY.md（合进 CLAUDE.md）。
