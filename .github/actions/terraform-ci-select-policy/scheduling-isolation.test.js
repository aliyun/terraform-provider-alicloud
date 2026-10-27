"use strict";

// Scheduling-isolation regression tests for the CI routing decoupling.
//
// Background: the Basic Policy Gate (acctest-terraform-basic-policy.yml, job
// `policy`) runs on the `terraform-ci-policy` runner pool. Five read-only
// consumer workflows (acctest-terraform-lint, consistency-check,
// testing-coverage-rate, collection-order-coverage, document-check) plus the
// inline-route consumer (acctest-terraform-basic) each have a `route` job that
// used to run on the SAME `terraform-ci-policy` pool and poll for the Gate's
// Check Run. Under batch concurrency the route jobs saturated the pool while
// waiting, queueing the Gate behind its own consumers and causing the route
// discovery budget to expire -> false "external/blocked" classification and
// spurious CI failures.
//
// The decoupling moves every consumer `route` job to `ubuntu-latest`, leaving
// `terraform-ci-policy` to the Gate. The route job only runs github-script
// (it reads the completed Check Run via the GitHub API; no fork code is
// checked out and the token stays fork-scoped), so ubuntu-latest is sufficient
// and does not widen any security boundary.
//
// These tests pin that invariant so a future regression (routing a consumer
// back onto the Gate's pool) fails fast in ci-policy-test.yml rather than
// reintroducing the batch race silently.

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");

const REPO_ROOT = path.resolve(__dirname, "../../..");
const WORKFLOWS_DIR = path.join(REPO_ROOT, ".github", "workflows");

const CONSUMER_WORKFLOWS = [
  { file: "acctest-terraform-lint.yml", job: "route", route: "composite" },
  { file: "consistency-check.yml", job: "route", route: "composite" },
  { file: "testing-coverage-rate.yml", job: "route", route: "composite" },
  { file: "collection-order-coverage.yml", job: "route", route: "composite" },
  { file: "document-check.yml", job: "route", route: "composite" },
  { file: "acctest-terraform-basic.yml", job: "route", route: "inline" },
];

const BASIC_POLICY_WORKFLOW = {
  file: "acctest-terraform-basic-policy.yml",
  job: "policy",
};

const GATE_POOL = "terraform-ci-policy";
const CONSUMER_POOL = "ubuntu-latest";

function readWorkflow(file) {
  const target = path.join(WORKFLOWS_DIR, file);
  assert.ok(fs.existsSync(target), `workflow file not found: ${file}`);
  return fs.readFileSync(target, "utf8");
}

