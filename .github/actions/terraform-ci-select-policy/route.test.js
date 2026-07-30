"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");

const { selectPolicy } = require("./route.js");

const OWNER = "aliyun";
const REPO = "terraform-provider-alicloud";
const HEAD_REPOSITORY = "api-tool-agent/terraform-provider-alicloud";
const HEAD_SHA = "1234567890abcdef1234567890abcdef12345678";
const OTHER_SHA = "abcdef1234567890abcdef1234567890abcdef12";
const CHECK_ID = 123456789;
const CHECK_URL =
  `https://github.com/${OWNER}/${REPO}/runs/${CHECK_ID}`;

function pullRequest() {
  return {
    number: 100,
    head: {
      repo: { full_name: HEAD_REPOSITORY },
      sha: HEAD_SHA,
    },
    base: {
      repo: { full_name: `${OWNER}/${REPO}` },
      ref: "master",
    },
  };
}

function pullRequestContext(overrides = {}) {
  const eventName = overrides.eventName ?? "pull_request";
  const payload = {
    action: overrides.action ?? "synchronize",
    pull_request: pullRequest(),
  };
  if (eventName === "pull_request_review") {
    payload.action = overrides.action ?? "submitted";
    payload.review = {
      state: overrides.reviewState ?? "approved",
      commit_id: overrides.reviewCommitId ?? HEAD_SHA,
      body: overrides.reviewBody ?? "",
    };
  }
  return {
    eventName,
    payload,
    repo: { owner: OWNER, repo: REPO },
    serverUrl: "https://github.com",
    sha: HEAD_SHA,
    ref: `refs/pull/${payload.pull_request.number}/merge`,
  };
}

function pushContext(overrides = {}) {
  const sha = overrides.sha ?? HEAD_SHA;
  const ref = overrides.ref ?? "refs/heads/master";
  const repository =
    overrides.repository ?? `${OWNER}/${REPO}`;
  return {
    eventName: "push",
    payload: {
      after: overrides.after ?? sha,
      deleted: overrides.deleted ?? false,
      ref,
      repository: { full_name: repository },
    },
    repo: { owner: OWNER, repo: REPO },
    serverUrl: "https://github.com",
    sha,
    ref,
  };
}

function binding(overrides = {}) {
  return {
    v: 1,
    p: 100,
    s: HEAD_SHA,
    r: HEAD_REPOSITORY,
    c: "trusted",
    q: true,
    ...overrides,
  };
}

function checkRun(overrides = {}) {
  const value = {
    id: CHECK_ID,
    name: "Basic Policy",
    head_sha: HEAD_SHA,
    external_id: JSON.stringify(binding()),
    status: "completed",
    conclusion: "success",
    details_url: CHECK_URL,
    html_url: CHECK_URL,
    app: {
      id: 15368,
      slug: "github-actions",
      name: "GitHub Actions",
      owner: { login: "github" },
    },
  };
  return Object.assign(value, overrides);
}

function harness(responses = [[]]) {
  let clock = 0;
  let index = 0;
  const calls = [];
  const outputs = {};
  const listForRef = async () => {
    throw new Error("paginate must own the request");
  };
  const github = {
    rest: { checks: { listForRef } },
    paginate: async (method, args) => {
      calls.push({ method, args });
      const response = responses[Math.min(index, responses.length - 1)];
      index += 1;
      return response;
    },
  };
  const core = {
    setOutput(name, value) {
      outputs[name] = value;
    },
  };
  const options = {
    github,
    core,
    now: () => clock,
    sleep: async (milliseconds) => {
      clock += milliseconds;
    },
    timeoutMs: 2,
    pollIntervalMs: 1,
  };
  return { calls, github, core, options, outputs };
}

async function execute(context, responses) {
  const testHarness = harness(responses);
  await selectPolicy({ context, ...testHarness.options });
  return testHarness;
}

