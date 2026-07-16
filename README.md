# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 配置（一次性）

```bash
bash bootstrap/install.sh                       # 装依赖（a1/aliyun/cloudspec/amp/gh）
cp bootstrap/.env.example bootstrap/.env        # 填 GH_TOKEN + 阿里云 AK
bin/a1id login jarvis                           # 至少登录 jarvis（编排层默认身份 WORKER_1782379562571）
bin/a1id login terraform-rd                     # Terraform 唯一公共身份；内部仍按 PD/RD/QA subagent 分工
# 个人身份仅在仓库主人本轮当面授权时切：
#   bin/a1id login <chenyi|guozai|linjun|shanye>
```

CloudSpec IDL 机器人能力由仓库内锁定快照提供；首次执行 Terraform release 的 CloudSpec
闭环前，运行 `bash bootstrap/cloudspec-core.sh doctor`，确认 amp 与 `aliyun cspec` 就绪。

钉钉数字员工还需 `cp bridge/jarvis.env.example bridge/jarvis.env`，填钉钉 appKey/appSecret/卡片模板 id。

## 启动

**当面用** —— 进目录起会话，CLAUDE.md 接管（preflight → 等任务，单条工单/即时任务按 loops 处理）：

```bash
claude
```

**数字员工** —— 常驻钉钉，收消息 → 跑 jarvis → 流式回卡：

```bash
bridge/run.sh start
```

## triage loop

批量扫派由 bridge 负责：`bridge/run.sh start` 起 ScanScheduler 定时扫池，对新单/外部更新单自动派发 headless jarvis 并发处理（`JARVIS_AUTO_DISPATCH=0` 可回退为钉钉授权前置模式），按 `loops/aone-triage.md` 全链跑到预发/CR。当面会话只处理用户直接给的单条工单或即时任务（`loops/adhoc-intake.md`）；`bootstrap/scan.sh` / `bootstrap/plan.sh` 保留作手动兜底与 bridge 内部计划步骤。**正式发布永停**，须人工介入。

## 任务看板

5 列卡片看板（任务池 / 待开始 / 进行中 / 审核中 / 已完成），真源=Aone 标签与状态。

```bash
bootstrap/serve.sh        # 起服务并打开 http://localhost:8787（未生成会自动 build）
bootstrap/serve.sh 9000   # 换端口
```

页内点「立即同步」或 `POST /refresh` 强制重扫 Aone 并重建；离线重建用 `bash bootstrap/board-html.sh --refresh`。

## HTML 报告在线预览

Jarvis 可把本地 HTML、zip 内 HTML、或 Aone 附件中的 HTML 报告上传到 AutomationAgent，并返回可贴回 Aone 的在线预览链接。预览链接走 AutomationAgent 的 `/reports/aone/*` 只读路径，Jarvis 上传走 `/api/reports/aone/*` 服务端 token 接口。

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
AutomationAgent 的上传接口需配置 `JARVIS_HTML_REPORT_TOKEN`，脚本会自动发送 `Authorization: Bearer <token>`；脚本会读取 `bootstrap/.env` 作为默认值，显式环境变量优先。
