# Provider-owned ACube review CI

The scripts and self-contained templates in this repository own CI behavior;
ACube continues to provide the review service. Generate templates with
`python3 scripts/ci/build_review_ci_templates.py` and run regressions with
`python3 -m unittest discover -s scripts/ci -p test_review_ci.py`.

## Behavior and trust boundary

- Both checks embed the same bridge, `review_ci.py`, and differ only in
  `ACUBE_REVIEW_KIND` (`formal` or `ai`) and their names. Each keeps exact MR/head
  binding, content reuse, ACube skips, P0/P1 blocking, P2 nonblocking, report
  links, a five-minute ACube unavailability window and a seven-hour polling
  deadline (eight-hour job timeout). Verified release return MRs skip both checks
  by policy.
- Formal review skips MRs without service implementation changes. Once
  [switched](#switching-pipeline-344310), AI review runs on the same MR events as
  formal review, and ACube decides whether Sonus actually runs; see
  [AI review on demand](#ai-review-on-demand).
- People dispose of findings on Code, by replying under the finding's MR comment.
  The check summary leads with those replies; the ACube report link stays as a
  secondary entry.
- Templates embed their scripts and perform no MR checkout. Publish only an
  approved immutable Provider commit, never load scripts from the reviewed branch.
- `ci/templates` deliberately does not register additional `.aoneci` pipelines.
  Existing checks continue during review; adding this code alone does not switch them.

## Pipelines

Both checks are existing Provider pipelines bound to Provider-owned code
templates. Keep their IDs, names and merge-gate associations:

| Pipeline | Name | Template | MR events | Source |
| --- | --- | --- | --- | --- |
| 344309 | ACube 形式化评审 | 5246 | opened, reopened | ci/templates/acube-formal-review.yaml |
| 344310 | ACube AI评审 | 5245 | accepted until the [switch](#switching-pipeline-344310), then opened, reopened | ci/templates/acube-ai-review.yaml |

Merging a change here activates nothing: a pipeline runs the YAML of the template
version it is pinned to. To activate a change, register the merged master commit
(full SHA) as a version of the template and pin the pipeline to it, for example
`a1 ci pipeline update 344309 --repo opensource-tools/terraform-provider-alicloud --template-version <SHA>`.
Read the pipeline back (`a1 ci pipeline get <ID> --repo opensource-tools/terraform-provider-alicloud -f json`):
`templateVersion` and `codeUrl` must show the new SHA. The value must really
change for the YAML to be regenerated, and a rerun keeps the YAML of the run it
repeats.

Preserve target branches `**`, concurrency group `${{git.merge_request.id}}`,
parallelism 1, cancel-in-progress true and the failure notification. Never
create a second automatically triggered pipeline for either check, and keep a
configuration snapshot before each change so a rollback restores it on the same
pipeline ID.

## P0/P1 gate update

The Provider bridge evaluates raw findings for both new and reused completed
reports. Any P0 or P1 blocks CI even if the ACube service returns success under
its earlier P0-only policy. Service failures remain failures; malformed reports
fail closed. Annotations do not waive findings. The summary shows the effective
CI conclusion, original ACube conclusion and repair guidance.

The gate was activated by pinning formal pipeline 344309 (Provider template
5246) to the merged SHA. The AI check applies the same gate once it runs Sonus;
see [AI review on demand](#ai-review-on-demand).

## Dispositions: rejected or ignored findings pass

ACube returns `dispositions` keyed by finding ID with the formal report. Each
P0/P1 finding passes only when its disposition `state` is `rejected` or
`ignored`; `verifying` keeps the check pending; anything else, including a
missing or malformed disposition, blocks. Fields inside a finding never waive it.
A pass also requires ACube's own `conclusion` to be `success`; a disagreement
fails closed. P2 still never blocks.

ACube decides the disposition. A person's decision always wins, whether given
by a reply under the finding's MR comment (see
[Disposing of findings on Code](#disposing-of-findings-on-code)) or in the ACube
workbench. Without one, the Agent verification verdict decides:
not reproduced rejects the finding, reproduced adopts it and posts it to the MR
as an ordinary review comment (so PR followers deliver it to RD), and a blocked
or failed verification waits for a person. ACube starts the verification for
undecided P0/P1 findings by itself while CI polls, so the check stays pending,
normally for under an hour, instead of failing early.

When a person changes a decision, or starts or retries an Agent verification, on
Code or in the ACube workbench, ACube reruns this check by itself shortly after,
using its dedicated review identity (Terraform 质量数字人) against AoneCI. It
reruns only the latest check of the current MR head, and only if that check had
already finished before the change. A check that is still running or queued is
left alone: it re-reads the dispositions while it waits. Completed reports are
reused for the same content, so the rerun reapplies the current dispositions
without a new review. While the reused report is still being verified, the same
run asks ACube again every minute. If the automatic rerun fails, ACube only logs
it; rerun the check by hand.

A rerun keeps the YAML of the run it repeats. Runs created before a pipeline is
pinned to a new version still execute the old bridge when rerun; start a fresh
run for those MRs (a new commit, or closing and reopening the MR).

The summary shows each finding's disposition, who made it, the reason and a link
to the Jarvis verification task. The log ends with stable lines for digital
humans that cannot open the BUC-only ACube page. A passed check lists every P0/P1:

```text
ACUBE_FINDING id=F-01 severity=P1 state=rejected actor=AGENT verdict=NOT_REPRODUCED evidence=<Jarvis task URL>
ACUBE_REVIEW_RESULT kind=PASSED
```

Digital humans read a failed check with the Jarvis SCM tool `read_ci_failure`,
which keeps only error-formatted lines (the last 8 per step, 360 characters
each). A failed check therefore prints only the findings that still block, at
most 5 lines with an overflow count, then the result line last, all prefixed with
AoneCI's failure marker:

```text
HIGHLIGHT__[ERROR] ACUBE_FINDING id=F-02 severity=P1 state=adopted actor=AGENT verdict=REPRODUCED comment=<MR comment ID>
HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=FINDINGS_BLOCKED|SERVICE_ERROR|TIMEOUT|MALFORMED_REPORT [reason=...]
```

Deploy the ACube change first. Until it is live, ACube returns no dispositions
and every P0/P1 keeps blocking, exactly as before.

## Waiting out ACube unavailability

ACube preproduction is redeployed often, and each redeploy makes it unavailable
for a few minutes. Every call the bridge makes to ACube therefore retries
temporary failures until ACube has been unavailable for five minutes in a row
(`UNAVAILABLE_WINDOW_SECONDS`). Temporary failures are HTTP 429, 500, 502, 503
and 504, refused or dropped connections, timeouts, invalid HTTP status lines,
and truncated response bodies. Delays grow from 5 to 20 seconds and each retry
is logged. A successful answer ends the outage: polling continues as usual and
a later outage gets a full window again. Waiting never extends past the
seven-hour review deadline.

Other HTTP statuses, redirects, oversized or unreadable answers and receipt or
report identity mismatches still fail at once. When the window runs out, the
check fails with the usual result line, for example:

```text
HIGHLIGHT__[ERROR] ACUBE_REVIEW_RESULT kind=SERVICE_ERROR reason=ACube_HTTP_502__check_task_in_ACube
```

The IntegrationTest and DocsExampleTest jobs apply the same five-minute window
to their safe ACube calls; see [README-scopes.md](README-scopes.md). They run the
script from the checked-out commit, so they pick it up once merged. The formal
review runs the YAML embedded in template 5246: its runs keep the old bridge
until pipeline 344309 is pinned to a Provider commit that contains this change.

## AI review on demand

After the [switch](#switching-pipeline-344310), pipeline 344310 runs the same
bridge as formal review on the same MR events (opened, reopened). Sonus is too
costly to run on every MR, so ACube decides for each head whether it actually
runs: hard rules first (for example, no service implementation change), then a
digital human's judgment and sampling. Anyone who can comment on the MR can
still ask for Sonus by commenting `/sonus`.

- When ACube decides not to run Sonus, the check passes as `SKIPPED`. The
  summary shows ACube's reason, for example `数字人判断无需 Sonus：……`, and the
  log prints it on one line before `ACUBE_REVIEW_RESULT kind=SKIPPED`. A skip is
  not an AI verdict.
- Otherwise the check stays pending while ACube waits for the formal review, for
  the judgment, for Sonus and for Agent verification. Sonus findings then follow
  formal review exactly: an Agent verifies each P0/P1, the verdict or a person
  disposes of it, and the same P0/P1 gate decides the check.
- All of this counts against the seven-hour deadline. A check that runs out
  fails as `TIMEOUT` and can be rerun.

### What the bridge expects from ACube

Paths and request fields do not change: `POST /api/v1/provider_runs/ci-reviews/ai`
with `mr`, `headSha`, `runId` and `contentReuseVersion`, then
`GET /api/v1/provider_runs/ci-reviews/ai/{taskId}?mr=<MR>&headSha=<head>`.

- The POST answers with a skip, a reused completed report for the same content,
  or a task receipt with `taskId`. A skip carries `status: skipped`,
  `conclusion: success`, the CI `headSha` and `mr`, and its reason in `summary`
  (a `reason` value is shown when `summary` is empty).
- Every poll answer carries the task's `taskId` and the CI `headSha`.
  `conclusion: pending` keeps the check waiting, whatever the `status`; ACube may
  describe the current stage in `summary`, which the check summary shows.
  `status: skipped` with `conclusion: success` ends the check as a skip, and
  `status: completed` with `findings` and `dispositions` goes through the gate.
- A skip without `conclusion: success`, or `conclusion: success` on a report that
  is neither skipped nor completed, fails closed as `SERVICE_ERROR`.
- As for formal review, ACube opens MR comment threads for P0/P1 findings,
  applies the Code replies below to them, reruns the AI check after a decision,
  and runs Sonus when anyone who can comment on the MR comments `/sonus`.

### Disposing of findings on Code

By default nobody needs to act: an Agent verification that does not reproduce a
finding rejects it, and one that reproduces it adopts it and posts it to the MR.
Fix adopted findings with a new commit. When a person must decide, one reply
under the finding's MR comment is enough:

| Reply | Meaning | Who may send it |
| --- | --- | --- |
| `采纳` | The finding is real and keeps blocking until fixed. | Anyone who can comment on the MR (the author included) |
| `驳回 理由` | False positive; releases the gate. | Anyone who can comment on the MR (the author included); a reason is required |
| `忽略 理由` | Real, but not fixed in this MR; releases the gate. | Anyone who can comment on the MR (the author included); a reason is required |
| `重试` | Run the Agent verification again. | Anyone who can comment on the MR (the author included) |

There is no reviewer or maintainer list. ACube records each `驳回` or `忽略`
reason together with who replied, answers under the thread and reruns the
check; the summary shows the decision, its author and the reason. Both checks' summaries lead
with these replies and `/sonus`; the ACube report link is a secondary entry, and
the workbench still accepts the same decisions, but nobody needs to open it.

### Switching pipeline 344310

The pipeline owner performs the switch personally, and only after the matching
ACube version is live on pre-acube: the AI endpoint accepts opened MRs (until
then it answers HTTP 409 `AI review requires reviewer approval` for MRs that
change service code), returns skips with a reason, and applies Code replies.
Switching earlier fails every AI check on MRs that change service code.

1. Snapshot both pipelines with `a1 ci pipeline get <ID> --repo opensource-tools/terraform-provider-alicloud -f json`.
2. Register the merged master commit (full SHA) as a version of template 5245,
   pin 344310 to it, and change its MR events from `accepted` to `opened` and
   `reopened`:
   `a1 ci pipeline update 344310 --repo opensource-tools/terraform-provider-alicloud --template-version <SHA> --trigger-mr-type opened --trigger-mr-type reopened --trigger-mr-target-branch '**'`
3. Register the same SHA for template 5246 and pin 344309 to it, so the formal
   check also leads with Code replies:
   `a1 ci pipeline update 344309 --repo opensource-tools/terraform-provider-alicloud --template-version <SHA>`
4. Read both pipelines back: `templateVersion` and `codeUrl` show the SHA;
   344310's MR events are `opened` and `reopened` with target branches `**`;
   concurrency, parallelism, cancel-in-progress and the failure notification are
   unchanged.
5. On a fresh run (a new commit, not a rerun), check that the AI check either
   skips with ACube's reason or polls an ACube task, that the formal and AI checks
   of one MR run side by side without cancelling each other, and that each event
   dispatches one run per check.

Rollback pins 344310 back to `6e74c600ada744206eb2d090d77cfca1cad36f6f` with MR
event `accepted`; that version passes at once without contacting ACube. 344309
goes back to its snapshot version, `01ea4e9723da0a5b0c7cfca700f3130e2e5375ae` at
the time of writing.
