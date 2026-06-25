#!/usr/bin/env python3
"""
DingTalk AI Card streaming sender.

Usage:
    python3 streaming.py --to <staffId> -m "message"
    python3 streaming.py --to <staffId> --stdin
    python3 streaming.py --to-group <openConversationId> -m "message"
"""

import argparse
import json
import os
import ssl
import subprocess
import sys
import time
import uuid
import urllib.request
import urllib.error
from pathlib import Path


# ---------------------------------------------------------------------------
# SSL context: auto-detect a usable CA bundle so the script works out-of-the-
# box on macOS Homebrew Python (which ships without bundled root certs).
# Priority: certifi > /etc/ssl/cert.pem (macOS) > ca-certificates.crt (Linux)
# The SSL_CERT_FILE env var is still honoured by ssl.create_default_context().
# ---------------------------------------------------------------------------
def _build_ssl_context():
    if os.environ.get("SSL_CERT_FILE"):
        return ssl.create_default_context()          # env var takes precedence
    try:
        import certifi
        return ssl.create_default_context(cafile=certifi.where())
    except ImportError:
        pass
    for candidate in ("/etc/ssl/cert.pem", "/etc/ssl/certs/ca-certificates.crt"):
        if os.path.isfile(candidate):
            return ssl.create_default_context(cafile=candidate)
    return ssl.create_default_context()

_SSL_CTX = _build_ssl_context()


def _urlopen(req, **kwargs):
    """urllib.request.urlopen wrapper that injects the auto-detected SSL ctx."""
    return urllib.request.urlopen(req, context=_SSL_CTX, **kwargs)


