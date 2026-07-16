# Jarvis execution architecture

Jarvis has two first-class work objects.  The distinction is recovery, not the
machine or command used to execute the work.

## Work objects

### Task

A Task represents one business outcome that must eventually converge after a
process or machine interruption.  It is persisted by the control plane and may
own multiple fenced Sessions over its lifetime.

Examples: Aone ticket delivery, revisit/wake continuation, persona handoff, and
PR CI/comment follow-up.  These triggers update the same Aone Task identity;
they do not introduce another execution mode.

Task states are `READY`, `LEASED`, `RUNNING`, `FINALIZING`, `SUSPENDED`,
`RETRY_WAIT`, `RECOVERY_REQUIRED`, `SUCCEEDED`, `FAILED_FINAL`, and `CANCELED`.
Session states are `CREATED`, `LEASED`, `RUNNING`, `SUSPENDED`, `RESUMABLE`,
`CLOSED`, `CORRUPTED`, and `CANCELED`.

### EphemeralJob

An EphemeralJob is a disposable local execution.  It has an execution id and
result/logs, but no Task row, Session, lease, fence, or cross-process recovery
promise.  Examples include probe scans, preflight/config checks, simple reads,
Tata/chat presentation, and subcommands already enclosed by a Task Session.

## Components

- `ExecutionRouter` maps `needs_recovery` to `TASK` or `EPHEMERAL_JOB`.
- `PersistenceExecutor` leases Tasks, owns `SessionController`, and reports terminal
  state.
- `EphemeralExecutor` executes disposable jobs.
- `CapacityManager` atomically allocates the shared machine slots before a Task
  lease or EphemeralJob process starts, preventing cross-executor oversubscription.
- `ExecutionRuntime` and `ProcessGuardian` provide common process launch,
  stdout/stderr capture, timeout, cancellation, process-group cleanup, and the
  fenced pre-exec gate.
- `ControlPlaneClient` covers Task, Session, Worker, Event, and Operation APIs.
- `SessionController` owns fenced Task Session transitions.

Every Task is persisted by the control plane and can only be executed by a
`PersistenceExecutor`.  A control-plane failure is fail-closed: Task work stays
unstarted instead of falling back to an untracked local process.  Pausing a
rollout means stopping PersistenceExecutors so Tasks remain queued; it does not
introduce another execution mode.

## Interactive Session ownership

Defaults:

```text
Lease           300 seconds
Heartbeat        30 seconds
Safety margin    90 seconds
Turn grace      600 seconds
Worker affinity 7200 seconds
```

The detached sidecar is the only remote heartbeat sender.  Each accepted
claim/start/heartbeat writes an atomic local `SessionPermit` containing the Task,
Session, Worker, fence, authoritative lease expiry, and sidecar health.  A
`PreToolUse` hook performs only local fail-closed permit validation; it does not
block a tool call on a network heartbeat.

After a Codex Desktop turn stops, the Session remains reusable for the
10-minute grace.  It is then suspended and the Task is parked without consuming
a slot.  The suspend detail carries a two-hour affinity expiry so the original
Worker is preferred before ordinary recovery is allowed.  Claude remains
process-scoped rather than using the Codex between-turn grace.  Explicit wait
states are not converted to turn-idle failures.

## External side effects

Operation receipts protect writes that are unsafe to replay blindly, such as
Aone state/comments, PR/CR mutations, notifications, and formal external
triggers.  Internal reads and ordinary subprocess steps are not Operations.
