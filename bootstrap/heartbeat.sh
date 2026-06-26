#!/usr/bin/env bash
# bootstrap/heartbeat.sh — heartbeat sidecar
#
# Usage: heartbeat.sh <instance_id> <follow_pid>
#
# Sends a heartbeat for <instance_id> every HB_INT seconds (default 60).
# Exits automatically when <follow_pid> is no longer alive.
#
# Environment:
#   HB_INT   — heartbeat interval in seconds (default 60)
d="$(cd "$(dirname "$0")" && pwd)"; while kill -0 "$2" 2>/dev/null; do bash "$d/coord.sh" heartbeat "$1"; sleep "${HB_INT:-60}"; done
