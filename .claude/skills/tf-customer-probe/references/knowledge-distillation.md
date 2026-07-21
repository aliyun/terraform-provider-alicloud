# 云产品知识蒸馏（KNOWLEDGE.md 契约）

> 本 reference 是**跨 skill 单点维护**——tf-customer-probe（Step E 收尾）、aone-triage（bookend 收尾）、
> provider-resource-dev（开发完成后）三处**都读它**（先例：provider-resource-dev/references/zhenyuan-verification.md
> 的跨 skill 复用模式）。

## 落点

**数据**在 `tf_playground/<product>/KNOWLEDGE.md`（独立数据仓 `terraflow/tf_playground`,直推 master + 工单报备,
非代码不走 MR）。playground 根解析优先级(与场景库同):
env `JARVIS_TF_PLAYGROUND` > config `paths.playground_dir` > `bootstrap/workspace.sh dir tf_playground` >
默认 `<jarvis 父目录>/terraform_playground`。

**契约**（本文件）在 jarvis 仓；三处消费方各加一行指针即可（不复制契约文本）。

## 五节条目结构

每个产品的 KNOWLEDGE.md 固定五节，缺节留空标题即可（不删小节，方便机读扫）：

```markdown
# <product> KNOWLEDGE

> jarvis 蒸馏的产品级可执行知识；条目按时序追加，格式见 `.claude/skills/tf-customer-probe/references/knowledge-distillation.md`。

## 命名与基本行为
- [YYYY-MM-DD][来源: <链接/路径>] <一条可执行的产品级事实>

## 参数约束与枚举 quirk
- ...

## 生命周期陷阱（ForceNew · 永久 diff · import）
- ...

## API 行为与废弃 action
- ...

## 报错 → 原因 → 解法
- ...
```

## 条目格式（硬规则）

```
- [YYYY-MM-DD][来源: <链接/路径>] <一条可执行的产品级事实>
```

- **日期**=蒸馏当日（UTC，`date -u +%F`）。
- **来源**=可回追的锚点：Aone 工单 URL / verdict 路径（`runs/probe/<YYYYMMDD>-<HHMMSS>-<sid>.json`）/
  上游 GitHub PR URL / provider 源码行号（`alicloud/resource_alicloud_xxx.go:LNN`）。
- **正文**=**一条可执行的产品级事实**（客户/开发者读到就能照做/避坑的一句话）；避免流程性/时序性叙述、
  避免"我们做了 X"式过程记录——那些进 verdict/工单评论，不进 KNOWLEDGE。

## 触发点（三处，缺一条链路都会漏收）

1. **probe 轮 Step E 收尾**（`.claude/skills/tf-customer-probe/SKILL.md` Step E）：
   跑完 tier-0/tier-1 汇报前，按下面**收录判据**过一遍本轮 findings + judgment 判定结果 + verdict 摘录，
   命中判据的条目追加进对应产品 `KNOWLEDGE.md`。
2. **aone-triage bookend 收尾**（`.claude/skills/aone-triage/SKILL.md` 主流程 bookend 段）：
   凡工单涉及某个 terraform 云产品（客户单 / 内部研发单 / probe 单皆算），在 `wrap.sh done` 之后、
   `claim.sh release/finish` 之前，把本轮**验证结论中可复用的产品级事实**蒸馏进 `<playground>/<product>/KNOWLEDGE.md`。
   **这是评审阻断项**——客户单场合的蒸馏钩子必须挂在 aone-triage 主流程，不能只挂 probe 侧。
3. **provider-resource-dev 完成开发后**（`.claude/skills/provider-resource-dev/SKILL.md` 步骤 8 PR 提交后）：
   把开发过程学到的产品 quirk（API 行为差异、schema 陷阱、必须的重试码等）蒸馏进 `<product>/KNOWLEDGE.md`，
   来源锚点写 upstream PR URL + provider 源码行号。

## 收录判据（可执行 / 跨场景复用 / 非文档已明示）

收进 KNOWLEDGE.md 的条目**必须同时满足三条**：

