---
name: terraform-provider-release
description: Release a Terraform Provider resource together with its data source for Alibaba Cloud. Covers both new-resource publishing (via generator) and updates to existing resources (auto-generated or hand-written). Anchored on an Aone work item end to end, runs in an isolated worktree, forbids local compilation, derives the concrete release scope by comparing AMP metadata against the local provider code (requirement gap analysis), verifies the underlying APIs are all OpenAPI before generation, and verifies via remote ACC tests before PR submission. The task is only complete after the PR is merged. NOT for 镇元/acube 生成器诊断链路（生成为空/损坏、resourceTypeCode 未命中、cspec 分支验证）— use provider-resource-dev. NOT for reviewing an existing PR — use terraform-pr-review.
metadata:
  version: "0.6.0"
  domain: terraform-provider
  triggers: 资源发布, provider 发布, terraform 发布, resource release, terraform provider release, publish resource, new resource, update resource, datasource release, data source release, 需求澄清, requirement clarification
---

# Terraform Provider Resource Release

Release a Terraform Provider resource together with its corresponding data source.

## Jarvis 自治边界（本仓部署时生效）

在 jarvis 无人值守语境下，本 skill 的自治边界如下（对应 `autonomy.md` 的 `stop: ["release_prod"]`）：

- **可自治**：全部开发/测试/开 PR/跑 ACC/回填 Aone 步骤，含轮询处理评审评论、按建议改码推分支。
- **必须 escalate**：最终 **PR merge** —— aliyun/terraform-provider-alicloud 是公共上游仓，merge 视同 `release_prod`，jarvis 不得自动合入。评审通过、CI 全绿、评论全清后停手，写 `escalation/` 并推送通知，等仓库主人人工合并。
- 因此原文"task is only complete after the PR is merged"是**对人跑者**的验收口径；对 jarvis，"PR 就绪（CI 绿+评论清）+ escalation 提交"即本轮收尾。

## 无人值守决策规则（替代"询问用户"）

jarvis 语境下，本 skill 所有"ask the user"节点按以下规则改走三通道：**规则内置**（有明确默认值自己决策+工单留痕）、**异步问工单**（在 Aone 评论提问，本轮释放，下轮 loop 续跑）、**escalation**（起草不发出）。

| 原询问点 | 无人值守行为 |
|---|---|
| 未提供工单号 | triage 场景工单天然存在；缺则按 `loops/adhoc-intake.md` 建单，不问人 |
| 镇元预发版 vs 线上版歧义 | AMP 两环境都查：线上存在用线上，否则用预发；判定理由评论回工单；两边元数据冲突 → low_conf escalate |
| 工单缺需求描述 | 在 Aone 评论 @提单人 补充描述，打 jarvis-idle 释放本单，本轮退出；下轮扫到更新自动续跑 |
| 修复方案需确认 | high_conf（OpenAPI+源码两层一致）→ 自动改+重跑 ACC，重试上限 3 次；low_conf → escalation |
| PR 评论无法自行解决 | 起草回复入 `escalation/`，不自动发出 |
| 等待 PR merge | 见"Jarvis 自治边界"——CI 绿+评论清即收尾，merge 是人工硬门 |
| Step 2 provider 仓本地路径 | 规则内置：`bootstrap/workspace.sh dir terraform_provider` 解析（CLAUDE.md #4），缺登记 → escalate(`missing_capability`)，不问人 |
| Step 9 镇元用例 vs 手写 | 工单已给镇元用例 ID 清单 → 走镇元生成；否则默认**手写 + 100% 属性覆盖**；拿不准 → 异步问工单 |
| Step 10.2 IDL source path / 请求 release IDL / generator 修复确认 | 异步问工单（@提单人/镇元侧），本轮 jarvis-idle 释放；generator 根因分析入 `escalation/`，不自动改生成器 |

**Applicable scenarios**
- **New resource release**: produce a new resource + data source + documentation via the generator
- **Existing resource update**: split into auto-generated and hand-written paths

**Key constraints**
- A single **Aone work item** MUST be linked; this skill is responsible for keeping the work item up to date
- All development, testing, and fixes happen inside an **isolated worktree** so concurrent tasks never pollute each other
- **Local compilation is strictly forbidden** — compiling the alicloud provider locally crashes the workstation. Run static checks only.
- For any generator-based path, **derive the concrete release scope by diffing AMP metadata against the local provider code** (requirement gap analysis) and **verify all orchestration APIs are OpenAPI** before generation; write both findings back to Aone
- Remote ACC acceptance tests must **all pass** before the PR is submitted
- **The task is only complete after the PR is merged** (CI passing + all comments closed)

