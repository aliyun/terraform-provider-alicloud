# Interactive Codex / Claude workers

Opening this repository in Codex or Claude registers that native conversation as
a one-slot `INTERACTIVE` Jarvis worker. The detached sidecar only heartbeats the
worker and its explicitly claimed session; it never leases the shared queue.
`bootstrap/claim.sh claim <aone-id> <project-id>` performs the database direct
claim before writing the Aone tag, so a conflict or unavailable control plane
fails closed without touching Aone. `release` suspends the persisted session and
`finish` completes it (or suspends it when finish is downgraded to idle).

The hook loads `bootstrap/.env` and `bridge/jarvis.env` from the main checkout.
`JARVIS_CONTROL_PLANE_TOKEN` falls back to `JARVIS_HTML_REPORT_TOKEN`; credentials
are never written to the local worker state. Use
`bash bootstrap/run-interactive-worker-hook.sh cli status` to query the worker's
persisted control-plane state for the current native conversation.

The sidecar starts only when the hook can verify the real Codex/Claude host
process. If process ancestry cannot be verified, the hook persists an `OFFLINE`
worker record but refuses claims. This is intentionally fail-closed: an arbitrary
terminal or desktop ancestor must never keep a task lease alive after the agent
has exited. The sidecar also verifies the host process start identity on every
heartbeat, so PID reuse cannot revive an old Worker incarnation.

## Codex project trust

Codex only runs repository hooks after the project has been trusted. On the first
open, approve trust for this checkout and start (or resume) the task again so the
`SessionStart` hook can register the worker. If trust was declined, `claim.sh`
intentionally fails closed with “interactive worker is not registered”; do not
bypass it by updating the Aone tag directly. Claude runs the matching
`SessionStart` and `SessionEnd` hooks. Codex Desktop reuses a global app-server
process and currently has no project `SessionEnd` event, so process liveness alone
cannot delimit one conversation. Its `UserPromptSubmit` and `Stop` hooks therefore
bracket each turn: the Worker keeps heartbeating while Codex is available, but the
current task session is explicitly suspended after
`JARVIS_INTERACTIVE_TURN_GRACE_SEC` (default 600 seconds), without consuming a
crash retry. The old `JARVIS_INTERACTIVE_STOP_GRACE_SEC` name remains a temporary
compatibility fallback. If Codex loses the `Stop` event entirely, an active-turn
inactivity limit (`JARVIS_INTERACTIVE_ACTIVE_TURN_TTL_SEC`, default 12 hours) is
refreshed by tool hooks and prevents the global app-server from renewing an
abandoned session forever without cutting off a turn that is still making
progress. If the suspend request cannot reach the control plane, the default
300-second lease remains the crash fallback. The detached sidecar heartbeats every
30 seconds by default, so an available control plane continuously renews the
session while the native conversation remains active. A prompt arriving before
suspension synchronously renews the same fence. A later continuation prompt
reclaims the suspended task through the normal control-plane claim/start path,
reuses the same durable session, and receives a new fence without consuming a
crash retry. A prompt that
explicitly names a different Aone item leaves the old task suspended so the new
item can follow its own `claim.sh` path. If another Worker already recovered the
task, the stale assignment is recorded as lost ownership and every later prompt
continues to warn/fail closed until a standard claim succeeds; it is never
silently forgotten. A Codex app-server restart similarly preserves only the
task/runtime lineage, never the dead process's fence or in-flight receipt, and
retries targeted recovery after the old lease has been resolved.
Codex sidecars with no current task exit and mark their Worker `OFFLINE` after
`JARVIS_INTERACTIVE_IDLE_TTL_SEC` (default eight hours); the next prompt registers
the same Worker again and restarts the sidecar.

Changing `.codex/hooks.json` changes Codex's project-trust hash. Re-trust the
checkout after pulling a hook update, then start or resume the task once so all
lifecycle hooks are active. If the first `SessionStart` registration fails, the
hook persists a fail-closed tombstone; a missing local Worker state is also
fail-closed. Tools remain blocked until a later trusted `SessionStart` registers
the conversation successfully.

While a task is attached, every Codex and Claude `PreToolUse` synchronously
renews the exact session fence before the tool may start. A stale fence,
control-plane timeout/conflict, delayed old turn, replaced host process, or
uncertain operation receipt blocks the tool with exit code 2. The only recovery
exception is the standalone command `/bin/bash
<absolute-current-repo>/bootstrap/claim.sh claim <same-aone> <same-project>`;
composed shell commands, environment prefixes, relative/lookalike scripts,
custom tools and a different target remain blocked. This closes the interval
between a sidecar discovering lost ownership and an old agent attempting another
code or Aone mutation.

Codex subagents share the root Worker's database session, but never inherit it
from `agent_id` alone. The root `PreToolUse` records an exact `spawn_agent`
receipt, and `SubagentStart` verifies that call against Codex's durable parent
transcript before binding the child to the root turn and current task/fence
epoch. A new root prompt, task claim, release, recovered fence, or replacement
process invalidates the old child. A successful explicit `followup_task` or
`send_message` may re-authorize only the registered target for the new epoch;
delayed children and children from an earlier Aone therefore cannot execute with
the current task's fence. Nested subagents follow the same parent-chain binding.

Claude launched through the local `idea` alias registers one Worker for the root
Claude session. Its `Agent` subagents keep the same Claude session identity and
run through the root Worker's `PreToolUse` fence; they do not create independent
Workers or assignments. Nested client launches are also isolated: if `idea` is
started from a Codex terminal (or Codex from Claude), the nearest real host
process selects the runtime context, so the inner client cannot accidentally use
the outer client's Worker or an already-ACKed claim receipt.
