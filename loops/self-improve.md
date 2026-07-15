# self-improve 自迭代 loop

> 遇到不会干 / 缺环境 / 缺上下文 / 路由未命中 → 记能力缺口 → worktree 起草补丁 → 开 MR/PR → 等人工合并。
> 绑 `autonomy.md` 的 `missing_capability` 触发；**永不自动合 master**。

---

## 一、触发

| 触发器 | 说明 |
|--------|------|
| 不会干 | 任务无对应技能/runbook |
| 缺环境 | bin/凭据/clone 缺失（verify 不过） |
| 缺上下文 | 找不到 schema/池/源码映射 |
| 路由未命中 | pools.json/workspaces.json 无匹配 |

对应 `autonomy.md` 的 `missing_capability` escalate trigger。

---

## 二、记缺口

```
escalation/cap-<slug>.md
```

内容：缺口类型、阻塞任务、缺什么、建议补丁、置信度。

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

补丁开 MR/PR 等人工评审，关联 `cap-<slug>.md`。

---

## 五、Done

| 结果 | 说明 |
|------|------|
| 入队 | cap 缺口 + 草案 PR 待人工合 |
| 永不 | **绝不自动合 master**（redline 硬门） |

---

## 六、工具链速查

| 工具 | 作用 |
|------|------|
| `escalation/cap-<slug>.md` | 缺口记录 |
| `autonomy.md missing_capability` | 触发器 |
| `bootstrap/log.sh escalate` | 审计上报 |
