# cap-repo-structure-cleanup

## 缺口类型

仓库结构累积 drift：40 个 bootstrap 脚本、双端 skill mirror 部分漂移、两个测试目录并存、topdir 有 79M 误落 provider clone、CLAUDE.md 与 loop/skill/cap 正文双写。2026-07-02 一轮只读审计整理清单如下，先记 cap，分批授权推进，避免只口头修。

## 本轮 CR 已修（不必再动）

1. 删除 `null/`（79M provider clone 误落 —— `main.go` 是 terraform-provider-alicloud package，含 .git，是误重定向落坑）
2. `.gitignore` 补 `.idea/`、`bootstrap/.env`
3. CLAUDE.md 精简：10 条纪律 → 7 条；开局动作 3 → 2 条主 + 1 条 bridge 说明；条 3/5/9/10 细节退到 loop / skill / cap 的 point-to
4. AGENTS.md 用 `bootstrap/skills-mirror-lib.sh` 的 `mirror_sed_claude_to_codex` 同步一次，防止已知漂移一行（bridge 定时扫段）继续扩散

## 未修待授权推进

### P0 · skill mirror drift（数据一致性风险，优先）

1. **反向 mirror `invoke-terraform-acc-test-remote` 到 `.claude/skills/`**：codex 侧独有；反向 sync 从没跑过；`bootstrap/sync-to-claude.sh` 跑一次即可复原。
2. **`skills-mirror-lib.sh` 补 `.claude/skills/` sed 规则**：目前 sed 只处理 `.claude/agents/` 路径，`dingtalk-ai-card` SKILL 里 `~/.claude/skills/` 硬编码路径没被转换，codex 侧被人工改成 `~/.Codex/skills/` 后进入真漂移。建议追加 `s|~/\.claude/skills/|~/.Codex/skills/|g` 与逆规则；或让相关 skill 用 `{SKILL_DIR}` 占位符脱敏路径。
3. **`aone-triage/SKILL.md` 的 SUSPEND 章节双端不一致**（.claude 侧多 21 行 SUSPEND 哨兵段 144–163，codex 侧完全缺失）：headless 是 bridge 特性，建议下沉到 CLAUDE.md/AGENTS.md 顶层 headless 模式纪律（本轮 CLAUDE.md 未加入，等评审），或让 mirror 允许双端保留。
4. **`skills-mirror-check.sh` while-piped-subshell 早退 bug**：首个 drift 就 `exit 1`，后续 drift 全部被吞；改成收集全部后统一报告。
5. **`terraform-pr-review` 单侧引用 `invoke-terraform-acc-test-remote`**：codex 侧多一行（71），.claude 侧缺；等反向 mirror 到位后校准。

### P1 · bootstrap 脚本合并（40 → ~27）

- **确定僵尸删除**：`status.sh`（0 外部引用）；如需终端视图加 `board.sh --tui`。
- **mirror 族合并**：`sync-to-codex.sh` + `sync-to-claude.sh` + `skills-mirror-check.sh` + `skills-mirror-lib.sh` → 单入口 `mirror.sh {to-codex|to-claude|check}`；`.claude/settings.json` PostToolUse hook 命令一并改。
- **升级族合并**：`sweep.sh` + `watchdog.sh` 归入 `reconcile.sh {stale|drift|orphan|all}`；同时补 `bootstrap/cron.example` 落地，若不打算跑 cron，直接删 sweep/watchdog 减少幻觉（两脚本目前无 cron/hook 实际调用）。
- **`refresh.sh` 折入**：15 行 shim → `board-html.sh --refresh` 或 `serve.sh /refresh` 内部逻辑；README/loops 同步改。
- **`scan.sh` 用 `cache.sh`**：目前 scan.sh 自己实现 mtime/TTL，与 cache.sh 完全同型；改 `cache.sh fresh scan.json <ttl> -- <query>`。
- **`serve.sh` 用 `lib.sh`**：用 `git rev-parse` 而非 `jarvis_root`，在 worktree 里会错定位到子仓根。
- **抽 `lib.sh` 助手**：`_xxx_repo_root/_xxx_script_dir` 助手函数在 plan/log/sweep/watchdog/wrap-check 5 处重复，`lib.sh` 应统一 `script_dir/source_log/require_pools_cfg`。
- **`wrap.sh` CODE_DIR 走 `workspace.sh`**：目前手拼路径，应经 workspace.sh 解析，对齐 CLAUDE.md 纪律 4。

