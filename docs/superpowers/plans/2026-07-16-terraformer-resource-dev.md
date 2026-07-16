# Terraformer Resource Development Skill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a first usable `terraformer-resource-dev` Skill to Jarvis and update the open design CR with the implementation.

**Architecture:** Keep `.claude/skills/terraformer-resource-dev` as the canonical source and generate the `.agents/skills` mirror with `bootstrap/mirror.sh`. Keep the entrypoint concise and move Terraformer-specific discovery, Import ID, pagination, file-selection, and validation details into one reference. Enforce the agreed boundaries with a shell contract test.

**Tech Stack:** Markdown Skill files, Bash contract tests, Jarvis mirror tooling, Python Skill validator.

## Global Constraints

- Final Skill layout contains only `SKILL.md` and `references/alicloud-resource-development.md`; do not add `agents/openai.yaml`.
- `.claude/skills/terraformer-resource-dev` is canonical; `.agents/skills/terraformer-resource-dev` is produced with `bootstrap/mirror.sh to-codex`.
- Treat parent-resource listing as only discovery pattern C; a multipart Import ID alone does not imply parent traversal.
- Derive Import ID segment count, order, and delimiter from Terraform Provider `d.SetId`, `ParseResourceId`, Import docs, and Import tests.
- Data Source may require a parent ID, but Terraformer must discover it only when the selected child List API requires parent scope.
- Do not produce, infer, or maintain resource relationships; only read and consume the unified relationship artifact when it explicitly contains the resource.
- Do not modify `/Users/shanye/programs/terraformer`; it is read-only evidence for this Skill implementation.
- Preserve the existing CR worktree and update CR 28627638; do not merge or delete the worktree.

---

### Task 1: Add the Skill contract and implementation

**Files:**
- Create: `test/terraformer_resource_dev_skill_rules_test.sh`
- Create: `.claude/skills/terraformer-resource-dev/SKILL.md`
- Create: `.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md`
- Create via mirror: `.agents/skills/terraformer-resource-dev/SKILL.md`
- Create via mirror: `.agents/skills/terraformer-resource-dev/references/alicloud-resource-development.md`
- Modify: `docs/superpowers/specs/2026-07-16-terraformer-resource-dev-design.md`

**Interfaces:**
- Consumes: `bootstrap/mirror.sh`, `bootstrap/skills-mirror-lib.sh`, workspace keys `terraformer` and `terraform_provider`, existing `provider-resource-dev` governance conventions.
- Produces: a discoverable `terraformer-resource-dev` Skill and a deterministic contract test named `test/terraformer_resource_dev_skill_rules_test.sh`.

- [ ] **Step 1: Write the failing contract test**

Create `test/terraformer_resource_dev_skill_rules_test.sh` with this content:

```bash
#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

for rel in \
  "SKILL.md" \
  "references/alicloud-resource-development.md"; do
  claude_file="$repo_root/.claude/skills/terraformer-resource-dev/$rel"
  codex_file="$repo_root/.agents/skills/terraformer-resource-dev/$rel"
  test -f "$claude_file"
  test -f "$codex_file"

  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$codex_file" > "$expected"
  diff -u "$expected" "$claude_file"
done

for unexpected in \
  "$repo_root/.claude/skills/terraformer-resource-dev/agents/openai.yaml" \
  "$repo_root/.agents/skills/terraformer-resource-dev/agents/openai.yaml"; do
  if [[ -e "$unexpected" ]]; then
    echo "terraformer_resource_dev_skill_rules_test: unexpected optional metadata $unexpected" >&2
    exit 1
  fi
done

for skill in \
  "$repo_root/.claude/skills/terraformer-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/SKILL.md"; do
  for term in \
    "description: Use when developing, diagnosing, or fixing an Alibaba Cloud resource in Terraformer" \
    "bootstrap/workspace.sh dir terraformer" \
    "references/alicloud-resource-development.md" \
    "terraform-rd" \
    "terraform-qa" \
    "InitResources" \
    "Do not produce or infer resource relationships"; do
    if ! grep -Fq -- "$term" "$skill"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
done

for reference in \
  "$repo_root/.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/references/alicloud-resource-development.md"; do
  for term in \
    "A. Direct full List" \
    "B. One List returns every composite-ID segment" \
    "C. Parent-child traversal" \
    "D. Complete enumeration is unavailable" \
    'd.SetId(...)' \
    'ParseResourceId(...)' \
    "A multipart Import ID does not by itself require parent traversal" \
    "A Data Source may require the parent ID" \
    "Reset pagination for every parent" \
    "Do not produce or infer connections" \
    "go test ./providers/alicloud" \
    "go test ./..." \
    "/tmp/terraformer"; do
    if ! grep -Fq -- "$term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $reference" >&2
      exit 1
    fi
  done
done

echo "terraformer_resource_dev_skill_rules_test: PASS"
```

