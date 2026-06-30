# cap-agent-skill-drift

## 缺口类型

skill / agent canonical source 不清晰,存在重复目录漂移。

## 阻塞任务

Jarvis 在 Codex / Claude 两类入口下都 vendored 了 Agent 和 Skill。Terraform 资源开发流程需要同时更新 `.agents/skills/provider-resource-dev` 与 `.claude/skills/provider-resource-dev`;如果只改一边,后续不同入口会执行不同流程。

## 当前发现

- `.agents/skills` 与 `.claude/skills` 有重复 skill 副本。
- `.agents/skills/provider-resource-dev` 与 `.claude/skills/provider-resource-dev` 本次已同步。
- `.agents/skills/dingtalk-ai-card` 与 `.claude/skills/dingtalk-ai-card` 存在脚本路径差异:一边指向 `~/.Codex`,一边指向 `~/.claude`。
- `.agents/skills/writing-weekly-report` 与 `.claude/skills/writing-weekly-report` 存在红线文案差异:Codex/Claude 署名描述不同。
- `.Codex/agents`、`.codex/agents`、`.claude/agents` 并存,目前没有一个测试或脚本声明 canonical source 与同步例外。

## 建议补丁

1. 明确 canonical source:建议 `.agents/skills` 作为 Codex 当前入口的主副本,`.claude/skills` 作为兼容副本。
2. 增加同步脚本或测试:核心 skill 默认两边一致;允许差异必须登记白名单和原因。
3. 对路径敏感 skill 使用相对路径或 skill 目录定位,避免写死 `~/.Codex` / `~/.claude`。
4. Agent 定义也做同样治理:明确 `.codex/agents`、`.Codex/agents`、`.claude/agents` 的生成关系和同步验证。

## 置信度

high:基于本仓 `diff -qr .agents/skills .claude/skills` 与 agent 文件盘点。

## 关联

Aone: https://project.aone.alibaba-inc.com/v2/project/2100304/req/83695664
