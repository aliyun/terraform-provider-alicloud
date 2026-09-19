# ArenaAI 页面操作与恢复

## 页面基线

先从当前页面可见状态确认控件，不盲目依赖历史选择器。当前 ArenaAI 页面通常可通过以下线索识别：

- 输入框：`textarea[placeholder="输入问题，即刻开始"]`
- 发送按钮：`button[class*="submitBtn"]`
- 用户消息：`.user-message`（注意：虚拟化下常被卸载，不可靠，见"幂等发送"与"状态定义"）
- 对战回答：`.battle_assistant`

站点更新后选择器可能变化。若选择器失效，先读取可见 DOM、标签、占位符和按钮状态，再做最小调整；不要通过读取站点存储或注入到无关页面来恢复。

## Chrome 控制顺序

1. 使用 `chrome:control-chrome` 提供的浏览器客户端完成选择、命名会话、导航和交互。
2. 如果 ArenaAI SPA 的 DOM/交互调用持续超时，先读取 Chrome 技能的故障排查文档，并使用它列出的替代能力。
3. 只使用当前 `chrome:control-chrome` 技能明确允许的浏览器控制接口，不切换到 Apple Events、独立自动化服务器或其他浏览器技能。
4. 若允许的接口无法创建独立 Chrome 窗口，向用户说明限制并请用户创建或激活一个新窗口后继续；不要静默换成内置浏览器或改造无关窗口。

## CSP 约束与原生操作（关键）

ArenaAI 页面的 Content Security Policy **不含 `'unsafe-eval'`**，因此任何以 `evaluate`/`eval`/`new Function` 形式执行的任意 JavaScript（包括"原生 setter + `dispatchEvent`"片段）**都会被 CSP 拦截，无法运行**。不要试图用 `evaluate` 注入脚本来填值或点击——这是已验证的死路。

> **例外：CDP 级 evaluate 不受 CSP 限制。** 上述封锁针对**页面上下文**的 eval（`chrome:control-chrome` 的 page eval、playwright 的 `page.evaluate`）。若改用 [AgentBridge 通路](agentbridge-method.md)（桌面扩展的 `chrome.debugger.attach`），其 `Runtime.evaluate` 是 CDP/devtools 协议层，**不受页面 CSP 限制**——即使页面禁 `unsafe-eval` 仍可执行任意表达式。**但**：仍应使用原生 fill（native setter + input/change events）与 `press Enter`，原因从 CSP 变为 React 受控组件兼容性，不是"evaluate 能跑了就改用 JS 点击"。AgentBridge 通路是 `chrome:control-chrome` 受阻（Chrome 153 封锁 `--remote-debugging-port` 或 relay 拥塞）时的备选，默认仍走本节的原生命令。

改用浏览器控制接口提供的**原生 DOM 一等命令**：这些命令在内容脚本层直接操作 DOM，绕过页面 CSP。经多轮 × 15 页实战验证（CSP 下稳定 15/15 提交）的命令与参数：

| 操作 | 命令 | 关键参数 |
|---|---|---|
| 填入文本 | `fill` | `strategy: native_setter`、`commit: none`、`clear: true` |
| 提交发送 | `press` | `key: Enter`、`strategy: auto`、`selector: textarea` |
| 读取文本 | `get_text` | `scope: element`（单元素）或 `scope: full`（全页），`maxChars` 限长 |

`press` 的 `strategy: auto` 会在调试器附加时优先用可信 `cdp_keyboard`，否则退回合成 `dom_keyboard`；两者均能触发 React 的 onSubmit。

## 稳定台账

不要仅使用"第几个标签页"作为身份，因为标签页会重排。为每个目标页维护：

```text
stable_tab_id
question_index
question_text
normalized_question_hash
pre_submit_user_message_count
pre_submit_answer_count
state
submit_method
verify
retry_count
attempt_started_at
```

每次重载或重新绑定标签页后，先用稳定 ID 和可见 URL/标题校验目标；无法证明是原目标页时停止，不猜测 ID。**同一窗口内 tabId 在多轮之间稳定不变**——下一轮直接从上一轮台账读取 tabId，不要重新发现。

## 跨轮台账纪律

每轮提交完成后，把该轮台账落盘为 `/tmp/arena_ledger_rN.json`（`N` 为轮次号），至少包含：`qidx`、`tabId`、`q`（问题全文）、`hash`、`theme`、`state`、`submit_method`、`verify`、`retry`。

生成**下一轮**问题时，必须先 dump **所有历史轮次**的台账（R1..R(N-1)），按 tabId 分组，逐个 tab 检查新问题的角度是否与**该 tab 上所有历史轮次**的问题都不同——不仅是与上一轮不同。"不同"指核心对象、决策目标和约束条件三者之一不同，而非仅改写措辞。tabId 从历史台账读取，不重新发现。

## 填入与提交（经验证序列）

按以下顺序执行，不要省略 sleep：

1. **填入**：`fill`，参数 `strategy: native_setter`、`commit: none`、`clear: true`，selector 指向 textarea。
2. **等待 React flush**：sleep 2.5 秒，让受控组件把 value 提交到内部状态。
3. **提交**：`press`，参数 `key: Enter`、`strategy: auto`、`selector: textarea`。
4. **等待提交落定**：sleep 3.0 秒。
5. **核验提交**：`get_text`，`scope: element`、`maxChars: 200`，selector 指向 textarea。

