# Tata 常驻进程保温设计（消冷启）

> 2026-06-25 · branch worktree-dingtalk-stream-bridge · 改 bridge/jarvis_dingtalk_bot.py

## 背景

Tata 现每轮 `claude -p` 冷起，闲聊也要 15-18s，全花在 CLI 起进程+加载。Jarvis
重型慢可接受，但 Tata 是门面要快。无 ANTHROPIC_API_KEY（CLI 走订阅登录），直连
API 不可行 → 用常驻 claude 进程保温。

## 验证（已实测）

单个 `claude -p --input-format stream-json --output-format stream-json` 进程喂两条
user JSON，依次回 A、B，一进程多轮保温。冷启只付一次，续轮 ~1-3s。

## 架构

- 每个 sender 一个长生 claude 子进程（=Tata 会话），cwd=tata_root，附 TATA_PROMPT，
  stdin 常开。每条消息写一行 `{"type":"user","message":{"role":"user","content":text}}`，
  读到该轮 `result` 为止，期间 text_delta 流式入卡。续轮不重启 = 保温。
- Jarvis 不变：仍每次 run_claude_stream 冷起；master 升级闸不变。

## 新组件 TataPool（bridge 内一类）

`procs: {staff -> {proc, last_used, lock}}`
- `send(staff, text) -> generator`：无/死进程→懒起；写 JSON、按 parse_stream_lines 读本轮、yield 累积；本轮 result/超时即止（不关进程）。
- 保温边界：空闲 >30min 回收；最多 10 个并发，超额 LRU 踢最旧。
- 进程崩→标记，下条自动重起。spawn 失败→回退一次性 run_tata_stream，不阻断。
- 关停 main 时全 kill。

## 配置

- 复用 tata_root/JARVIS_TATA_STAFF/JARVIS_MASTER_STAFF。
- 新增 `JARVIS_TATA_IDLE_MIN`=30、`JARVIS_TATA_MAX`=10。

## 数据流

消息 → audience 闸 → TataPool.send(staff,text) 流式入卡 → 全文 extract_jarvis_task
→ master 且有任务→第二卡 run_claude_stream(Jarvis)。Tata 不再用 uuid/resume（进程即会话）。

## 错误处理

- 进程死/管道断→重起重试一次,再败回退一次性。永不崩 WS。
- 单测：framing 正确、读到 result 即止、idle/LRU 回收、崩溃重起。

## 收尾

spec+plan 落 docs/；建 Aone(2100304)+MR 待辰羿合,不 push/merge。