### P2 · 测试目录合并

- `bootstrap/tests/` 10 文件（用 `<name>.sh` 与被测脚本同名，`lib.sh` 有 source 误伤风险） + `test/` 18 文件（用 `<name>_test.sh`），互不重叠。
- 方案：全部集中到顶级 `test/`，一律 `_test.sh` 后缀。要移动+重命名 10 个：`bootstrap/tests/{aone_comment_format,coord,github_identity,heartbeat,lib,reconcile,triage_one,watchdog,wrap_check,wrap_done}.sh` → `test/<name>_test.sh`；完成后 `rmdir bootstrap/tests`。
- 合并前 grep `bootstrap/tests/` 所有引用面（README/CI/settings/hooks/loops/cron）一并改。

### P3 · skill 拆分与去重

- **`aone-triage/references/tf-customer-request-routing.md`（631 行）拆 3 份**：决策树 + 团队分工 + 镇元查证；把「provider 代码类型判定 + 生成器/手写路由」这段与 `provider-resource-dev` 打通（单点维护），避免两个 skill 双写同一 provider 代码路径判定规则。
- **`terraform-changelog` SKILL.md（322 行）** 接近 400 阈值：Step 9（cut release）+ Step 10（submit PR）可外链 references。
- **`terraform-pr-review` vs `provider-resource-dev` 各自 `sync-provider.sh` 副本**（989B vs 1382B）：是两份不同的 provider 同步脚本，已废弃 or 差异化未标注；查明后收敛。

### P4 · `loops/plugin-dev.md` 归档

- 全仓仅 `config/workspaces.json.agent_toolkit.loop` + `bootstrap/wrap.sh` 一句注释引用，近 30 commit 无活动。
- 方案：降级为 agent-toolkit skill reference（需要新起 `agent-toolkit` skill），或归档到 `escalation/archived/`；同步改 `workspaces.json` 的 `loop` 字段指向新落点。

### P5 · 现有 `escalation/cap-*.md` 处理

- **`cap-auto-memory-save-policy.md`**：CLAUDE.md 纪律 7 已 point-to，原正文可归档 `escalation/archived/` 或降级为 skill reference（cap 保留 minimal）。
- **`cap-agent-skill-drift.md` / `cap-bootstrap-test-drift.md`**：未闭环无 Aone 关联，属"记完躺尸"。挂 Aone 工作项，由 [loops/self-improve.md](../loops/self-improve.md) 例行跟踪，缺回环等于没记。

### P6 · 顶层 & 长期

- **AGENTS.md 由 hook 生成**（而非 tracked）：当前每次 CLAUDE.md 更新要手工同步 AGENTS.md，已发生过一次漂移。长期方案是 `.gitignore` 排除 AGENTS.md，由 `bootstrap/preflight.sh` 或 pre-commit 用 `mirror_sed_claude_to_codex` 从 CLAUDE.md 生成。**风险**：codex 侧首次拉仓时无 AGENTS.md，需要 preflight 兜底生成才能开工。本 CR 未开门，先记路线图，等评审。
- **tools/ 归属**：`tools/acube_terraform_generate.py` + `tools/terraform_generated_diff.py` 目前只服务 `provider-resource-dev`，理论上更合适 `.agents/skills/provider-resource-dev/scripts/`，但涉及双端 mirror；当前 skill README 加注释即可。
- **`.worktrees/` 空目录 rmdir**（可选）：worktree-guard 会按需重建。

## 置信度

各条方向 high_conf（只读审计交叉验证过），待授权分批执行避免一轮 PR 太大。P0 是数据一致性风险，最先做。

## 关联

- CLAUDE.md 自我迭代段已 point-to 本 cap
- [loops/self-improve.md](../loops/self-improve.md)（结构性重构的沉淀路径）
- 4 份只读审计 subagent 报告（会话上下文，未落文件）
