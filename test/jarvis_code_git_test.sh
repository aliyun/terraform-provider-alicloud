#!/usr/bin/env bash
set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
helper="$repo_root/.agents/skills/cloudspec-amp-workflow/scripts/jarvis-code-git.sh"
tmp_dir="$(mktemp -d)"
auth_file="$tmp_dir/auth.yaml"
git_log="$tmp_dir/git.log"
fake_git="$tmp_dir/git"
fake_repo="$tmp_dir/repo"
out="$tmp_dir/out"
err="$tmp_dir/err"
pass=0
fail=0

cleanup() { rm -rf "$tmp_dir"; }
trap cleanup EXIT

mkdir -p "$fake_repo/.git"
cat >"$auth_file" <<'YAML'
version: 1
platforms:
  code:
    host: code.alibaba-inc.com
    auth_type: private_token
    token: fake-private-token
    user: open_jarvis
YAML
chmod 600 "$auth_file"

cat >"$fake_git" <<'EOF'
#!/usr/bin/env bash
set -uo pipefail
{
  printf 'ARGS=%s\n' "$*"
  printf 'PROMPT=%s\n' "${GIT_TERMINAL_PROMPT:-}"
  printf 'USER=%s\n' "$("$GIT_ASKPASS" 'Username for HTTPS')"
  printf 'PASS=%s\n' "$("$GIT_ASKPASS" 'Password for HTTPS')"
} >>"$JARVIS_CODE_TEST_LOG"
exit "${FAKE_GIT_STATUS:-0}"
EOF
chmod 700 "$fake_git"

run_helper() {
  : >"$out"
  : >"$err"
  : >"$git_log"
  JARVIS_CODE_AUTH_FILE="$auth_file" \
  JARVIS_CODE_GIT_BIN="$fake_git" \
  JARVIS_CODE_TEST_LOG="$git_log" \
    bash "$helper" "$@" >"$out" 2>"$err"
  status=$?
}

ok() { printf 'PASS %s\n' "$1"; pass=$((pass + 1)); }
bad() { printf 'FAIL %s\n' "$1" >&2; fail=$((fail + 1)); }
contains() { grep -Fq "$2" "$1"; }

run_helper check
[ "$status" -eq 0 ] && contains "$out" 'auth=private_token' \
  && ! contains "$out" 'fake-private-token' && ok 'check validates without leaking token' \
  || bad 'check validates without leaking token'

clone_dest="$tmp_dir/clone-target"
run_helper clone cloudspec-model/VPC_pop_Vpc_2016-04-28 \
  feature/85289668-tunnel-bandwidth "$clone_dest"
[ "$status" -eq 0 ] \
  && contains "$git_log" 'ARGS=clone --branch feature/85289668-tunnel-bandwidth --single-branch https://code.alibaba-inc.com/cloudspec-model/VPC_pop_Vpc_2016-04-28.git' \
  && contains "$git_log" 'PROMPT=0' \
  && contains "$git_log" 'USER=open_jarvis' \
  && contains "$git_log" 'PASS=fake-private-token' \
  && ! contains "$out" 'fake-private-token' \
  && ! contains "$err" 'fake-private-token' \
  && ok 'clone uses clean HTTPS URL and askpass' \
  || bad 'clone uses clean HTTPS URL and askpass'

run_helper clone cloudspec-model/VPC_pop_Vpc_2016-04-28 master "$tmp_dir/master"
[ "$status" -ne 0 ] && contains "$err" 'only feature/* branches are allowed' \
  && [ ! -s "$git_log" ] && ok 'clone rejects master' \
  || bad 'clone rejects master'

run_helper push "$fake_repo" cloudspec-model/VPC_pop_Vpc_2016-04-28 HEAD main
[ "$status" -ne 0 ] && contains "$err" 'only feature/* branches are allowed' \
  && [ ! -s "$git_log" ] && ok 'push rejects master' \
  || bad 'push rejects master'

run_helper push "$fake_repo" cloudspec-model/VPC_pop_Vpc_2016-04-28 HEAD \
  refs/heads/feature/85289668-tunnel-bandwidth
expected_push="ARGS=-C $fake_repo push https://code.alibaba-inc.com/cloudspec-model/VPC_pop_Vpc_2016-04-28.git HEAD:refs/heads/feature/85289668-tunnel-bandwidth"
[ "$status" -eq 0 ] \
  && contains "$git_log" "$expected_push" \
  && ok 'push uses private-token HTTPS transport' \
  || bad 'push uses private-token HTTPS transport'

chmod 644 "$auth_file"
run_helper check
[ "$status" -ne 0 ] && contains "$err" 'must not be group/world accessible' \
  && ok 'insecure auth permissions fail closed' \
  || bad 'insecure auth permissions fail closed'

