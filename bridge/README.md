# Jarvis DingTalk bridge (two-tier: Tata × Jarvis)

Long-running process: holds a DingTalk **Stream** WebSocket, receives messages
you send to the bot, and streams the answer back **live** into one AI card with a
real typewriter effect. Two tiers: a lightweight **Tata** front 陪聊 everyone, and
heavy **Jarvis** only 辰羿 can escalate to.

## Two-tier flow
```
DingTalk DM ─► 受众闸(JARVIS_TATA_STAFF, 空=全员)
            └► Tata 卡: claude -p + 轻量人设, cwd=空目录(不吃 CLAUDE.md), 秒回流式
               扫回复末尾 [[JARVIS]] <任务> 哨兵
                 ├ sender == JARVIS_MASTER_STAFF(辰羿) ► Tata 卡定格"交给 Jarvis 处理…"
                 │                                       第二张卡: 重型 Jarvis(cwd=仓库根)流式接力
                 └ 其他人/无哨兵 ► 只闲聊, 哨兵剥掉不上屏
```
1. **受众闸**: `JARVIS_TATA_STAFF` 空=全员可跟 Tata 聊; 填名单则收窄。
2. **Tata 一层**: 建卡即跑 `claude -p <text> --append-system-prompt <Tata人设>
   --output-format stream-json --include-partial-messages --verbose`, cwd=空目录,
   秒回; per-sender session。回复里 `[[JARVIS]]` 哨兵剥掉再上屏。
3. **master 闸 + 升级**: 仅辰羿(`JARVIS_MASTER_STAFF`)且 Tata 发了哨兵, 第二张卡跑重型
   Jarvis(cwd=`JARVIS_ROOT` 仓库根, 独立 session)流式接力。300s cap, 截 2KB。
4. Card helpers imported from the skill's `streaming.py`. master 闸在 bridge 端硬判, Tata 人设不可信。

One round per sender at a time; second message → "请稍候". Other senders independent.

## Setup
1. **DingTalk console**: app → enable **Stream mode** (not webhook). Grant the
   robot message read/send. Note appKey/appSecret. Create an AI card template, note its id.
2. `cp jarvis.env.example jarvis.env`, fill APP_KEY/SECRET/TEMPLATE_ID, set
   `JARVIS_TATA_STAFF` (空=全员), `JARVIS_MASTER_STAFF` (你的 staffId, 默认 320687),
   `JARVIS_TATA_ROOT`/`JARVIS_ROOT`.
3. Ensure deps: `python3 -c "import dingtalk_stream"` and `claude` on PATH (`~/.local/bin`).
4. `./run.sh start` → `./run.sh status` → `./run.sh logs`. Stop: `./run.sh stop`.

## Autostart
- mac: fill placeholders in `com.jarvis.dingtalk.plist.example`, copy to
  `~/Library/LaunchAgents/`, `launchctl load …`.
- linux: systemd note at bottom of the plist example.

## Migrate to another machine
Copy `bridge/` + `jarvis.env`, install `dingtalk-ai-card` skill + claude, then
`./run.sh start`. No webhook/public IP needed — Stream is outbound. State (pid/log)
lives in `${JARVIS_STATE:-~/.jarvis}`.
