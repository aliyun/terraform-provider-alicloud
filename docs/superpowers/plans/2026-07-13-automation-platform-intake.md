# Automation Platform Intake Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Register the Automation Platform Aone pool and complete the workspace, triage, delivery, SLS, and board context needed for Jarvis to process its workitems.

**Architecture:** Keep `config/pools.json` and `config/workspaces.json` as the factual sources. Existing generic scan/dispatch code consumes the new pool without a bridge branch; `aone-triage` selects a dedicated delivery reference, and the HTML board derives its filters from pool configuration instead of a hard-coded list.

**Tech Stack:** Bash, jq, JSON, Python embedded in shell, Markdown skills, Aone CLI metadata.

## Global Constraints

- Automation Platform Aone project is exactly `1091779`.
- Scheduled scanning uses only assignee `WORKER_1782379562571` (`open-jarvis`); direct Aone URL/ID handling is not assignee-restricted.
- Main application is `aliyun-automation-platform`, app ID `172823`, repo ID `2156624`.
- Delivery pipelines are prestage `66` and production `67`; production remains human-gated.
- Register exactly these six repositories: backend, frontend, runtime, integration tests, public IaCService API, and internal IaCService API.
- Include SLS contexts `systemlog-prod`, `systemlog-pre`, and `aliyun-automation-platform-1252907582134651-pop-aliyun-cn` through the existing `sls-log-query-aliyun-automation-platform` skill.
- Do not register or deliver Agent portal, AgentRuntime, PlayGround, FC sandbox, WebSocket/STS orchestration, or `aliyun-automation-agent` as part of this pool.
- Do not store historical feature branches or temporary flow IDs.
- Modify only the isolated branch `worktree-automation-platform-context`; never push or merge master.
- Write a failing regression test before every production/configuration behavior change and observe the expected failure.

---

### Task 1: Register the Automation Platform pool and scan contract

**Files:**
- Modify: `test/pools_config_test.sh`
- Modify: `test/scan_test.sh`
- Modify: `config/pools.json`

**Interfaces:**
- Consumes: `scan.sh` dynamic `.pools | to_entries[]` traversal and per-pool `assignee` semantics.
- Produces: pool key `automation_platform`, project `1091779`, and scan items with `pool=automation_platform` and `pool_project=1091779`.

- [ ] **Step 1: Write failing pool configuration assertions**

Update the count and ID/assignee loops in `test/pools_config_test.sh`, then add exact assertions:

```bash
[ "$(jq '.pools|length' "$POOLS_JSON")" = "5" ] && ok "5 pools" || bad "expected 5 pools"
[ "$(jq '.lines|length' "$POOLS_JSON")" = "3" ] && ok "3 lines" || bad "expected 3 lines"

for id in 1086837 528766 2124589 2100304 1091779; do
  jq -e --argjson p "$id" '[.pools[].project]|index($p)' "$POOLS_JSON" >/dev/null \
    && ok "project $id present" || bad "project $id missing"
done

for pool in tf_customer tf_provider mcp_server api_toolkit automation_platform; do
  [ "$(jq -r ".pools.$pool.assignee" "$POOLS_JSON")" = "WORKER_1782379562571" ] \
    && ok "$pool assignee" || bad "$pool assignee mismatch"
done

jq -e '.pools.automation_platform.line=="automation_platform"' "$POOLS_JSON" >/dev/null \
  && ok "automation platform independent line" || bad "automation platform line"
jq -e '.pools.automation_platform.apps[0] == {
  "app":172823,
  "repo_id":2156624,
  "name":"aliyun-automation-platform",
  "repo":"aliyun-automation-platform",
  "pipelines":{"prestage":66,"prod":67},
  "delivery":"delivery-aliyun-automation-platform.md"
}' "$POOLS_JSON" >/dev/null && ok "automation platform app facts" || bad "automation platform app facts"

[ "$(jq -r '.pools.api_toolkit.done_status["产品类需求"]' "$POOLS_JSON")" = "已发布" ] \
  && ok "api_toolkit.done_status.产品类需求" || bad "api_toolkit done key"
[ "$(jq -r '.pools.automation_platform.progress_status["产品类需求"]' "$POOLS_JSON")" = "开发中" ] \
  && ok "platform req progress" || bad "platform req progress"
[ "$(jq -r '.pools.automation_platform.done_status["任务"]' "$POOLS_JSON")" = "已完成" ] \
  && ok "platform task done" || bad "platform task done"
```

