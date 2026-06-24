# 设计：纪律工程化 + subagent 一条一派

> Aone: 83491649（api_toolkit / 2100304）。把"提示词纪律"换成"代码兜底"，解决两个根因问题。

## 问题

1. **收尾约束全是提示词，无强制。** "回填 Aone / 更新 tag/status" 只写在 `CLAUDE.md 纪律5`、`loops/aone-triage.md`，没有一行代码强制。`wrap.sh done` 改 status 是可选参数，且所有 `a1` 调用 `|| echo` warn-only 不报错——漏调或调失败流程仍退 0。结果：tag/status 漂，工单统计失真。
2. **无 subagent 隔离。** 无 `.claude/agents/`，triage/dev 全在主会话内联跑，上下文越堆越脏，"编码交子代理"只是一句提示词。

## 目标 / 非目标

- 目标：tag/status 收尾从"靠模型自觉"变"代码兜底 + 钩子拦截 + 事后对账"三层冗余；每条工单跑在独立 subagent，主会话只编排。
- 非目标：不改 scan/plan/claim 既有契约；不动 release_prod 硬门；不引入外部服务。

## 架构：4 层

### 第1层 · 编排：一条工单一 subagent
- 新增 `bootstrap/triage-one.sh <id> <pool> <project>`：bookend 全在脚本，主 Agent 不靠记忆。
  顺序：`claim.sh claim` → 派 subagent 跑单 → `wrap.sh done <id> <summary> <status>`（status 必填）→ `claim.sh release`。任一步失败 → escalate + 不 release。
- 新增 `.claude/agents/`：`triager`（默认全流程，内部可再派）、`dev`（编码调试）、`pr-review`、`verify`（OpenAPI+源码查证）。主会话只编排，每条派一个，上下文干净独立。

### 第2层 · wrap.sh done 收紧（直接改硬）
- `status` 从可选改**必填**，缺则 `exit 1`。
- `a1 comment/update` 失败从 warn-only 改 `exit 1`（漂的根因）。本地 `run_done` 仍先写（不丢账），但 Aone 未达即非零退出，收尾视为失败。

### 第3层 · Stop-hook 硬闸
- 新增项目级 `.claude/settings.json` 注册 Stop hook → `bootstrap/wrap-check.sh`。
- 认领台账 `.my-day/claims-<date>.json`（claim 时追加 id，release 时勾掉）。Stop 时校验每条已认领 id 是否齐：`runs/<date>-<id>.md` + `jarvis-done` + status。缺一 `exit 2` 拦住会话，附 reason 逼收齐。

### 第4层 · reconcile 对账器
- 新增 `bootstrap/reconcile.sh`：扫 `runs/` vs Aone，jarvis-claimed 无 jarvis-done / 状态没推 → 自动补；接进 loop 尾。漏网兜底。

## 数据流
scan → plan → 主Agent 逐条 → `triage-one.sh`（claim→subagent→done→release）→ loop 尾 `reconcile.sh` → 会话结束触发 Stop-hook wrap-check。

## 验收
- 漏传 status / a1 失败 → `wrap.sh done` 报错；认领未收齐 → 会话被拦；triage 跑在 subagent，主会话不堆代码细节；reconcile 能把已知漂复位。

## 测试
- `wrap.sh done` 缺 status / 模拟 a1 失败 → 非零。`wrap-check.sh` 台账缺项 → exit 2。`triage-one.sh` dry-run 编排顺序。bats 或 shell 断言。