test("pull_request selects the trusted runner from a valid binding", async () => {
  const result = await execute(pullRequestContext(), [[checkRun()]]);

  assert.deepEqual(result.outputs, {
    classification: "trusted",
    relevant: "true",
    runner: JSON.stringify("terraform-ci-trusted"),
    heavy_runner: JSON.stringify("terraform-ci-heavy"),
  });
  assert.equal(result.calls.length, 1);
  assert.deepEqual(result.calls[0].args, {
    owner: OWNER,
    repo: REPO,
    ref: HEAD_SHA,
    check_name: "Basic Policy",
    filter: "all",
    per_page: 100,
  });
});

test("pull_request selects the hosted runner for an external binding", async () => {
  const external = checkRun({
    external_id: JSON.stringify(binding({ c: "external", q: false })),
  });
  const result = await execute(pullRequestContext(), [[external]]);

  assert.deepEqual(result.outputs, {
    classification: "external",
    relevant: "false",
    runner: JSON.stringify("ubuntu-latest"),
    heavy_runner: JSON.stringify("ubuntu-latest"),
  });
});

test("pull_request_review accepts only an exact approved review", async () => {
  const result = await execute(
    pullRequestContext({ eventName: "pull_request_review" }),
    [[checkRun()]]
  );
  assert.equal(result.outputs.classification, "trusted");
  assert.equal(result.outputs.relevant, "true");
});

test("pull_request_review rejects edited reviews", async () => {
  await assert.rejects(
    execute(
      pullRequestContext({
        eventName: "pull_request_review",
        action: "edited",
      }),
      [[checkRun()]]
    ),
    /review metadata/i
  );
});

test("pull_request_review ignores approval text in the body", async () => {
  await assert.rejects(
    execute(
      pullRequestContext({
        eventName: "pull_request_review",
        reviewState: "commented",
        reviewBody: "approved",
      }),
      [[checkRun()]]
    ),
    /review metadata/i
  );
});

test("pull_request_review rejects approval for a stale head", async () => {
  await assert.rejects(
    execute(
      pullRequestContext({
        eventName: "pull_request_review",
        reviewCommitId: OTHER_SHA,
      }),
      [[checkRun()]]
    ),
    /review metadata/i
  );
});

test("push to the repository master branch is trusted and relevant", async () => {
  const result = await execute(pushContext(), []);
  assert.deepEqual(result.outputs, {
    classification: "trusted",
    relevant: "true",
    runner: JSON.stringify("terraform-ci-trusted"),
    heavy_runner: JSON.stringify("terraform-ci-heavy"),
  });
  assert.equal(result.calls.length, 0);
});

test("push rejects another repository", async () => {
  await assert.rejects(
    execute(pushContext({ repository: "someone/else" }), []),
    /push metadata/i
  );
});

test("push rejects another ref", async () => {
  await assert.rejects(
    execute(pushContext({ ref: "refs/heads/feature" }), []),
    /push metadata/i
  );
});

test("push rejects a stale or deleted commit", async () => {
  await assert.rejects(
    execute(pushContext({ after: OTHER_SHA }), []),
    /push metadata/i
  );
  await assert.rejects(
    execute(pushContext({ deleted: true }), []),
    /push metadata/i
  );
  await assert.rejects(
    execute(pushContext({ sha: "0".repeat(40) }), []),
    /push metadata/i
  );
});

test("zero Basic Policy checks times out closed", async () => {
  await assert.rejects(
    execute(pullRequestContext(), [[]]),
    /timed out/i
  );
});

test("duplicate Basic Policy checks fail closed", async () => {
  await assert.rejects(
    execute(pullRequestContext(), [[checkRun(), checkRun({ id: 2 })]]),
    /ambiguous/i
  );
});