- [ ] **Step 2: Write a failing real-pool scan regression**

Append Test 14 to `test/scan_test.sh`. Copy the repository's current `config/pools.json` into a temporary `JARVIS_ROOT`, so the test fails before the production pool is registered and passes afterward. Its stub must record all list arguments, return a workitem only for project `1091779`, and return an empty array for every other project.

Assert that all three req/bug/task calls contain `--project 1091779`, contain `assignedTo=WORKER_1782379562571`, do not contain `--assignee`, and that returned rows contain:

```bash
jq -e 'all(.[]; .pool=="automation_platform" and .pool_project=="1091779")' "$scan_output"
```

- [ ] **Step 3: Run tests and verify RED**

Run:

```bash
bash test/pools_config_test.sh
bash test/scan_test.sh
```

Expected: `pools_config_test.sh` fails because the fifth pool and third line do not exist; Test 14 fails because the real pool is absent from the repository configuration contract.

- [ ] **Step 4: Add the minimal pool configuration**

Add `lines.automation_platform = "Automation Platform"` and a pool containing the exact app block from Step 1, plus:

```json
"project": 1091779,
"name": "自动化服务台",
"line": "automation_platform",
"desc": "自动化服务台 / IaCService 产品研发与交付；不含 Agent 链路",
"assignee": "WORKER_1782379562571",
"dev": true,
"done_status": {
  "产品类需求": "已发布",
  "功能缺陷": "Fixed",
  "线上问题": "Fixed",
  "任务": "已完成"
},
"progress_status": {
  "产品类需求": "开发中",
  "功能缺陷": "Open",
  "线上问题": "Open",
  "任务": "处理中"
},
"exclude_status": [
  "已发布", "已取消", "已完成", "Fixed", "Closed", "Won'tfix",
  "Worksforme", "Duplicate", "Invalid", "External", "ByDesign"
],
"exclude_title": []
```

Add a routing entry matching `自动化服务台`, `aliyun-automation-platform`, and `IaCService` to `automation_platform`.

- [ ] **Step 5: Run tests and verify GREEN**

Run the two commands from Step 3. Expected: both exit 0, Test 14 reports the fixed-assignee scan contract, and the pre-existing `api_toolkit` key mismatch is gone.

- [ ] **Step 6: Commit**

```bash
git add config/pools.json test/pools_config_test.sh test/scan_test.sh
git commit -m "feat: register automation platform pool"
```

---

### Task 2: Register all Automation Platform workspaces

**Files:**
- Modify: `test/workspaces_config_test.sh`
- Modify: `config/workspaces.json`

**Interfaces:**
- Consumes: `bootstrap/workspace.sh dir "$WORKSPACE_KEY"` and `JARVIS_WORKSPACE_ROOT`.
- Produces: six workspace keys sharing pool `automation_platform`.

- [ ] **Step 1: Write failing workspace assertions**

Add jq assertions for these exact records:

```text
automation_platform:
  repo=aliyun-automation-platform
  git_url=git@gitlab.alibaba-inc.com:aliyun-automation-platform/aliyun-automation-platform.git
  project=1091779 app=172823 repo_id=2156624
  pipelines.prestage=66 pipelines.prod=67
  delivery=delivery-aliyun-automation-platform.md

automation_platform_frontend:
  repo=iac-service
  git_url=git@gitlab.alibaba-inc.com:aliyun-api/iac-service.git

automation_platform_runtime:
  repo=iac-service-runtime
  git_url=git@gitlab.alibaba-inc.com:opensource-tools/iac-service-runtime.git

automation_platform_function_test:
  repo=automation-function-test
  git_url=git@gitlab.alibaba-inc.com:aliyun-automation-platform/automation-function-test.git

automation_platform_api:
  repo=IaCService_pop_IaCService_2021-08-06
  git_url=git@gitlab.alibaba-inc.com:cloudspec-model/IaCService_pop_IaCService_2021-08-06.git

automation_platform_api_inner:
  repo=IaCService-inner_pop_IaCService-inner_2021-09-01
  git_url=git@gitlab.alibaba-inc.com:cloudspec-model/IaCService-inner_pop_IaCService-inner_2021-09-01.git
```

Assert each record has `default_branch == "master"` and `pools | index("automation_platform") != null`. Create matching directories beneath the test `tmpdir`, call `workspace.sh dir` for every key, and compare the resolved path to `"$tmpdir/$repo_dir"`.

