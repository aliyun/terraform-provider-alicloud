# cap-sync-auto-memory-to-skill

## 缺口类型

自动生成的 project auto-memory(`~/.claude/projects/-Users-gzzz-Workplace-jarvis/memory/`)是**机器本地文件,不落 git**,不能跨设备共享。多台设备用 jarvis 时(主开发机、钉钉桥接机、cron/routine 机),memory 只能靠人肉 copy 或重复采集,违背"仓库主人认知一次落库,处处可用"的自迭代原则。

## 阻塞任务

- 换机器 / 工作站后 auto-memory 全丢,已积累的 user / feedback / project / reference 全部重来
- 同一工单在 A 机学到的教训(如 memory `wrap-done-single-comment`、`read-description-last-paragraph`)在 B 机不生效,B 机重复踩坑
- 跨设备联调时不同 jarvis 实例上下文错位:钉钉桥接机的实例看不到主开发机学到的团队分工表 / 反馈 / 用户偏好,导致回复漂移

## 现状(已存在的邻近机制)

- `bootstrap/sync-to-codex.sh` + `bootstrap/sync-to-claude.sh`:PostToolUse hook,skill 层文件在 `.claude/skills/` 与 `.agents/skills/` 双向镜像,codex↔claude 共享 skill 已解决
- **未解决**:上述 sync 只覆盖 `.claude/skills/` 与 `.agents/skills/`,不覆盖 memory(`~/.claude/projects/...` 路径,不在仓内也不在 sync 白名单)

## 建议补丁

### 方向一(推荐):hook-based sync + preflight 校对

- **PostToolUse hook**:检测到 Write/Edit 目标在 `~/.claude/projects/-Users-gzzz-Workplace-jarvis/memory/` 时,rsync 或 cp 到 `.claude/skills/memory/`(仓内新目录,随仓入 git),同时镜像到 `.agents/skills/memory/`(codex 侧)
- **preflight.sh 开局校对**:比对 `~/.claude/projects/.../memory/` 与 `.claude/skills/memory/`,按 mtime 更新旧的一侧;冲突留 `<name>.conflict.<hostname>.md` 让人工裁决
- **`MEMORY.md` 索引也同步**:MEMORY.md 是 auto-memory 的目录,一并镜像
- 优点:hook 已有基座(`sync-to-codex.sh` 是现成范式);colocate 到 `.claude/skills/memory/` 后 skill trigger 逻辑不受影响(memory 目录不带 SKILL.md,不当 skill 触发)

### 方向二:符号链接

- `ln -sfn <repo>/.claude/skills/memory ~/.claude/projects/-Users-gzzz-Workplace-jarvis/memory`
- 简单粗暴但要求 preflight 每台机第一次跑校对/建立符号链接;codex 侧路径不同,同一 memory 实体给 claude 和 codex 都能读需要多个 symlink;跨账号(claude vs codex 分属不同用户目录)不友好

### 方向三:cron / wrap.sh 收尾时增量提交

- 每 N 分钟或每次 bookend 后 diff memory 目录,commit 到 refs/memory 分支(不进 master)
- 门槛低但产生噪声 commit;历史被自动 commit 淹没

## 建议

**方向一优先**(hook-based + preflight 校对)。粒度合适 + 无需改用户环境 + 与现有 sync-to-codex/claude 同构。

## 关联

- CLAUDE.md 工作纪律 #4 自我迭代:本 cap 是"跨设备共享认知"能力补齐,直接对应
- `bootstrap/sync-to-codex.sh` / `bootstrap/sync-to-claude.sh`:现成模型可复用
- feat/master-guard-and-aone-triage-refactor(本 CR):不冲突,可独立实现

## 置信度

low_conf:方案未实现验证。先建单等仓库主人拍板方向,再起草补丁(应该在独立 CR 中实现)。