- [ ] **Step 2: Run the contract test and verify RED**

Run:

```bash
bash test/terraformer_resource_dev_skill_rules_test.sh
```

Expected: non-zero because `.claude/skills/terraformer-resource-dev/SKILL.md` does not exist. Confirm the failure is caused by the missing Skill, not shell syntax.

- [ ] **Step 3: Run the official initializer in scratch space**

Run:

```bash
scratch_root="$(mktemp -d /tmp/terraformer-skill-scaffold.XXXXXX)"
python3 /Users/shanye/.codex/skills/.system/skill-creator/scripts/init_skill.py \
  terraformer-resource-dev \
  --path "$scratch_root" \
  --resources references \
  --interface display_name="Terraformer Resource Development" \
  --interface short_description="Develop and repair Alibaba Cloud Terraformer resources" \
  --interface 'default_prompt=Use $terraformer-resource-dev to implement or diagnose an Alibaba Cloud Terraformer resource.'
```

Expected: initializer succeeds in `/tmp`. Do not copy its optional `agents/openai.yaml` into Jarvis; the user explicitly selected the minimal repository layout.

- [ ] **Step 4: Write the canonical Skill entrypoint**

Create `.claude/skills/terraformer-resource-dev/SKILL.md` with this content:

```markdown
---
name: terraformer-resource-dev
description: Use when developing, diagnosing, or fixing an Alibaba Cloud resource in Terraformer, including unsupported resources, incomplete discovery, incorrect or multipart Import IDs, parent-scoped listing, pagination defects, endpoint failures, or invalid generated Terraform state or HCL.
---

# Terraformer Resource Development

## Core model

Treat a Terraformer resource as a discovery adapter, not a second Terraform Provider resource implementation:

```text
InitResources
  -> enumerate remote objects and emit Provider-compatible Import IDs
  -> ProviderWrapper.Refresh uses the installed Alicloud Provider Import/Read
  -> ConvertTFstate produces Terraform state and HCL
```

Diagnose failures at the correct layer: discovery, Import ID, Provider Read, or state/HCL conversion. Do not copy Provider CRUD logic into `InitResources`.

## Start every task

1. Resolve the repository with `bash bootstrap/workspace.sh dir terraformer`; resolve Provider evidence with `bash bootstrap/workspace.sh dir terraform_provider`.
2. Preserve dirty files in the Terraformer checkout and create an isolated worktree before modifying tracked files.
3. Use the existing Aone claim/bookend flow. Assign implementation to `terraform-rd` and acceptance verification to `terraform-qa`.
4. Read [references/alicloud-resource-development.md](references/alicloud-resource-development.md) before choosing an API or writing code.
5. Classify the request as a new resource or a repair. For a repair, change only files required by the demonstrated root cause and add a regression test.

## Evidence order

Use this order and record the decisive evidence:

1. Terraform Provider Resource source.
2. Provider Import documentation and Import acceptance tests.
3. Provider Data Source source for List/filter/pagination behavior only.
4. Provider service/client implementation.
5. Terraformer resources with the same discovery pattern.
6. OpenAPI metadata or official API documentation.
7. Read-only live API/export results when credentials and existing resources are available.

The Provider's `d.SetId(...)`, `ParseResourceId(...)`, Import docs, and Import tests define the Import ID. Do not infer it from names or Data Source arguments.

## Select one discovery pattern

Choose exactly one primary `InitResources` pattern:

- **A. Direct full List:** the List API enumerates resources without parent scope.
- **B. One List returns every composite-ID segment:** one response includes all parent and child ID pieces.
- **C. Parent-child traversal:** the child List API requires a parent ID, so enumerate parents and then children; reset pagination for each parent.
- **D. Complete enumeration is unavailable:** use an existing explicit scope/filter input or report the unsupported boundary.

A multipart Import ID does not imply pattern C. A Data Source may require a parent ID because its caller supplies a scope; Terraformer must discover that parent only when the child List API requires it.

## Change only applicable files

- Always add or repair `providers/alicloud/resource_alicloud_<name>.go`.
- Update `providers/alicloud/alicloud_provider.go` only when registration in `SupportedResourceByProduct` or the global-resource list is required.
- Add client/service or endpoint support only when the current product client cannot issue the required API call.
- Add resource-level tests that lock Import ID construction, pagination, empty results, and error propagation.
- Do not modify Terraform Provider code as part of a Terraformer task; split Provider defects into `provider-resource-dev` work.
- Do not produce or infer resource relationships. Read the unified relationship artifact and consume only an explicit matching declaration.

## Validation gates

Run target checks before broad checks:

1. Verify `gofmt` reports no target files.
2. Run the resource regression test and `go test ./providers/alicloud`.
3. Build the binary to `/tmp/terraformer` so the repository stays clean.
4. Confirm the resource is visible through the Terraformer CLI registration path.
5. Run or record `go test ./...`; compare failures with the baseline instead of hiding existing unrelated failures.
6. When an account and an existing resource are available, perform a read-only export, inspect state/HCL, run `terraform validate`, and run `terraform plan -refresh-only`.

If live validation is unavailable, report "static validation only" and list the missing acceptance evidence. Never create cloud resources merely to make a Terraformer discovery check possible unless the user explicitly authorizes it.

## Delivery

Keep the worktree after opening a CR/MR, link it to Aone immediately, and do not merge or release. Report the selected discovery pattern, Import ID evidence, files changed, tests run, existing baseline failures, and any live-validation gap.
```

