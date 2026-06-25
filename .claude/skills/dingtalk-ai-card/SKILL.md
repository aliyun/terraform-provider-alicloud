---
version: 1.0.0
description: >
  通过钉钉 AI 卡片发送流式消息（打字机效果）。当用户需要发送流式钉钉消息、
  AI 卡片消息、打字机效果消息、或通过钉钉卡片推送内容时使用此技能。
  触发关键词：流式消息、AI卡片、打字机效果、streaming card、发卡片、钉钉卡片推送。
allowed-tools: Bash, Read
name: dingtalk-ai-card
x-source: aone-open
---

# 钉钉 AI Card 流式消息技能

通过钉钉 AI 卡片 API 发送带打字机效果的流式消息。

## 前置条件

### 1. 安装 am CLI

```bash
curl -fsSL https://am.io.alibaba-inc.com/install.sh | bash
```

### 2. 绑定机器人

凭证全走环境变量，绑定可自动幂等完成（已绑则跳过）：

```bash
# 设好 DINGTALK_APP_KEY / DINGTALK_APP_SECRET / DINGTALK_STAFF_ID（DINGTALK_ROBOT_CODE 默认=appKey）
bash ~/.claude/skills/dingtalk-ai-card/scripts/ensure-bind.sh
```

手动等价：

```bash
am bind --type=bot \
  --access-key-id=$DINGTALK_APP_KEY \
  --access-key-secret=$DINGTALK_APP_SECRET \
  --account-id=$DINGTALK_STAFF_ID \
  --robot-code $DINGTALK_APP_KEY
```

### 3. 创建 AI 卡片模板

在 [钉钉卡片平台](https://open-dev.dingtalk.com/fe/card) 创建：

1. 新建模板 → 卡片类型选「消息卡片」→ 场景选「AI 卡片」
2. 关联你的机器人应用
3. 在「输入中」状态下，确认 Markdown 组件已开启「流式组件」开关
4. 记住绑定的 **markdown 变量名**（默认通常是 `content`）
5. 保存并发布，记录 **模板 ID**（形如 `xxxxxxxx.schema`）

## 使用方法

脚本位于 `scripts/streaming.py`，自动从 `am bind` 配置读取凭据。

### 发送流式消息（打字机效果）

```bash
python3 ~/.claude/skills/dingtalk-ai-card/scripts/streaming.py \
  --to <staffId> \
  --template-id <templateId> \
  -m $'消息内容\n支持换行和**Markdown**'
```

### 一次性发送（无打字机效果）

```bash
python3 ~/.claude/skills/dingtalk-ai-card/scripts/streaming.py \
  --to <staffId> \
  --template-id <templateId> \
  --no-stream \
  -m "直接显示完整内容"
```

### 从 stdin 读取

```bash
echo "内容" | python3 ~/.claude/skills/dingtalk-ai-card/scripts/streaming.py \
  --to <staffId> \
  --template-id <templateId> \
  --stdin
```

### 发送到群聊

```bash
python3 ~/.claude/skills/dingtalk-ai-card/scripts/streaming.py \
  --to-group <openConversationId> \
  --template-id <templateId> \
  -m "群消息"
```

### 参数说明

| 参数 | 缩写 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--to` | | 二选一 | | 目标用户 staffId |
| `--to-group` | | 二选一 | | 目标群 openConversationId |
| `--message` | `-m` | 是* | | 消息内容（与 `--stdin` 二选一） |
| `--stdin` | | 是* | | 从标准输入读取消息 |
| `--template-id` | | 是 | `$DINGTALK_TEMPLATE_ID` | AI 卡片模板 ID |
| `--chunk-size` | | 否 | 2 | 每次流式更新的字符数 |
| `--delay` | | 否 | 0.15 | 每次更新间隔秒数 |
| `--no-stream` | | 否 | | 跳过打字机效果，直接显示完整内容 |
| `--key` | | 否 | content | 卡片模板中流式 markdown 变量名 |
| `--bot` | | 否 | | am 的命名 bot profile |
| `--app-key` | | 否 | am 配置 | 覆盖 appKey |
| `--app-secret` | | 否 | am 配置 | 覆盖 appSecret |
| `--robot-code` | | 否 | am 配置 | 覆盖 robotCode |

### 环境变量

| 变量 | 说明 |
|------|------|
| `DINGTALK_APP_KEY` | appKey（优先级高于 am 配置） |
| `DINGTALK_APP_SECRET` | appSecret |
| `DINGTALK_ROBOT_CODE` | robotCode |
| `DINGTALK_TEMPLATE_ID` | 默认模板 ID |

## 重要提示

- shell 中传含换行的消息必须用 `$'...\n...'` 语法，不要用双引号中的 `\n`
- markdown 类型的流式变量，`isFull` 必须为 `true`（脚本已默认设置）
- 单次流式内容不超过 1KB，总大小建议不超过 3KB
- 卡片模板必须已发布且关联了正确的机器人应用
- 机器人应用需把出口 IP 加进白名单（IPv4/IPv6 都加，或临时关闭），否则 `am`/卡片 403
- macOS python.org 的 Python 3.13 缺根证书：跑 `streaming.py` 前 `export SSL_CERT_FILE=$(python3 -c 'import certifi;print(certifi.where())')`，否则 SSL 校验失败
- 常用模板可写进 `DINGTALK_TEMPLATE_ID` 免每次带 `--template-id`

## 底层 API 参考

| 步骤 | API | 方法 |
|------|-----|------|
| 获取 token | `POST /v1.0/oauth2/accessToken` | appKey + appSecret |
| 创建并投放卡片 | `POST /v1.0/card/instances/createAndDeliver` | openSpaceId + templateId |
| 流式更新 | `PUT /v1.0/card/streaming` | outTrackId + content |

### openSpaceId 格式

| 场域 | 格式 |
|------|------|
| 机器人单聊 | `dtv1.card//IM_ROBOT.<userId/staffId>` |
| 群聊 | `dtv1.card//IM_GROUP.<openConversationId>` |

### 卡片状态流转

```
处理中（投放后） → 输入中（首次流式更新后） → 完成（isFinalize=true）
                                             → 失败（isError=true）
```
