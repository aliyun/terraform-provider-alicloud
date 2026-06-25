# Tata × Jarvis 双层桥接 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax.

**Goal:** 把单层重型桥接拆成 Tata 门面陪聊 + 仅辰羿可升级重型 Jarvis 的两层。

**Architecture:** 复用 bridge/jarvis_dingtalk_bot.py 的真流式建卡/流推；process() 先跑 Tata（轻 system-prompt，空 cwd），扫回复里的 `[[JARVIS]]` 哨兵，仅当 sender==master 才起第二张卡跑重型 Jarvis。

**Tech Stack:** Python 3.13, dingtalk_stream, headless claude CLI, dingtalk-ai-card streaming.py（import）。

## Global Constraints

- 只改 `bridge/jarvis_dingtalk_bot.py`，连带 `bridge/README.md`、`bridge/jarvis.env.example`；不碰 master。
- 卡片截断 2KB（`MAX_REPLY=2000`），CARD_KEY="content"，节流 `PUT_MIN_INTERVAL=0.4`/`PUT_MIN_GROWTH=40` 不变。
- 永不崩 WS 循环：任何轮异常 finalize is_error。
- commit 中文、结尾保 Co-Authored-By；语法 `ast.parse` + `import dingtalk_stream` 保持绿。
- master 闸在 bridge 端硬判，Tata 提示词不可信。

---

### Task 1: 哨兵解析 + 闸门 helper（含单测）

**Files:**
- Modify: `bridge/jarvis_dingtalk_bot.py`（加 `extract_jarvis_task`、`tata_audience`、`master_staff`、`tata_root`）
- Test: scratch 单测文件，验完即删

**Interfaces:**
- Produces: `extract_jarvis_task(text)->(clean:str, task:str|None)`；`tata_audience()->set`；`master_staff()->str`；`tata_root()->str`

- [ ] **Step1 写失败测试**：sentinel 末行 `[[JARVIS]] 部署预发`→clean 去掉该行、task=="部署预发"；无 sentinel→(原文,None)；多空行容忍。
- [ ] **Step2 跑测试见 RED**（function 未定义）。
- [ ] **Step3 实现**：正则 `^\s*\[\[JARVIS\]\]\s*(.+)$` 多行扫，剥行返干净文本+任务；`tata_audience` 读 `JARVIS_TATA_STAFF`(空→空集=全员)；`master_staff` 默认 320687；`tata_root` 读 `JARVIS_TATA_ROOT` 默认 `~/.jarvis/tata-cwd`（建空目录）。
- [ ] **Step4 跑测试见 GREEN**。
- [ ] **Step5 commit**（feat: 哨兵+闸门 helper）。

### Task 2: Tata 轻量 runner

**Files:** Modify `bridge/jarvis_dingtalk_bot.py`（加 `run_tata_stream`）
**Interfaces:** Consumes `parse_stream_lines`；Produces `run_tata_stream(text,sid,resume)` 生成器，cwd=tata_root，附 `--append-system-prompt <人设>`，复用 stream-json 流式，不吃 jarvis CLAUDE.md。

- [ ] **Step1** 复制 run_claude_stream 改 cwd=tata_root()、加 `--append-system-prompt` TATA_PROMPT（spec 草案人设）。
- [ ] **Step2** `ast.parse` 过；空 root 目录确保创建。
- [ ] **Step3 commit**。

### Task 3: process() 两层接力 + 文档

**Files:** Modify bot（process/_stream_round 返回终文）、README、jarvis.env.example
**Interfaces:** Consumes Task1/2。

- [ ] **Step1** `_stream_round` 返回最终 acc；新增 `JARVIS_PROMPT`(可空)。
- [ ] **Step2** process：audience 闸(空→全员)；Tata session per-sender 跑 `run_tata_stream`→拿全文→`extract_jarvis_task`；clean 文上屏；task 且 `staff==master_staff()`→第二卡跑重型 Jarvis(run_claude_stream, jarvis session)。非 master 丢 task。
- [ ] **Step3** README 画两层流程；env.example 补 `JARVIS_TATA_STAFF/JARVIS_MASTER_STAFF/JARVIS_TATA_ROOT`。
- [ ] **Step4** 语法+依赖绿；commit。

## Self-Review
覆盖 spec 全部点（受众/master 闸/哨兵剥离/两 session/handoff 单卡/错误不崩）。无占位。helper 命名一致。