function escapeRegex(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function stripInlineComment(value) {
  const hash = value.indexOf("#");
  if (hash >= 0) value = value.slice(0, hash);
  return value.replace(/^['"]|['"]$/g, "").trim();
}

// Locate the literal runs-on value for a named job. The route and policy jobs
// use a scalar label (no ${{ }} expression), so a line scan is sufficient and
// avoids pulling in a YAML dependency.
function runsOnForJob(yaml, jobName) {
  const lines = yaml.split(/\r?\n/);
  const headerRe = new RegExp(`^  ${escapeRegex(jobName)}:\\s*$`);
  let inside = false;
  for (const line of lines) {
    if (!inside) {
      if (headerRe.test(line)) inside = true;
      continue;
    }
    // A new 2-space-indented job header or a top-level key ends this job.
    if (/^[A-Za-z]/.test(line) || /^  [A-Za-z]/.test(line)) break;
    const match = line.match(/^    runs-on:\s*(.+?)\s*$/);
    if (match) return stripInlineComment(match[1]);
  }
  return null;
}

// Collect the permission lines scoped under a named job, so we can assert the
// read-only token posture is preserved alongside the runner-label move.
function jobPermissions(yaml, jobName) {
  const lines = yaml.split(/\r?\n/);
  const headerRe = new RegExp(`^  ${escapeRegex(jobName)}:\\s*$`);
  const permissions = new Set();
  let inside = false;
  let inPermissions = false;
  for (const line of lines) {
    if (!inside) {
      if (headerRe.test(line)) inside = true;
      continue;
    }
    if (/^[A-Za-z]/.test(line) || /^  [A-Za-z]/.test(line)) break;
    if (/^    permissions:\s*$/.test(line)) {
      inPermissions = true;
      continue;
    }
    if (inPermissions) {
      // A sibling 4-space key (not 6-space permission entry) closes the block.
      if (/^    [A-Za-z]/.test(line)) {
        inPermissions = false;
        continue;
      }
      const match = line.match(/^\s+([A-Za-z-]+):\s*([A-Za-z-]+)\s*$/);
      if (match) permissions.add(`${match[1]}:${match[2]}`);
    }
  }
  return permissions;
}

test("every consumer route job runs on ubuntu-latest (decoupled from the Gate pool)", () => {
  for (const { file, job } of CONSUMER_WORKFLOWS) {
    const yaml = readWorkflow(file);
    const value = runsOnForJob(yaml, job);
    assert.ok(
      value !== null,
      `${file}: ${job} job has no runs-on`
    );
    assert.equal(
      value,
      CONSUMER_POOL,
      `${file}: consumer ${job} job must run on ${CONSUMER_POOL} to avoid saturating the Basic Policy Gate pool; found "${value}"`
    );
  }
});

test("the Basic Policy Gate keeps terraform-ci-policy (producer pool unchanged)", () => {
  const yaml = readWorkflow(BASIC_POLICY_WORKFLOW.file);
  const value = runsOnForJob(yaml, BASIC_POLICY_WORKFLOW.job);
  assert.ok(
    value !== null,
    `${BASIC_POLICY_WORKFLOW.file}: ${BASIC_POLICY_WORKFLOW.job} job has no runs-on`
  );
  assert.equal(
    value,
    GATE_POOL,
    `${BASIC_POLICY_WORKFLOW.file}: Basic Policy Gate must keep ${GATE_POOL}; found "${value}"`
  );
});

test("consumer route jobs and the Basic Policy Gate are on different pools (no scheduling race)", () => {
  const basic = runsOnForJob(
    readWorkflow(BASIC_POLICY_WORKFLOW.file),
    BASIC_POLICY_WORKFLOW.job
  );
  assert.ok(basic !== null);
  assert.equal(basic, GATE_POOL);
  for (const { file, job } of CONSUMER_WORKFLOWS) {
    const consumer = runsOnForJob(readWorkflow(file), job);
    assert.ok(consumer !== null, `${file}: ${job} job missing runs-on`);
    assert.notEqual(
      consumer,
      basic,
      `${file}: consumer ${job} job (${consumer}) must not share the Basic Policy pool (${basic})`
    );
    assert.equal(
      consumer,
      CONSUMER_POOL,
      `${file}: consumer ${job} job must be on ${CONSUMER_POOL}`
    );
  }
});

test("consumer route jobs keep the read-only policy token (checks:read, pull-requests:read)", () => {
  // Moving the route job to ubuntu-latest must not widen its token scope. The
  // route job does not check out fork code; it only reads the Basic Policy
  // Check Run via the GitHub API. This guards against an accidental privilege
  // change riding on the runner-label move.
  for (const { file, job } of CONSUMER_WORKFLOWS) {
    const yaml = readWorkflow(file);
    const permissions = jobPermissions(yaml, job);
    assert.ok(
      permissions.has("checks:read"),
      `${file}: ${job} job must keep permissions.checks=read`
    );
    assert.ok(
      permissions.has("pull-requests:read"),
      `${file}: ${job} job must keep permissions.pull-requests=read`
    );
  }
});

test("the Basic Policy Gate keeps checks:write (producer authority unchanged)", () => {
  const yaml = readWorkflow(BASIC_POLICY_WORKFLOW.file);
  const permissions = jobPermissions(yaml, BASIC_POLICY_WORKFLOW.job);
  assert.ok(
    permissions.has("checks:write"),
    `${BASIC_POLICY_WORKFLOW.file}: ${BASIC_POLICY_WORKFLOW.job} job must keep checks:write to create the Basic Policy Check Run`
  );
  assert.ok(
    permissions.has("pull-requests:read"),
    `${BASIC_POLICY_WORKFLOW.file}: ${BASIC_POLICY_WORKFLOW.job} job must keep pull-requests:read`
  );
});

test("consumer downstream jobs still resolve their runner from route outputs (binding preserved)", () => {
  // The decoupling only moves the route job's own runner; downstream jobs must
  // keep consuming needs.route.outputs.runner / heavy_runner so the route
  // decision (trusted -> terraform-ci-trusted, external -> ubuntu-latest) is
  // unchanged. This catches a regression where a downstream job is hardcoded
  // to a pool instead of honoring the route classification.
  for (const { file } of CONSUMER_WORKFLOWS) {
    const yaml = readWorkflow(file);
    // Every consumer workflow has at least one downstream job that consumes
    // route.outputs.runner or route.outputs.heavy_runner via fromJSON.
    assert.ok(
      /fromJSON\(needs\.route\.outputs\.(runner|heavy_runner)\)/.test(yaml),
      `${file}: downstream jobs must resolve runner from route outputs via fromJSON(needs.route.outputs.runner|heavy_runner)`
    );
  }
});

test("all consumer workflows keep permissions:{} at the workflow level (no default token widening)", () => {
  for (const { file } of CONSUMER_WORKFLOWS) {
    const yaml = readWorkflow(file);
    assert.ok(
      /^permissions:\s*\{\s*\}\s*$/m.test(yaml),
      `${file}: workflow-level permissions:{} must remain so the default GITHUB_TOKEN is not widened`
    );
  }
});
