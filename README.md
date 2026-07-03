# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 配置（一次性）

```bash
bash bootstrap/install.sh                       # 装依赖（a1/aliyun/cloudspec/gh）
cp bootstrap/.env.example bootstrap/.env        # 填 GH_TOKEN + 阿里云 AK
bin/a1id login jarvis                           # 至少登录 jarvis（默认身份）；其他身份按需 `bin/a1id login <chenyi|guozai|linjun>`
```

钉钉数字员工还需 `cp bridge/jarvis.env.example bridge/jarvis.env`，填钉钉 appKey/appSecret/卡片模板 id。

## 启动

**当面用** —— 进目录起会话，CLAUDE.md 接管 scan→plan→triage：

```bash
claude
```

**数字员工** —— 常驻钉钉，收消息 → 跑 jarvis → 流式回卡：

```bash
bridge/run.sh start
```

## triage loop

人确认凭证后，Claude 自动执行以下流程：`bootstrap/scan.sh` 拉取当前 Aone 入箱清单 → `bootstrap/plan.sh` 生成执行计划（含动作、置信度、auto/stop、不可逆点）→ supervised 模式下等你逐条授权 → 按 `loops/aone-triage.md` 处理授权项，全链跑到预发/CR。**正式发布永停**，须人工介入。

## 任务看板

5 列卡片看板（escalated / merged / done / inflight / 任务池），真源=Aone 标签与状态。

```bash
bootstrap/serve.sh        # 起服务并打开 http://localhost:8787（未生成会自动 build）
bootstrap/serve.sh 9000   # 换端口
```

页内点「立即同步」或 `POST /refresh` 强制重扫 Aone 并重建；离线重建用 `bash bootstrap/board-html.sh --refresh`。

## HTML 报告在线预览

Jarvis 可把本地 HTML、zip 内 HTML、或 Aone 附件中的 HTML 报告上传到 AutomationAgent，并返回可贴回 Aone 的在线预览链接。接口走 AutomationAgent 的 `/buc/reports/aone/*` BUC 组。

```bash
# 上传本地 HTML / zip / 目录
bash bootstrap/html-report-preview.sh upload 83843879 ./report.html
bash bootstrap/html-report-preview.sh upload 83843879 ./reports.zip --comment

# 从 Aone 附件中选择最新的 .html/.htm/.zip，提取 HTML 后上传并回贴评论
bash bootstrap/html-report-preview.sh from-aone 83843879 --comment

# 上传所有匹配附件；默认预发，可通过 --base-url 切换环境
bash bootstrap/html-report-preview.sh from-aone 83843879 --all --base-url https://pre-agent.aliyun-inc.com
```

默认目标是 `https://pre-agent.aliyun-inc.com`；也可用 `JARVIS_HTML_REPORT_BASE_URL` 或 `--base-url` 覆盖。