**Final deliverables**
- resource code + data source code
- resource documentation + data source documentation
- test cases (generated or hand-written, 100% attribute coverage)
- all remote ACC cases PASS
- PR is **merged** (CI passed + all comments resolved)
- Aone work item updated to the final completed state, with the requirement-gap analysis and any non-OpenAPI APIs logged

---

## Workflow

### Step 1: Link an Aone Work Item

Every provider resource release MUST be linked to an Aone work item.

- **If the user has not provided a work item** → ask for the Aone link / ID. Do not proceed to later steps until the work item is confirmed.
- The work item is used end-to-end for:
  - extracting `product_code` / `resource_code` (Step 4)
  - identifying the target Zhenyuan version, staging or production (Step 5)
  - reading the business requirement description (Step 7 hand-written branch)
  - writing back progress at key milestones (Step 5 requirement-gap analysis, Step 5 non-OpenAPI log, Step 11, Step 12)

读工单：`bootstrap/aone-get.sh <id>`（或有 coop MCP 时用 `mcp__coop__query_workitem_detail`）读字段；回填进展/更新状态：`bootstrap/wrap.sh sync <id> "<进展>"`（或 coop MCP `add_comment`）。**a1/bootstrap 路径为主**（数字机器人无 coop MCP），MCP 为备。

### Step 2: Prepare the Provider Repo Worktree

**Do NOT ask the user for the repo path** — resolve it from the workspace registry:
`PROV=$(bootstrap/workspace.sh dir terraform_provider)`（本机覆盖走 `workspaces.local.json` / `JARVIS_WORKSPACE_ROOT`，见 CLAUDE.md 工作纪律 #4）。

Every release task creates a new worktree based on that repo, aligned with the upstream master（登记的 remote 布局里 **`origin` 即上游 aliyun**，`fork` 为 api-tool-agent fork，见 `config/workspaces.json` remotes）:

```bash
cd "$PROV"
git fetch origin
git worktree add -b feat/<aone-id>-<brief> ../<worktree_dir> origin/master
cd ../<worktree_dir>
```

**Why a worktree**: each Aone work item gets its own worktree. All subsequent development, testing, and fixes for this work item run inside that worktree, isolated from other parallel tasks.

### Step 3: Resolve the Generator Working Directory

**Do NOT ask the user** — the Terraform generator is a registered workspace `terraform_generator_v4`:
- Repo: `git@gitlab.alibaba-inc.com:opensource-tools/terraform-generator-v4.git` (branch `main`)
- Resolve its local path: `GEN=$(bootstrap/workspace.sh dir terraform_generator_v4)`
- **If absent, clone it there first**: `git clone git@gitlab.alibaba-inc.com:opensource-tools/terraform-generator-v4.git "$GEN"`
- **Before every generation use (Step 6 code-gen / Step 9 test-case gen), pull fresh** so you run the latest generator: `git -C "$GEN" pull --ff-only`

This path is used by the generation steps (Step 6) and by the test-case generation step (Step 9).

### Step 4: Decide Release Scope (new vs update)

Extract from the Aone work item:
- **product_code**
- **resource_code**

Check whether the resource already exists in the provider source tree:

```bash
ls alicloud/resource_alicloud_<product>_<resource>.go 2>/dev/null
```

- File **does not exist** → **New resource release** (path: Step 5 → Step 6)
- File **exists** → **Existing resource update**, decide auto-generated vs hand-written by inspecting the **first line** of the source file:
  - Auto-generated marker present → auto-generated update path (Step 5 → Step 6)
  - Marker absent → hand-written update path (skip Step 5, go to Step 7)

### Step 5: Pre-Generation Analysis (Requirement Gap + OpenAPI Check)

> **When to run**: any generator-based path (new resource, or auto-generated existing-resource update).
> **Skip when**: the resource is hand-written (Step 4 routed to Step 7).

This step uses the `amp-resource-metadata` skill. A single AMP Meta fetch fuels two parallel analyses — **5.3 requirement gap** and **5.4 OpenAPI check** — and both findings are written back to the Aone work item as separate comments.

#### 5.1 Determine the target Zhenyuan version

From the Aone work item, identify whether the resource to release is the **Zhenyuan staging (pre)** version or the **Zhenyuan production (online)** version. If the work item is ambiguous, ask the user.

#### 5.2 Fetch the resource Meta

Use `amp-resource-metadata` to fetch the resource schema (which contains attribute definitions and the CRUD API mappings):

```bash
python3 {AMP_SKILL_DIR}/scripts/get_resource_type.py \
  --service-code <product_code> \
  --resource-code <resource_code> \
  [--env pre]   # add --env pre for the staging version; omit for production
```

