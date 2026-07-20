---
name: terraform-provider-live-repro
description: Reproduce Terraform Provider issues against real Alibaba Cloud resources with a minimal local HCL case, Provider debug traces, OpenAPI RequestIds, refresh/drift evidence, controlled scene preservation or cleanup, and an optional AutomationAgent HTML report. Use when a user asks to actually reproduce a provider bug, create resources from supplied Terraform, investigate create/read round-trip drift or unexpected replacement, preserve a live cloud scene for a product team, or produce a complete API timeline; prefer invoke-terraform-acc-test-remote for ordinary acceptance-test execution unless the reported behavior requires the caller's exact HCL and live API responses.
---

# Terraform Provider Live Repro

Turn a reported Provider symptom into a repeatable, auditable live reproduction. Keep the reported input exact, separate cloud API behavior from Provider state behavior, and leave either a verified live scene or a verified clean account.

## Routing boundary

- Use this skill for direct Terraform templates, create/read round trips, refresh drift, unexpected ForceNew replacement, exact reporter inputs, RequestId timelines, or live-scene handoff.
- Use `invoke-terraform-acc-test-remote` for normal `TestAcc*` execution, PR regression coverage, or broad remote validation.
- If existing AccTests pass but do not contain the reporter's exact trigger fields, treat them as regression coverage rather than reproduction proof.
- If an Aone ID is present, load `aone-triage` first and verify the item title/description. Uploading a report does not authorize comments, status changes, or closure.

## Hard safety rules

1. Create real resources only when the user explicitly asks to create/apply/reproduce, or when an already-approved workflow clearly includes live creation.
2. Before apply, record the intended disposition: `preserve` or `destroy-after-evidence`. Preserve only when requested or needed for downstream investigation.
3. Never apply a drift/replacement plan. A plan showing replacement is evidence, not an execution artifact.
4. Never print, persist, commit, or upload AK/SK, bearer tokens, signatures, passwords, or provider configuration values.
5. Treat `.tfplan` files and raw `TF_LOG` traces as sensitive. Extract an allowlisted timeline, then delete the originals.
6. Never destroy a preserved scene. For disposable scenes, destroy only the exact resources in that reproduction state and verify removal.
7. Do not silently replace inaccessible user-supplied VPC/VSwitch IDs. Either stop or use an isolated equivalent network and state that deviation prominently in the report.

## Workflow

### 1. Establish the case

Record:

- reporter HCL and exact provider version;
- expected behavior and observed symptom;
- region, zone, resource type, and lifecycle phase;
- Aone ID when available;
- requested working directory and scene disposition.

Default to an isolated directory under `~/Workspace/troubleshoot/<slug>-<aone-id-or-date>`. Do not overwrite another reproduction or reuse unrelated Terraform state.

### 2. Preflight identity and dependencies

Verify without exposing credentials:

```bash
terraform version
aliyun sts GetCallerIdentity
```

Verify the requested region/zone and referenced VPC/VSwitch with read-only APIs before planning. Report only identity labels/IDs that are necessary for internal diagnosis; never show credential values.

Pin the reported provider version exactly. Pass credentials through the environment or an existing CLI profile. Do not write plaintext secrets to `.tfvars`.

### 3. Build the minimal reproduction

Preserve all trigger fields and their exact values/relationships. Remove only unrelated resources. If the original network is inaccessible and an isolated equivalent is approved, keep the same region, zone, resource shape, and dependency graph.

Run:

```bash
terraform fmt -recursive
terraform init -backend=false -input=false -no-color
terraform validate -no-color
```

### 4. Gate the first plan

Generate a saved create plan and inspect its JSON:

```bash
terraform plan -input=false -no-color -out=create.tfplan
terraform show -json create.tfplan | jq '{changes: [.resource_changes[] | {address, actions: .change.actions}]}'
```

Proceed only when the plan contains exactly the expected creates and no update, delete, or replacement. If this gate fails, stop and diagnose configuration, identity, region, or stale-state problems.

### 5. Apply with trace evidence

Record local and UTC start/end times. Enable Provider trace logging only for the reviewed apply:

```bash
DEBUG=terraform \
TF_LOG=DEBUG \
TF_LOG_PROVIDER=DEBUG \
TF_LOG_PATH="$evidence_dir/terraform-apply.raw.log" \
terraform apply -input=false -auto-approve -no-color create.tfplan
```

