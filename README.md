# jarvis

可被 Claude 完全接替日常工作的自包含起点母版。

## 新机起步

首次在一台新机器上 clone 后，按序一次性设置：

1. `git clone … && cd jarvis`
2. 装依赖：`bash bootstrap/install.sh`（按 `bootstrap/deps.lock` 装 a1/aliyun/cloudspec/gh，已装即跳过）
3. `cp bootstrap/.env.example bootstrap/.env`，填 `GH_TOKEN` + 阿里云 AK
4. a1 登录：`bin/a1id login jarvis`（新机无凭证；不登 scan/wrap/claim 全挂。多身份见 `bin/README.md`）
5. 进目录启动 `claude`——开局动作（`bootstrap/preflight.sh` 日级自检 → scan→plan→triage）由 CLAUDE.md 接管

之后每次干活只需第 5 步；`install.sh`/登录已就位则跳过。

> 钉钉双向桥（可选）：见 `bridge/README.md`，含装依赖、起停与迁移到新机的步骤。

## triage loop

人确认凭证后，Claude 自动执行以下流程：`bootstrap/scan.sh` 拉取当前 Aone 入箱清单 → `bootstrap/plan.sh` 生成执行计划（含动作、置信度、auto/stop、不可逆点）→ supervised 模式下等你逐条授权 → 按 `loops/aone-triage.md` 处理授权项，全链跑到预发/CR。**正式发布永停**，须人工介入。