- [ ] **Step 2: Run test and verify RED**

Run `bash test/workspaces_config_test.sh`.

Expected: FAIL on the missing frontend/runtime/test/API workspace assertions.

- [ ] **Step 3: Add minimal workspace records**

Use the exact records from Step 1. Set operations as follows:

```json
"automation_platform": {
  "ops": {"build":"./mvnw -q -DskipTests package","test":"./mvnw -q test"}
},
"automation_platform_frontend": {
  "ops": {"build":"npm run build","test":"npm test -- --runInBand","lint":"npm run lint"}
},
"automation_platform_runtime": {
  "ops": {"build":"go build ./...","test":"go test ./..."}
},
"automation_platform_function_test": {
  "ops": {"test":"mvn -q test"}
},
"automation_platform_api": {
  "ops": {"build":"cloudspec build"}
},
"automation_platform_api_inner": {
  "ops": {"build":"cloudspec build"}
}
```

Descriptions must state each repository's role. Only the backend record carries project/app/repo_id/pipelines/delivery.

- [ ] **Step 4: Run test and verify GREEN**

Run `bash test/workspaces_config_test.sh`. Expected: `workspaces_config_test: PASS`.

- [ ] **Step 5: Commit**

```bash
git add config/workspaces.json test/workspaces_config_test.sh
git commit -m "feat: register automation platform workspaces"
```

---

### Task 3: Add Automation Platform triage and delivery guidance

**Files:**
- Modify: `test/aone_triage_templates_sync_test.sh`
- Modify: `.claude/skills/aone-triage/SKILL.md`
- Create: `.claude/skills/aone-triage/references/delivery-aliyun-automation-platform.md`
- Modify: `.agents/skills/aone-triage/SKILL.md` via mirror script
- Create: `.agents/skills/aone-triage/references/delivery-aliyun-automation-platform.md` via mirror script
- Modify: `loops/adhoc-intake.md`

**Interfaces:**
- Consumes: `space=1091779`, workspace keys from Task 2, pipelines 66/67, and the existing SLS skill.
- Produces: explicit routing to `delivery-aliyun-automation-platform.md` for scheduled and direct-ID workitems.

- [ ] **Step 1: Write failing mirror and semantic assertions**

Add `references/delivery-aliyun-automation-platform.md` to the mirror loop in `test/aone_triage_templates_sync_test.sh`. Add grep assertions requiring the canonical SKILL/reference to contain all of:

```text
references/delivery-aliyun-automation-platform.md
1091779
automation_platform
172823
prestage 66
prod 67
systemlog-prod
systemlog-pre
aliyun-automation-platform-1252907582134651-pop-aliyun-cn
sls-log-query-aliyun-automation-platform
```

Also assert the reference states that Agent portal, AgentRuntime, PlayGround, and `aliyun-automation-agent` are excluded, and that production requires explicit human approval.

- [ ] **Step 2: Run test and verify RED**

Run `bash test/aone_triage_templates_sync_test.sh`.

Expected: FAIL because the canonical reference and its Codex mirror do not exist.

- [ ] **Step 3: Create the canonical delivery reference**

Create these sections with exact coordinates and commands:

```markdown
# 交付链路：aliyun-automation-platform

## 路由与边界
project `1091779`; app `172823`; backend repo ID `2156624`.
This path includes backend/frontend/runtime/integration/public API/internal API workspaces.
It explicitly excludes Agent portal, AgentRuntime, PlayGround, FC sandbox, WebSocket/STS orchestration, and `aliyun-automation-agent`.

## Repository selection
Map backend, frontend, runtime, integration, public API, and internal API changes to the six Task 2 workspace keys.
Request/response/POP Action/resource definition changes must check the relevant IaCService API repository and run `cloudspec build`.

## Worktree and CR/MR
Set `WORKSPACE_KEY` to the selected Task 2 key, resolve it with `bash bootstrap/workspace.sh dir "$WORKSPACE_KEY"`, create a branch worktree from the registered default branch, run registered ops, push only the feature branch, and sync the CR/MR link to Aone.

## Prestage and production
`bin/a1id -- app link 172823`
`bin/a1id -- app cr submit "$CR_ID" --pipeline-id 66`
Stop after prestage deployment and business/SLS verification.
Only after explicit human approval: `bin/a1id -- app cr submit "$CR_ID" --pipeline-id 67`.

## SLS
Use `sls-log-query-aliyun-automation-platform` for `systemlog-pre`, `systemlog-prod`, and `aliyun-automation-platform-1252907582134651-pop-aliyun-cn`.
```

