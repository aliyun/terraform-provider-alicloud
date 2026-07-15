# cap: bridge 新 checkout 缺 bootstrap/.env → JARVIS_GITHUB_TOKEN 未注入 → PR 推不出，发布类工单卡 jarvis-idle

- **缺口类型**: 部署/配置缺陷（missing_capability + 模板文档误导）
- **阻塞任务**: Terraform 资源发布自动审核类工单（Redis TairInstance，工单 84291978）——headless jarvis 已完成代码开发，但一直卡在 jarvis-idle 不开 PR，表面「无进展」
- **置信度**: high_conf（实测复现 + 修复验证：`github-identity.sh check` 从 "JARVIS_GITHUB_TOKEN is required" 转为返回 `api-tool-agent`）

## 现象

发布类工单被 headless jarvis 正常认领、PD 分诊 high_conf、RD gap 分析零差异并完成代码开发（commit 已落本地分支），但收尾一直停在 jarvis-idle、不开 PR。run 审计里写着：

> GitHub push 受阻: JARVIS_GITHUB_TOKEN 未注入 dispatch 环境, 需 bridge 注入 token 后续跑

此后每轮 scan 判 `skip (idle_no_human)`，无人工介入则不再重派，工单表面「无进展」。

## 根因

1. bridge 迁到新 checkout（`.../workspace/jarvis` → `.../workspace/jarvis-preview`）运行，但 gitignored 的 `bootstrap/.env`（含 `JARVIS_GITHUB_TOKEN`）没随迁——git 不带、需每机 / 每 checkout 手建。
2. `bridge/run.sh` start 只在 `bootstrap/.env` 存在时才 source；文件缺失即静默无 token（其余 JARVIS_* 变量来自 `bridge/jarvis.env`，照常运行，唯独漏了 token 这一项，故障隐蔽）。
3. `~/.zshrc` 虽 `export JARVIS_GITHUB_TOKEN`，但 bridge 经 launchd / nohup 起、非交互式 login shell，拿不到 `.zshrc` 的 export。
4. **`bootstrap/.env.example` 模板误导**：写的是 `GH_TOKEN=`，而代码（`github-identity.sh`、`claim.sh`）实际读的是 `JARVIS_GITHUB_TOKEN`（`GH_TOKEN` 只在内部 `GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh …` 临时赋值）。照模板填 `GH_TOKEN` 也过不了身份门。

结果：每个被派发的 headless jarvis 都继承一个无 token 的环境，`github-identity.sh check` 硬失败 "JARVIS_GITHUB_TOKEN is required"，push / 开 PR 全断，发布类工单只能卡在 idle。

## 修复（本 PR）

1. `bootstrap/.env.example`：把误导的 `GH_TOKEN=` 改为真正需要的 `JARVIS_GITHUB_TOKEN=`，并加醒目注释——`.env` 不随 git 迁移、换 checkout / 机器必重建，`run.sh` 启动 source 它，`.zshrc` export 对 launchd/nohup 起的 bridge 无效，`github-identity.sh check` 可自检。
2. 当值机器（本次 jarvis-preview checkout）：已补 `bootstrap/.env` 含有效 `JARVIS_GITHUB_TOKEN`（api-tool-agent，权限 600、gitignored）。

## 复发防护

- 换 checkout / 机器起 bridge 前，先跑 `bootstrap/github-identity.sh check`，返回 `api-tool-agent` 才算就绪；报 "JARVIS_GITHUB_TOKEN is required" 即缺 `.env`。
- 排障信号：run 审计出现「GitHub push 受阻: JARVIS_GITHUB_TOKEN 未注入」= 本问题。
- 卡在 jarvis-idle 且始终不开 PR 的发布类工单，优先怀疑 token 注入。
- 关联同族 GitHub 身份陷阱：`escalation/archived/cap-github-commit-identity.md`（那条管 commit 作者 → CLA；本条管 push token 注入 → 能否 push）。
