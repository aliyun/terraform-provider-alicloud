---
name: screenshot-evidence
description: 可视化查证截图取证——aone-triage / terraform-pr-review / provider-resource-dev 查证后的浏览器截图证据流程。Terraform Aone 三层查证时必须使用：PD 生成本地 OpenAPI、CloudSpec/ACube、Provider 截图与 manifest，最终 RD 统一上传报告并只在唯一聚合回复中贴链接。其它需要可视化证据时也使用；纯文字且无需人工复核的场景不触发。
---

# screenshot-evidence — 可视化查证截图取证

> 在 aone-triage 查证环节补充浏览器截图证据，让人工审核一目了然。

## 触发条件

以下场景应在文字查证之外**追加截图取证**：

- Terraform aone-triage 三层查证：**强制**生成 OpenAPI、CloudSpec/ACube 映射、Provider 源码三层截图或逐层 N/A 说明
- 非 Terraform aone-triage 查证阶段：存在需要人工复核的网页、字段或对比结果
- terraform-pr-review 需要对比 OpenAPI 文档与 PR diff 时
- provider-resource-dev 需要展示当前文档 bug 与修复对比时
- 任何需要在 Aone 工单/评论中提供可视化证据的场景

## 前置条件

- **截图通道（任一即可，按优先级探测）**：
  1. **Playwright MCP** 已连接（`mcp__playwright__*` 工具可用）——交互会话常见
  2. **仓库内 `capture.sh`** 通道可用（Playwright Python 绑定 + 本地 chromium，或 headless Chrome/Chromium 二进制）——headless 必走这条，因为交互态 Playwright MCP **不会**注入进 bridge 拉起的 headless 会话（根因见 `references/headless-screenshot-channels.md`）
  - 两者都不可用 → 逐层写 `n-a + missing_capability`，继续文字查证和最终评论，**绝不静默跳过截图或因此挂起任务**（见 Step 0）
- **JARVIS_HTML_REPORT_TOKEN** 已通过运行时配置注入（截图与 HTML 都走 server-token）
- **JARVIS_HTML_REPORT_BASE_URL** 可选；默认 `https://pre-agent.aliyun-inc.com`
- 截图存储由 AutomationAgent 服务端负责：私有 bucket `jarvis-upload-files`，owner
  account `1983056807138283`
- **html-report-preview.sh** 内建 `X-Request-Context` WAF header（env `JARVIS_HTML_REPORT_WAF_HEADER` 可覆盖），无需额外处理

## 取证流程

### Step 0: 能力探测（执行前必做）

在截图**之前**先探测通道，避免运行到中途才发现无浏览器而失败：

```bash
bash .claude/skills/screenshot-evidence/scripts/capture.sh probe
# stdout: <channel-name>          → exit 0，通道可用（playwright_python / chrome_binary）
# stdout: missing_capability: ... → exit 3，无可用通道
```

- 探测命中任一通道 → 记下通道名，Step 2 用对应方式截图。
- 探测为 `missing_capability`（exit 3）→ **不静默跳过**：在 manifest（Step 2.1）把该层写 `n-a`，
  `note` 写 `missing_capability: <probe 给出的原因>`，继续完成文字查证。截图属于证据增强项，
  不得仅因缺少截图把 finalizer 标为 blocked/missing_capability，也不得输出 SUSPEND。
  `validate-manifest.py` 接受带原因的 `n-a` 行，finalizer 在最终聚合评论中说明截图已降级。
- 探测只读、幂等、无副作用，每个 headless run 开头调用一次后可复用结论。
- 交互会话若 `mcp__playwright__*` 可用，可跳过本步直接走 Step 2 的 Playwright 路径；headless 必跑。

### Step 1: 确定截图目标

根据查证对象选择截图页面：

| 证据层 | 目标 URL | 截取内容 |
|--------|---------|---------|
| OpenAPI 定义 | `https://next.api.alibabacloud.com/document/{Product}/{Version}/{Action}` | 参数表中目标字段行（含 Valid Values） |
| CloudSpec/ACube 映射 | AMP 资源元数据页、ACube mapping 响应页，或将只读查询结果渲染为本地 HTML | resourceTypeCode、Terraform resource type、映射状态与关键字段 |
| Provider 文档（当前） | `https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/{resource_name}` | 目标字段的描述和 Valid values |
| 兄弟资源文档（对比） | 同上，替换为兄弟资源名 | 同字段的正确描述 |
| GitHub PR diff | `https://github.com/aliyun/terraform-provider-alicloud/pull/{N}/files` | Files changed 标签页的 diff |
| Provider 源码 | `https://github.com/aliyun/terraform-provider-alicloud/blob/master/alicloud/{file}.go` | 目标字段的 schema 定义 |

### Step 2: 截图（交互用 Playwright MCP，headless 用 capture.sh）

交互会话用 `mcp__playwright__browser_run_code_unsafe` 定位元素截图（绕过滚动问题）：

