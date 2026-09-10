"use strict";

const POLICY_VERSION = 1;
const POLICY_CHECK_NAME = "Basic Policy";
const POLICY_TIMEOUT_MS = 90000;
const POLICY_POLL_INTERVAL_MS = 3000;
const TRUSTED_RUNNER = "terraform-ci-trusted";
const TRUSTED_HEAVY_RUNNER = "terraform-ci-heavy";
const EXTERNAL_RUNNER = "ubuntu-latest";
const SHA_PATTERN = /^[0-9a-f]{40}$/;
const REPOSITORY_PATTERN =
  /^[A-Za-z0-9](?:[A-Za-z0-9-]{0,38})\/[A-Za-z0-9._-]{1,100}$/;
const PULL_REQUEST_ACTIONS = new Set([
  "opened",
  "reopened",
  "synchronize",
  "ready_for_review",
]);
const SUPPORTED_BASE_REFS = new Set(["master", "release/v2"]);
const SUPPORTED_PUSH_REFS = new Set(
  Array.from(SUPPORTED_BASE_REFS, (branch) => `refs/heads/${branch}`)
);
const ACTIONS_APP = Object.freeze({
  id: 15368,
  slug: "github-actions",
  name: "GitHub Actions",
  owner: "github",
});

function isRecord(value) {
  return value !== null && typeof value === "object" && !Array.isArray(value);
}

function isCommitSha(value) {
  return (
    typeof value === "string" &&
    SHA_PATTERN.test(value) &&
    value !== "0".repeat(40)
  );
}

function expectedRepository(context) {
  const owner = context?.repo?.owner;
  const repo = context?.repo?.repo;
  const fullName = `${owner}/${repo}`;
  if (
    typeof owner !== "string" ||
    typeof repo !== "string" ||
    !REPOSITORY_PATTERN.test(fullName)
  ) {
    throw new Error("Repository metadata is invalid.");
  }
  return { fullName, owner, repo };
}

function expectedServerOrigin(context) {
  if (typeof context?.serverUrl !== "string") {
    throw new Error("Server metadata is invalid.");
  }
  let server;
  try {
    server = new URL(context.serverUrl);
  } catch {
    throw new Error("Server metadata is invalid.");
  }
  if (
    server.protocol !== "https:" ||
    server.username !== "" ||
    server.password !== "" ||
    (server.pathname !== "" && server.pathname !== "/") ||
    server.search !== "" ||
    server.hash !== ""
  ) {
    throw new Error("Server metadata is invalid.");
  }
  return server.origin;
}

function pullRequestMetadata(context, repository) {
  const payload = context?.payload;
  const pullRequest = payload?.pull_request;
  if (!isRecord(payload) || !isRecord(pullRequest)) {
    throw new Error("Pull request metadata is invalid.");
  }

  const pullNumber = pullRequest.number;
  const headRepository = pullRequest.head?.repo?.full_name;
  const headSha = pullRequest.head?.sha;
  const baseRepository = pullRequest.base?.repo?.full_name;
  const baseRef = pullRequest.base?.ref;
  if (
    !Number.isSafeInteger(pullNumber) ||
    pullNumber <= 0 ||
    typeof headRepository !== "string" ||
    !REPOSITORY_PATTERN.test(headRepository) ||
    !isCommitSha(headSha) ||
    baseRepository !== repository.fullName ||
    typeof baseRef !== "string" ||
    !SUPPORTED_BASE_REFS.has(baseRef)
  ) {
    throw new Error("Pull request metadata is invalid.");
  }

  return { headRepository, headSha, pullNumber };
}

function eventMetadata(context, repository) {
  if (context?.eventName === "push") {
    const payload = context.payload;
    const ref = context.ref;
    if (
      !isRecord(payload) ||
      typeof ref !== "string" ||
      !SUPPORTED_PUSH_REFS.has(ref) ||
      payload.action !== undefined ||
      payload.repository?.full_name !== repository.fullName ||
      payload.ref !== ref ||
      payload.deleted !== false ||
      !isCommitSha(context.sha) ||
      payload.after !== context.sha
    ) {
      throw new Error("Push metadata is invalid.");
    }
    return { kind: "push" };
  }

  const pullRequest = pullRequestMetadata(context, repository);
  if (context?.eventName === "pull_request") {
    if (!PULL_REQUEST_ACTIONS.has(context.payload.action)) {
      throw new Error("Pull request metadata is invalid.");
    }
    return { kind: "pull_request", ...pullRequest };
  }

  if (context?.eventName === "pull_request_review") {
    const review = context.payload.review;
    if (
      context.payload.action !== "submitted" ||
      !isRecord(review) ||
      review.state !== "approved" ||
      review.commit_id !== pullRequest.headSha
    ) {
      throw new Error("Pull request review metadata is invalid.");
    }
    return { kind: "pull_request_review", ...pullRequest };
  }

  throw new Error("Unsupported event.");
}

