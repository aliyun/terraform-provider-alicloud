# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 使用

1. `git clone` 本仓库
2. `cp bootstrap/.env.example bootstrap/.env` 填好 `GH_TOKEN`/阿里云密钥（a1 用容器内凭证）
3. 进目录启动 `claude` 开始使用

## triage loop

人确认凭证后，Claude 自动执行以下流程：`bootstrap/scan.sh` 拉取当前 Aone 入箱清单 → `bootstrap/plan.sh` 生成执行计划（含动作、置信度、auto/stop、不可逆点）→ supervised 模式下等你逐条授权 → 按 `loops/aone-triage.md` 处理授权项，全链跑到预发/CR。**正式发布永停**，须人工介入。