```javascript
async (page) => {
  await page.goto('<URL>', {waitUntil: 'domcontentloaded', timeout: 60000});
  await page.waitForTimeout(3000);
  
  // 定位包含目标文本的元素
  const items = await page.$$('li, td, tr');
  for (const item of items) {
    const text = await item.innerText();
    if (text.includes('<目标字段名>')) {
      await item.screenshot({path: '.my-day/screenshots/<aone-id>/<name>.png'});
      return 'Captured';
    }
  }
  return 'Not found';
}
```

**关键规则：**
- 截图保存到 `.my-day/screenshots/<aone-id>/`（gitignored）
- **不压缩**：保留原始 PNG 质量
- OpenAPI 页面是 SPA，需等待渲染完成
- Terraform Registry 有右侧 sidebar 遮挡，用元素级截图（`element.screenshot()`）而非 viewport

#### headless / 无 Playwright MCP：走 capture.sh

headless bridge 会话**没有** `mcp__playwright__*`（根因见 `references/headless-screenshot-channels.md`）。
当 Step 0 探测返回 `chrome_binary` 或 `playwright_python` 时，用仓库内 `capture.sh` 全页截图：

```bash
bash .claude/skills/screenshot-evidence/scripts/capture.sh capture \
  "<URL>" .my-day/screenshots/<aone-id>/<name>.png \
  --wait 3000 --full-page --width 1280 --height 2000 \
  --text "<目标字段名>"
# stdout: <channel-name>；exit 0 成功，exit 3 missing_capability，exit 1 capture_error
```

- `capture.sh` 是仓库内可控、可降级通道（Playwright Python + 本地 chromium 优先，否则 headless Chrome/Chromium 二进制），不依赖交互态 MCP。
- 两个通道都支持 `--text` 定位 `li/td/tr` 中包含目标文本的元素并截图；未命中目标文本返回
  `capture_error`，不得用无关全页图冒充元素证据。不需要元素定位时去掉 `--text`。
- `chrome_binary` 通过 Chrome DevTools Protocol 获取页面内容尺寸并执行
  `Page.captureScreenshot(captureBeyondViewport=true)`，因此 `--full-page` 是真全页截图，不再受
  `--height` viewport 高度限制。超长页面命中安全上限时改用 `--text`。
- SPA 等渲染需 `--wait`（毫秒）让其稳定。
- `chrome_binary` 对启动、CDP 就绪和页面捕获类瞬时错误自动使用全新 profile 最多尝试 3 次；
  可用 `JARVIS_SCREENSHOT_CHROME_ATTEMPTS`（1–5）调整。
- exit 1（capture_error）或 exit 3（missing_capability）→ 该层 manifest 写 `n-a` 并记录
  原始分类与脱敏原因，然后继续文字查证和最终评论；不得把 exit 1 改写为“无浏览器通道”。

### Step 2.1: Terraform PD 本地交接

Terraform PD 只生成本地文件，不执行 OSS、pre-agent 或 Aone 写入。在
`.my-day/screenshots/<aone-id>/evidence-manifest.md` 记录：

```markdown
| layer | result | screenshot | source | note |
|---|---|---|---|---|
| OpenAPI | pass/fail/n-a | <本地绝对路径或 N/A> | <来源 URL/命令> | <字段结论或 N/A 原因> |
| CloudSpec/ACube | pass/fail/n-a | <本地绝对路径或 N/A> | <来源 URL/命令> | <映射结论或 N/A 原因> |
| Provider | pass/fail/n-a | <本地绝对路径或 N/A> | <来源 URL/文件:行> | <实现结论或 N/A 原因> |
```

返回结构中增加：

```yaml
visual_evidence_manifest: /absolute/path/.my-day/screenshots/<aone-id>/evidence-manifest.md
```

浏览器能力、登录态或页面不可达时，不要把该层悄悄删掉；保留 `n-a` 行并写清
`missing_capability` 或 `capture_error`。这样 finalizer 能区分“不适用”“截图降级”和“漏做”；
只要文字查证与业务验收已完成，截图降级不改变业务 outcome。

finalizer 上传前必须执行确定性校验：

```bash
python3 .claude/skills/screenshot-evidence/scripts/validate-manifest.py \
  .my-day/screenshots/<aone-id>/evidence-manifest.md
```

校验器要求三层各一行；pass/fail 的截图必须是实际存在的绝对路径，n-a 必须使用 `N/A`
并给出原因。校验失败时不上传报告，在最终聚合评论中写明“截图报告降级：manifest 无效”，
继续输出文字查证结果；不得因此阻断评论或挂起任务。

### Step 3: 通过服务端上传私有图片并获取签名 URL

客户端不接触 OSS 凭据。辅助脚本先 source `bootstrap/runtime-config.sh` 并调用
`jarvis_load_runtime_config`，再使用 `JARVIS_HTML_REPORT_TOKEN` 逐张调用：

```text
POST /api/reports/aone/<aone-id>/images
Authorization: Bearer <server-token>
Content-Type: multipart/form-data
file=@<screenshot>
```

AutomationAgent 服务端负责把对象写入 `jarvis-upload-files`。服务端必须在首次实际 OSS
操作前调用 STS 校验当前账号严格等于 `1983056807138283`；账号不符、无法取得身份、
bucket region/endpoint 无法发现、private 写入失败或签名失败时都要 fail-closed。对象保持
private，只向客户端返回限时 signed GET URL。

