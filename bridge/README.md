# Jarvis DingTalk bridge (inbound, two-way)

Long-running process: holds a DingTalk **Stream** WebSocket, receives messages
you send to the bot, runs a headless `claude` round, and streams the answer
back **live** into one AI card. Turns the one-way `dingtalk-ai-card` push into a
two-way conversation with a real typewriter effect.

## Flow (true streaming)
1. You DM the bot → `jarvis_dingtalk_bot.py` checks whitelist (`JARVIS_ALLOW_STAFF`).
2. Creates one AI card immediately, then reads `claude -p <text> --output-format
   stream-json --include-partial-messages --verbose --session-id <per-sender uuid>`
   token-by-token and PUTs the growing text onto that card (throttled), so you
   watch the answer build up. 300s cap. Finalizes when done; truncated to 2KB.
3. Card helpers are imported from the skill's `streaming.py` (no subprocess回放).

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
