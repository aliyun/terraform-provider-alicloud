# 模型分层（两条车道 / model lanes）

Jarvis headless 派发按工单类型选不同模型后端，兼顾质量与成本。

## 车道划分

| 车道 | 判定 | 档链（env） | 后端 |
|------|------|-------------|------|
| **Terraform 线** | `_is_terraform_ticket(pool, title)`：池 line=`terraform_provider`（`tf_customer` 1086837 / `tf_provider` 528766）**或**标题命中 `alicloud_` / `terraform-provider` / `tf-provider` 等 | `JARVIS_SETTINGS_TF` | ideamo → ideamore → glm5.2（兜底） |
| **其他工作**（默认） | 其余全部 | `JARVIS_SETTINGS` | glm5.2（百炼 key） |

- `JARVIS_SETTINGS_TF` **未设时自动回退** `JARVIS_SETTINGS`——分层是 opt-in；只配一条 `JARVIS_SETTINGS` 时两条车道等价旧版单链，向后兼容。
- `JARVIS_CC`（重型启动脚本）设了则**全覆盖、不分车道**。
- Tata 门面（`JARVIS_TATA_SETTINGS`）属「其他工作」，独立配置；如需与其他线同走 glm5.2，指到同一档即可。
- qwen3.7 全线退役，glm5.2 成为两条车道的统一兜底。

## 车道如何随会话持久化（resume 正确性）

`jarvis_cmd` 按 `md5(session_id) % pool` 在**选中车道内**粘网关。若 resume 走错车道会落到另一网关/token。因此车道随会话外化：

- **Task Session** 的不可变 `inputPayload.terraform` 是可恢复任务的唯一车道真源；换机、重启或重新 lease 后仍使用同一车道。
- **in-flight 记录**（`.my-day/bridge/inflight.json`）只服务本地 EphemeralJob 的进程诊断，不参与 Task 恢复。
- **suspend 记录**（`.my-day/suspended/*.json`）带 `terraform` 字段 → 挂起唤醒 `_wake` 复原车道。
- Task Session 缺少不可变 `inputPayload` 时执行器 fail-closed，不读取当前 Task payload 猜测恢复输入。

派发点判定车道（均在 `bridge/jarvis_dingtalk_bot.py`）：`_dispatch` / 手动授权 `处理 #id` / `全部处理` / revisit / persona → 按 `_is_terraform_ticket`；probe（tf-probe）→ 恒 `terraform=True`；Tata 委派 handoff-exec → 默认车道。

---

## 后端接入

| 车道成员 | 接入方式 | 后端 |
|----------|----------|------|
| ideamo / ideamore | idea-cc-fix 网关 `:9090`（`idealab.alibaba-inc.com`，直通） | IdeaLab |
| **glm5.2** | **直连百炼 anthropic 端点，不经任何网关** | 百炼(DashScope) `glm-5.2-fast-preview` |
| ~~qwen3.7~~ | ~~idea-cc-fix-qwen `:9091`~~ | 退役 |

> **glm 不走 idea 网关**：百炼的 anthropic-compatible 端点原生说 Anthropic 协议、直接认 `glm-5.2-fast-preview` 这个 model id，claude 客户端把 `ANTHROPIC_BASE_URL` 指过去即可，**无需 idea-cc-fix 网关 / 端口 / model_map / agent-sandbox 改动**。ideamo/ideamore 仍走 9090（IdeaLab 用 claude 别名，需网关翻译）；两者在 failover 档链里各自独立指向自己的后端，互不影响。

### 已验证（本机实测，非推断）

对 `https://dashscope.aliyuncs.com/apps/anthropic` + `BAILIAN_KEY` + `glm-5.2-fast-preview` 全链路打通：

| 路径 | 结果 |
|------|------|
| `POST /v1/messages`（Bearer / x-api-key 两种鉴权） | 200，原生 Anthropic message 格式 |
| 流式 `stream:true` | 200，正确 SSE（`message_start`/`content_block_*`）——`run_claude_stream` 依赖 |
| 1-token 健康探测（`max_tokens:1`, content `"."`） | 200——`_probe_settings` failover 探活依赖 |
| `POST /v1/messages/count_tokens` | 200 `{"input_tokens":N}` |
| **真实 `claude -p` headless 跑一轮** | `subtype:success`, `result:"pong"`, `modelUsage.glm-5.2-fast-preview` |