**辅助脚本**批量上传 PNG/JPG：

```bash
bash .claude/skills/screenshot-evidence/scripts/upload-screenshots.sh <aone-id> <screenshot-dir>
# 输出：name|signed_url（stdout 每文件一行，需要落盘时自行重定向）
```

脚本错误只输出固定脱敏错误码，不回显 token 或服务端响应。缺 token 返回 exit code 3；
任一文件上传失败则不输出任何部分结果。禁止绕过该接口在本地配置、解密或调用个人
AK/SK，也禁止把任何明文/密文凭据放进命令参数、日志、报告或工单。

### Step 4: 组装 HTML 报告

用 Python 组装 HTML，要点：

1. **图片用 OSS 签名 URL**（`<img src="<signed-url>">`），不用 base64（WAF 拦截）
2. **来源文字加超链接**：如 `来源：<a href="...">next.api.alibabacloud.com · Cloudfw / CreateNatFirewallControlPolicy</a>`
3. **对比用 grid 布局**：左红（❌ 错误）右绿（✅ 正确）
4. **附交叉验证表**：验证层 / 结果 / 证据链接

报告模板见 `references/report-template.html`（可变量替换）。

### Step 5: 上传到 pre-agent 预览

```bash
# 非 Terraform：允许当前处理者直接评论
bash bootstrap/html-report-preview.sh upload <aone-id> <report.html> --comment

# Terraform finalizer：只取预览 URL，禁止产生第二条 Aone 评论
bash bootstrap/html-report-preview.sh upload <aone-id> <report.html>
```

Terraform finalizer 必须先校验 manifest 中三层都存在（或有明确 N/A 原因），再上传截图和报告。
存在 `n-a` 时可以上传包含降级说明的部分报告；若截图、manifest 或上传链路失败，则跳过报告，
把失败分类和原因写入唯一聚合评论，继续按文字证据收口业务结论。
上传命令不得传 `--comment`；将返回的预览 URL 写进唯一聚合回复。executor 托管的 headless
run 必须放入 `AONE_RESULT.reply_body`，由 executor 单次落账，run 内不得调用 `wrap.sh`；
仅非 executor 托管的独立 finalizer 才按 bookend 调用一次 `wrap.sh done`。
非 Terraform 流程可保留 `--comment` 的端到端行为。

### Step 6: 更新 Aone 工单详情

以下 description 更新只适用于非 Terraform 流程：

用 `bin/a1id -- project workitem update <id> --body-file <path>` 在 description 中：
- 添加「可视化查证报告」章节，含超链接指向在线预览
- 证据来源用超链接（description 支持 markdown 链接渲染）
- **不要在 description 中嵌入 `<img>` 标签**——Aone 渲染器剥离 img src 的 query 参数

Terraform 主处理 run 不单独更新 description，也不调用 `--comment`；finalizer 把报告链接放入
唯一聚合回复，避免截图能力绕过单写者边界。

## 平台限制（已验证）

| 限制 | 影响 | 应对 |
|------|------|------|
| Aone 评论区不 autolink 裸 URL、不渲染 `<a href>` | 裸 URL/HTML 锚都是不可点的死文本 | 评论与详情一律用 markdown `[text](url)`（唯一可点格式，84307546 评论 124870464 实测） |
| Aone 渲染器剥离 `<img src>` 的 query 参数 | OSS 签名 URL 失效 → 403 | 图片只在 pre-agent 在线报告中展示 |
| 账号级 Block Public Access | OSS 对象无法 public-read | 服务端保持 private 并返回签名 URL |
| pre-agent WAF 拦截 base64 data URI | HTML 中不能内嵌图片 | 图片走 OSS 签名 URL |

## OSS Bucket 规范

- **服务端唯一 Bucket**: `jarvis-upload-files`
- **Owner account**: `1983056807138283`（STS 必须精确匹配，否则 fail-closed）
- **Region/endpoint**: 服务端从 bucket 元数据发现，不在客户端硬编码
- **ACL**: bucket 与对象均为 private
- **服务端对象前缀**: `reports/aone/<aone-id>/images/<uuid>-<original-filename>`
- **访问方式**: 服务端返回限时 signed GET URL；无签名直链不得公开可读
- **禁止**: public-read、其他 bucket、本地 OSS 凭据、个人 AK/SK、把凭据传给 agent

## 与 aone-triage 集成

在 aone-triage 查证阶段（SKILL.md §3）完成后，追加调用本 skill：

```
非 Terraform：查证（文字）→ 截图取证 → 组装报告 → 上传预览并回贴 → 可选更新详情
Terraform：PD 三层查证 → 尽力截图 + 完整 manifest → RD finalizer 校验并尽力上传（无 --comment）→ 唯一聚合回复；截图失败只降级报告，不阻断评论
```

aone-triage 的 wrap.sh done 草稿中增加一行（评论区仅 markdown `[text](url)` 可点，裸 URL 不 autolink，见 aone-triage SKILL §4 渲染 quirk）：

```
📊 可视化查证报告：[在线查看](<preview-url>)
```
