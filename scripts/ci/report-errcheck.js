"use strict";

function markdown(value) {
  return String(value).replace(/[\r\n]/g, " ").replace(/[\\`*_\[\]<>|]/g, "\\$&");
}

async function reportErrcheck({ github, context, core }) {
  const event = context.payload.workflow_run;
  if (
    context.eventName !== "workflow_run" || context.payload.action !== "completed" ||
    !Number.isSafeInteger(event?.id) || event.id <= 0 ||
    !Number.isSafeInteger(event.run_attempt) || event.run_attempt <= 0
  ) {
    throw new Error("Invalid completed workflow event.");
  }

  const repository = context.repo;
  const fullName = `${repository.owner}/${repository.repo}`;
  const { data: run } = await github.rest.actions.getWorkflowRun({
    ...repository, run_id: event.id,
  });
  if (run.run_attempt !== event.run_attempt) {
    core.info("Ignoring an obsolete workflow attempt.");
    return;
  }
  if (
    run.id !== event.id || run.event !== "pull_request" ||
    run.name !== "Golangci-Lint" ||
    run.path !== ".github/workflows/acctest-terraform-lint.yml" ||
    run.repository?.full_name !== fullName || run.status !== "completed" ||
    !Number.isSafeInteger(run.check_suite_id) || run.check_suite_id <= 0 ||
    !/^[0-9a-f]{40}$/.test(run.head_sha ?? "")
  ) {
    throw new Error("Unexpected source workflow.");
  }

  const jobs = await github.paginate(github.rest.actions.listJobsForWorkflowRun, {
    ...repository, run_id: run.id, filter: "latest", per_page: 100,
  });
  const matches = jobs.filter((job) => job.name === "runner / errcheck");
  if (matches.length !== 1) {
    throw new Error("Expected exactly one errcheck job.");
  }
  const job = matches[0];
  const checkId = Number(job.check_run_url?.split("/").pop());
  if (
    job.run_id !== run.id || job.run_attempt !== run.run_attempt ||
    !Number.isSafeInteger(checkId) || checkId <= 0
  ) {
    throw new Error("Errcheck job does not belong to this workflow attempt.");
  }
  const { data: source } = await github.rest.checks.get({
    ...repository, check_run_id: checkId,
  });
  if (
    source.name !== job.name || source.head_sha !== run.head_sha ||
    source.check_suite?.id !== run.check_suite_id ||
    source.app?.id !== 15368 || source.app?.slug !== "github-actions" ||
    source.status !== "completed" ||
    !["success", "failure", "cancelled", "timed_out", "skipped", "neutral", "action_required"].includes(source.conclusion)
  ) {
    throw new Error("Errcheck Check Run source is invalid.");
  }

  const annotations = await github.paginate(github.rest.checks.listAnnotations, {
    ...repository, check_run_id: checkId, per_page: 100,
  });
  if (annotations.length !== source.output?.annotations_count) {
    throw new Error("Errcheck annotations are incomplete.");
  }
  const findings = annotations.filter((item) =>
    [item.message, item.title].some((text) =>
      typeof text === "string" && text.startsWith("[errcheck] reported by reviewdog")
    )
  );
  let conclusion = source.conclusion;
  let title = `execution ${source.conclusion}`;
  if (source.conclusion === "success") {
    conclusion = findings.length ? "neutral" : "success";
    title = findings.length
      ? `⚠ ${findings.length} reported (non-blocking)`
      : "no reported findings";
  } else if (source.conclusion === "skipped") {
    conclusion = "neutral";
    title = "skipped (not checked)";
  }

  const repositoryUrl = `${context.serverUrl}/${fullName}`;
  const sourceUrl = `${repositoryUrl}/actions/runs/${run.id}/job/${checkId}`;
  const lines = [
    `## Errcheck: ${title}`, "",
    `[Source job and full log](${sourceUrl})`, "",
    "Errcheck findings are advisory. The original job result is unchanged.", "",
  ];
  for (const finding of findings.slice(0, 20)) {
    if (
      typeof finding.path !== "string" || typeof finding.message !== "string" ||
      !Number.isSafeInteger(finding.start_line) || finding.start_line <= 0
    ) {
      throw new Error("Invalid errcheck annotation.");
    }
    const location = markdown(`${finding.path.slice(0, 256)}:${finding.start_line}`);
    const path = finding.path.split("/").map((part) =>
      encodeURIComponent(part).replace(/[!'()*]/g, (char) =>
        `%${char.charCodeAt(0).toString(16).toUpperCase()}`
      )
    ).join("/");
    const url = path.length <= 1500
      ? `${repositoryUrl}/blob/${run.head_sha}/${path}#L${finding.start_line}`
      : sourceUrl;
    const message = finding.message.replace(/^\[errcheck\] reported by reviewdog[^\n]*\n?/, "").split("\n")[0];
    lines.push(`- [${location}](${url}) — ${markdown(message.slice(0, 300))}`);
  }
  if (findings.length > 20) {
    lines.push("", "Showing the first 20 reported findings; see the full log for details.");
  }
  // shortcut: Actions can cap annotations; use reported counts and link the full log.
  lines.push("", "Counts cover available annotations; GitHub may limit annotations. The full log can contain more findings.");
  const summary = lines.join("\n");
  const externalId = `errcheck-report:${run.id}`;
  const checks = await github.paginate(github.rest.checks.listForRef, {
    ...repository, ref: run.head_sha, filter: "all", per_page: 100,
  });
  const reports = checks.filter((check) =>
    check.external_id === externalId && check.head_sha === run.head_sha &&
    check.app?.id === 15368 && check.app?.slug === "github-actions" &&
    check.name?.startsWith("errcheck findings / ")
  );
  if (reports.length > 1) {
    throw new Error("Errcheck report is ambiguous.");
  }
  const report = {
    ...repository, name: `errcheck findings / ${title}`, status: "completed",
    conclusion, external_id: externalId, details_url: sourceUrl,
    output: { title: `Errcheck: ${title}`, summary },
  };
  if (reports.length) {
    await github.rest.checks.update({ ...report, check_run_id: reports[0].id });
  } else {
    await github.rest.checks.create({ ...report, head_sha: run.head_sha });
  }
  await core.summary.addRaw(summary).write();
}

module.exports = { reportErrcheck };
