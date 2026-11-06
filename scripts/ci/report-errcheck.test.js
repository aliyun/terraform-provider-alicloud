"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { reportErrcheck } = require("./report-errcheck.js");

const SHA = "d2ecfc7387b2725fd381acb36c4cb7264db71791";
const APP = { id: 15368, slug: "github-actions" };

function harness() {
  const writes = [];
  const run = {
    id: 40010646353, run_attempt: 1, event: "pull_request", name: "Golangci-Lint",
    path: ".github/workflows/acctest-terraform-lint.yml", head_sha: SHA,
    repository: { full_name: "aliyun/terraform-provider-alicloud" },
    status: "completed", check_suite_id: 108451975031,
    // Fork PR runs can have no pull_requests entries.
    pull_requests: [],
  };
  const job = {
    name: "runner / errcheck", run_id: run.id, run_attempt: 1,
    check_run_url: "https://api.github.com/repos/aliyun/terraform-provider-alicloud/check-runs/120610082736",
  };
  const source = {
    name: job.name, head_sha: SHA, app: APP, status: "completed",
    conclusion: "success", check_suite: { id: run.check_suite_id },
    output: { annotations_count: 0 },
  };
  const annotations = [];
  const reports = [];
  const actions = {
    getWorkflowRun: async () => ({ data: run }),
    listJobsForWorkflowRun() {},
  };
  const checks = {
    get: async () => ({ data: source }),
    listAnnotations() {}, listForRef() {},
    create: async (value) => writes.push({ method: "create", ...value }),
    update: async (value) => writes.push({ method: "update", ...value }),
  };
  const github = {
    rest: { actions, checks },
    paginate: async (method, args) => {
      assert.equal(args.owner, "aliyun");
      assert.equal(args.repo, "terraform-provider-alicloud");
      assert.equal(args.per_page, 100);
      if (method === actions.listJobsForWorkflowRun) {
        assert.equal(args.filter, "latest");
        return [job];
      }
      if (method === checks.listAnnotations) return annotations;
      if (method === checks.listForRef) {
        assert.equal(args.ref, SHA);
        assert.equal(args.filter, "all");
        return reports;
      }
      throw new Error("Unexpected API call.");
    },
  };
  const context = {
    eventName: "workflow_run", serverUrl: "https://github.com",
    repo: { owner: "aliyun", repo: "terraform-provider-alicloud" },
    payload: { action: "completed", workflow_run: { id: run.id, run_attempt: 1 } },
  };
  const core = {
    info() {},
    summary: { addRaw() { return this; }, async write() {} },
  };
  return { run, job, source, annotations, reports, writes, github, context, core };
}

function finding(path = "alicloud/resource_alicloud_direct_mail_config_set.go", line = 176) {
  return {
    path, start_line: line, title: "", annotation_level: "failure",
    message: "[errcheck] reported by reviewdog 🐶\nError return value is not checked (errcheck)\n\nRaw Output:\nignored",
  };
}

function addFindings(h, items) {
  h.annotations.push(...items);
  h.source.output.annotations_count = h.annotations.length;
}

