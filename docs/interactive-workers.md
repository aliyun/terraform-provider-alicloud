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
`SessionStart` and `SessionEnd` hooks. Codex currently has no project `SessionEnd`
event, so its sidecar watches the real Codex host PID and marks the worker offline
when that process exits.
