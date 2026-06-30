# GitHub Token PR Flow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enforce Jarvis-owned GitHub PR creation through `JARVIS_GITHUB_TOKEN` and document the Terraform Provider internal Aone split-workflow.

**Architecture:** Add one focused shell guard that verifies the GitHub token identity and can either check-only or execute a `gh` command under `GH_TOKEN=$JARVIS_GITHUB_TOKEN`. Wire the guard into preflight verification and update the Jarvis runbooks/skills that describe GitHub PR and Terraform Provider resource handoff behavior.

**Tech Stack:** Bash, GitHub CLI `gh`, existing shell test harness, Markdown runbooks, Aone workitem process docs.

## Global Constraints

- Work happens only in a linked worktree branch, never directly on `master`.
- Jarvis GitHub PR actions and branch pushes must use `JARVIS_GITHUB_TOKEN`; ambient `gh auth` or git credentials are not acceptable.
- The token login must be exactly `api-tool-agent`.
- Terraform Provider PR heads created by Jarvis must be under `api-tool-agent:<branch>`.
- Missing token, failing `gh api user`, or a non-`api-tool-agent` login must block and escalate.
- `.agents/skills` and `.claude/skills` copies must remain in sync where both exist.
- External PR/MR/Aone text must not contain AI signatures or watermarks.

---

### Task 1: GitHub Identity Guard

**Files:**
- Create: `bootstrap/github-identity.sh`
- Create: `bootstrap/tests/github_identity.sh`

**Interfaces:**
- Produces: `bootstrap/github-identity.sh check`
- Produces: `bootstrap/github-identity.sh gh <args...>`
- Produces: `bootstrap/github-identity.sh push <owner/repo> <local-ref> <remote-ref>`

- [ ] **Step 1: Write the failing test**

Create `bootstrap/tests/github_identity.sh` with cases for:
- missing `JARVIS_GITHUB_TOKEN` exits non-zero and mentions `JARVIS_GITHUB_TOKEN`
- wrong login from `gh api user --jq .login` exits non-zero and mentions `api-tool-agent`
- correct login passes and invokes `gh api user`
- `gh` subcommand runs the supplied command with `GH_TOKEN` set to `JARVIS_GITHUB_TOKEN`
- `push` subcommand runs `git push` with a token-backed askpass helper and no ambient git credentials

- [ ] **Step 2: Run test to verify it fails**

Run: `bash bootstrap/tests/github_identity.sh`
Expected: fail because `bootstrap/github-identity.sh` does not exist.

- [ ] **Step 3: Write minimal implementation**

Implement `bootstrap/github-identity.sh` with:
- expected login hard-coded to `api-tool-agent`
- `check` mode: verify env var, call `GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh api user --jq .login`, compare to expected login
- `gh` mode: call check first, then `GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh "$@"`
- `push` mode: call check first, then push `local-ref:remote-ref` to `https://github.com/<owner/repo>.git` using `GIT_ASKPASS` with username `x-access-token` and password `JARVIS_GITHUB_TOKEN`
- clear stderr messages and non-zero exits on failure

- [ ] **Step 4: Run test to verify it passes**

Run: `bash bootstrap/tests/github_identity.sh`
Expected: pass all cases.

### Task 2: Preflight Coverage And PR Flow Docs

**Files:**
- Modify: `bootstrap/verify.sh`
- Modify: `test/verify_test.sh`
- Modify: `AGENTS.md`
- Modify: `loops/adhoc-intake.md`
- Modify: `.agents/skills/terraform-pr-review/SKILL.md`
- Modify: `.claude/skills/terraform-pr-review/SKILL.md`
- Modify: `.agents/skills/provider-resource-dev/SKILL.md`
- Modify: `.claude/skills/provider-resource-dev/SKILL.md`
- Modify: `config/workspaces.json`

**Interfaces:**
- Consumes: `bootstrap/github-identity.sh check`

- [ ] **Step 1: Write/update failing tests**

Extend `test/verify_test.sh` so a stubbed `bootstrap/github-identity.sh` failure causes `bootstrap/verify.sh` to print `FAIL jarvis-github-token` and exit non-zero.

- [ ] **Step 2: Run test to verify it fails**

Run: `bash test/verify_test.sh`
Expected: fail until `bootstrap/verify.sh` calls the new guard.

- [ ] **Step 3: Wire verification**

Update `bootstrap/verify.sh` to run `bootstrap/github-identity.sh check` as `PASS/FAIL jarvis-github-token`.

- [ ] **Step 4: Update GitHub PR flow docs**

Document that GitHub PR create/comment actions must use:

```bash
bootstrap/github-identity.sh check
bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud HEAD <branch>
bootstrap/github-identity.sh gh pr create --repo aliyun/terraform-provider-alicloud ...
```

Replace stale personal-fork / ambient `gh pr` wording for Jarvis-owned PR paths with `api-tool-agent` fork/head requirements.

- [ ] **Step 5: Run tests**

Run:
- `bash test/verify_test.sh`
- `bash bootstrap/tests/github_identity.sh`
- `bash test/provider_resource_dev_skill_sync_test.sh`

Expected: all pass.

### Task 3: Terraform Provider Internal Aone Split Workflow

**Files:**
- Modify: `loops/adhoc-intake.md`
- Modify: `loops/aone-triage.md`
- Modify: `.agents/skills/aone-triage/SKILL.md`
- Modify: `.claude/skills/aone-triage/SKILL.md`
- Modify: `.agents/skills/provider-resource-dev/SKILL.md`
- Modify: `.claude/skills/provider-resource-dev/SKILL.md`
- Modify: `config/pools.json` only if the existing tf_provider pool does not already encode project `528766`

**Interfaces:**
- Consumes: existing `tf_provider` pool project `528766`

- [ ] **Step 1: Add a doc sync test**

Add a test that asserts both provider-resource-dev skill copies mention:
- `terraform-alicloud`
- `528766`
- `WORKER_1782379562571`
- `双向关联`
- `客户主单`

- [ ] **Step 2: Run test to verify it fails**

Run the new test.
Expected: fail before docs are updated.

- [ ] **Step 3: Update runbooks and skills**

Add the rule:
- for non-automated Terraform Provider resource development, create or reuse an internal `terraform-alicloud` Aone item in project `528766`
- assign it to `WORKER_1782379562571`
- link it bidirectionally to the customer main item
- put detailed development, validation, PR, CI, and acceptance notes on the internal item
- keep the customer main item to brief key milestones and blockers
- sync detailed dependency questions to cloudspec_gap or product upstream items when those exist

- [ ] **Step 4: Run tests**

Run:
- new doc sync test
- `bash test/provider_resource_dev_skill_sync_test.sh`
- `bash test/aone_triage_templates_sync_test.sh`
- `bootstrap/verify.sh`

Expected: all pass.
