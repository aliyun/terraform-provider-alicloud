#!/usr/bin/env python3
"""Local Cloud DNS group recording pilot; no remote service is required."""
import argparse
import errno
import hashlib
import json
import os
from pathlib import Path
import re
import secrets
import shutil
import signal
import socket
import subprocess
import sys
import tempfile
import time

REPO = Path(__file__).resolve().parent.parent
CASE = "TestRecordingAlidnsGroupLiveTerraformCLI"
PACKAGE = "./alicloud/service/alidns"
PROFILE = '''(version 1)
(allow default)
(deny network*)
(allow network-outbound (remote ip "localhost:*"))
(allow network-inbound (local ip "localhost:*"))
(allow network* (local unix-socket) (remote unix-socket))
'''


def write_json(path, value):
    path.write_text(json.dumps(value, indent=2) + "\n")


def base_env():
    # Never inherit proxy, debug, profile, Go execution flags or credential vars.
    env = {k: os.environ[k] for k in ("PATH", "HOME", "TMPDIR") if k in os.environ}
    env.update(ALIBABA_CLOUD_ACCOUNT_TYPE="Domestic", CHECKPOINT_DISABLE="1",
               ALICLOUD_REGION="cn-hangzhou", ALIBABA_CLOUD_REGION_ID="cn-hangzhou")
    return env


def record_env():
    if os.environ.get("ALICLOUD_ACCESS_KEY") or os.environ.get("ALICLOUD_SECRET_KEY"):
        ak, sk = os.environ.get("ALICLOUD_ACCESS_KEY"), os.environ.get("ALICLOUD_SECRET_KEY")
    else:
        ak, sk = os.environ.get("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.environ.get("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
    if not ak or not sk:
        raise RuntimeError("record requires ALICLOUD_ACCESS_KEY/ALICLOUD_SECRET_KEY (or ALIBABA_CLOUD equivalents)")
    if os.environ.get("ALICLOUD_SECURITY_TOKEN") or os.environ.get("ALIBABA_CLOUD_SECURITY_TOKEN"):
        raise RuntimeError("this pilot requires static test credentials; STS/profile credentials are not supported")
    return dict(base_env(), ALICLOUD_ACCESS_KEY=ak, ALICLOUD_SECRET_KEY=sk,
                ALIBABA_CLOUD_ACCESS_KEY_ID=ak, ALIBABA_CLOUD_ACCESS_KEY_SECRET=sk,
                ALIBABA_CLOUD_USER_AGENT="AlibabaCloud-Agent-Skills")


def run_log(args, env, log, timeout=600):
    with log.open("w") as out:
        process = subprocess.Popen(args, env=env, cwd=REPO, stdout=out,
                                   stderr=subprocess.STDOUT, start_new_session=True)
        try:
            code = process.wait(timeout=timeout)
        except BaseException:
            # Stop Terraform descendants before the independent cleanup query.
            try:
                os.killpg(process.pid, signal.SIGKILL)
            except ProcessLookupError:
                pass
            process.wait()
            raise
    if code:
        raise RuntimeError(f"command failed; private log: {log}")


def cloud(env, action, **params):
    args = ["aliyun", "alidns", action, "--region", "cn-hangzhou"]
    for key, value in params.items():
        args += ["--" + key, str(value)]
    result = subprocess.run(args, env=env, capture_output=True, timeout=45)
    if result.returncode:
        # CLI errors may contain request/credential details; do not echo them.
        raise RuntimeError(f"Cloud DNS {action} failed; cleanup may need manual verification")
    return json.loads(result.stdout)


def groups(env, name):
    found = []
    page = 1
    while True:
        data = cloud(env, "DescribeDomainGroups", KeyWord=name, PageSize=100, PageNumber=page)
        rows = data.get("DomainGroups", {}).get("DomainGroup", [])
        found.extend(g for g in rows if g["GroupName"] in (name, name + "-u"))
        if page * 100 >= int(data["TotalCount"]):
            return found
        if not rows or page >= 100:
            raise RuntimeError("incomplete group inventory; refusing to assume cleanup")
        page += 1


def cleanup(env, name, output):
    for group in groups(env, name):
        if int(group.get("DomainCount", -1)) != 0:
            raise RuntimeError("test group contains domains; refusing to delete it")
        cloud(env, "DeleteDomainGroup", GroupId=group["GroupId"])
    for _ in range(12):
        remaining = groups(env, name)
        if not remaining:
            write_json(output / "cleanup.json", {"remainingGroups": 0})
            return
        time.sleep(2)
    raise RuntimeError("test groups remain; inspect the recorded group name before retrying")


def baseline_hashes(output):
    return {name: hashlib.sha256((output / "cassette" / name).read_bytes()).hexdigest()
            for name in ("cassette.json", "seed")}


def result(output, mode):
    report = json.loads((output / (mode + ".json")).read_text())
    if (report.get("success") is not True or report.get("mode") != mode
            or report.get("testCase") != CASE or report.get("sanitized") is not True
            or report.get("interactionCount", 0) <= 0
            or report.get("unmatchedRequests") != 0 or report.get("unusedInteractions") != 0):
        raise RuntimeError(f"{mode} report did not pass")
    log = (output / (mode + ".log")).read_text()
    if f"--- PASS: {CASE} (" not in log:
        raise RuntimeError("selected test did not PASS (a skipped test is not sufficient)")
    return report


def compile_test(work, log):
    for command in ("go", "terraform"):
        if not shutil.which(command):
            raise RuntimeError(f"install {command} first")
    binary = work / "alidns.test"
    run_log(["go", "test", "-p=1", "-c", "-o", str(binary), PACKAGE],
            base_env(), log)
    return binary


def test_env(output, name, mode, env):
    return dict(env, TF_ACC="1", TF_ACC_LIVE_RECORDING="1", TF_RECORDING_GROUP_NAME=name,
                TF_ACC_RECORDING_CASE=CASE, TF_ACC_RECORDING_MODE=mode,
                TF_ACC_RECORDING_PATH=str(output / "cassette"),
                TF_ACC_RECORDING_RESULT_PATH=str(output / (mode + ".json")))


def test_args(binary):
    return [str(binary), "-test.run=^" + CASE + "$", "-test.v", "-test.timeout=9m"]


def record(output):
    env = record_env()
    if not shutil.which("aliyun"):
        raise RuntimeError("install aliyun CLI first (used for independent cleanup)")
    output.mkdir(parents=True, mode=0o700, exist_ok=False)
    name = "tf-rr-" + secrets.token_hex(5)
    manifest = {"version": 1, "testCase": CASE, "groupName": name, "state": "recording"}
    write_json(output / "manifest.json", manifest)
    print(f"Recording empty DNS group; private output: {output}", flush=True)
    with tempfile.TemporaryDirectory(prefix="tf-recording-build-") as temporary:
        binary = compile_test(Path(temporary), output / "build-record.log")
        if groups(env, name):
            raise RuntimeError("generated group name already exists; refusing to run")
        try:
            run_log(test_args(binary), test_env(output, name, "record", env), output / "record.log")
        finally:
            cleanup(env, name, output)
    report = result(output, "record")
    blob = (output / "cassette" / "cassette.json").read_bytes()
    for key, value in os.environ.items():
        if any(part in key for part in ("ACCESS_KEY", "SECRET_KEY", "SECURITY_TOKEN")) and len(value) > 8:
            if value.encode() in blob:
                raise RuntimeError("credential scan failed; recording is not ready to share")
    manifest.update(state="ready", hashes=baseline_hashes(output), interactionCount=report["interactionCount"],
                    sourceCommit=subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=REPO, text=True).strip())
    write_json(output / "manifest.json", manifest)
    print(f"PASS: recorded {report['interactionCount']} interactions; remaining test groups: 0")
    print(f"Replay: python3 scripts/recording.py replay {str(output)!r}")


