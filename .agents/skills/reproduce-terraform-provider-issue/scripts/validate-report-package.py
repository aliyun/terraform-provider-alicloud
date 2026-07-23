#!/usr/bin/env python3
"""Validate a Terraform reproduction report package without creating resources."""

from __future__ import annotations

import argparse
import html
import json
import os
import re
import subprocess
import tempfile
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import NoReturn


EXIT_INVALID = 2
EXIT_BLOCKED = 3
REQUIRED_FILES = (
    "REPORT.md",
    "REPORT.html",
    "template/main.tf",
    "template/README.md",
)
QUERY_KEYS = ("spm", "activeTab", "source", "sourcePath", "params")
DEFAULT_PREVIEW_ORIGIN = "https://pre-agent.aliyun-inc.com"
LOOPBACK_HOSTS = {"127.0.0.1", "::1", "localhost"}
FORBIDDEN_MARKUP_RE = re.compile(
    r"<\s*(?:script|iframe|object|embed|base|form|input|button)\b"
    r"|\bon[a-z]+\s*="
    r"|(?:href|src)\s*=\s*['\"]?\s*javascript:",
    flags=re.IGNORECASE,
)
DATA_URI_RE = re.compile(r"\bdata:[^,\s\"']+[,;]", flags=re.IGNORECASE)
DIRECT_SECRET_RE = re.compile(
    r"-----BEGIN [A-Z ]*PRIVATE KEY-----"
    r"|Authorization\s*:\s*Bearer\s+\S+"
    r"|Cookie\s*:\s*\S+"
    r"|\bLTAI[A-Za-z0-9]{12,}\b",
    flags=re.IGNORECASE,
)
SECRET_ASSIGNMENT_RE = re.compile(
    r"""(?ix)
    (?P<key>
        access[_-]?key(?:[_-]?(?:id|secret))?
        |secret[_-]?key
        |security[_-]?token
        |api[_-]?token
        |client[_-]?secret
        |password
    )
    \s*["']?\s*[:=]\s*
    (?:
        "(?P<double>[^"]*)"
        |'(?P<single>[^']*)'
        |(?P<bare>[^\s,}\]]+)
    )
    """,
)
RAW_TF_LOG_RE = re.compile(
    r"(?m)^(?:\d{4}-\d{2}-\d{2}T[^\n]*)?"
    r"\[(?:TRACE|DEBUG|INFO|WARN|ERROR)\][^\n]*"
    r"(?:terraform|provider|tf_req_id|tf_rpc|rpc=)"
    r"|^\s*TF_LOG(?:_PROVIDER|_PATH)?\s*=",
    flags=re.IGNORECASE,
)
INIT_EXTERNAL_RE = re.compile(
    r"failed to query available provider packages"
    r"|failed to install provider"
    r"|could not retrieve the list of available versions"
    r"|registry\.terraform\.io"
    r"|(?:network|connection|request) (?:error|failed|refused|reset|timed out)"
    r"|no such host|temporary failure in name resolution|tls handshake timeout",
    flags=re.IGNORECASE,
)


class InvalidPackage(Exception):
    """The package violates a deterministic validation contract."""


class BlockedValidation(Exception):
    """Validation cannot complete because an external gate is unavailable."""


class Result:
    def __init__(self, output_format: str) -> None:
        self.output_format = output_format
        self.checks: list[dict[str, str]] = []
        self.preview: dict[str, str] = {"status": "not_requested"}
        self.viewer_copy: dict[str, str] = {
            "status": "platform_blocked",
            "detail": "viewer copy is owned by the AutomationAgent platform",
        }

    def add(self, name: str, status: str, detail: str = "") -> None:
        item = {"name": name, "status": status}
        if detail:
            item["detail"] = detail
        self.checks.append(item)

    def emit(self, *, success: bool, status: str, exit_code: int) -> None:
        payload = {
            "success": success,
            "status": status,
            "exit_code": exit_code,
            "checks": self.checks,
            "preview": self.preview,
            "viewer_copy": self.viewer_copy,
        }
        if self.output_format == "json":
            print(json.dumps(payload, ensure_ascii=False, sort_keys=True))
            return
        for item in self.checks:
            suffix = f": {item['detail']}" if item.get("detail") else ""
            print(f"{item['name']}: {item['status']}{suffix}")
        print(f"preview: {self.preview['status']}")
        print(f"viewer_copy: {self.viewer_copy['status']}")
        print(f"result: {status} (exit {exit_code})")


def fail(message: str) -> NoReturn:
    raise InvalidPackage(message)