- [ ] **Step 5: Write the technical reference**

Create `.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md` with the following required structure and content. Keep the headings and exact contract sentences because the rule test relies on them.

```markdown
# Alicloud Terraformer resource development

## Contents

1. Runtime architecture
2. Source-of-truth checklist
3. InitResources discovery patterns
4. Multipart Import IDs
5. Pagination and errors
6. File selection
7. Tests and validation
8. Common mistakes

## 1. Runtime architecture

`Generator.InitResources()` loads the Alicloud client, calls one or more read-only APIs, converts each discovered object into `terraformutils.Resource`, and appends it to `g.Resources`. Terraformer then delegates Import/Read to the installed Provider through `ProviderWrapper.Refresh`; `ConvertTFstate` converts the returned Provider state to state and HCL.

Keep `InitResources` limited to discovery and Provider-compatible IDs. Do not reproduce Create, Update, Delete, schema flattening, or drift logic from the Provider.

## 2. Source-of-truth checklist

Read sources in this order:

1. Provider Resource: find `d.SetId(...)`, every `ParseResourceId(...)`, the Importer, and Read lookup parameters.
2. Import docs/tests: confirm segment order, delimiter, and import round trip.
3. Provider Data Source: reuse only the List API choice, filters, response path, and pagination semantics.
4. Provider service/client: confirm product endpoint, API version, RPC/ROA style, retryable errors, and response normalization.
5. Terraformer same-pattern resources: reuse repository conventions, not identity assumptions.
6. OpenAPI: verify request/response fields when Provider code is indirect or generated.
7. Live read-only call: validate only when credentials and an existing resource are available.

When sources conflict, Provider Import/Read behavior wins for the ID contract. Record the conflict rather than guessing.

## 3. InitResources discovery patterns

### A. Direct full List

Use when one List API enumerates all resources without a parent identifier. Paginate until the API's explicit completion signal, or until a short page when no stronger signal exists. Build each final ID from fields returned by that same item.

### B. One List returns every composite-ID segment

Use when one response item contains every segment required by the Provider Import ID. Preserve the Provider-defined order and delimiter. Do not add a parent List merely because the ID has multiple segments.

### C. Parent-child traversal

Use only when the child List API requires parent scope and Terraformer must enumerate the whole account/region scope:

1. List all parents with complete pagination.
2. For each parent, create a fresh child request.
3. Reset pagination for every parent.
4. List every child page.
5. Join parent and child segments once, at the leaf, using the Provider contract.
6. Return errors with parent ID and page/token context; never silently skip one parent.

A Data Source may require the parent ID because the Terraform caller supplies a query scope. Terraformer cannot impose that Data Source input on a full export; it discovers parents only for this pattern.

The following is pseudocode for the loop shape, not a copy-ready SDK call:

```go
for _, parentID := range parentIDs {
    pageNumber := 1
    nextToken := ""

    for {
        children, page, err := listChildren(parentID, pageNumber, nextToken, pageSize)
        if err != nil {
            return nil, fmt.Errorf("list children for parent %s: %w", parentID, err)
        }
        for _, child := range children {
            importID, err := buildProviderImportID(parentID, child.ID)
            if err != nil {
                return nil, err
            }
            ids = append(ids, importID)
        }
        if page.NextToken == "" && len(children) < pageSize {
            break
        }
        pageNumber++
        nextToken = page.NextToken
    }
}
```

In real code, choose either the API's token contract or page-number contract; do not combine both unless that API actually returns both.

### D. Complete enumeration is unavailable

Use when the service offers only exact lookup, the parent cannot be enumerated, or permissions make account-wide discovery impossible. Reuse an existing Terraformer scope/filter mechanism when it can express the missing input. Otherwise stop and report the limitation; do not claim complete support or guess IDs.

## 4. Multipart Import IDs

The only valid evidence for segment count, order, and delimiter is the Provider Resource's `d.SetId(...)`, its `ParseResourceId(...)` calls, Import docs, and Import tests.

A multipart Import ID does not by itself require parent traversal. All segments may already be present in one List response (pattern B), or earlier segments may require parent discovery (pattern C).

Implementation rules:

- Carry parent, child, attachment, or account segments as separate variables while traversing.
- Validate every required segment before joining.
- Join exactly once when creating the leaf `terraformutils.Resource` ID.
- Do not trim, encode, reorder, or change delimiters without Provider evidence.
- Test the normal ID, missing segment, order, delimiter, and special-character boundary.

## 5. Pagination and errors

- Prefer `NextToken`, `TotalCount`, `IsTruncated`, or an equivalent explicit signal.
- When using returned item count, compare it with the exact page-size value sent in the request.
- Reset pagination for every parent; initialize page number/token inside the parent loop.
- Cover empty first page, short last page, exactly full last page, and multiple pages.
- Include action, resource type, parent ID, page number, or token in wrapped errors.
- Treat permission, endpoint, decode, and single-parent failures as errors, not empty results.
- Follow the repository's retry helpers and product client conventions; do not invent a second retry framework.

## 6. File selection

| File | Change when |
|---|---|
| `providers/alicloud/resource_alicloud_<name>.go` | Always for a new resource; normally for a repair |
| `providers/alicloud/alicloud_provider.go` | Registration or global-resource classification is missing |
| Product client/service files | No existing client can issue the API call |
| Endpoint configuration | The current endpoint resolution is proven insufficient |
| Resource `_test.go` | Lock ID, pagination, empty-result, and error behavior |
| Unified relationship consumer | The shared artifact explicitly declares this resource |

Do not produce or infer connections from Provider schema, Data Source arguments, or API field names. The unified producer owns relationship semantics.

Do not modify `cmd`, module entrypoints, README, Provider source, or unrelated shared code unless repository evidence proves the resource cannot work without that change.

## 7. Tests and validation

Use TDD for repairs: demonstrate the current failure, add the smallest regression, then implement the fix.

Static gates:

```bash
RESOURCE_FILE=providers/alicloud/resource_alicloud_example.go
gofmt -l "$RESOURCE_FILE"
go test ./providers/alicloud
go build -o /tmp/terraformer .
```

Confirm registration through the CLI's supported-resource listing or equivalent code path. Run or record `go test ./...`; the current repository has existing unrelated failures, so compare the broad result with the captured baseline while requiring all target-package checks to pass.

When live read-only validation is possible:

1. Export only the target product/resource.
2. Compare discovered count and IDs with the API response.
3. Inspect generated state and HCL.
4. Run `terraform init` and `terraform validate` in the generated directory.
5. Run `terraform plan -refresh-only` and investigate any read/import drift.

When credentials or an existing resource are unavailable, report static validation only and list the unverified live steps.

## 8. Common mistakes

| Mistake | Correct action |
|---|---|
| Treating every multipart ID as parent-child discovery | Select pattern B when one response already contains all segments |
| Copying a Data Source's required parent argument | Enumerate parents only when the child List API requires it |
| Initializing page number outside the parent loop | Reset pagination for every parent |
| Requesting one page size and testing termination with another | Use the same page-size variable for request and termination |
| Guessing an Import ID from API primary keys | Read Provider `d.SetId(...)`, `ParseResourceId(...)`, and Import evidence |
| Editing connection maps by inspection | Read only an explicit declaration from the unified relationship artifact |
| Treating `go test ./...` baseline failures as success or as a new regression | Report the baseline delta and require target tests to pass |
| Creating cloud resources for convenience | Use existing resources or obtain explicit authorization |
```

