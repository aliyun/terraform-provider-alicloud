---
name: writing-weekly-report
description: 写周报/月报/工作汇报技能。覆盖触发：用户说"写周报/做周报/总结这周/近 X 周的工作/做个工作汇报/月报"，或要求把工作汇总更新到钉钉文档（且内容是周报性质）。本技能负责跨仓数据采集（jarvis git+runs/+week.sh+a1 events user）、按"对外价值方向"归纳（不是按仓库分块）、生成接近辰羿本人风格的简洁周报；钉钉文档读写能力（dingtalk-doc-rw）当前未落地，涉及写钉钉文档时按 missing_capability escalate，草稿留对话内。
---

# 写周报

## 核心原则

### 1. 价值视角，不是活动视角

- ❌ 不按仓库分块（cloudspec / automation_agent / acube 各一段）
- ✅ 按对外可交付价值分块（MCP Server RunIaC 上线 / 工具体系对内门户 PlayGround 上线 / Terraform AutoPilot 能力建设）
- 同一价值方向横跨多个仓的工作合并写

### 2. 基础设施工具不当主线

- 内部支撑工具（Jarvis 自治、a1id 多身份、钉钉桥）不要展开技术细节
- 写"用它交付了多少需求"，不写"我搭了什么模块"

### 3. 粒度纪律

| 项 | 上限 |
|---|---|
| 一级方向 | 3-5 个，绝不超过 7 |
| 每方向子项 | 2-4 个 |
| 每子项 | 1 行 |
| 杂项工单（一次性修复、版本升级、小蜜单条） | 默认不写 |

### 4. 状态词标准化

- 已上线 / 已合入 / 已 merged / 已通过预发验收 / 已闭环
- 开发中 / 验收中 / 评审中 / 联调中
- 待合 / 待 prod / 待预发腾空

不要用模糊词："基本完成"、"差不多了"、"在推进"。

### 5. 关键标识——摘要标题用可点击链接

- 摘要中的 Aone 工单 / MR / CR：写成 `[Aone 83346651](url)` / `[MR !28231014](url)` / `[CR !28205668](url)` 的**简洁可点击链接**形式，不要裸号
- 对外发布物（GitHub PR）：带 markdown 链接
- **不必给出分支名**

## 数据采集

### 本仓

```bash
git log --since="YYYY-MM-DD" --pretty=format:"%h %ad %s" --date=short --no-merges
ls -lat runs/ escalation/
bootstrap/week.sh   # 优先,基于 Aone jarvis-done 真源
```

### 跨仓

```bash
bin/a1id -- staff list <花名>                                                          # 拿工号
bin/a1id -- events user <工号> --since YYYY-MM-DD --until YYYY-MM-DD --per-page 100 --type push
bin/a1id -- events user <工号> --since YYYY-MM-DD --until YYYY-MM-DD --per-page 100 --type merge-request
```

只关心 `config/workspaces.json` 登记的仓。

### 身份去重

- 辰羿 / ChenHanZhang / jarvis / chenhanzhang.chz — 同一人（工号 320687）
- open-jarvis（`WORKER_xxx`）— 独立服务账号，写在"服务账号代发 CR"子项
- Jarvis 代码账号（蚂蚁内部 `admin.for.jarvis@antfin.com`）— 不相关，不写

## 输出格式（钉钉文档）

- 一级方向用 `**xxx**` 段落标题（不用 `#`/`##`/`###`）
- 二级用标准 `-` bullet（不用 `●`/`○`/`■`）
- 整篇用 `:::` 包钉钉容器块
- 表格慎用（钉钉 `[...]`/`{...}` 会被吞）

### 范本（gold standard）

```
:::
# **重要事项**

**MCP Server RunIaC 上线**
- **RunIaC 上线**：6-13 ~ 6-25 持续推进，已通过预发验收。
- **MCP token bizid 隔离**：开发中，已补 token / McpToken 单测覆盖。
- **RunIaC 错误去重**：6-26 新开，验收中。

**工具体系对内门户 / PlayGround 上线**
- **PlayGround 体验优化**：工具卡对齐 / 长表格横滚 / 输入框增高 / 新会话引导 / sessionId 复制 / 文件树折叠 / tool 事件回放 / 首页转静态，已上线。PlayGround 主动探测 MCP token 健康状态，连接失败 / 不可用 / 需认证三态透出，已合入。
- **沙箱 / PrivateLink / 凭证策略**：放开预发 MCPServer、PrivateLink 文档、凭证前缀去状态复述改 MCP-only + 明文凭证轮转提示。
- **Playground 访问日志**：已开发，[Aone 83346651](https://project.aone.alibaba-inc.com/v2/project/2100304/req/83346651)。

**Terraform AutoPilot 能力建设**
- **Jarvis 自治系统从 0 到 1**（terraflow/jarvis）。从 6-22 起仓库从零搭建，5 天内落 200+ 提交，覆盖 supervised 授权关口、Aone 唯一真源、并发认领锁、子代理编排、状态看板、worktree+PR 纪律全链路；具备无人值守 triage 到预发的闭环能力。
- 目前近一周多个需求已通过 Jarvis 交付。
:::
```

## 写入流程

1. 采集 → 归纳价值方向
2. 草稿先在对话内呈现，**不直接写文档**
3. 用户确认/调整后，如指定钉钉文档 URL：钉钉文档读写能力（`dingtalk-doc-rw`）**当前未落地**——按 `missing_capability` escalate（loops/self-improve.md），草稿留对话内交用户手动粘贴
4. 该能力落地后：`update_doc` 覆盖式写入，写入后必须 `read_doc` 回读校验渲染

## 对照表：AI 草稿易写错 vs 终稿保留

每次写完后按此自检；有系统性新差异回写本表。

| 草稿（AI 易写） | 终稿（保留） |
|---|---|
| 按仓库分块，每仓一段 | 按价值方向，跨仓合并 |
| 7-8 个一级方向 | 3-5 个 |
| Jarvis 自治写 6 子项展开 | 一句话 + "近一周多个需求已通过 Jarvis 交付" |
| 列所有杂项工单 | 删杂项，只留主线 |
| 加"未完待办"段 | 不加 |
| 加"其他工单交付"段 | 不加 |
| 用 ●/○ 字符 bullet | `**粗体**` + `-` |
| 一级用 `### xxx` | 用 `**xxx**` |
| 裸 Aone 号 + 分支名 | `[Aone xxx](url)` 简洁链接，不带分支名 |

## 红线

1. **不带 AI 署名**：周报对外，禁出现「🤖 Generated with Codex」「Co-Authored-By: Codex」等任何 AI 水印（AGENTS.md 工作纪律 #5）；commit / PR / 钉钉文档正文都不带
2. **回读校验**：钉钉文档写入后必须 `read_doc` 回读，确认渲染正常
3. **草稿先于写入**：除非用户当面授权直接写文档，否则草稿先在对话内呈现