function parseBinding(externalId, pullRequest) {
  if (
    typeof externalId !== "string" ||
    externalId.length === 0 ||
    externalId.length > 255
  ) {
    throw new Error("Basic Policy binding is invalid.");
  }

  let binding;
  try {
    binding = JSON.parse(externalId);
  } catch {
    throw new Error("Basic Policy binding is invalid.");
  }
  if (
    !isRecord(binding) ||
    Object.keys(binding).sort().join(",") !== "c,p,q,r,s,v" ||
    binding.v !== POLICY_VERSION ||
    binding.p !== pullRequest.pullNumber ||
    binding.s !== pullRequest.headSha ||
    binding.r !== pullRequest.headRepository ||
    !["trusted", "external"].includes(binding.c) ||
    typeof binding.q !== "boolean"
  ) {
    throw new Error("Basic Policy binding is stale or invalid.");
  }
  return binding;
}

function validateSource(check, repository, pullRequest, origin) {
  const checkId = check?.id;
  const canonicalUrl =
    `${origin}/${repository.owner}/${repository.repo}/runs/${checkId}`;
  if (
    !Number.isSafeInteger(checkId) ||
    checkId <= 0 ||
    check?.name !== POLICY_CHECK_NAME ||
    check?.head_sha !== pullRequest.headSha ||
    check?.app?.id !== ACTIONS_APP.id ||
    check?.app?.slug !== ACTIONS_APP.slug ||
    check?.app?.name !== ACTIONS_APP.name ||
    check?.app?.owner?.login !== ACTIONS_APP.owner ||
    check?.details_url !== canonicalUrl ||
    check?.html_url !== canonicalUrl
  ) {
    throw new Error("Basic Policy Check Run source is invalid.");
  }
}

function setOutputs(core, classification, relevant) {
  if (!core || typeof core.setOutput !== "function") {
    throw new Error("Actions output interface is invalid.");
  }
  if (
    !["trusted", "external"].includes(classification) ||
    typeof relevant !== "boolean"
  ) {
    throw new Error("Actions routing output is invalid.");
  }
  const runner =
    classification === "trusted" ? TRUSTED_RUNNER : EXTERNAL_RUNNER;
  const heavyRunner =
    classification === "trusted" ? TRUSTED_HEAVY_RUNNER : EXTERNAL_RUNNER;
  const outputs = {
    classification,
    relevant: String(relevant),
    runner: JSON.stringify(runner),
    heavy_runner: JSON.stringify(heavyRunner),
  };
  for (const [name, value] of Object.entries(outputs)) {
    core.setOutput(name, value);
  }
  return outputs;
}

async function selectPolicy({
  github,
  context,
  core,
  now = Date.now,
  sleep = (milliseconds) =>
    new Promise((resolve) => setTimeout(resolve, milliseconds)),
  timeoutMs = POLICY_TIMEOUT_MS,
  pollIntervalMs = POLICY_POLL_INTERVAL_MS,
}) {
  const repository = expectedRepository(context);
  const origin = expectedServerOrigin(context);
  const event = eventMetadata(context, repository);

  if (event.kind === "push") {
    return setOutputs(core, "trusted", true);
  }

  if (
    !github ||
    typeof github.paginate !== "function" ||
    typeof github.rest?.checks?.listForRef !== "function"
  ) {
    throw new Error("GitHub Checks interface is invalid.");
  }
  if (
    typeof now !== "function" ||
    typeof sleep !== "function" ||
    !Number.isSafeInteger(timeoutMs) ||
    timeoutMs <= 0 ||
    timeoutMs > POLICY_TIMEOUT_MS ||
    !Number.isSafeInteger(pollIntervalMs) ||
    pollIntervalMs <= 0 ||
    pollIntervalMs > timeoutMs
  ) {
    throw new Error("Polling configuration is invalid.");
  }

  const deadline = now() + timeoutMs;
  while (now() < deadline) {
    const checkRuns = await github.paginate(
      github.rest.checks.listForRef,
      {
        owner: repository.owner,
        repo: repository.repo,
        ref: event.headSha,
        check_name: POLICY_CHECK_NAME,
        filter: "all",
        per_page: 100,
      }
    );
    if (!Array.isArray(checkRuns)) {
      throw new Error("Basic Policy Check Run lookup is invalid.");
    }

    const candidates = checkRuns.filter(
      (check) => check?.name === POLICY_CHECK_NAME
    );
    if (candidates.length === 0) {
      await sleep(Math.min(pollIntervalMs, Math.max(0, deadline - now())));
      continue;
    }
    if (candidates.length !== 1) {
      throw new Error("Basic Policy Check Run is ambiguous.");
    }

    const check = candidates[0];
    validateSource(check, repository, event, origin);
    const binding = parseBinding(check.external_id, event);
    if (check.status !== "completed") {
      if (!["queued", "in_progress"].includes(check.status)) {
        throw new Error("Basic Policy Check Run status is invalid.");
      }
      await sleep(Math.min(pollIntervalMs, Math.max(0, deadline - now())));
      continue;
    }
    if (check.conclusion !== "success") {
      throw new Error("Basic Policy rejected this event.");
    }
    return setOutputs(core, binding.c, binding.q);
  }

  throw new Error("Timed out waiting for a unique Basic Policy result.");
}

module.exports = {
  selectPolicy,
};
