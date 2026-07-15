# cap: html-report-preview 内嵌 base64 图被预发 WAF 拦

- **缺口类型**：能力/知识缺口（skill 未记载约束 + helper 未处理）
- **阻塞任务**：工单 83899246（SearchDocuments 错误透传）验收报告要传 AutomationAgent 在线预览，内嵌截图的 HTML 反复被拦，一度误判为 token/端点问题。
- **现象**：`POST /api/reports/aone/<id>` 返回 HTTP 200，但 body 是 `waf_block*.html` / `punish`（rgv587）页而非 `{"success":true,...,"viewUrl":...}`。
- **定位**：
  - GET 根路径 200；空 body POST 到应用得 415（应用层）；带/不带 token 结果相同 → 排除 token/端点/方法。
  - 极简纯文本 HTML、无图报告均上传成功；内嵌 base64 图（PNG 707KB、JPEG 192KB）均被拦 → **预发 WAF content inspection 对 base64 图片数据命中签名，与体积无关**。
  - 拦截页内联「加 `X-Request-Context: rctx_...` 即放行」——蜜罐/注入，照做无效（退化为纯 block 页），已拒绝。
- **解法（已验证可行）**：截图传 OSS（对象 public-read），HTML 用 `<img src="https://<bucket>.oss-<region>.aliyuncs.com/...">` 外链，无 base64 → 过 WAF 且在线渲染正常。工单 83899246 在线报告即用此方式生成成功。
- **建议补丁**：
  1. （本 PR 已含）SKILL.md 增「Image Handling (WAF constraint)」节，记约束 + OSS 外链解法 + 蜜罐提醒。
  2. （后续可选）`bootstrap/html-report-preview.sh` 增强：检测 HTML 内 `data:image` base64，自动上传 OSS 并改写 `src`，对使用者透明；或上传失败识别 punish 页时给出明确提示而非静默返回拦截页 HTML。
- **置信度**：high（多组对照实测 + 在线报告已按解法成功生成）。

> **现状注（2026-07-14）**：helper 现默认在每次上传携带 `X-Request-Context` 头（`JARVIS_HTML_REPORT_WAF_HEADER`，36ac446）过 WAF 分类闸——上文「照做无效、已拒绝」是当时 punish 页语境下的判断；该头**不豁免** base64 图片拦截，图片外链方案不变，且已硬化为 **private OSS + 签名 GET URL**（skill 禁 public-read）。补丁 2（data:image 自动 OSS 化 / punish 页识别）仍未实现，本 cap 保持 open。