From the returned Meta, extract:
- the **attribute schema** (used by 5.3)
- the list of **APIs used in the Terraform orchestration** — Create, Read, Update, Delete (and List for the data source) (used by 5.4)

#### 5.3 Requirement Gap Analysis (Requirement Clarification)

Derive the concrete release scope by diffing the AMP Meta against the existing local Terraform provider code.

- **Existing-resource auto-gen update** (Step 4 routed file exists + auto-gen marker):
  - Read `alicloud/resource_alicloud_<product>_<resource>.go` and the data source file
  - List the attributes, valid values, and behaviors currently supported by the provider
  - Compute the gap: capabilities described in the AMP Meta that are **missing from the provider code**
  - Each gap is one concrete requirement of this release
- **New resource** (Step 4 routed file does not exist):
  - There is no existing provider code; the **entire AMP Meta** is the requirement scope
  - List every attribute / behavior from the Meta as a requirement entry

Write the gap analysis back to the Aone work item as a structured comment. Recommended columns per entry:
- `attribute / behavior`
- `type` (string / int / bool / list / object / API behavior / valid-value addition)
- `category` (new attribute, new valid value, new CRUD branch, etc.)
- `source in AMP Meta` (path / pointer)
- `recommended provider field name` (snake_case)

回填此需求集：`bootstrap/wrap.sh sync <id> "<gap 分析>"`（或有 coop MCP 时用 `mcp__coop__add_comment`）。**a1/bootstrap 路径为主**（数字机器人无 coop MCP），MCP 为备。This is the **concrete requirement set** that downstream Steps (6 generation, 9 tests, etc.) implement and verify against.

> If the gap is empty for an existing-resource update, there is nothing to release — report this to the user and stop early.

#### 5.4 OpenAPI Check

For every API extracted from the Meta in 5.2, use `amp-resource-metadata` to inspect its gateway / backend info and decide whether it is an OpenAPI:

```bash
python3 {AMP_SKILL_DIR}/scripts/get_runtime_api.py \
  --pop-code <PopCode> \
  --pop-version <PopVersion> \
  --api-name <ApiAction> \
  [--env pre]
```

Use the OpenAPI determination rule provided by `amp-resource-metadata` (PopCode + GatewayType + runtime backend signature) to classify the API as **OpenAPI** or **non-OpenAPI**.

#### 5.5 Log non-OpenAPI APIs to Aone

For every API classified as non-OpenAPI in 5.4, record a separate comment on the Aone work item with the columns: `API action`, `product / version`, `CRUD role` (Create / Read / Update / Delete / List), `reason it is not OpenAPI`.

> 5.3 (requirement gap) and 5.5 (non-OpenAPI APIs) are written as **two separate comments** so they can be reviewed independently. Neither blocks downstream steps when findings are present — the release proceeds, and the visibility signal is preserved on the work item.

### Step 6: Generator Path — Generate Code and Docs

> **Applies to**: new resource release **and** auto-generated existing-resource update (path decided in Step 4).

Run the generator using the version identified in Step 5.1 (Zhenyuan staging or production). The generator produces in a single pass:
- resource code
- **data source code**
- documentation (resource + data source)

**This step does NOT produce test cases.** Test cases are handled separately in Step 9.

### Step 7: Hand-written Update Path

> **Applies to**: existing-resource update routed to hand-written (Step 4 detected no auto-generated marker).

Check whether the Aone work item contains a **business requirement description** for the release:

- **Description present** → develop the provider code by hand according to that requirement
- **Description absent** → ask the user to contact the cloud product developer to add the requirement into the Aone work item, **pause subsequent steps** until the requirement is in place

### Step 8: Static Check (Local Compilation Forbidden)

> **Critical warning**: **compiling the alicloud provider locally crashes the workstation**.
> NEVER run `go build`, `go install`, `make`, or any other command that triggers a full provider build.

Run static checks only on the generated or hand-written code（工作区登记命令见 `config/workspaces.json` ops）:

```bash
gofmt -l alicloud/        # ops.fmt —— 输出为空即过
go vet ./...              # ops.vet —— 静态分析
```

### Step 9: Generate / Write Test Cases

**Ask the user** whether to use **Zhenyuan platform (镇元)** test cases to generate Terraform test cases:

- **Use Zhenyuan platform cases**
  - The user must provide a **list of test case IDs**
  - Run the generator against those specific IDs to produce the Terraform test cases

