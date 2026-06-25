# adhoc-intake 即时任务接入 loop

> 实时任务接入：用户随手丢一个请求（带不带 Aone 都行），意图分诊 → 关联 Aone → 解析工作区 → cd 进项目 → 动手（默认只读，dev 需授权）→ 审计 → 缺口上报。
> 与 `loops/aone-triage.md` 互补：triage 是入箱批扫的工单池，本 loop 是单条即时插入的任意任务。

---

## 一、触发

| 方式 | 说明 |
|------|------|
| 手动触发 | 用户在会话中直接发请求：「看下 alicloud 这个资源支不支持 X」「跑下 mcp_server 构建」「这个 PR 帮我评审」 |
| 链接触发 | 贴 Aone 工作项 / PR / next.api 链接，无固定工单池 |
| 无 Aone | 纯本地任务（读代码、跑 build/test、查证），不一定有工作项 |

与 triage 的区别：triage 走 `scan.sh` 批量扫池；本 loop 接单条即时任务，**不扫池**。

---

## 二、输入分诊（intent classify）

由 Claude 对请求做意图分类，不靠脚本扫表：

| 意图 | 例 | 默认动作面 |
|------|----|-----------|
| 查证 / 问答 | 「provider 支不支持 X」 | 只读 |
| 评审 | 「评审这个 PR」 | 只读 |
| 开发 | 「补字段 / 改 bug / 加资源」 | dev（需授权） |
| 运维 | 「跑 build/test/vet/fmt」 | 只读（ops 白名单） |

---

## 三、关联 Aone（supervised 授权关口）

1. **找现有工单**：按意图/关键字检索是否已有对应 Aone 工作项。
2. **命中** → 复用其 id；记录池（`config/pools.json` 路由）。
3. **未命中** → **反问选池**（不擅自创建）：列出候选池（tf_provider / tf_customer / mcp_server / cloudspec / api_toolkit），等用户选。
4. **授权后建 + 双向关联**：建工作项 → 工单挂本任务链接，本任务记 Aone id（双向）。
   - **统一建为需求（`--category req`）**：adhoc 接入默认开需求单，不开 task/bug，便于后续走需求→变更→发布链路。确属缺陷再用 bug。
   - ad-hoc PR 无明确归属 → 默认落 **tf_provider (528766)**。
5. **凡要写工单(新建 or 复用)都开 bookend**：拿到 id 后(无论第 2 步复用还是第 4 步新建)，**开局即 `bootstrap/claim.sh claim <id> <project>`**(打 jarvis-claimed、入台账)，收尾走 bookend 收口(`triage-one.sh` 或 `wrap.sh done`+`claim.sh release`)。漏 claim = 后续标签/状态/对账全失灵。纯本地只读、不动任何工单的任务可跳过本节直接进四。

保留 **supervised 门**：建/写 Aone 前逐项等授权。

---

## 四、解析工作区 + cd

读 `config/workspaces.json` 取 `workspaces.<key>`（terraform_provider | mcp_server）：`repo` / `path` / remotes / `default_branch` / `pools` / `ops`。

- 评审 alicloud PR：看 `upstream_remote=alicloud`，改在 `origin=ChenHanZhang` fork。
- cd 进 `path`，dev 先开 worktree 切分支（CLAUDE.md 工作纪律）。

---

## 五、动手（只读默认 / dev 授权）

| 面 | 操作 | 授权 |
|----|------|------|
| 只读 | 读码、查证、`ops.build/test/vet/fmt` | 默认 |
| dev | worktree → 改 → 验 → MR/PR；mcp_server 到预发 | 授权后 |

红线（推 master / 零差异 CR / 正式发布）→ escalate，禁止自动。

---

## 六、审计

```bash
bootstrap/wrap.sh done <id> "<任务+落点>" "<status>"
```

收尾回填 Aone（评论+改状态）并落 `runs/`；dev 中途用 `wrap.sh sync <id> "<进展>"` 报进展。临时数据走 `.my-day/`，禁往仓库根甩 scratch。Aone 唯一真源——禁止只在本地推进不落 Aone。

> **收尾必走 bookend，禁裸 `log.sh run_done`**：凡 claim 过的工单(新建/复用都算)，收尾用 `bootstrap/triage-one.sh <id> <pool> <project> "<summary>" <status>`(claim→done→release 一把成对)，或手跑 `wrap.sh done` + `claim.sh release` 两步。裸 run_done 只落本地审计：标签停 claimed/状态不动/release 漏。claim 开头、release 收尾闭合，`wrap-check.sh`+`reconcile.sh` 才兜得住。

---

## 七、Done

| 结果 | 说明 |
|------|------|
| 完成 | 只读出结论 / dev 到 MR·预发，`wrap.sh done` 回填 Aone+入 `runs/` |
| escalation | 缺口/低置信/红线 → 写 `escalation/`，触发 `self-improve` |

---

## 八、工具链速查

| 工具 | 作用 |
|------|------|
| `config/workspaces.json` | 工作区 canonical schema → repo/path/remotes/ops |
| `config/pools.json` | 池路由（ad-hoc PR→tf_provider 528766，客户 1086837） |
| `bootstrap/wrap.sh sync/done` | 进展实时回填 Aone（真源）+收尾审计 |
| `loops/self-improve.md` | 缺口→escalation→补丁 |