# AMP config projection is an exact internal ASKPASS prompt, not a public command.
chmod 600 "$auth_file"
run_helper amp-config
[ "$status" -ne 0 ] && [ ! -s "$git_log" ] \
  && ! contains "$out" 'fake-private-token' && ok 'amp-config is not a public command' \
  || bad 'amp-config is not a public command'

fake_home="$tmp_dir/home"
mkdir -p "$fake_home/.amp"
# Keep installed Python packages available without exposing the real auth tree.
python3 - "$fake_home" <<'PY'
import site
import sys
from pathlib import Path
source = Path(site.getusersitepackages())
site_path = Path(sys.argv[1]) / source.relative_to(Path.home())
site_path.parent.mkdir(parents=True, exist_ok=True)
site_path.symlink_to(source, target_is_directory=True)
PY
amp_config="$fake_home/.amp/config.yaml"
cat >"$amp_config" <<'YAML'
profiles:
  default:
    auth:
      type: buc
      token_file: fake-private-token
    credentials:
      access_key_secret: fake-private-token
YAML
HOME="$fake_home" JARVIS_CODE_ASKPASS_MODE=1 run_helper AMPConfig
[ "$status" -eq 0 ] && [ ! -s "$err" ] && [ ! -s "$git_log" ] \
  && python3 - "$out" <<'PY' && ok 'AMP config replaces BUC and supplies required safe defaults' \
  || bad 'AMP config replaces BUC and supplies required safe defaults'
import json
import sys
from pathlib import Path
assert json.loads(Path(sys.argv[1]).read_text()) == {
    "current_profile": "default", "profiles": {"default": {"auth": {"type": "private_token"}}},
    "output": {"default": "json"},
    "http": {"timeout_seconds": 30, "debug": False}, "upgrade": {}}
PY

cat >"$amp_config" <<'YAML'
current_profile: custom
profiles:
  default:
    token: fake-private-token
  custom:
    endpoint: https://custom-amp.example.invalid
    openapi_version: 2026-04-20
    auth: {type: buc, token_file: fake-private-token}
    credentials:
      source: oauth
      oauth_profile: existing-signing-profile
      oauth_site: https://oauth.example.invalid
      token: fake-private-token
      access_key_id: fake-private-token
      access_key_secret: fake-private-token
      secret_file: /nonexistent-signing-secret
output: {default: yaml}
http:
  timeout_seconds: 73
  debug: true
  headers: {token: fake-private-token}
upgrade: {token: fake-private-token}
tokens: {private_token: fake-private-token}
YAML
HOME="$fake_home" AMP_HOME="$tmp_dir/ignored" AMP_PROFILE=default \
  JARVIS_CODE_ASKPASS_MODE=1 run_helper AMPConfig
[ "$status" -eq 0 ] && [ ! -s "$err" ] && [ ! -s "$git_log" ] \
  && python3 - "$out" <<'PY' && ok 'AMP config preserves routing and timeout with fixed safe settings' \
  || bad 'AMP config preserves routing and timeout with fixed safe settings'
import json
import sys
from pathlib import Path
assert json.loads(Path(sys.argv[1]).read_text()) == {
    "current_profile": "custom", "profiles": {"custom": {
        "endpoint": "https://custom-amp.example.invalid", "openapi_version": "2026-04-20",
        "auth": {"type": "private_token"},
        "credentials": {"source": "oauth", "oauth_profile": "existing-signing-profile",
                        "oauth_site": "https://oauth.example.invalid"}}},
    "output": {"default": "json"},
    "http": {"timeout_seconds": 73, "debug": False}, "upgrade": {}}
PY

HOME="$fake_home" JARVIS_CODE_ASKPASS_MODE=1 run_helper 'AMPConfig extra'
[ "$status" -eq 0 ] && ! contains "$out" 'profiles' && [ ! -s "$err" ] \
  && ok 'AMP config prompt must match exactly' || bad 'AMP config prompt must match exactly'

printf 'profiles: [\nfake-private-token:\n' >"$amp_config"
HOME="$fake_home" JARVIS_CODE_ASKPASS_MODE=1 run_helper AMPConfig
[ "$status" -ne 0 ] && [ ! -s "$out" ] \
  && [ "$(<"$err")" = 'AMP config is unavailable or invalid' ] \
  && ok 'invalid AMP YAML fails closed without parser secrets' \
  || bad 'invalid AMP YAML fails closed without parser secrets'

rm "$amp_config"
HOME="$fake_home" JARVIS_CODE_ASKPASS_MODE=1 run_helper AMPConfig
[ "$status" -ne 0 ] && [ ! -s "$out" ] \
  && [ "$(<"$err")" = 'AMP config is unavailable or invalid' ] \
  && ok 'missing AMP config fails closed' || bad 'missing AMP config fails closed'

printf '%s passed, %s failed\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