Use the named variables `CR_ID` and `WORKSPACE_KEY` exactly as shown so the committed document contains executable command patterns rather than unresolved placeholders.

- [ ] **Step 4: Update canonical triage routing and adhoc intake**

In `.claude/skills/aone-triage/SKILL.md`:

- add the new reference to frontmatter;
- add `1091779 automation_platform` to the space routing table before generic app rows;
- add the reference to the self-delivery list.

In `loops/adhoc-intake.md`, change the selectable pool list to `tf_provider / tf_customer / mcp_server / automation_platform / api_toolkit`, and describe workspace resolution as configuration-driven rather than listing only two keys.

- [ ] **Step 5: Generate Codex mirrors**

Run:

```bash
bash bootstrap/mirror.sh to-codex .claude/skills/aone-triage/SKILL.md
bash bootstrap/mirror.sh to-codex .claude/skills/aone-triage/references/delivery-aliyun-automation-platform.md
```

- [ ] **Step 6: Run tests and verify GREEN**

Run:

```bash
bash test/aone_triage_templates_sync_test.sh
bash bootstrap/mirror.sh check
```

Expected: both commands exit 0 and the sync test prints `aone_triage_templates_sync_test: PASS`.

- [ ] **Step 7: Commit**

```bash
git add .claude/skills/aone-triage .agents/skills/aone-triage loops/adhoc-intake.md test/aone_triage_templates_sync_test.sh
git commit -m "docs: add automation platform delivery flow"
```

---

### Task 4: Render board pools dynamically and run integration verification

**Files:**
- Modify: `test/board_probe_test.sh`
- Modify: `bootstrap/board-html.sh`
- Modify: `docs/flow.md`

**Interfaces:**
- Consumes: ordered `.pools` object and pool names from `config/pools.json`.
- Produces: a visible filter row and counts for every configured pool, including `automation_platform`.

- [ ] **Step 1: Write failing board assertions**

Extend the existing board HTML test after `HTMLF` is generated:

```bash
grep -q 'data-pf="automation_platform"' "$HTMLF" \
  && ok "html includes automation_platform filter" \
  || bad "html omits automation_platform filter"
grep -q '<span class="nm">自动化服务台</span>' "$HTMLF" \
  && ok "html uses configured platform name" \
  || bad "html omits configured platform name"
```

- [ ] **Step 2: Run test and verify RED**

Run `bash test/board_probe_test.sh`.

Expected: FAIL because `board-html.sh` still hard-codes four pool keys.

- [ ] **Step 3: Derive pools from configuration**

In the embedded Python, load configuration before calculating totals and replace the hard-coded pool list:

```python
cfg=json.load(open(sys.argv[3]))
POOLS=list(cfg.get("pools",{}).keys())
NAMES={k:v.get("name",k) for k,v in cfg.get("pools",{}).items()}
```

Keep the existing `COLOR.get(k,"#98a2b3")` fallback so a newly configured pool renders without a code change. Remove the duplicate later `cfg` load. Change `docs/flow.md` from a numeric pool count to “并行扫描 `config/pools.json` 中配置的工作池”.

- [ ] **Step 4: Run targeted tests and verify GREEN**

Run:

```bash
bash test/board_probe_test.sh
bash test/pools_config_test.sh
bash test/scan_test.sh
bash test/workspaces_config_test.sh
bash test/aone_triage_templates_sync_test.sh
bash bootstrap/mirror.sh check
```

Expected: every command exits 0; board output includes the new filter and configured name.

- [ ] **Step 5: Validate syntax and diff hygiene**

Run:

```bash
jq -e . config/pools.json
jq -e . config/workspaces.json
bash -n bootstrap/board-html.sh
bash -n test/pools_config_test.sh
bash -n test/scan_test.sh
bash -n test/workspaces_config_test.sh
bash -n test/aone_triage_templates_sync_test.sh
bash -n test/board_probe_test.sh
git diff --check master...HEAD
```

Expected: all commands exit 0 with no diff-check output.

- [ ] **Step 6: Commit**

```bash
git add bootstrap/board-html.sh test/board_probe_test.sh docs/flow.md
git commit -m "fix: render board pools from config"
```
