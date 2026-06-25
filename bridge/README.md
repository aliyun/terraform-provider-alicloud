# Jarvis DingTalk bridge (two-tier: Tata × Jarvis)

Long-running process: holds a DingTalk **Stream** WebSocket, receives messages
you send to the bot, and streams the answer back **live** into one AI card with a
real typewriter effect. Two tiers: a lightweight **Tata** front 陪聊 everyone, and
heavy **Jarvis** only 辰羿 can escalate to.

## Two-tier flow
```
DingTalk DM ─► 受众闸(JARVIS_TATA_STAFF, 空=全员)
            └► Tata 卡: 常驻 idea 进程保温(每 staff 一个), cwd=空目录(不吃 CLAUDE.md), 秒回流式
               扫回复末尾 [[JARVIS]] <任务> 哨兵
                 ├ sender == JARVIS_MASTER_STAFF(辰羿) ► Tata 卡定格"交给 Jarvis 处理…"
                 │                                       第二张卡: 重型 Jarvis(cc/cwd=仓库根)流式接力
                 └ 其他人/无哨兵 ► 只闲聊, 哨兵剥掉不上屏
```
1. **受众闸**: `JARVIS_TATA_STAFF` 空=全员可跟 Tata 聊; 填名单则收窄。
2. **Tata 一层(保温)**: 每 sender 一个长生 **idea** 进程(`claude --settings ~/.claude/idea_settings.json`,
   走 idealab 网关)跑双向 stream-json, 逐轮喂一行 `{"type":"user",...}` 读到 result, 进程不关 →
   冷启只付一次, 续轮 ~1-3s, cwd=空目录。启动即预热 `JARVIS_TATA_PREWARM`(默认3)个未绑定 generic
   进程待命, 首批 sender 直接领走免冷启、后台补满。空闲 30min 回收, >10 并发 LRU 踢旧; 起不来回退一次性冷起。
   回复里 `[[JARVIS]]` 哨兵剥掉再上屏。
3. **master 闸 + 升级**: 仅辰羿(`JARVIS_MASTER_STAFF`)且 Tata 发了哨兵, 第二张卡跑重型
   Jarvis = **cc**(`~/claude-start.sh`, 预检后 exec claude, 挂 stdin</dev/null 防 IP 不符卡死),
   cwd=`JARVIS_ROOT` 仓库根, 独立 session 流式接力。300s cap, 截 2KB。
4. Card helpers imported from the skill's `streaming.py`. master 闸在 bridge 端硬判, Tata 人设不可信。

One round per sender at a time; second message → "请稍候". Other senders independent.

## Setup
1. **DingTalk console**: app → enable **Stream mode** (not webhook). Grant the
   robot message read/send. Note appKey/appSecret. Create an AI card template, note its id.
2. `cp jarvis.env.example jarvis.env`, fill APP_KEY/SECRET/TEMPLATE_ID, set
   `JARVIS_TATA_STAFF` (空=全员), `JARVIS_MASTER_STAFF` (你的 staffId, 默认 320687),
   `JARVIS_TATA_ROOT`/`JARVIS_ROOT`. 保温调参: `JARVIS_TATA_IDLE_MIN`(30)/`JARVIS_TATA_MAX`(10)/`JARVIS_TATA_PREWARM`(3);
   启动档: `JARVIS_TATA_SETTINGS`(idea, 默认 ~/.claude/idea_settings.json)、`JARVIS_CC`(默认 ~/claude-start.sh)。
3. Ensure deps: `python3 -c "import dingtalk_stream"`, `claude` on PATH (`~/.local/bin`),
   idea 设置档 `~/.claude/idea_settings.json` + `~/claude-start.sh` 在位。
4. `./run.sh start` → `./run.sh status` → `./run.sh logs`. Stop: `./run.sh stop`.

## Autostart
- mac: fill placeholders in `com.jarvis.dingtalk.plist.example`, copy to
  `~/Library/LaunchAgents/`, `launchctl load …`.
- linux: systemd note at bottom of the plist example.

## Migrate to another machine
Copy `bridge/` + `jarvis.env`, install `dingtalk-ai-card` skill + claude, then
`./run.sh start`. No webhook/public IP needed — Stream is outbound. State (pid/log)
lives in `${JARVIS_STATE:-~/.jarvis}`.
