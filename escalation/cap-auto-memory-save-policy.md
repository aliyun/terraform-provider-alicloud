# cap-auto-memory-save-policy

## 缺口类型

auto-memory 的**保存决策**错位:技术性 / 团队性 / 项目性知识被存入 `~/.claude/projects/-Users-gzzz-Workplace-jarvis/memory/`(per-machine,不落 git,不跨设备),而这类知识本应入相关 skill/reference(落 git,跨设备,并且能通过 skill trigger 自然出现在使用场景)。

## 走过的弯路(留个记录)

1. **最初误判**:以为缺口是"auto-memory 没同步到仓内",提方向一"新建 `.claude/skills/memory/` + hook sync"。
2. **仓库主人指正**:即使 sync 到 `.claude/skills/memory/`,claude/codex 默认不会自动读该目录——auto-memory 的自动加载只覆盖 `~/.claude/projects/<proj>/memory/MEMORY.md`(前 200 行注入),skill trigger 只加载命中 skill 的 `SKILL.md`,其它 md 都要显式 Read。sync 一份放到无人读的地方是纯冗余。
3. **真正的问题**:save memory 时应先看内容能否补入相关 skill,能则补(去重),不能则不落。auto-memory 只留个人/机器/临时上下文,这类内容 per-machine 本来就合理,跨设备没意义。

## 现状核对(佐证)

save 至今的 3 条 auto-memory,内容 100% 已在 skill 中出现:

| memory 文件 | 内容 | 已覆盖处 |
|---|---|---|
| `feedback_wrap_done_single_comment.md` | wrap.sh done 别先手动 comment | `aone-triage/SKILL.md` 反模式段 |
| `feedback_read_description_last_paragraph.md` | 小蜜工单末段是真实诉求 | `aone-triage/references/tf-customer-request-routing.md` Step 1 |
| `reference_team_roster_tf_alicloud.md` | 团队分工表(11 云产品专属 + 4 通用路由) | `aone-triage/references/tf-customer-request-routing.md` 团队分工速查 |

即 sync 需求本身是幻觉——技术知识本就该只在一处存(skill),重复维护两处才是问题。

## 补丁(本 CR 已实现)

1. **CLAUDE.md 新增工作纪律 #10「auto-memory 只存 personal/machine 上下文,技术知识入 skill」**:save memory 前扫 `.claude/skills/**/*.md`,已覆盖则不写;属技术/团队/项目类且 skill 未覆盖 → 补入 skill/reference,不落 memory;仅 personal 偏好 / 临时状态才走 auto-memory。
2. **清理现有 3 条冗余 auto-memory**:全部内容已在 skill 中,per-machine 删除,MEMORY.md 留说明注释。此步不入 CR(auto-memory 文件不在仓内),但 CR 合入后各机第一次跑到应用新纪律即可自然收敛。
3. **不需要方向一/二/三**:不建 memory 目录、不加 sync hook、不做 symlink、不做 cron 提交。

## 后续可能

- 若观察到"确有内容不属技术知识但需要跨设备"(如仓库主人偏好在多机通用),再评估轻量方案(单独 config 文件 + 手动 commit,而不是给 auto-memory 加通用 sync)。
- 未来 skill 数量膨胀后,save 前的"扫 skill 判重"如果人工太累,可考虑 preflight 加个校对脚本或 subagent 辅助——现阶段 skill 数(10 上下)可手扫。

## 关联

- CLAUDE.md 工作纪律 #4 自我迭代 + 新增 #10 一并落地
- feat/master-guard-and-aone-triage-refactor(本 CR)
- 与 `bootstrap/sync-to-codex.sh` 的 skill 双向镜像模型解耦:memory 完全不进 sync 范围

## 置信度

high_conf:方案已在本 CR 内实现;实际观察(有无新的 skill-worthy 内容误落 memory)在下一次 save memory 事件才能验证,若发现回归再补脚本级校对。