- **Hand-written cases**
  - Write test cases by hand based on the provider changes
  - **Attribute coverage MUST reach 100%** (hard requirement)
  - **All critical logic introduced by the provider change MUST be covered**

> **Do not skip data source tests**: if a data source was generated or written, the corresponding data source test cases MUST also be added.

### Step 10: Remote ACC Acceptance Tests

Run Terraform remote acceptance tests and execute **every** acceptance case once (use the `invoke-terraform-acc-test-remote` skill).

If any case fails, diagnose the root cause and self-fix using the branches below. This step is only complete when **every** acceptance case passes.

#### Step 10.1: Fixing failures for hand-written resources

Locate the provider bug → fix the provider code → re-run the tests → PASS.

#### Step 10.2: Fixing failures for auto-generated cases

First determine the failure category — is it a **resource definition** problem or a **generator** problem? Report the conclusion and supporting analysis to the user.

- **Resource definition problem**
  - Ask the user for the **IDL source path** of the resource definition
  - Attempt to fix the IDL resource definition yourself
  - Push the fix to the remote
  - Ask the user to **release the IDL**

- **Generator problem**
  - Provide a detailed analysis of the root cause
  - Recommend a fix
  - Wait for the user to confirm and apply the fix

### Step 11: Submit the PR and Update Progress

#### 11.1 Submit the PR

> **⚠️ CRITICAL — Sanitize all GitHub-facing content**: the target repo is **public on GitHub**. NEVER include any internal information in the PR **title**, **body**, **commit messages**, or **code comments**. Forbidden categories:
> - **Aone** work item IDs, links, or any mention of Aone / 内部工单
> - **Claude / AI assistant** attribution — no `Co-Authored-By: Claude`, no "generated by Claude", no "AI-assisted", no tool-bot signatures
> - **Internal personnel names** — developers, reviewers, PMs, or anyone surfaced from internal systems (Aone assignee, Code Review reviewers, etc.); 内部聊天记录 / OKR / 项目代号同禁
> - **Customer information** — customer names, account UIDs, opportunity / contract IDs, ticket numbers, business context tied to a specific customer
> - **Diagnostic details** — customer instance IDs (`r-xxx` / `i-xxx` / `lb-xxx` / `s-xxx` 等), `RequestId` literals, machine names / RAM user names from error output
>
> When summarizing the change for the PR, describe only the **public-facing technical capability** (resource, attributes, behavior). If the Aone work item contains internal motivation or customer context, restate the change in product-neutral language. **Scrub before every push** — including iteration commits during the comment-resolution loop in Step 12: run `bash <jarvis仓>/bootstrap/pre-push-sanitize.sh`（禁品正则真源）from the worktree, then read `git log -p origin/master..HEAD` to catch customer names the script cannot enumerate.

> **Identity gate（CLAUDE.md 工作纪律 #6）**: all GitHub writes go through `bootstrap/github-identity.sh` — run `check` first (token account MUST be `api-tool-agent`); commit via `github-identity.sh commit`, push via `github-identity.sh push <owner/repo> <local-ref> <remote-ref>` (PR head MUST be `api-tool-agent:<branch>`), `gh` writes via `github-identity.sh gh ...`. Never fall back to ambient `gh auth` / local git identity.

> **Note**: do **NOT** write the CHANGELOG when submitting the PR. The CHANGELOG is written by the release engineer in a later step — do not edit `CHANGELOG.md`.

The PR body MUST list the **passing test cases** in this format:

```
=== RUN TestAccAliCloudxxxx_basic0
--- PASS: TestAccAliCloudxxxx_basic0 (144.21s)

=== RUN TestAccAliCloudxxxx_basic1
--- PASS: TestAccAliCloudxxxx_basic1 (108.57s)

=== RUN TestAccAliCloudxxxx_basic2
--- PASS: TestAccAliCloudxxxx_basic2 (3.88s)
```

#### 11.2 Ensure PR CI Passes

After the PR is submitted, monitor all CI tasks:
- **All pass** → proceed to 11.3
- **Any failure** → fix and re-push until every CI task passes

#### 11.3 Update Aone Progress

Once the PR CI is fully green, **write the current progress back into the Aone work item** (PR link, test results, status, etc.).

### Step 12: Poll PR Comments Until Merge

After the PR is submitted, **continuously poll PR comments** until the PR is **merged**. Only then is the release task complete.

Fetch comments periodically via `gh pr view --comments`（只读可直用 gh）。

#### 12.1 Poll for comments

Periodically check for new PR comments (reviewer comments, AI bot comments, CI bot comments, etc.). A new comment → go to 12.2.

#### 12.2 Resolve comments

