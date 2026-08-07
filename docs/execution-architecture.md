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
bridge restart Scheduler, Bot, and Worker all quiesce before sharing one
30-second shutdown deadline. The executor reports `DRAINING`, stops local task
processes, relinquishes Sessions as `RESUMABLE` without increasing Task retry
count, and reports `OFFLINE` when handoff completes. The next process
re-registers the same Worker and resumes the same runtime session on that
machine. If the bounded handoff cannot reach the control plane, shutdown logs a
failed handoff and lease expiry/recovery reconciles the existing
`runtimeSessionId` and fence; replacement never bypasses that ownership. A
crash cannot be taken over until the old heartbeat/leases are stale and no
Session remains `LEASED` or `RUNNING`.

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

## RECOVERY_REQUIRED and the desired-state upsert boundary

Discovery and execution are deliberately separate. When the scan runner's
`_decide` selects `dispatch force=True` for a `jarvis-idle` ticket with a new
human comment, `_dispatch` still goes through `execution_router.enqueue`
(`bridge/jarvis_task_router.py:392-398`) → `_persist_task`
(`bridge/jarvis_task_router.py:354-376`) → `client.upsert_desired_task`. That
call advances the Task's **desiredRevision** (a `TASK_UPSERTED` event).
`force=True` is scoped to the EphemeralExecutor local dedup ledger
(`scan.py:657-661`, gated by `not execution_router.is_task(envelope)`); on the
control-plane Task path it does not itself flip an execution state.

**Whether the upsert moves the execution status depends on the status it finds**
— an earlier revision of this section claimed it "never" does, which is true
only for the states below and is exactly why the zombie loop went unnoticed for
weeks:

| Task status when the upsert lands | Effect |
|---|---|
| `RUNNING` / `LEASED` | `TASK_UPSERTED RUNNING→RUNNING`. Inert for status. The running Session finishes as `SESSION_COMPLETED RUNNING→READY` instead of `→SUCCEEDED` when desired advanced mid-run, so the newer revision still gets executed. |
| `RECOVERY_REQUIRED` / `RETRY_WAIT` | Inert (`RECOVERY_REQUIRED→RECOVERY_REQUIRED`). Scan cannot self-heal these; see the Task 471 case below. |
| `READY` | `READY→READY`. Already queued. |
| **resting** — `SUCCEEDED` / `FAILED` / `CANCELED` | **Re-arms the Task to `READY`.** A Worker then leases it and runs a full generation. |

That last row is a feature when the revision genuinely advanced (it is how a
human comment on a finished ticket gets answered) and a defect when it did not.
The control plane's idempotency cannot tell the two apart, because
`TaskEnvelope.request_id` hashes the whole envelope rather than the revision:
`sourceStatus` moving as the ticket advances, an edited title or a reworded
prompt each mint a fresh key for byte-identical desired state.

So discovery must not re-publish a revision the control plane already holds.
`ScanRunner._desired_revision_already_published` compares
`envelope.desired_revision` against the live Task's `desiredRevision` before
`enqueue` and returns `noop_desired_revision_unchanged` on a match — accepted
(nothing left to do, so the decision comes to rest and is not replayed) but not
counted as a dispatch. An unreadable Task falls through to the upsert: one
duplicate run costs minutes, a suppressed dispatch can lose a human comment
permanently. Re-arming a genuinely parked Task stays the `recovery` runner's
job, not discovery's.

Observed case (Task 665, aone 1086837/84202521, 2026-07-19~08-07): `generation`
reached 12 with `desiredRevision` frozen at
`modified:2026-07-29-18:22|policy:v6|input:15fc9fc6130092f8`. The timeline holds
seven `TASK_UPSERTED SUCCEEDED→READY` re-arms plus one `CANCELED→READY` — the
last of these on 08-07T09:38Z, *after* `TASK_SOURCE_TERMINAL_CONVERGED` had
already cancelled the Task on a terminal source status. Each re-arm consumed a
lease slot for 20-40 minutes to re-reach the same conclusion. Five such zombies
(`gen` 6/8/12/8 and one real run) filled all five slots of the only online
worker, and new tickets queued behind `ready → NO_FREE_SLOTS`.

What actually controls leaseability:

- `lease_task` (`bridge/jarvis_persistence_executor.py:1104-1161`,
  `client.lease_task` at `bridge/jarvis_task_client.py:407`) only pulls Tasks
  in `READY`. Tasks in `RECOVERY_REQUIRED`, `RETRY_WAIT` or `SUSPENDED` are
  never returned, regardless of how many times scan re-upserts the desired
  revision — unlike a resting Task, which an upsert first moves to `READY` and
  which is leasable from then on (see the table above). There is no `lease`
  INFO line for the
  skipped Task because the server simply returns the next matching READY Task
  (or none); only `task lease conflict` / `task lease failed` are logged.
- Once `retryCount` exceeds `maxRetries` after a `SESSION_FAILED`, the Task
  parks in `RECOVERY_REQUIRED` and every later `TASK_UPSERTED` becomes a
  `RECOVERY_REQUIRED→RECOVERY_REQUIRED` no-op. Scan cannot self-heal it;
  recovery requires the `recovery` runner, or an explicit `force-release` /
  `force-redispatch` via `bootstrap/control-plane-status.sh`.
- `field_repair` is a pre-execution gate inside `PersistentTaskExecution.execute`
  (`bridge/persistent_tasks.py:1079-1086`, `field_repair_worker.repair_only`).
  A `field_repair_transient` / `apply_timeout` here fails the Session **before**
  `dispatch_item` is reached — so a bot.log without `dispatch_item #<id> start`
  does not mean "never leased"; it means execution stopped at the field-repair
  gate. Always corroborate with `bootstrap/control-plane-status.sh task
  <aone_id>` (Task status / generation / retry / lastError / session
  RESUMABLE / event timeline) before concluding "Worker did not lease".

Observed case (Task 471, aone 528766/84103836, 2026-07-30~31): a human comment
on a `jarvis-idle` ticket caused scan to upsert `desired=comment:<id>` three
times; all events were `TASK_UPSERTED RECOVERY_REQUIRED→RECOVERY_REQUIRED`,
because the Worker had already leased the Task, hit
`field_repair_transient / apply_timeout` at the repair gate
(`SESSION_FAILED RUNNING→RECOVERY_REQUIRED`), and parked at `retry=4/3`
(over `maxRetries=3`). The human comment went unanswered until an operator
manually closed the Aone status. `force=True` dispatching was inert
throughout, and the `lease=- worker=-` RESUMABLE session was never re-leased
because the Task never returned to READY.

Separately, `_point_read_source_status` (`bridge/scheduler/runners/scan.py:146
-182`) runs a full `a1id project workitem get <id>` subprocess to reconcile the
control-plane `sourceStatus` display field, with
`SOURCE_STATUS_POINT_TIMEOUT_SECONDS` defaulting to 30s. Heavy tickets (many
comments / relations / long description) exceed this steadily for hours,
leaving `sourceStatus` stale and emitting `point-read #<id> failed:
TimeoutExpired` noise. This is display-only and does not affect dispatch or
lease; the same `a1id` slowness is the likely root of the field-repair
`apply_timeout` above. Raising the timeout or point-reading a lightweight
status-only endpoint clears the noise and may reduce repair-gate timeouts.

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