---

## settings 档（bridge 宿主机 `~/.claude/`）

每个「档」是一份 Claude Code settings JSON。**ideamo / ideamore 档不变**（走 9090 IdeaLab）。新增 glm5.2 档——直连百炼：

`~/.claude/glm5.2.json`：
```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://dashscope.aliyuncs.com/apps/anthropic",
    "ANTHROPIC_MODEL": "glm-5.2-fast-preview",
    "ANTHROPIC_SMALL_FAST_MODEL": "glm-5.2-fast-preview",
    "ANTHROPIC_AUTH_TOKEN": "<粘 ~/.zshrc 里 BAILIAN_KEY 的值——切勿提交进任何仓库>"
  }
}
```

> - `ANTHROPIC_BASE_URL` 到 `.../apps/anthropic` 的**根**（不带 `/v1`）——claude CLI 与 `_probe_settings` 都会自动补 `/v1/messages`。
> - `ANTHROPIC_MODEL` = 百炼真实 model id `glm-5.2-fast-preview`（无需 model_map）；
>   `ANTHROPIC_SMALL_FAST_MODEL` 也设成它，否则 claude 后台/haiku 调用会发默认 claude-haiku id 被百炼拒。
> - `ANTHROPIC_AUTH_TOKEN` = 百炼 key。该档是**唯一**持有明文 key 的地方，`~/.claude/` 不入库。

`ideamo.json` / `ideamore.json` 保持现状（示意）：
```json
{ "env": { "ANTHROPIC_BASE_URL": "http://127.0.0.1:9090", "ANTHROPIC_MODEL": "<ideamo 模型 id>", "ANTHROPIC_AUTH_TOKEN": "<IDEALAB_API_KEY>" } }
```

`bridge/jarvis.env` 接线：
```sh
# 其他工作 → glm5.2（直连百炼）
JARVIS_SETTINGS=~/.claude/glm5.2.json
# Terraform 线 → ideamo → ideamore → glm5.2 兜底
JARVIS_SETTINGS_TF=~/.claude/ideamo.json,~/.claude/ideamore.json,~/.claude/glm5.2.json
```

---

## 宿主机落地 runbook（`jarvis@<bridge-host>`）

> 密码不入库、不入命令行日志；SSH 由仓库主人本人执行，或改配公钥免密后再交由自动化。
> 百炼 key = `~/.zshrc` 的 `BAILIAN_KEY`（`sk-…`），只手动粘进宿主机 claude 档，切勿写入任何仓库文件。

1. **profile**：新建 `~/.claude/glm5.2.json`（上），`ANTHROPIC_AUTH_TOKEN` 填 `BAILIAN_KEY` 的值；确认 `ideamo.json`/`ideamore.json` 仍健康。**无需动 agent-sandbox / idea-cc-fix**。
2. **接线**：`bridge/jarvis.env` 写入两条 `JARVIS_SETTINGS*`（上）。
3. **重启桥**：`bridge/run.sh restart`。
4. **冒烟**：
   - 非 Terraform 单：确认落 glm5.2（百炼）。
   - Terraform 单：确认落 ideamo（主档探活通过时）；停掉 9090 验证自动顶到 ideamore→glm5.2。
   - resume：挂起一条 terraform 单再唤醒，确认仍落 terraform 车道（`.my-day/suspended/<id>.json` 里 `terraform:true`）。
5. **退役 qwen3.7**：确认无 env 再引用 qwen 档后，agent-sandbox 可摘 `idea-cc-fix-qwen`（或先保留观察一轮再摘）。

## 相关代码/配置

- 车道选择：`bridge/jarvis_dingtalk_bot.py` `jarvis_cmd()`
- 判定：`_is_terraform_ticket()` / 池 line：`config/pools.json`
- env 契约：`bridge/jarvis.env.example`
- 测试：`test/bridge_dispatch_test.sh`（`SettingsLaneTest` / `DispatchLanePersistTest`）
