#!/usr/bin/env bash
# bootstrap/cloudspec-build-linux.sh — build the cloudspec CLI from source on Linux.
#
# WHY: upstream ships NO Linux artifacts — the acube download endpoint returns
# 500 「下载失败:<ver>」 for every cloudspec-linux-*.zip (all versions/arches)
# while mac zips download fine. The official Linux path per
# https://aliyuque.antfin.com/cloudspec/model/cli-install is: get read access
# to code.alibaba-inc.com/cloudspec/cloudspec-cli and build locally.
#
# This script automates that path on a worker:
#   clone (the worker's git credential store already carries a
#   code.alibaba-inc.com token with read access — verified for open_jarvis)
#   → uv venv → uv pip install -r requirements.txt → pyinstaller (same flags
#   as the repo's own build_linux.sh) → install to ~/.local/opt/cloudspec
#   (binary + _internal/ must stay adjacent — PyInstaller onedir) → symlink
#   ~/.local/bin/cloudspec.
#
# The produced binary's glibc baseline = this build host's glibc (AliOS 7.2 =
# 2.17), so a zip built on one worker runs on the whole fleet. The script
# always leaves ~/cloudspec-linux-amd64.zip behind and prints an OSS upload
# hint so the other workers can install from OSS instead of rebuilding.
#
# Env overrides:
#   CLOUDSPEC_CLI_REPO  git URL      (default: HTTPS repo on code.alibaba-inc.com)
#   CLOUDSPEC_PY        python ver   (default: 3.12)
#   UV_INDEX_URL        pip index    (default: aliyun mirror — workers are in-region)
#
# The whole pipeline (venv + deps + pyinstaller + smoke) was dry-run end-to-end
# on macOS before this script landed; the Linux run uses identical steps.
set -euo pipefail

REPO_URL="${CLOUDSPEC_CLI_REPO:-https://code.alibaba-inc.com/cloudspec/cloudspec-cli.git}"
PY_VER="${CLOUDSPEC_PY:-3.12}"
export UV_INDEX_URL="${UV_INDEX_URL:-https://mirrors.aliyun.com/pypi/simple}"

die()  { printf '[ERR] %s\n' "$*" >&2; exit 1; }
info() { printf '[..]  %s\n' "$*"; }
ok()   { printf '[OK]  %s\n' "$*"; }
step() { printf '\n===== %s =====\n' "$*"; }

command -v git >/dev/null 2>&1 || die "git not found"
command -v zip >/dev/null 2>&1 || die "zip not found (sudo yum install -y zip)"
# uv: self-install when absent (hosts whose yum python sufficed never needed
# the step-2 uv fallback, but the build always does).
if ! command -v uv >/dev/null 2>&1 && [ ! -x "$HOME/.local/bin/uv" ]; then
  info "uv missing; installing from astral.sh"
  curl -LsSf https://astral.sh/uv/install.sh | sh || die "uv install failed"
fi
command -v uv >/dev/null 2>&1 || export PATH="$HOME/.local/bin:$PATH"
command -v uv >/dev/null 2>&1 || die "uv still not on PATH after install"

t=$(mktemp -d -t cloudspec-build.XXXXXX)
# shellcheck disable=SC2064
trap "rm -rf '$t'" EXIT

step "1. Clone $REPO_URL"
git clone --depth 1 "$REPO_URL" "$t/src" \
  || die "clone failed — this token may lack read access; request it on the repo page (yuque cli-install doc)"
cd "$t/src"

step "2. venv (python $PY_VER) + requirements"
uv venv "$t/venv" --python "$PY_VER"
# shellcheck source=/dev/null
source "$t/venv/bin/activate"
uv pip install -q -r requirements.txt
version=$(python -c "from version.version_manager import VersionManager; print(VersionManager.version())")
ok "deps installed; source version = $version"

step "3. pyinstaller (flags from upstream build_linux.sh)"
pyinstaller --log-level ERROR \
  --hidden-import=queue --hidden-import=tqdm --hidden-import=tqdm.auto \
  --hidden-import=tqdm.std --hidden-import=tqdm.utils --hidden-import=tqdm.gui \
  --hidden-import=tqdm.contrib --hidden-import=urllib3 --hidden-import=certifi \
  --hidden-import=ssl --hidden-import=json --hidden-import=platform \
  --hidden-import=subprocess --add-data="config.json:." cloudspec.py
[ -x dist/cloudspec/cloudspec ] || die "pyinstaller produced no dist/cloudspec/cloudspec"

step "4. Install → ~/.local/opt/cloudspec (+ ~/.local/bin symlink)"
rm -rf "$HOME/.local/opt/cloudspec"
mkdir -p "$HOME/.local/opt" "$HOME/.local/bin"
cp -Rp dist/cloudspec "$HOME/.local/opt/cloudspec"
ln -sfn "$HOME/.local/opt/cloudspec/cloudspec" "$HOME/.local/bin/cloudspec"
"$HOME/.local/bin/cloudspec" version || die "installed binary failed smoke test"
ok "cloudspec $version installed"

step "5. Fleet zip"
zip_out="$HOME/cloudspec-linux-amd64.zip"
(cd dist/cloudspec && zip -qr "$zip_out" .)
sha=$(sha256sum "$zip_out" | awk '{print $1}')
ok "left behind: $zip_out ($(wc -c <"$zip_out" | tr -d ' ') bytes)"
cat <<EOF

Fleet distribution (run where ossutil lives, e.g. the scheduler host):
  ossutil cp "$zip_out" oss://cc-packet/cloudspec-linux-amd64-$version.zip
  sha256: $sha
Other workers then install with the deps.lock flow pointed at that OSS URL,
or simply: unzip to ~/.local/opt/cloudspec + symlink ~/.local/bin/cloudspec.
EOF