- **Can resolve yourself** → modify the code → push to the PR branch（经 `github-identity.sh push`，push 前重跑 sanitize 自查）→ reply to the comment confirming the fix
- **Cannot resolve yourself** → **ask the user**, wait for the decision, then act

#### 12.3 Loop until merged

Repeat 12.1 / 12.2 until the PR is **merged**.

#### 12.4 Final wrap-up

Once the PR is merged:
- Update the Aone work item to the final completed state (status, merged PR link, commit info, etc.)
- The release task is officially complete

---

## Important Notes

1. **Never compile locally** — building the alicloud provider locally will crash the workstation; static checks only.
2. **Aone work item is both source and sink** — every task must be linked to a work item, and progress MUST be written back at each milestone (requirement-gap analysis, non-OpenAPI APIs detected, development complete, tests passed, PR CI green, PR merged).
3. **Worktree isolation** — one worktree per Aone work item; concurrent tasks must not contaminate each other.
4. **Step 5 is mandatory for generator paths** — never skip the requirement-gap analysis (5.3) or the OpenAPI check (5.4) for new resources or auto-generated existing-resource updates. Both produce visibility signals on the Aone work item even when they do not block progress.
5. **Requirements come from the AMP-vs-code diff (generator paths) or the Aone work item (hand-written path)** — never invent requirements from outside these sources.
6. **Generation ≠ test generation** — the generator path in Step 6 produces only code + docs; test cases are handled independently in Step 9.
7. **100% attribute coverage** — a hard requirement when writing test cases by hand. Do not enter Step 10 without it.
8. **Do not skip data source tests** — if a data source exists, it must have its own tests.
9. **Submitting the PR is not the finish line** — comments MUST be polled and resolved until the PR is merged.
10. **Do not touch CHANGELOG.md** — the release engineer writes it later.
11. **Sanitize all GitHub-facing content (CRITICAL)** — the alicloud provider repo is **public on GitHub**. Never expose internal information in PR title, PR body, commit messages, or code comments. Forbidden: Aone references, Claude/AI attribution, internal personnel names, customer information. See Step 11.1 for full rules. This applies to **every** push, including iteration commits during the Step 12 comment loop.

## Acceptance Criteria

- [ ] Aone work item is linked and updated to the merged/completed state
- [ ] Worktree is created and based on the latest upstream master（`origin/master`，登记布局 origin=上游 aliyun）
- [ ] For generator-based paths: requirement-gap analysis (Step 5.3) and OpenAPI check (Step 5.4) executed; both findings logged to Aone as separate comments
- [ ] Resource + data source code and documentation are all produced
- [ ] Test cases are complete (generated or hand-written), with 100% attribute coverage
- [ ] Data source test cases are in place
- [ ] **All** remote ACC acceptance cases PASS
- [ ] PR is submitted; the body contains the passing test list in the required format
- [ ] All CI tasks on the PR pass
- [ ] All PR comments are resolved
- [ ] PR is **merged**
- [ ] All GitHub-facing content (PR title/body, commit messages, code comments) is sanitized — no Aone references, no Claude/AI attribution, no internal personnel names, no customer information

---

## Rules

手写资源/手写用例的硬编码规范（`GetOkExists`、数组类型断言、命名约定等）单点维护在
[`provider-resource-review/SKILL.md`](../provider-resource-review/SKILL.md) §6 Code Bug Patterns —— 开发与评审共用同一套，此处不重复。

---

## 发版前强化门禁（可选,additive）

> **本节为 jarvis 侧 additive 指引,不改动上述任何既有步骤;是否启用由跑者按发版风险自定。**

提 PR / 合并前,可选跑一条**发版前全量过闸**门禁(jarvis F3 探测线),把 tf-probe 的 tier-0（三方一致性全量）
+ tier-1（全场景生命周期）串成一次红/黄/绿判定:

```bash
bootstrap/rc-gate.sh <provider-dir> [--quick]   # 绿/黄=退0过闸, 红=退1阻断, 不可判=退2
```

- 🟢绿/🟡黄 → 照常进入本 SOP 的远程 ACC / PR / 人工合并环节（黄项知情放行）。
- 🔴红（api_gap S3+ / 场景 fail / destroy 残留）→ 回「修复方案」修红项,复跑转绿/黄再提交。
- 门禁**不替代**远程 ACC（互补:ACC 深验单资源,门禁广验全量+全场景）,**不触碰** `release_prod`。

详见 [`references/rc-gate.md`](references/rc-gate.md);完整读法与实现见 jarvis 仓 `loops/tf-probe.md`「四点五、RC 门禁」+ `bootstrap/rc-gate.sh`。
