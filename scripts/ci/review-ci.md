# Provider-owned ACube review CI

The scripts and self-contained templates in this repository own CI behavior;
ACube continues to provide the review service. Generate templates with
`python3 scripts/ci/build_review_ci_templates.py` and run regressions with
`python3 -m unittest discover -s scripts/ci -p test_review_ci.py`.

## Behavior and trust boundary

- Formal review retains exact MR/head
  binding, content reuse, non-service-change skips, P0/P1 blocking, P2 nonblocking,
  report links and finding actions, HTTP retries and seven-hour polling deadline.
- AI review exits successfully with `status=skipped`, `reason=ai-review-disabled`.
  It does not request credentials or submit a remote task; success is not an AI verdict.
- Templates embed their scripts and perform no MR checkout. Publish only an
  approved immutable Provider commit, never load scripts from the reviewed branch.
- `ci/templates` deliberately does not register additional `.aoneci` pipelines.
  Existing checks continue during review; adding this code alone does not switch them.

## Cutover after this MR is merged

Keep the existing Provider pipeline IDs, names and merge-gate associations:

| Pipeline | Name | Existing template | MR events | Provider template |
| --- | --- | --- | --- | --- |
| 344309 | ACube 形式化评审 | 5224 | opened, reopened | ci/templates/acube-formal-review.yaml |
| 344310 | ACube AI评审 | 5225 | accepted | ci/templates/acube-ai-review.yaml |

1. Publish/register the two templates from this repository at the approved merged
   SHA. Change each existing pipeline's template binding to the corresponding
   Provider-owned template and immutable version; updating only the version of
   an ACube-owned template does NOT migrate its source repository.
2. Preserve target branches `**`, concurrency group `${{git.merge_request.id}}`,
   parallelism 1, cancel-in-progress true and existing failure notifications.
   Do not create a second automatically triggered pipeline for either check.
3. Read back each pipeline: its code URL must point to terraform-provider-alicloud,
   its version to the approved commit, and events/IDs must match the table.
4. Verify formal review on an exact-head MR, including remote report identity and
   conclusion. Verify AI reports disabled/success without creating an ACube task.
   Check that only one run per check is dispatched for each triggering event.

Before switching, retain pipeline configuration snapshots. If updating an existing
template binding is unsupported, stop the cutover and resolve platform support;
do not enable replacement triggers alongside the old ones. Existing runs are not
cancelled by this source migration. Rollback restores the saved bindings on the
same pipeline IDs, never adds duplicate checks.

The CLI currently exposes pipeline template-version updates but not template
publication/source registration. Cutover requires the supported platform template
management interface; source merge and live activation must be reported separately.

## P0/P1 gate update

The Provider bridge evaluates raw findings for both new and reused completed
reports. Any P0 or P1 blocks CI even if the ACube service returns success under
its earlier P0-only policy. Service failures remain failures; malformed reports
fail closed. Annotations do not waive findings. The summary shows the effective
CI conclusion, original ACube conclusion and repair guidance.

After merging this update, pin existing formal pipeline 344309 (Provider
template 5246) to the approved merged SHA. Preserve its triggers and strategy;
do not change AI pipeline 344310. Source merge alone does not activate the gate.
A rerun can retain the old YAML snapshot: verify a new run's template version.

## Dispositions: rejected or ignored findings pass

ACube returns `dispositions` keyed by finding ID with the formal report. Each
P0/P1 finding passes only when its disposition `state` is `rejected` or
`ignored`; `verifying` keeps the check pending; anything else, including a
missing or malformed disposition, blocks. Fields inside a finding never waive it.
A pass also requires ACube's own `conclusion` to be `success`; a disagreement
fails closed. P2 still never blocks.

ACube decides the disposition. A person's reject, ignore or adopt in the ACube
workbench always wins. Without one, the Agent verification verdict decides:
not reproduced rejects the finding, reproduced adopts it and posts it to the MR
as an ordinary review comment (so PR followers deliver it to RD), and a blocked
or failed verification waits for a person. ACube starts the verification for
undecided P0/P1 findings by itself while CI polls, so the check stays pending,
normally for under an hour, instead of failing early.

When a person changes a decision, or starts or retries an Agent verification in
the ACube workbench, ACube reruns this check by itself about 15 seconds later,
using its dedicated review identity (Terraform 质量数字人) against AoneCI. It
reruns only the latest check of the current MR head, and only if that check had
already finished before the change. A check that is still running or queued is
left alone: it re-reads the dispositions while it waits. Completed reports are
reused for the same content, so the rerun reapplies the current dispositions
without a new review. While the reused report is still being verified, the same
run asks ACube again every minute. If the automatic rerun fails, ACube only logs
it; rerun the check by hand.

A rerun keeps the YAML of the run it repeats. Runs created before pipeline 344309
is pinned to this version still execute the old bridge when rerun; start a fresh
run for those MRs (a new commit, or closing and reopening the MR).

The summary shows each finding's disposition, who made it, the reason and a link
to the Jarvis verification task. The log ends with stable lines for digital
humans that cannot open the BUC-only ACube page:

```text
ACUBE_REVIEW_RESULT kind=PASSED|FINDINGS_BLOCKED|SERVICE_ERROR|TIMEOUT|MALFORMED_REPORT|SKIPPED
ACUBE_FINDING id=F-01 severity=P1 state=rejected actor=AGENT verdict=NOT_REPRODUCED evidence=<Jarvis task URL>
ACUBE_FINDING id=F-02 severity=P1 state=adopted actor=AGENT verdict=REPRODUCED comment=<MR comment ID>
```

Deploy the ACube change first. Until it is live, ACube returns no dispositions
and every P0/P1 keeps blocking, exactly as before.