def read_utf8(path: Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except UnicodeDecodeError as error:
        fail(f"{path}: file is not valid UTF-8: {error}")


def decode_html_entities(value: str) -> str:
    decoded = value
    for _ in range(3):
        candidate = html.unescape(decoded)
        if candidate == decoded:
            break
        decoded = candidate
    return decoded


def contains_credential_value(text: str) -> bool:
    if DIRECT_SECRET_RE.search(text):
        return True
    placeholders = {
        "",
        "null",
        "none",
        "redacted",
        "<redacted>",
        "***",
        "xxxxx",
        "example",
        "example-value",
    }
    for match in SECRET_ASSIGNMENT_RE.finditer(text):
        value = next(
            (
                item
                for item in (
                    match.group("double"),
                    match.group("single"),
                    match.group("bare"),
                )
                if item is not None
            ),
            "",
        ).strip()
        lowered = value.lower()
        if (
            lowered not in placeholders
            and not lowered.startswith(("your_", "your-", "${", "var."))
        ):
            return True
    return False


def scan_package(package_dir: Path) -> None:
    for relative in REQUIRED_FILES:
        path = package_dir / relative
        if not path.is_file():
            fail(f"required file is missing: {relative}")

    for path in package_dir.rglob("*"):
        relative = path.relative_to(package_dir)
        if path.is_symlink():
            fail(f"symbolic links are forbidden in report packages: {relative}")
        lower_name = path.name.lower()
        lower_parts = {part.lower() for part in relative.parts}
        if ".terraform" in lower_parts:
            fail(f"forbidden Terraform data directory: {relative}")
        if path.is_dir():
            continue
        if (
            lower_name.endswith((".tfplan", ".tfstate", ".tfstate.backup"))
            or lower_name == "crash.log"
            or lower_name.startswith("crash.")
            or lower_name.endswith(".log")
            or (
                lower_name.endswith(".txt")
                and re.search(
                    r"(?:terraform|provider|tf)[_-]?(?:debug|trace|raw)",
                    lower_name,
                )
            )
            or lower_name.endswith((".tfvars", ".auto.tfvars"))
            or lower_name in {".env", "credentials.tfrc.json", ".terraformrc"}
        ):
            fail(f"forbidden sensitive artifact: {relative}")
        text = read_utf8(path)
        decoded_text = decode_html_entities(text)
        if DATA_URI_RE.search(decoded_text):
            fail(f"data URI or base64 payload is forbidden: {relative}")
        if contains_credential_value(decoded_text):
            fail(f"credential-like content is forbidden: {relative}")
        if RAW_TF_LOG_RE.search(decoded_text):
            fail(f"raw Terraform/provider debug log is forbidden: {relative}")
        if (
            path.suffix.lower() in {".html", ".htm", ".md", ".txt"}
            and FORBIDDEN_MARKUP_RE.search(decoded_text)
        ):
            fail(f"executable HTML is forbidden: {relative}")


def extract_hcl(package_dir: Path) -> tuple[str, str, str]:
    hcl_path = package_dir / "template/main.tf"
    hcl = read_utf8(hcl_path)
    if not hcl.endswith("\n"):
        fail("template/main.tf must end with one or more newline bytes")

    report_md = read_utf8(package_dir / "REPORT.md")
    md_blocks = re.findall(
        r"(?ms)^```hcl[ \t]*\n(.*?)^```[ \t]*$",
        report_md,
    )
    if len(md_blocks) != 1:
        fail("REPORT.md must contain exactly one complete ```hcl fenced block")

    report_html = read_utf8(package_dir / "REPORT.html")
    html_blocks = re.findall(
        r'(?is)<code\b[^>]*\bclass=["\'][^"\']*\blanguage-hcl\b[^"\']*["\'][^>]*>(.*?)</code>',
        report_html,
    )
    if len(html_blocks) != 1:
        fail("REPORT.html must contain exactly one language-hcl code block")
    html_hcl = html.unescape(html_blocks[0])

    expected = hcl.encode("utf-8")
    if md_blocks[0].encode("utf-8") != expected:
        fail("REPORT.md HCL bytes differ from template/main.tf")
    if html_hcl.encode("utf-8") != expected:
        fail("REPORT.html HCL bytes differ from template/main.tf")
    return hcl, report_md, report_html


def java_urlencode(value: str) -> str:
    return urllib.parse.quote_plus(value, safe="*-._", encoding="utf-8").replace(
        "+", "%20"
    )


def extract_online_url(text: str, source: str, *, decode_html: bool = False) -> str:
    if decode_html:
        text = html.unescape(text)
    urls = re.findall(
        r"https://api\.aliyun\.com/terraform\?[^)\s<>\"']+",
        text,
    )
    if len(urls) != 1:
        fail(f"{source} must contain exactly one api.aliyun.com/terraform link")
    return urls[0]


def validate_online_url(url: str, hcl: str) -> str:
    split = urllib.parse.urlsplit(url)
    if (
        split.scheme != "https"
        or split.netloc != "api.aliyun.com"
        or split.path != "/terraform"
        or split.fragment
    ):
        fail("online Terraform URL origin/path must be https://api.aliyun.com/terraform")

    pairs = split.query.split("&")
    if len(pairs) != len(QUERY_KEYS):
        fail("online Terraform URL query must contain exactly five fields")
    raw: dict[str, str] = {}
    ordered_keys: list[str] = []
    for pair in pairs:
        if "=" not in pair:
            fail("online Terraform URL query contains a field without '='")
        key, value = pair.split("=", 1)
        if key in raw:
            fail(f"online Terraform URL query repeats {key}")
        raw[key] = value
        ordered_keys.append(key)
    if tuple(ordered_keys) != QUERY_KEYS:
        fail("online Terraform URL query order or fixed field set is invalid")
    if raw["spm"] != "XToCode.TerraformAI.QA.0":
        fail("online Terraform URL spm is invalid")
    if raw["activeTab"] != "code" or raw["source"] != "PlayGround":
        fail("online Terraform URL activeTab/source is invalid")

    source_match = re.fullmatch(
        r"TerraformAI/([0-9]{13})::([0-9]{13})", raw["sourcePath"]
    )
    if not source_match or source_match.group(1) != source_match.group(2):
        fail("sourcePath must use the same 13-digit timestamp twice")
    if raw["params"] != java_urlencode(hcl):
        fail("params must be Java URLEncoder UTF-8 output encoded exactly once")
    return source_match.group(1)


def extract_variable_block(hcl: str, variable: str) -> str:
    match = re.search(rf'variable\s+"{re.escape(variable)}"\s*\{{', hcl)
    if not match:
        fail(f'variable "{variable}" is required')
    start = match.end()
    depth = 1
    index = start
    while index < len(hcl) and depth:
        if hcl[index] == "{":
            depth += 1
        elif hcl[index] == "}":
            depth -= 1
        index += 1
    if depth:
        fail(f'variable "{variable}" block is incomplete')
    return hcl[start : index - 1]


def validate_profile(hcl: str) -> None:
    block = extract_variable_block(hcl, "profile")
    requirements = (
        (r"(?m)^\s*type\s*=\s*string\s*$", "type = string"),
        (r"(?m)^\s*default\s*=\s*null\s*$", "default = null"),
        (r"(?m)^\s*nullable\s*=\s*true\s*$", "nullable = true"),
    )
    for pattern, label in requirements:
        if not re.search(pattern, block):
            fail(f'variable "profile" must declare {label}')


def run_terraform(package_dir: Path, result: Result) -> None:
    template_dir = package_dir / "template"
    commands = [
        ("terraform_fmt", ["terraform", "fmt", "-check", "-recursive"]),
        (
            "terraform_init",
            [
                "terraform",
                "init",
                "-backend=false",
                "-input=false",
                "-no-color",
            ],
        ),
        ("terraform_validate", ["terraform", "validate", "-no-color"]),
    ]
    if (template_dir / ".terraform.lock.hcl").is_file():
        commands[1][1].append("-lockfile=readonly")

    allowed_env = {
        "PATH",
        "LANG",
        "LC_ALL",
        "LC_CTYPE",
        "TMPDIR",
        "TEMP",
        "TMP",
        "SSL_CERT_FILE",
        "SSL_CERT_DIR",
        "HTTPS_PROXY",
        "HTTP_PROXY",
        "NO_PROXY",
        "https_proxy",
        "http_proxy",
        "no_proxy",
    }
    with tempfile.TemporaryDirectory(prefix="report-package-runtime-") as runtime_dir:
        runtime_path = Path(runtime_dir)
        home_dir = runtime_path / "home"
        data_dir = runtime_path / "terraform-data"
        home_dir.mkdir()
        data_dir.mkdir()
        env = {key: value for key, value in os.environ.items() if key in allowed_env}
        env.update(
            {
                "HOME": str(home_dir),
                "TF_DATA_DIR": str(data_dir),
                "TF_IN_AUTOMATION": "1",
                "CHECKPOINT_DISABLE": "1",
            }
        )
        for name, command in commands:
            try:
                completed = subprocess.run(
                    command,
                    cwd=template_dir,
                    env=env,
                    text=True,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                    timeout=180,
                    check=False,
                )
            except FileNotFoundError as error:
                raise BlockedValidation("terraform executable is unavailable") from error
            except subprocess.TimeoutExpired as error:
                raise BlockedValidation(f"{name} timed out") from error
            detail = completed.stdout.strip()[-2000:]
            if completed.returncode == 0:
                result.add(name, "passed")
                continue
            result.add(name, "failed", detail or f"exit {completed.returncode}")
            if name == "terraform_init" and INIT_EXTERNAL_RE.search(detail):
                raise BlockedValidation(
                    f"terraform init could not complete: {detail or completed.returncode}"
                )
            raise InvalidPackage(
                f"{name} failed: {detail or completed.returncode}"
            )


def parse_preview_origin(value: str) -> urllib.parse.SplitResult:
    try:
        origin = urllib.parse.urlsplit(value)
        port = origin.port
    except ValueError as error:
        fail(f"preview origin is invalid: {error}")
    if (
        not origin.hostname
        or origin.username is not None
        or origin.password is not None
        or "%" in origin.netloc
        or origin.path not in {"", "/"}
        or origin.query
        or origin.fragment
        or "?" in value
        or "#" in value
    ):
        fail("preview origin must contain only a scheme and host")
    host = origin.hostname.lower()
    if origin.scheme != "https" and not (
        origin.scheme == "http" and host in LOOPBACK_HOSTS
    ):
        fail("preview origin must use HTTPS (HTTP is allowed only for loopback tests)")
    if port is not None and not (1 <= port <= 65535):
        fail("preview origin port is invalid")
    return origin


def validate_absolute_preview_url(
    value: object,
    view_url: str,
    expected_origin: urllib.parse.SplitResult,
) -> str:
    if not isinstance(value, str):
        fail("preview record must include an absolute url")
    try:
        absolute = urllib.parse.urlsplit(value)
        absolute_port = absolute.port
        expected_port = expected_origin.port
    except ValueError as error:
        fail(f"preview absolute url is invalid: {error}")
    if (
        not absolute.hostname
        or absolute.username is not None
        or absolute.password is not None
        or "%" in absolute.netloc
    ):
        fail("preview absolute url authority is invalid")
    if (
        absolute.scheme != expected_origin.scheme
        or absolute.hostname.lower() != expected_origin.hostname.lower()
        or absolute_port != expected_port
    ):
        fail("preview absolute url origin does not match --preview-origin")
    if absolute.path != view_url or "%" in absolute.path:
        fail("preview absolute url path must exactly equal the readonly viewUrl")
    if absolute.query or absolute.fragment or "?" in value or "#" in value:
        fail("preview absolute url must not contain a query or fragment")
    return value


def load_preview_record(path: Path, preview_origin: str) -> dict[str, object]:
    if not path.is_file():
        raise BlockedValidation(f"preview JSON is unavailable: {path}")
    records: list[dict[str, object]] = []
    for line_number, line in enumerate(read_utf8(path).splitlines(), start=1):
        if not line.strip():
            continue
        try:
            record = json.loads(line)
        except json.JSONDecodeError as error:
            fail(f"preview JSON line {line_number} is invalid: {error}")
        if not isinstance(record, dict):
            fail(f"preview JSON line {line_number} must be an object")
        records.append(record)
    if len(records) != 1:
        fail("preview JSON must contain exactly one upload record")
    record = records[0]
    report_id = record.get("reportId")
    view_url = record.get("viewUrl")
    view_match = (
        re.fullmatch(
            r"/reports/aone/([A-Za-z0-9_-]+)/([A-Za-z0-9_-]+)/view",
            view_url,
        )
        if isinstance(view_url, str)
        else None
    )
    if (
        record.get("success") is not True
        or record.get("status") != "uploaded"
        or not isinstance(report_id, str)
        or not re.fullmatch(r"[A-Za-z0-9_-]+", report_id)
        or not view_match
        or view_match.group(2) != report_id
    ):
        fail(
            "preview record must be success=true/status=uploaded with reportId "
            "and a readonly /reports/aone/.../view route"
        )
    origin = parse_preview_origin(preview_origin)
    record["url"] = validate_absolute_preview_url(record.get("url"), view_url, origin)
    return record


def verify_preview(
    record: dict[str, object],
    report_html: str,
    expected_hcl: str,
    markers: list[str],
) -> None:
    url = record.get("url")
    if not isinstance(url, str) or not re.match(r"^https?://", url):
        raise BlockedValidation("preview record has no absolute anonymous GET url")
    request = urllib.request.Request(
        url,
        method="GET",
        headers={
            "Accept": "text/html",
            "User-Agent": "jarvis-report-package-validator/1",
        },
    )
    class NoRedirect(urllib.request.HTTPRedirectHandler):
        def redirect_request(self, *_args: object, **_kwargs: object) -> None:
            return None

    opener = urllib.request.build_opener(
        urllib.request.ProxyHandler({}),
        NoRedirect(),
    )
    try:
        with opener.open(request, timeout=20) as response:
            status = response.status
            content_type = response.headers.get_content_type()
            final_url = response.geturl()
            body = response.read().decode("utf-8")
    except (urllib.error.URLError, UnicodeDecodeError, TimeoutError) as error:
        raise BlockedValidation(f"anonymous preview GET failed: {error}") from error
    if status != 200 or content_type != "text/html":
        raise BlockedValidation(
            f"anonymous preview GET returned {status} {content_type}"
        )
    if final_url != url:
        raise BlockedValidation("anonymous preview GET final URL changed")

    title_match = re.search(r"(?is)<title>(.*?)</title>", report_html)
    if not title_match:
        fail("REPORT.html has no title")
    expected_title = html.unescape(title_match.group(1)).strip()
    remote_title_match = re.search(r"(?is)<title>(.*?)</title>", body)
    if (
        not remote_title_match
        or html.unescape(remote_title_match.group(1)).strip() != expected_title
    ):
        raise BlockedValidation("preview title does not match REPORT.html")
    remote_blocks = re.findall(
        r'(?is)<code\b[^>]*\bclass=["\'][^"\']*\blanguage-hcl\b[^"\']*["\'][^>]*>(.*?)</code>',
        body,
    )
    if len(remote_blocks) != 1 or html.unescape(remote_blocks[0]) != expected_hcl:
        raise BlockedValidation("preview HCL does not match template/main.tf")
    decoded_body = html.unescape(body)
    for marker in markers:
        if marker not in decoded_body:
            raise BlockedValidation(f"preview marker is missing: {marker}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("package_dir", type=Path)
    parser.add_argument("--require-preview", action="store_true")
    parser.add_argument("--preview-json", type=Path)
    parser.add_argument(
        "--preview-origin",
        default=DEFAULT_PREVIEW_ORIGIN,
        help=f"expected preview origin (default: {DEFAULT_PREVIEW_ORIGIN})",
    )
    parser.add_argument("--preview-marker", action="append", default=[])
    parser.add_argument("--require-viewer-copy", action="store_true")
    parser.add_argument("--format", choices=("text", "json"), default="text")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    result = Result(args.format)
    try:
        package_dir = args.package_dir.resolve()
        if not package_dir.is_dir():
            fail(f"package directory does not exist: {package_dir}")
        scan_package(package_dir)
        result.add("package_safety", "passed")

        hcl, report_md, report_html = extract_hcl(package_dir)
        result.add("hcl_byte_identity", "passed")

        md_url = extract_online_url(report_md, "REPORT.md")
        html_url = extract_online_url(report_html, "REPORT.html", decode_html=True)
        if md_url != html_url:
            fail("REPORT.md and REPORT.html online Terraform URLs differ")
        timestamp = validate_online_url(md_url, hcl)
        result.add("online_prefill", "passed", f"timestamp={timestamp}")

        validate_profile(hcl)
        result.add("profile_variable", "passed")
        run_terraform(package_dir, result)

        if args.preview_json:
            record = load_preview_record(args.preview_json, args.preview_origin)
            result.preview = {
                "status": "recorded",
                "reportId": str(record["reportId"]),
                "viewUrl": str(record["viewUrl"]),
            }
            result.add("preview_record", "passed")
        else:
            record = None

        if args.require_preview:
            if record is None:
                raise BlockedValidation(
                    "--require-preview needs --preview-json from a successful upload"
                )
            if not args.preview_marker or any(
                not marker.strip() for marker in args.preview_marker
            ):
                fail("--require-preview needs at least one non-empty --preview-marker")
            verify_preview(record, report_html, hcl, args.preview_marker)
            result.preview["status"] = "verified"
            result.add("preview_anonymous_get", "passed")
        elif args.preview_marker:
            fail("--preview-marker requires --require-preview")

        if args.require_viewer_copy:
            raise BlockedValidation(
                "viewer copy remains platform_blocked; this report package cannot fix it"
            )
    except InvalidPackage as error:
        result.add("validation", "failed", str(error))
        result.emit(success=False, status="invalid", exit_code=EXIT_INVALID)
        return EXIT_INVALID
    except BlockedValidation as error:
        result.add("validation", "blocked", str(error))
        result.emit(success=False, status="blocked", exit_code=EXIT_BLOCKED)
        return EXIT_BLOCKED

    result.emit(success=True, status="validated", exit_code=0)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
