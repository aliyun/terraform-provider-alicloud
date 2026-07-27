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
- The Scheduler `reply` runner reconstructs current `AONE_REPLY` waits from the
  control plane and persists wake Tasks when a human response arrives.

Every Task is persisted by the control plane and can only be executed by a
`PersistenceExecutor`.  A control-plane failure is fail-closed: Task work stays
unstarted instead of falling back to an untracked local process.  An already
running Session tolerates transient heartbeat transport/5xx failures while its
last successful fenced renewal proves more than the 90-second safety margin.
Heartbeat retries remain on the 30-second Session loop; a stale fence stops
immediately, and continued unavailability stops the process when the proof
reaches the safety boundary.  Pausing a rollout means stopping
PersistenceExecutors so Tasks remain queued; it does not introduce another
execution mode.

The queue-pull Worker is installation-scoped, not process-scoped. Its stable
`workerKey` is stored in `.my-day/bridge/worker-id` (override with
`JARVIS_WORKER_ID_FILE`), while every bridge start creates a new `processUuid`.
Ordinary mutations carry both identities. During a planned stop or explicit
bridge restart the executor reports `DRAINING`, synchronously stops local task
processes, relinquishes each Session as `RESUMABLE` without increasing Task
retry count, and finally reports `OFFLINE`. The next process re-registers the
same Worker and resumes the same runtime session on that machine. A crash
cannot be taken over until the old heartbeat/leases are stale and no Session
remains `LEASED` or `RUNNING`.

## Managed wait wake-up

A managed Session persists `waitType`, `waitKey`, `waitCursor`, and expiry in
the control plane before relinquishing its Worker slot. The bridge pages over
current Aone reply waits and polls only comments newer than that durable cursor.
It advances the Task with a deterministic `comment:<id>` revision only after
the control plane accepts the observation.

Bridge restart or machine replacement therefore loses only an in-memory poll
throttle, never a wait. A control-plane 503 or Aone read failure advances no
cursor and is retried on a later sensor tick. Concurrent sensors are safe
because Aone wake cursors cannot move backwards. The accepted wake atomically
changes the suspended Session to `RESUMABLE` and binds its continuation input;
the next lease reuses the original runtime session and receives the new reply.

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

## Dead interactive Session recovery

Interrupted Task recovery is owned by the control plane. Its reaper expires
stale leases, fences the old Session and returns replay-safe work to the
durable Task lifecycle. The bridge does not keep a second local redispatch
ledger or spawn replacement work from Scheduler state.

The `claim_health` runner independently reads the current
`jarvis-claimed` inventory and corroborates each item with `tasks/by-aone` and
the Task timeline. It treats advancing RUNNING/LEASED Sessions and unexpired
waits as healthy, allows 15 minutes for heartbeat/lease convergence, and
requires two observations at least five minutes apart for missing Tasks,
terminal residue or malformed control-plane structure. It only publishes an
idempotent alert; it never releases a claim or executes business work.

## External side effects

Operation receipts protect writes that are unsafe to replay blindly, such as
Aone state/comments, PR/CR mutations, notifications, and formal external
triggers.  Internal reads and ordinary subprocess steps are not Operations.

Required Aone receipts left `UNKNOWN` are recovered from the control-plane
`external-recovery-candidates` read model, never from local `orphanOperations`.
Only the scheduler-role bridge runs `ExternalOperationRecoveryScheduler`. It
leases one receipt with a stable per-worker/per-operation token, renews around
the Aone read, and checks comment digest, exact status, or the complete tag
postcondition. A definite found/not-found result is reconciled; unavailable or
ambiguous readback releases the lease without any Aone write.