def load_manifest(output):
    manifest = json.loads((output / "manifest.json").read_text())
    if (manifest.get("version") != 1 or manifest.get("state") != "ready"
            or manifest.get("testCase") != CASE
            or not re.fullmatch(r"tf-rr-[0-9a-f]{10}", manifest.get("groupName", ""))):
        raise RuntimeError("unsupported or incomplete recording manifest")
    if manifest.get("hashes") != baseline_hashes(output):
        raise RuntimeError("recording files changed; refusing replay")
    return manifest


def replay(output):
    # Fail before executing any test on systems without the verified isolation.
    if sys.platform != "darwin" or not shutil.which("sandbox-exec"):
        raise RuntimeError("automatic isolated replay currently supports macOS only; no unisolated fallback")
    manifest = load_manifest(output)
    with tempfile.TemporaryDirectory(prefix="tf-replay-") as temporary:
        work = Path(temporary)
        binary = compile_test(work, output / "build-replay.log")
        home = work / "home"
        home.mkdir()
        env = test_env(output, manifest["groupName"], "replay", dict(
            base_env(), HOME=str(home), ALICLOUD_ACCESS_KEY="recording-access-key",
            ALICLOUD_SECRET_KEY="recording-secret-key",
            ALIBABA_CLOUD_ACCESS_KEY_ID="recording-access-key",
            ALIBABA_CLOUD_ACCESS_KEY_SECRET="recording-secret-key"))
        profile = work / "replay.sb"
        profile.write_text(PROFILE)
        # Probe and exec share the same sandbox and credential-free environment.
        run_log(["sandbox-exec", "-f", str(profile), sys.executable,
                 str(Path(__file__).resolve()), "_isolated"] + test_args(binary),
                env, output / "replay.log")
    report = result(output, "replay")
    if baseline_hashes(output) != manifest["hashes"]:
        raise RuntimeError("replay changed the baseline")
    print(f"PASS: replayed {report['interactionCount']} interactions; external network denied; baseline unchanged")


def isolated(args):
    try:
        connection = socket.create_connection(("1.1.1.1", 443), timeout=3)
    except OSError as error:
        if error.errno not in (errno.EPERM, errno.EACCES):
            raise RuntimeError("network probe must fail with permission denied") from None
    else:
        connection.close()
        raise RuntimeError("external networking is available; refusing replay")
    print("External TCP probe: permission denied", flush=True)
    os.execve(args[0], args, os.environ)


def terminate(signum, frame):
    raise KeyboardInterrupt


def main():
    signal.signal(signal.SIGTERM, terminate)
    os.umask(0o077)
    if len(sys.argv) > 1 and sys.argv[1] == "_isolated":
        isolated(sys.argv[2:])
        return
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("mode", choices=("record", "replay"))
    parser.add_argument("directory", nargs="?", help="private recording directory")
    args = parser.parse_args()
    if args.mode == "replay" and not args.directory:
        parser.error("replay requires a recording directory")
    output = Path(args.directory).expanduser().resolve() if args.directory else REPO / ".terraform-recordings" / secrets.token_hex(6)
    try:
        {"record": record, "replay": replay}[args.mode](output)
    except (RuntimeError, OSError, ValueError, subprocess.SubprocessError) as error:
        parser.exit(1, f"ERROR: {error}\n")


if __name__ == "__main__":
    main()