test("PR 10601's five findings produce a visible neutral check on the fork head", async () => {
  const h = harness();
  addFindings(h, [
    finding("alicloud/data_source_alicloud_direct_mail_config_sets.go", 244),
    finding("alicloud/data_source_alicloud_direct_mail_dedicated_ip_pools.go", 206),
    finding(undefined, 176), finding(undefined, 177), finding(undefined, 178),
  ]);
  await reportErrcheck(h);
  const report = h.writes[0];
  assert.equal(report.conclusion, "neutral");
  assert.equal(report.head_sha, SHA);
  assert.match(report.name, /5 reported \(non-blocking\)/);
  assert.match(report.output.summary, new RegExp(`/blob/${SHA}/alicloud/`));
  assert.match(report.output.summary, /data_source_alicloud_direct_mail_config_sets\.go#L244/);
  assert.match(report.output.summary, /#L178/);
  assert.equal(h.source.conclusion, "success");
  assert.equal(h.writes.length, 1);
});

test("a clean run reports success", async () => {
  const h = harness();
  await reportErrcheck(h);
  assert.equal(h.writes[0].conclusion, "success");
  assert.match(h.writes[0].name, /no reported findings/);
});

test("tool failures and timeouts stay failed even with findings", async () => {
  for (const conclusion of ["failure", "timed_out", "cancelled"]) {
    const h = harness();
    h.source.conclusion = conclusion;
    addFindings(h, [finding()]);
    await reportErrcheck(h);
    assert.equal(h.writes[0].conclusion, conclusion);
    assert.match(h.writes[0].name, /execution/);
  }
});

test("skipped checks are distinguished from clean checks", async () => {
  const h = harness();
  h.source.conclusion = "skipped";
  await reportErrcheck(h);
  assert.equal(h.writes[0].conclusion, "neutral");
  assert.match(h.writes[0].name, /not checked/);
});

test("obsolete attempts cannot overwrite newer reports", async () => {
  const h = harness();
  h.run.run_attempt = 2;
  await reportErrcheck(h);
  assert.equal(h.writes.length, 0);
});

test("rejects wrong workflow and source bindings before writing", async () => {
  const mutations = [
    (h) => { h.context.eventName = "pull_request"; },
    (h) => { h.run.event = "push"; },
    (h) => { h.run.path = ".github/workflows/untrusted.yml"; },
    (h) => { h.run.repository.full_name = "other/repository"; },
    (h) => { delete h.run.check_suite_id; delete h.source.check_suite.id; },
    (h) => { h.job.run_attempt = 2; },
    (h) => { h.job.check_run_url = "invalid"; },
    (h) => { h.source.head_sha = "a".repeat(40); },
    (h) => { h.source.check_suite.id += 1; },
    (h) => { h.source.app = { id: 1, slug: "other" }; },
    (h) => { h.source.status = "in_progress"; },
    (h) => { h.source.conclusion = null; },
    (h) => { h.source.output.annotations_count = 1; },
  ];
  for (const mutate of mutations) {
    const h = harness();
    mutate(h);
    await assert.rejects(reportErrcheck(h));
    assert.equal(h.writes.length, 0);
  }
});

test("unrelated annotations do not become errcheck findings", async () => {
  const h = harness();
  addFindings(h, [{ ...finding(), message: "Unrelated deprecation warning" }]);
  await reportErrcheck(h);
  assert.equal(h.writes[0].conclusion, "success");
});

test("reruns update only the matching report and can clear old findings", async () => {
  const h = harness();
  h.reports.push(
    { id: 2, name: "Basic Policy", app: APP, head_sha: SHA },
    { id: 3, name: "errcheck findings / ⚠ 5 reported (non-blocking)", app: APP,
      head_sha: SHA, external_id: `errcheck-report:${h.run.id}` },
  );
  await reportErrcheck(h);
  assert.equal(h.writes[0].method, "update");
  assert.equal(h.writes[0].check_run_id, 3);
  assert.equal(h.writes[0].conclusion, "success");
  assert.match(h.writes[0].name, /no reported findings/);
});

test("renders untrusted text as text and bounds large summaries", async () => {
  const h = harness();
  const item = finding("alicloud/a[link](x)<script>.go", 1);
  item.message = "[errcheck] reported by reviewdog 🐶\n<script>bad</script> [link]\n::error::not a command";
  addFindings(h, [item, ...Array.from({ length: 24 }, () => finding("文".repeat(512)))]);
  await reportErrcheck(h);
  const summary = h.writes[0].output.summary;
  assert.match(summary, /25 reported/);
  assert.match(summary, /%3Cscript%3E/);
  assert.match(summary, /%28x%29/);
  assert.ok(!summary.includes("<script>"));
  assert.ok(!summary.includes("::error::"));
  assert.match(summary, /first 20/);
  assert.ok(summary.length < 60000);
});