- [ ] **Step 6: Remove optional metadata from the approved design**

Update `docs/superpowers/specs/2026-07-16-terraformer-resource-dev-design.md` so section 3 shows only `SKILL.md` and `references/alicloud-resource-development.md` in both mirrors. Remove the `agents/openai.yaml` responsibility and rule-test requirement. In the completion criteria, replace `agent metadata` with `technical reference`. Add this sentence after the layout:

```markdown
初版不包含可选的 `agents/openai.yaml`；Jarvis 依靠 `SKILL.md` frontmatter 发现和触发该 Skill，后续只有在需要 Codex UI 展示元数据时才单独增加。
```

- [ ] **Step 7: Generate the Codex mirror**

Run:

```bash
bash bootstrap/mirror.sh to-codex \
  .claude/skills/terraformer-resource-dev/SKILL.md \
  .claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md
```

Expected: both `.agents/skills/terraformer-resource-dev` files are created with the repository's Claude-to-Codex token transformation.

- [ ] **Step 8: Run the contract test and verify GREEN**

Run:

```bash
bash test/terraformer_resource_dev_skill_rules_test.sh
```

Expected: `terraformer_resource_dev_skill_rules_test: PASS`.

- [ ] **Step 9: Validate the Skill and mirror**

Run:

```bash
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  .claude/skills/terraformer-resource-dev
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py \
  .agents/skills/terraformer-resource-dev
bash bootstrap/mirror.sh check
```