- **可执行**：读到就能照做（避坑、选参、认错、走 workaround），不是背景/时序叙述。
- **跨场景复用**：同产品的其他资源 / 其他客户 / 未来某类场景大概率也会命中。
- **非文档已明示**：官方 TF 文档 / OpenAPI 文档已经明确写清的照抄事实不收（读者会自己去查），
  收的是文档**没写清**、**没提**、或**与代码/API 实际行为不一致**的产品级事实。

**不收**（一律进 verdict / 工单评论 / draft，别污染 KNOWLEDGE）：

- 流程性一次性内容（"本轮我们跑了 X 场景发现 Y"）；
- 单次事件叙述（"某月某日某工单里 A 客户遇到 B"）；
- 客户/账号/工单标识（含 Aone 号）——那类信息只在**内部仓**（jarvis）的过程材料里出现，**KNOWLEDGE 追求可复用**；
- 已 CHANGELOG Unreleased 记过的临时修复注记（那是发版说明的活）。

## 消费约定（三处开发前都读）

- **`.claude/skills/tf-customer-probe`**：Step B 挑场景前，若涉及产品 X，先读 `<playground>/<X>/KNOWLEDGE.md`（存在即读）。
- **`.claude/skills/aone-triage`**：查证阶段（SKILL.md 第 3 步）若命中某个 terraform 产品，先读 `<playground>/<X>/KNOWLEDGE.md`。
- **`.claude/skills/provider-resource-dev`**：开发某资源前（步骤 6 手改），先读 `<playground>/<X>/KNOWLEDGE.md`。

playground 路径解析：`bootstrap/workspace.sh dir tf_playground` 或 env `JARVIS_TF_PLAYGROUND`；
文件不存在即跳过（无信息 ≠ 阻断），存在则读全文（文件按小节结构组织，五节均需扫）。

## sanitize 预埋（对外流转必过 CLAUDE.md #5 禁品清单）

- **内部仓（jarvis / tf_playground）内**：KNOWLEDGE.md 条目可带 Aone 工单链接作为来源锚点（数据仓 `terraflow/tf_playground`
  gitlab 私仓，非对外），便于回追。
- **对外产物**（GitHub 公开 PR 评论 / 上游 issue / 公开文档 / 未来毕业成产品级 skill 发到 registry）**必过 sanitize**：
  按 CLAUDE.md 工作纪律 #5 禁品清单逐条过滤——AI 署名 / 客户名或账号 UID / Aone URL 或工单号 / 客户实例 ID
  (`r-xxx`/`i-xxx` 等) / RequestId / 花名+工号引用一律剥掉。
- **实操**：需要把 KNOWLEDGE.md 某条目搬到对外产物时，只留"可执行的产品级事实"正文，来源锚点改成**公开可达**
  的替代（OpenAPI 文档 URL / provider PR URL / 上游 issue URL），把内部锚点整段删掉；宁可失一次内部回追、
  不可泄一次内部信息。

## 毕业标准（进入产品级 skill）

- **某产品 KNOWLEDGE.md 条目 ≥ 15 条** 且
- **被 dev/review/probe 消费 ≥ 3 次**（在 aone-triage 查证 / provider-resource-dev 开发 / tf-customer-probe 判定的
  会话中被引用），

即触发**起草 `.claude/skills/<product>-*` 产品级 skill**（如 `alicloud-vpc-quirks` / `alicloud-oss-lifecycle`），
把稳定的产品级事实沉淀为可搜索、可 trigger 的 skill；KNOWLEDGE.md 继续作为**未稳定条目的写入队列**。
毕业动作建 Aone 跟踪，并在相关产品 skill/reference 中记录（"KNOWLEDGE → 产品级 skill 毕业"）。

## 自检（写完看一遍）

- 条目落在正确产品目录（跨产品复用不明显的写主产品；例：VPC + OSS 组合场景的 KNOWLEDGE 落主产品目录，
  与场景 authoring 的跨产品组合场景落主产品目录约定一致）。
- 来源锚点可回追（内部仓可含 Aone；对外流转前按上述 sanitize 段处理）。
- 每条正文都能通过"读到就能照做/避坑"的可执行检验，不是过程叙述。
- 条目格式正确：`- [YYYY-MM-DD][来源: ...] <正文>`（保证机读扫描 grep 得动）。
