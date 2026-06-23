# 抽离 Aone 项目信息为结构化 config

> 状态：设计 v0 · 2026-06-22 · 待用户评审

## 目标

把散在 `.claude/skills/aone-triage` 里的 Aone 项目事实(池ID/产品cfs/指派/pipeline/认领标签/routing)抽成**唯一结构化数据源**,技能/scan/plan 引用它。解耦后能逐池视检状态。

## 单一数据源 `config/pools.json`

```yaml
pools:
  terraform:    {project: 2100304, product_cfs: "107239=906688", assignee: 320687, dev: false}
  agent_portal: {project: 2124589, app: 283346, assignee: 320687, dev: true,  pipelines: {prestage: 66, prod: 67}, version_cfs: 100340, delivery: delivery-aliyun-automation-agent.md}
  cloudspec_gap:{project: 2165097, assignee: 479782, dev: false}
routing:
  - {match: ["terraform","provider"], pool: terraform}
  - {match: ["AgentRuntime","PlayGround","MCP"], pool: agent_portal}
  - {match: ["cloudspec 缺口"], pool: cloudspec_gap}
claim: {tag: jarvis-claimed, fallback: title}
```

数据=配置；流程与坑留 SKILL.md/delivery-*.md(后者用文件名指过去)。

## 改动
1. 建 `config/pools.json`(verbatim 抽现有 ID,不臆造)
2. SKILL.md/routing.md/delivery 的硬编码 ID → 改"见 config/pools.json",留流程
3. scan/plan 读 `claim.tag`、`pools[].project`,不再 grep prose
4. `verify.sh` 加一条:config 可解析、池数≥3
5. 视检:`pools.sh` 逐池跑 `a1 list --project` 计数,输出每池 open 数

## 验收
- 改池/标签/指派=改 1 行；scan 跳过 `NOT tag=jarvis-claimed`;`pools.sh` 列各池状态；技能去重指 config 无散落 ID。

## YAGNI
不做：池权限管理、自动建标签、迁全部坑。
