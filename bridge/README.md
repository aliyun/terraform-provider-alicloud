# Jarvis DingTalk bridge (inbound, two-way)

Long-running process: holds a DingTalk **Stream** WebSocket, receives messages
you send to the bot, runs a headless `claude` round, streams the answer back as
an AI card. Turns the one-way `dingtalk-ai-card` push into a two-way conversation.

## Flow
1. You DM the bot → `jarvis_dingtalk_bot.py` checks whitelist (`JARVIS_ALLOW_STAFF`).
2. Acks `🟡 收到, 处理中…` (no-stream) so you see it instantly.
3. `claude -p <text> --output-format text --session-id <per-sender uuid>` in `JARVIS_ROOT`, 300s cap.
4. Streams the answer back (truncated <3KB) via the skill `streaming.py` card.

One claude round per sender at a time; second message → "请稍候". Other senders independent.

## Setup
1. **DingTalk console**: app → enable **Stream mode** (not webhook). Grant the
   robot message read/send. Note appKey/appSecret. Create an AI card template, note its id.
2. `cp jarvis.env.example jarvis.env`, fill APP_KEY/SECRET/TEMPLATE_ID, set
   `JARVIS_ALLOW_STAFF` (your staffId, default 320687), `JARVIS_ROOT`.
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