Record Terraform resource IDs immediately. Do not publish the raw log.

### 6. Capture direct API truth

Call the product's read APIs with the created IDs. Capture for every call:

- timestamp and API name;
- RequestId;
- instance/resource ID;
- relevant status and fields;
- create-request fields needed to prove the trigger was sent.

For state round-trip issues, compare at least:

1. create request parameters and create response;
2. the primary describe/read response;
3. related-resource APIs that may retain the missing relationship;
4. Terraform state after Provider Read.

For NAS VPC/VSwitch cases, compare `CreateFileSystem`, `DescribeFileSystems`, and `DescribeMountTargets`. Explicitly distinguish `VpcId=""`, `QuorumVswId=null`, and a missing `QuorumVswId` field.

### 7. Extract a sanitized apply timeline

Use the bundled parser; it emits only allowlisted request/response fields:

```bash
python3 {SKILL_DIR}/scripts/extract-api-timeline.py \
  "$evidence_dir/terraform-apply.raw.log" \
  --format markdown > "$evidence_dir/apply-api-timeline.md"
```

Read [references/evidence-contract.md](references/evidence-contract.md) before sharing evidence or writing the final report.

### 8. Prove refresh drift

Run a normal refresh/plan with a separate raw trace:

```bash
set +e
DEBUG=terraform \
TF_LOG=DEBUG \
TF_LOG_PROVIDER=DEBUG \
TF_LOG_PATH="$evidence_dir/terraform-drift-plan.raw.log" \
terraform plan -input=false -no-color -detailed-exitcode -out=drift.tfplan
plan_rc=$?
set -e
```

Interpret `0` as empty plan, `2` as changes, and other codes as failure. Extract replacement paths:

```bash
terraform show -json drift.tfplan | jq '{changes: [
  .resource_changes[]
  | select(.change.actions == ["delete", "create"])
  | {address, replace_paths: .change.replace_paths}
]}'
```

Extract the refresh API timeline with the same parser. Never run `terraform apply drift.tfplan`.

### 9. Preserve or clean the scene

For `preserve`:

- keep cloud resources and Terraform state;
- verify the state resource count and live resource statuses;
- list IDs, region/zone, mount endpoints, and the explicit preservation warning;
- delete sensitive plan files and raw traces after sanitized evidence is complete;
- provide read-only verification commands and a clear `do not destroy` notice.

For `destroy-after-evidence`:

- review that destroy contains exactly the reproduction resources;
- run `terraform destroy` using the same credential source;
- verify state count is zero and cloud lookups return NotFound/zero count;
- state what was removed and that the resources are no longer recoverable.

### 10. Produce and publish the report

Copy [assets/report-template.md](assets/report-template.md) and fill every applicable section. Keep conclusions layered:

1. cloud API behavior;
2. Provider schema/Read behavior;
3. Terraform plan consequence;
4. safe fix direction and product-team questions.

Render a self-contained HTML file without base64 images:

```bash
python3 {SKILL_DIR}/scripts/render-report-html.py report.md report.html
```

When the user requests an online report, invoke `html-report-preview` and use the repository helper with a verified Aone ID:

```bash
bash bootstrap/html-report-preview.sh upload <aone-id> report.html \
  --base-url https://agent.aliyun-inc.com --format jsonl
```

Do not use browser cookies or a personal BUC session if the server token is missing. Verify the returned `/reports/aone/.../view` URL without credentials: HTTP 200, `text/html`, report title, one instance ID, one API name, and one key RequestId.

### 11. Remove sensitive artifacts

After the sanitized report is complete, resolve and delete only the exact generated files:

- `create.tfplan`;
- `drift.tfplan`;
- raw apply/refresh debug logs.

Keep the `.tf` files, lockfile, sanitized evidence, report, and state when preserving the scene. Scan the shareable artifacts for actual credential values before delivery.

## Completion gates

Do not claim success until all applicable gates pass:

- initial plan reviewed and apply completed;
- created IDs and API RequestIds captured;
- direct API response compared with Terraform state;
- drift plan exit code and replacement paths recorded;
- scene preservation or cleanup verified;
- raw sensitive artifacts removed;
- report contains deviations and exact timestamps;
- online report, when requested, returns 200 and contains distinctive evidence markers.
