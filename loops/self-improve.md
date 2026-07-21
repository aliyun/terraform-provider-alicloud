# self-improve 自迭代 loop

> 遇到不会干 / 缺环境 / 缺上下文 / 路由未命中 → 在 Aone 记录能力缺口 → worktree 起草补丁 → 开 MR/PR → 等人工合并。
> 绑 `autonomy.md` 的 `missing_capability` 触发；**永不自动合 master**。

---

## 一、触发

| 触发器 | 说明 |
|--------|------|
| 不会干 | 任务无对应技能/runbook |
| 缺环境 | bin/凭据/clone 缺失（verify 不过） |
| 缺上下文 | 找不到 schema/池/源码映射 |
| 路由未命中 | pools.json/workspaces.json 无匹配 |

对应 `autonomy.md` 的 `missing_capability` trigger。

---

## 二、记缺口

优先复用当前 Aone 工作项；没有可承载的工作项时，按 `loops/adhoc-intake.md` 在所属团队新建能力改造单。记录缺口类型、阻塞任务、缺什么、建议补丁、置信度和复现证据。需要人工决策时将控制面 Task 置为 `SUSPENDED`，由幂等事件发布器把问题写回 Aone；不落本地草稿文件。

---

## 三、起草补丁（worktree）

worktree 切分支，补 loops/ 或 .claude/skills/ 或 config/ 之一：

| 缺口 | 补丁 |
|------|------|
| 缺 runbook | 新 `loops/*.md` |
| 缺技能 | `.claude/skills/*` |
| 缺路由/工作区 | `config/pools.json` / `config/workspaces.json` |

编码交子代理（CLAUDE.md）。

---

## 四、MR/PR

补丁开 MR/PR 等人工评审，关联 Aone 工作项。技术结论同步进相关 skill/reference，让其他机器和后续会话自然加载。

---

## 五、Done

| 结果 | 说明 |
|------|------|
| 入队 | Aone 缺口记录 + 草案 PR 待人工合 |
| 永不 | **绝不自动合 master**（redline 硬门） |

---

## 六、工具链速查

| 工具 | 作用 |
|------|------|
| Aone 工作项 | 缺口、证据、决策与进度真源 |
| `autonomy.md missing_capability` | 触发器 |
| 相关 `.claude/skills/*` / `loops/*` | 可复用技术知识与执行流程 |
