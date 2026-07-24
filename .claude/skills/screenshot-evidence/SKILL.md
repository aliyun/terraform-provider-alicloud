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

- **Playwright MCP** 已连接（`mcp__playwright__*` 工具可用）
- **aliyun CLI** 已认证（`aliyun oss` 可用）
- **OSS bucket**: `jarvis-report-images`（cn-hangzhou, private）
- **JARVIS_HTML_REPORT_TOKEN** 已设置（pre-agent 上传）
- **html-report-preview.sh** 内建 `X-Request-Context` WAF header（env `JARVIS_HTML_REPORT_WAF_HEADER` 可覆盖），无需额外处理

## 取证流程

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

### Step 2: Playwright 截图

用 `mcp__playwright__browser_run_code_unsafe` 定位元素截图（绕过滚动问题）：

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
`missing_capability`。这样 finalizer 能区分“不适用”和“漏做”。

finalizer 上传前必须执行确定性校验：

```bash
python3 .claude/skills/screenshot-evidence/scripts/validate-manifest.py \
  .my-day/screenshots/<aone-id>/evidence-manifest.md
```

校验器要求三层各一行；pass/fail 的截图必须是实际存在的绝对路径，n-a 必须使用 `N/A`
并给出原因。校验失败即 blocked/missing_capability，不得继续生成“完整证据”报告。

### Step 3: 上传 OSS + 生成签名 URL

```bash
BUCKET="oss://jarvis-report-images"
ENDPOINT="oss-cn-hangzhou.aliyuncs.com"
PREFIX="reports/<aone-id>"
TIMEOUT=15768000  # 6 个月

# 上传（private ACL）
aliyun oss cp <file>.png "${BUCKET}/${PREFIX}/<file>.png" --acl private -e "$ENDPOINT" -f

# 生成签名 URL
aliyun oss sign "${BUCKET}/${PREFIX}/<file>.png" --timeout "$TIMEOUT" -e "$ENDPOINT"
```

**辅助脚本**批量处理上传 + 签名：

```bash
bash .claude/skills/screenshot-evidence/scripts/upload-screenshots.sh <aone-id> <screenshot-dir>
# 输出: name|signed_url 行打到 stdout(每文件一行,需要落盘自行重定向)
```

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
| 账号级 Block Public Access | OSS 对象无法 public-read | 必须用签名 URL |
| pre-agent WAF 拦截 base64 data URI | HTML 中不能内嵌图片 | 图片走 OSS 签名 URL |

## OSS Bucket 规范

- **Bucket**: `jarvis-report-images`（cn-hangzhou）
- **ACL**: 对象级 private
- **路径规则**: `reports/<aone-id>/<screenshot-name>.png`
- **签名有效期**: 6 个月（`15768000` 秒）
- **禁止**: public-read（账号策略阻止 + skill 规范禁止）
- **禁止**: 使用个人 AKSK 或其他 bucket

## 与 aone-triage 集成

在 aone-triage 查证阶段（SKILL.md §3）完成后，追加调用本 skill：

```
非 Terraform：查证（文字）→ 截图取证 → 组装报告 → 上传预览并回贴 → 可选更新详情
Terraform：PD 三层查证 → 本地截图 + manifest → RD finalizer 校验并上传（无 --comment）→ 唯一聚合回复
```

aone-triage 的 wrap.sh done 草稿中增加一行（评论区仅 markdown `[text](url)` 可点，裸 URL 不 autolink，见 aone-triage SKILL §4 渲染 quirk）：

```
📊 可视化查证报告：[在线查看](<preview-url>)
```