**提交成功判据（主信号，稳健，不受虚拟化影响）**：textarea 被清空回占位符 `输入问题，即刻开始`，或为空。占位符 vs 空字符串是 `get_text` 的读取差异，**不**代表提交失败。

**用户消息匹配（旁证，虚拟化脆弱）**：`.user-message` 出现与分配问题精确匹配的文本。ArenaAI 会**虚拟化**对话流，`get_text scope=full` 只读到当前挂载的视口切片，用户消息常被卸载而不可见；因此用户消息匹配**仅在可达时作确认**，**不可作为提交成功的必要条件**。主信号（textarea 清空）已足够。

## 幂等发送

发送前按以下顺序判断：

1. `get_text` 读 textarea（scope=element）：若内容为占位符或空，且页面已有该问题的用户消息（可达时）→ 视为已提交，不点击。
2. textarea 内容与分配问题精确匹配、按钮启用：可以 `press` 提交一次。
3. textarea 已清空但没有可见用户消息：**以 textarea 清空为主信号视为已提交**，不推断失败、不重发（用户消息可能被虚拟化卸载）。
4. textarea 非空但按钮禁用：等待一次页面状态更新；仍禁用则诊断页面，不强制点击。
5. 没有 textarea 或按钮：页面尚未就绪或 UI 已变化，重新检查可见状态。

不要因为一批操作超时就再次点击整批。超时可能发生在点击已经生效之后。

## 状态定义

为每页记录以下状态：

- `prepared`：`fill` 后 textarea 内容与分配问题完全一致，发送按钮可用。
- `submitted`：**textarea 被清空回占位符或为空**（主信号）；若用户消息可达且精确匹配，作确认旁证。
- `generating`：已提交，页面仍显示停止生成控件或"生成中，请等待"等明显状态。
- `completed`：已提交，回答基线之后新增至少一个回答容器且含有意义正文，且生成控件已消失或恢复为普通禁用发送按钮。**稳健旁证**：真实模型名出现（流式生成期间模型匿名显示"A 模型"/"B 模型"，对战结束后才显名）+ 投票按钮启用。
- `failed`：出现 `Failed to fetch`、明确错误提示、空白结果且不再生成，或长时间无进展。

"有意义的正文"应明显超过模型标签或占位符。可将 100 个可见字符作为普通问答的保守下限，但应结合页面内容判断；不能仅凭回答容器存在、投票按钮出现或 `A/B` 标签出现认定成功。

生成中的按钮目前可能含有 `Stop Loading` 标题；完成后按钮通常回到普通发送图标并因输入为空而禁用。站点变化时以可见语义为准。

## 轮询与恢复（含虚拟化噪声与稳健轮询）

**虚拟化噪声地板**：ArenaAI 虚拟化对话，`get_text scope=full` 的 `caps.totalTextChars` 只统计当前挂载视口切片，因此**即使已停止生成，连续两次读取也会有 ±10–30 字符的抖动**，且随对话轮次累积（8 轮 Q&A 后更明显）。真实模型流式生成速度为 10–50+ 字符/秒；**低于约 1 字符/秒的增量是虚拟化抖动，不是活跃生成**，尤其出现在长对话页面上时。

- 正常轮询间隔为 15–30 秒。读取 `caps.totalTextChars`；连续两次相等视为字面稳定。
- 但对长对话页面，字面稳定可能被噪声掩盖：若增量 < ~1 字符/秒且持续，按"噪声内、实际已完成"处理，并用上述稳健旁证（真实模型名 + 投票按钮）交叉确认。
- 生成超过 90 秒且连续两次长度不增长（或增量落在噪声地板内），可视为完成/卡住，以旁证区分。
- 每次提交或恢复尝试设置 6 分钟墙钟期限；即使文字仍增长，到期也停止无限等待，保留页面并报告状态。

**稳健轮询（防 daemon 拥塞）**：15 个对战同时流式时，浏览器控制 daemon 会变慢，`get_text` 在 30 秒超时下易失败并**整批崩溃**。轮询脚本必须：

- 单命令超时设为 **60 秒**；
- 捕获 `TimeoutExpired` 异常，返回哨兵值（如 `-2`）标记"不明确"，**不向上抛出**；
- 每个超时 tab 重试一次；
- 把超时 tab 标记为 `unclear` 而非 `failed`，不中断整批轮询。

- 若同一批中至少 3 页或 20% 页面出现相同的鉴权、限流或服务错误，立即暂停新的发送和重试，等待 60 秒后复查一次；仍为系统性错误时停止批量重放并报告。
- `Failed to fetch`：等待其他页面完成以降低负载，然后在原页重新 `fill` 同一问题并重发一次。
- 空白结果、发送按钮无法重新启用或生成卡死：只重载该页，等待输入框恢复，再 `fill` 原问题并重发。
- 重试前必须检查是否已经提交（textarea 清空）或已有有效回答，避免重复。
- 每页最多 2 次恢复性重试；超过上限后保留页面并向用户报告，不循环刷新。

## 最终验收

最后一次全量扫描必须重新读取每个目标页面，而不是依赖早先的点击返回值。至少汇总：

```text
目标页面数
唯一分配问题数
已提交数（textarea 清空为主信号）
含有效回答数
已停止生成数（字面稳定 或 增量在噪声地板内 + 旁证）
待处理页码
发生过恢复性重试的页码
不明确页码（轮询超时/噪声未收敛）
```

核心计数等于目标页面数且待处理为空，才能报告全部成功；存在不明确页时如实标注（例如"字面稳定 N/15，其余在噪声地板内、实际已完成"），不夸大也不漏报。
