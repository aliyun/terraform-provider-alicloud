# Tata 常驻进程保温 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development. Steps use checkbox (`- [ ]`).

**Goal:** Tata 用常驻 idea 进程多轮保温消冷启；Jarvis 改用 cc 启动。

**Architecture:** 新增 TataPool：每 sender 一个长生 `idea`(=claude --settings idea_settings) 子进程跑 stream-json 双向，逐轮喂 user JSON 读到 result。process() 用 pool 发 Tata；Jarvis 用 cc(~/claude-start.sh)。

**Tech Stack:** Python3.13, dingtalk_stream, claude CLI(idea/cc), parse_stream_lines。

## Global Constraints

- 只改 `bridge/jarvis_dingtalk_bot.py`+README+jarvis.env.example，不碰 master/不 push。
- Tata=`claude --settings ~/.claude/idea_settings.json`(env JARVIS_TATA_SETTINGS)；Jarvis=`~/claude-start.sh`(env JARVIS_CC)，挂 stdin</dev/null。
- idle 30min(JARVIS_TATA_IDLE_MIN)、max 10 LRU(JARVIS_TATA_MAX)。
- 永不崩 WS；起不来回退一次性 run_tata_stream。MAX_REPLY=2000/CARD_KEY/节流不变。
- commit 中文+Co-Authored-By；ast.parse+import dingtalk_stream 绿。

---

### Task 1: 启动命令 helper

**Files:** Modify bridge/jarvis_dingtalk_bot.py
**Interfaces:** Produces `tata_cmd()->list`(=[claude,--settings,设置档]) ；`jarvis_cmd()->list`(=[~/claude-start.sh])。

- [ ] Step1 测：JARVIS_TATA_SETTINGS 覆盖时 tata_cmd 末元素=该路径；jarvis_cmd[0] 以 claude-start.sh 结尾。
- [ ] Step2 RED。
- [ ] Step3 实现：tata_cmd 用 claude_bin()+["--settings",os.environ.get JARVIS_TATA_SETTINGS or ~/.claude/idea_settings.json]；jarvis_cmd=[os.environ.get JARVIS_CC or ~/claude-start.sh]。
- [ ] Step4 GREEN。Step5 commit。

### Task 2: TataPool 常驻保温

**Files:** Modify bot（加 class TataPool）
**Interfaces:** Consumes tata_cmd/parse_stream_lines；Produces `TataPool(max_n,idle_min)`，`send(staff,text)->generator(累积文本)`，`shutdown()`。

- [ ] Step1 测：mock 进程 framing——写一行 user JSON,读到 result 停,返累积;空闲>idle 回收;>max LRU 踢旧;崩溃下条重起;spawn 失败抛特定异常。
- [ ] Step2 RED。
- [ ] Step3 实现：Popen(tata_cmd+stream-json双向+partial+verbose) cwd=tata_root,stdin常开;send 写 `{"type":"user","message":{"role":"user","content":text}}\n`+flush,逐行 parse_stream_lines 读到 result yield;last_used 更新;锁串行;reap 线程/惰性扫。
- [ ] Step4 GREEN。Step5 commit。

### Task 3: 接 process + Jarvis 用 cc + 文档

**Files:** Modify bot/README/jarvis.env.example
**Interfaces:** Consumes Task1/2。

- [ ] Step1 process Tata 轮改 self.pool.send(staff,text)（替 run_tata_stream），起不来回退一次性；run_claude_stream 用 jarvis_cmd 起;handler init 建 pool,main 退出 shutdown。
- [ ] Step2 README 注 idea/cc+保温;env.example 补 JARVIS_TATA_IDLE_MIN/MAX/SETTINGS/CC。
- [ ] Step3 语法+依赖绿;commit。

## Self-Review
覆盖:保温/idle30/max10LRU/崩溃重起/回退/idea/cc/stdin</dev/null。无占位。命名一致 tata_cmd/jarvis_cmd/TataPool.send。