Expected: both Skill validators succeed and mirror check exits zero.

- [ ] **Step 10: Run Jarvis regression tests**

Run these separately:

```bash
bash test/aone_comment_format_test.sh
bash test/wrap_done_test.sh
bash test/provider_resource_dev_skill_sync_test.sh
```

Expected: every command exits zero.

- [ ] **Step 11: Commit the implementation**

```bash
git add \
  docs/superpowers/specs/2026-07-16-terraformer-resource-dev-design.md \
  docs/superpowers/plans/2026-07-16-terraformer-resource-dev.md \
  test/terraformer_resource_dev_skill_rules_test.sh \
  .claude/skills/terraformer-resource-dev \
  .agents/skills/terraformer-resource-dev
git commit -m "feat: add terraformer resource development skill"
```

### Task 2: Forward-test and deliver the Skill

**Files:**
- Modify only if forward tests expose a concrete gap: `.claude/skills/terraformer-resource-dev/SKILL.md`
- Modify only if forward tests expose a concrete gap: `.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md`
- Regenerate after any canonical edit: `.agents/skills/terraformer-resource-dev/**`

**Interfaces:**
- Consumes: Task 1 Skill, `/Users/shanye/programs/terraformer` as read-only evidence, three evaluation prompts below.
- Produces: forward-test evidence, any minimal refinement, and an updated CR 28627638.

- [ ] **Step 1: Run three fresh-context application scenarios with the Skill**

Use fresh `terraform-rd` evaluation agents. Each agent must read the Skill first and must not modify Terraformer:

```text
Scenario A: A new resource has a global List API whose items contain the complete Provider Import ID. Explain the InitResources approach and files to inspect.

Scenario B: A child List API requires workspace_id; Provider Import ID is workspace_id:member_id. Explain discovery, pagination, and ID construction.

Scenario C: Provider schema implies a parent relation, but the unified relationship artifact has no declaration. Explain what the Terraformer change should do.
```

Expected:

- Scenario A selects pattern A or B and does not add parent traversal.
- Scenario B selects pattern C, discovers parents, resets child pagination per parent, and uses Provider evidence for `workspace_id:member_id`.
- Scenario C refuses to infer a relationship and records/consumes only unified producer output.

- [ ] **Step 2: Refine only demonstrated gaps**

If an evaluator violates an expected behavior, add the smallest explicit sentence or example that closes that gap, regenerate the mirror, and rerun the failed scenario. Do not add hypothetical material.

- [ ] **Step 3: Run final verification**

```bash
bash test/terraformer_resource_dev_skill_rules_test.sh
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py .claude/skills/terraformer-resource-dev
uvx --with pyyaml python /Users/shanye/.codex/skills/.system/skill-creator/scripts/quick_validate.py .agents/skills/terraformer-resource-dev
bash bootstrap/mirror.sh check
bash test/aone_comment_format_test.sh
bash test/wrap_done_test.sh
git diff --check origin/master...HEAD
```

Expected: every command exits zero.

- [ ] **Step 4: Commit refinements when present**

```bash
git add .claude/skills/terraformer-resource-dev .agents/skills/terraformer-resource-dev
git commit -m "docs: tighten terraformer resource development guidance"
```

Skip this commit when forward tests require no file changes.

- [ ] **Step 5: Push and update the existing CR**

```bash
git push origin worktree-terraformer-resource-dev-design
```

Verify CR 28627638 contains the design, plan, Skill, reference, mirror, and rule test; then sync the CR link and verification summary to Aone 84375416 and release the claim without changing status.
