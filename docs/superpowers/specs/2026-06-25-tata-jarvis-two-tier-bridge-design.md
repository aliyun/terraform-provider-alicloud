# Tata × Jarvis 双层钉钉桥接设计

> 2026-06-25 · branch worktree-dingtalk-stream-bridge · 改造 bridge/jarvis_dingtalk_bot.py

## 背景

当前每条钉钉消息都直接拉起重型 Jarvis（headless claude，cwd=仓库根，吃整套
CLAUDE.md+autonomy+triage bootstrap），首轮 ~100s。太重、太慢，且任何白名单
用户都能调起 Jarvis 干真活。要拆出一个轻量门面 **Tata** 陪聊，**仅辰羿
(staffId 320687)** 能透过 Tata 升级到 Jarvis。

## 目标

- 全员能跟 Tata 闲聊（轻、快，首轮秒回）。
- 只有辰羿能让 Tata 喊起重型 Jarvis；别人诱导也调不动。
- 复用已就绪的真流式建卡/流推，不重写卡片逻辑。

## 架构

| | Tata（门面） | Jarvis（重型） |
|---|---|---|
| 受众 | 全员（`JARVIS_TATA_STAFF` 空=不限；可填名单收窄） | 仅辰羿 `JARVIS_MASTER_STAFF=320687` |
| 实现 | `claude -p` + 精简 `--append-system-prompt`，cwd=空目录，**不吃 jarvis CLAUDE.md** | `claude` cwd=`JARVIS_ROOT` 仓库根（现状） |
| 速度 | 首轮秒回 | 重，~100s |
| session | per-sender uuid（Tata 套） | 仅辰羿（Jarvis 套，独立） |

## 数据流

1. 钉钉消息进 `process()` → **Tata 闸**：sender ∈ `JARVIS_TATA_STAFF`（空=放行全员）；
   空文本/锁忙照旧。
2. 建卡 → 跑 **Tata 一轮**（真流式秒回）：`claude -p <text> --append-system-prompt <人设>
   --output-format stream-json --include-partial-messages --verbose`，cwd=空目录，
   `--session-id/--resume` 用 Tata session。
3. Tata 人设要求：干真活(查证/开发/运维)时在回复**末尾单起一行** `[[JARVIS]] <任务>`。
4. bridge 收尾扫 Tata 全文哨兵 + sender==`JARVIS_MASTER_STAFF` → 起**第二张卡**跑
   Jarvis(cwd=仓库根, Jarvis session)流式接力；Tata 卡定格"交给 Jarvis 处理…"。
   非辰羿哨兵丢弃，只当闲聊。

## 人设（Tata system-prompt 草案）

> 你是 Tata，钉钉里的轻量助手。日常陪聊、答疑、查资料，语气简洁友好。
> 你**不能**直接动仓库/发布/调 IaC。若辰羿要你干真活(查证/开发/运维/碰工单)，
> 在回复最后单起一行 `[[JARVIS]] <一句话任务>`，由系统转交 Jarvis；其余人只闲聊。

## 配置

- `JARVIS_TATA_STAFF`：Tata 受众名单，逗号分隔；空=全员。
- `JARVIS_MASTER_STAFF`：可调 Jarvis 的 staffId，默认 320687。
- `JARVIS_TATA_ROOT`：Tata cwd，默认空临时目录（不加载 jarvis bootstrap）。
- 复用：`JARVIS_ROOT`(Jarvis cwd)、`CLAUDE_TIMEOUT`、卡片三件套。

## 哨兵 & 闸门

- 哨兵 `[[JARVIS]]` 只是路由信号，正则剥掉再不展示给用户。
- 两道独立闸：Tata 受众闸（聊不聊）+ master 闸（升不升级）。master 闸在 bridge 端
  做，Tata 提示词不可信，越权一律拦。

## 错误处理

- Tata/Jarvis 任一轮异常 → finalize is_error，不崩 WS 循环（现状保留）。
- 升级时 Jarvis 卡独立失败不影响 Tata 已发卡。

## 测试

- `parse_stream_lines` 已有单测；新增哨兵剥离/master 闸：非辰羿带哨兵不升级、
  辰羿带哨兵升级、无哨兵不升级。
- 语法 `ast.parse`、依赖 `import dingtalk_stream` 保持绿。

## 收尾纪律

spec commit；按 CLAUDE.md 建 Aone 单(池2100304)+MR 待辰羿合，不 push/merge master。