def ensure_bind():
    """Best-effort idempotent am bind from env vars; never fatal."""
    helper = Path(__file__).resolve().parent / "ensure-bind.sh"
    if not helper.exists():
        return
    try:
        subprocess.run(["bash", str(helper)], check=False,
                       stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except Exception:
        pass


def load_am_config(bot_name=None):
    if bot_name:
        config_path = Path.home() / ".config" / "aone-message-cli" / "profiles" / bot_name / "config.properties"
    else:
        config_path = Path.home() / ".config" / "aone-message-cli" / "config.properties"
    config = {}
    if config_path.exists():
        for line in config_path.read_text().splitlines():
            line = line.strip()
            if line and not line.startswith("#") and "=" in line:
                k, v = line.split("=", 1)
                config[k.strip()] = v.strip()
    return config


def get_access_token(app_key, app_secret):
    data = json.dumps({"appKey": app_key, "appSecret": app_secret}).encode()
    req = urllib.request.Request(
        "https://api.dingtalk.com/v1.0/oauth2/accessToken",
        data=data,
        headers={"Content-Type": "application/json"},
    )
    resp = _urlopen(req)
    result = json.loads(resp.read())
    return result["accessToken"]


def create_and_deliver_card(token, template_id, robot_code, target, target_type="user"):
    out_track_id = str(uuid.uuid4())
    body = {
        "cardTemplateId": template_id,
        "outTrackId": out_track_id,
        "cardData": {"cardParamMap": {}},
        "callbackType": "STREAM",
        "userIdType": 1,
    }

    if target_type == "user":
        body["openSpaceId"] = f"dtv1.card//IM_ROBOT.{target}"
        body["imRobotOpenDeliverModel"] = {"spaceType": "IM_ROBOT", "robotCode": robot_code}
        body["imRobotOpenSpaceModel"] = {"supportForward": True}
    elif target_type == "group":
        body["openSpaceId"] = f"dtv1.card//IM_GROUP.{target}"
        body["imGroupOpenDeliverModel"] = {"robotCode": robot_code}
        body["imGroupOpenSpaceModel"] = {"supportForward": True}

    data = json.dumps(body).encode()
    req = urllib.request.Request(
        "https://api.dingtalk.com/v1.0/card/instances/createAndDeliver",
        data=data,
        headers={
            "Content-Type": "application/json",
            "x-acs-dingtalk-access-token": token,
        },
    )
    resp = _urlopen(req)
    result = json.loads(resp.read())

    if result.get("success"):
        for dr in result.get("result", {}).get("deliverResults", []):
            if dr.get("success"):
                return out_track_id
            raise RuntimeError(f"Card delivery failed: {dr.get('errorMsg')}")
    raise RuntimeError(f"Card creation failed: {result}")


def streaming_update(token, out_track_id, key, content, is_full=True, is_finalize=False, is_error=False):
    body = {
        "outTrackId": out_track_id,
        "guid": str(uuid.uuid4()),
        "key": key,
        "content": content,
        "isFull": is_full,
        "isFinalize": is_finalize,
        "isError": is_error,
    }
    data = json.dumps(body).encode()
    req = urllib.request.Request(
        "https://api.dingtalk.com/v1.0/card/streaming",
        data=data,
        method="PUT",
        headers={
            "Content-Type": "application/json",
            "x-acs-dingtalk-access-token": token,
        },
    )
    resp = _urlopen(req)
    return json.loads(resp.read())


def stream_text(token, out_track_id, text, chunk_size=2, delay=0.15, key="content"):
    chars = list(text)
    accumulated = ""
    for i in range(0, len(chars), chunk_size):
        accumulated += "".join(chars[i : i + chunk_size])
        is_last = (i + chunk_size) >= len(chars)
        streaming_update(token, out_track_id, key, accumulated, is_finalize=is_last)
        if not is_last:
            time.sleep(delay)


def send_full(token, out_track_id, text, key="content"):
    streaming_update(token, out_track_id, key, text, is_finalize=True)


def main():
    parser = argparse.ArgumentParser(description="DingTalk AI Card streaming sender")
    parser.add_argument("--to", help="Target user staffId")
    parser.add_argument("--to-group", help="Target group openConversationId")
    parser.add_argument("--message", "-m", help="Message content")
    parser.add_argument("--stdin", action="store_true", help="Read message from stdin")
    parser.add_argument("--chunk-size", type=int, default=2, help="Chars per streaming update (default: 2)")
    parser.add_argument("--delay", type=float, default=0.15, help="Delay between updates in seconds (default: 0.15)")
    parser.add_argument("--no-stream", action="store_true", help="Send as one shot without streaming effect")
    parser.add_argument("--key", default="content", help="Card template streaming variable name (default: content)")
    parser.add_argument("--template-id", help="AI card template ID")
    parser.add_argument("--app-key", help="DingTalk app key")
    parser.add_argument("--app-secret", help="DingTalk app secret")
    parser.add_argument("--robot-code", help="Robot code")
    parser.add_argument("--bot", help="am bot profile name")
    args = parser.parse_args()

    ensure_bind()
    am_config = load_am_config(args.bot)
    app_key = args.app_key or os.environ.get("DINGTALK_APP_KEY") or am_config.get("aliding.access-key-id")
    app_secret = args.app_secret or os.environ.get("DINGTALK_APP_SECRET") or am_config.get("aliding.access-key-secret")
    robot_code = (args.robot_code
                  or os.environ.get("DINGTALK_ROBOT_CODE")
                  or am_config.get("aliding.robot-code")
                  or app_key)          # robotCode defaults to appKey
    template_id = args.template_id or os.environ.get("DINGTALK_TEMPLATE_ID")

    if not template_id:
        print("Error: --template-id or DINGTALK_TEMPLATE_ID required.", file=sys.stderr)
        sys.exit(1)
    if not app_key or not app_secret:
        print("Error: credentials required. Use --app-key/--app-secret, env vars, or `am bind`.", file=sys.stderr)
        sys.exit(1)
    if not args.to and not args.to_group:
        print("Error: --to or --to-group required.", file=sys.stderr)
        sys.exit(1)

    message = args.message
    if args.stdin:
        message = sys.stdin.read()
    if not message:
        print("Error: --message or --stdin required.", file=sys.stderr)
        sys.exit(1)

    token = get_access_token(app_key, app_secret)
    target = args.to or args.to_group
    target_type = "user" if args.to else "group"

    out_track_id = create_and_deliver_card(token, template_id, robot_code, target, target_type)
    print(json.dumps({"outTrackId": out_track_id, "status": "created"}))

    if args.no_stream:
        send_full(token, out_track_id, message, key=args.key)
    else:
        stream_text(token, out_track_id, message, chunk_size=args.chunk_size, delay=args.delay, key=args.key)

    print(json.dumps({"outTrackId": out_track_id, "status": "done"}))


if __name__ == "__main__":
    main()
