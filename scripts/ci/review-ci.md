# Provider-owned ACube review CI

The scripts and self-contained templates in this repository own CI behavior;
ACube continues to provide the review service. Generate templates with
`python3 scripts/ci/build_review_ci_templates.py` and run regressions with
`python3 -m unittest discover -s scripts/ci -p test_review_ci.py`.

## Behavior and trust boundary

- Formal review preserves the deployed bridge from ACube commit
  `14da0fb400ff4b25b3656d0269dfd5e1ce2706fc` without source changes: exact MR/head
  binding, content reuse, non-service-change skips, P0 blocking, P1/P2 nonblocking,
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