for (const [name, mutate] of [
  ["id", (check) => { check.app.id = 1; }],
  ["slug", (check) => { check.app.slug = "other"; }],
  ["name", (check) => { check.app.name = "Other"; }],
  ["owner", (check) => { check.app.owner.login = "other"; }],
]) {
  test(`wrong GitHub Actions app ${name} fails closed`, async () => {
    const candidate = checkRun();
    mutate(candidate);
    await assert.rejects(
      execute(pullRequestContext(), [[candidate]]),
      /source/i
    );
  });
}

for (const [name, mutate] of [
  ["details URL", (check) => { check.details_url += "?spoofed=1"; }],
  ["HTML URL", (check) => { check.html_url += "?spoofed=1"; }],
]) {
  test(`wrong canonical ${name} fails closed`, async () => {
    const candidate = checkRun();
    mutate(candidate);
    await assert.rejects(
      execute(pullRequestContext(), [[candidate]]),
      /source/i
    );
  });
}

test("wrong Check Run head SHA fails closed", async () => {
  await assert.rejects(
    execute(
      pullRequestContext(),
      [[checkRun({ head_sha: OTHER_SHA })]]
    ),
    /source/i
  );
});

for (const [name, externalId] of [
  ["invalid JSON", "not-json"],
  ["extra key", JSON.stringify({ ...binding(), extra: true })],
  ["stale SHA", JSON.stringify(binding({ s: OTHER_SHA }))],
  ["wrong pull request", JSON.stringify(binding({ p: 101 }))],
  ["wrong repository", JSON.stringify(binding({ r: "someone/else" }))],
  ["blocked classification", JSON.stringify(binding({ c: "blocked" }))],
  ["non-boolean relevance", JSON.stringify(binding({ q: "true" }))],
]) {
  test(`external_id ${name} fails closed`, async () => {
    await assert.rejects(
      execute(
        pullRequestContext(),
        [[checkRun({ external_id: externalId })]]
      ),
      /binding/i
    );
  });
}

for (const classification of ["blocked", "invalid"]) {
  test(`${classification} classification emits no executable routes`, async () => {
    const candidate = checkRun({
      external_id: JSON.stringify(binding({ c: classification })),
    });
    const result = harness([[candidate]]);

    await assert.rejects(
      selectPolicy({
        context: pullRequestContext(),
        ...result.options,
      }),
      /binding/i
    );
    assert.deepEqual(result.outputs, {});
  });
}

test("an overlong external_id fails closed", async () => {
  await assert.rejects(
    execute(
      pullRequestContext(),
      [[checkRun({ external_id: "x".repeat(256) })]]
    ),
    /binding/i
  );
});

test("queued and in-progress checks are polled before success", async () => {
  const queued = checkRun({ status: "queued", conclusion: null });
  const inProgress = checkRun({ status: "in_progress", conclusion: null });
  const result = harness([[queued], [inProgress], [checkRun()]]);
  result.options.timeoutMs = 4;
  await selectPolicy({
    context: pullRequestContext(),
    ...result.options,
  });
  assert.equal(result.calls.length, 3);
  assert.equal(result.outputs.classification, "trusted");
});

test("an unknown Check Run status fails closed", async () => {
  await assert.rejects(
    execute(
      pullRequestContext(),
      [[checkRun({ status: "waiting", conclusion: null })]]
    ),
    /status/i
  );
});

test("a completed non-success Check Run fails closed", async () => {
  await assert.rejects(
    execute(
      pullRequestContext(),
      [[checkRun({ conclusion: "failure" })]]
    ),
    /rejected/i
  );
});

test("malformed pull request metadata fails closed", async () => {
  const context = pullRequestContext();
  context.payload.pull_request.number = "100";
  await assert.rejects(
    execute(context, [[checkRun()]]),
    /pull request metadata/i
  );
});

test("unsupported events fail closed", async () => {
  const context = pullRequestContext();
  context.eventName = "workflow_dispatch";
  await assert.rejects(
    execute(context, [[checkRun()]]),
    /unsupported event/i
  );
});
