---
name: screenshot-evidence
description: 可视化查证截图取证——aone-triage / terraform-pr-review / provider-resource-dev 查证后需要浏览器截图证据时触发：Playwright 元素级截图 → OSS private + 签名 URL → 组装 HTML 对比报告 → pre-agent 预览回贴 Aone。NOT for 纯文字查证或无需可视化证据的场景。
---

# screenshot-evidence — 可视化查证截图取证

> 在 aone-triage 查证环节补充浏览器截图证据，让人工审核一目了然。

## 触发条件

以下场景应在文字查证之外**追加截图取证**：

- aone-triage 查证阶段（证据 1-4 层完成后）
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
bash .Codex/skills/screenshot-evidence/scripts/upload-screenshots.sh <aone-id> <screenshot-dir>
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
bash bootstrap/html-report-preview.sh upload <aone-id> <report.html> --comment
```

返回预览 URL 并自动发 Aone 评论。

### Step 6: 更新 Aone 工单详情

用 `bin/a1id -- project workitem update <id> --body-file <path>` 在 description 中：
- 添加「可视化查证报告」章节，含超链接指向在线预览
- 证据来源用超链接（description 支持 markdown 链接渲染）
- **不要在 description 中嵌入 `<img>` 标签**——Aone 渲染器剥离 img src 的 query 参数

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
查证（文字）→ 截图取证（本 skill）→ 组装报告 → 上传预览 → 更新工单详情
```

aone-triage 的 wrap.sh done 草稿中增加一行（评论区仅 markdown `[text](url)` 可点，裸 URL 不 autolink，见 aone-triage SKILL §4 渲染 quirk）：

```
📊 可视化查证报告：[在线查看](<preview-url>)
```
