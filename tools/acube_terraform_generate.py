#!/usr/bin/env python3
"""Generate Terraform provider files through Acube, with offline fixture support."""

from __future__ import annotations

import argparse
import json
import re
import ssl
import sys
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path, PurePosixPath
from typing import Any


DEFAULT_GENERATOR_PREFIX = "/api/v1/terraform/generator"
MAPPING_PATH = "/api/v1/terraform/resource/mapping/createMapping"
BUILD_PATH = "/api/v1/terraform_vendor_build/createLocalBuildTask"
PRODUCT_PREFIXES = [
    "resource_manager",
    "cloud_sso",
    "cloud_storage_gateway",
    "cloud_firewall",
    "cloud_monitor",
    "cloud_control",
    "api_gateway",
    "log_service",
    "event_bridge",
    "private_link",
    "direct_mail",
    "action_trail",
    "data_works",
    "open_search",
    "alikafka",
    "alb",
    "arms",
    "bastionhost",
    "cdn",
    "cen",
    "cms",
    "cr",
    "cs",
    "dbs",
    "dns",
    "dts",
    "ecs",
    "edas",
    "ehpc",
    "emr",
    "ens",
    "ess",
    "fc",
    "ga",
    "hbr",
    "kms",
    "mns",
    "nas",
    "oos",
    "oss",
    "ots",
    "polardb",
    "ram",
    "rds",
    "ros",
    "sag",
    "slb",
    "vpc",
    "vpn",
]
FILE_MAP_KEYS = {
    "files",
    "fileMap",
    "file_map",
    "generatedFiles",
    "generated_files",
    "resultFiles",
    "fileContentMap",
}
CONTENT_KEYS = ("content", "fileContent", "body", "text", "source")
LOG_KEYS = {"logs", "log", "output", "stdout", "stderr", "message", "msg", "errorMessage"}


def camelize(value: str) -> str:
    parts = [part for part in re.split(r"[_\-\s]+", value) if part]
    return "".join(part[:1].upper() + part[1:] for part in parts)


def snake_case(value: str) -> str:
    value = re.sub(r"(.)([A-Z][a-z]+)", r"\1_\2", value)
    value = re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", value)
    return value.replace("-", "_").lower()


def token_prefix(tokens: list[str], prefix: str) -> bool:
    prefix_tokens = [part for part in prefix.split("_") if part]
    return tokens[: len(prefix_tokens)] == prefix_tokens


def infer_product_and_resource_code(
    resource: str | None,
    product_override: str | None,
    resource_code_override: str | None,
) -> tuple[str, str]:
    if not resource and not (product_override and resource_code_override):
        raise ValueError("--resource is required unless both --product and --resource-code are provided")

    base = resource.removeprefix("alicloud_") if resource and resource.startswith("alicloud_") else resource
    tokens = [part for part in (base or "").split("_") if part]

    if product_override:
        product = product_override
        product_prefix = snake_case(product_override)
        if token_prefix(tokens, product_prefix):
            resource_tokens = tokens[len(product_prefix.split("_")) :]
        else:
            resource_tokens = tokens
    else:
        matched_prefix = next((prefix for prefix in PRODUCT_PREFIXES if token_prefix(tokens, prefix)), None)
        if matched_prefix:
            product = camelize(matched_prefix)
            resource_tokens = tokens[len(matched_prefix.split("_")) :]
        elif tokens:
            product = camelize(tokens[0])
            resource_tokens = tokens[1:]
        else:
            raise ValueError(f"cannot infer product/resourceCode from resource: {resource!r}")

    if resource_code_override:
        resource_code = resource_code_override
    else:
        if not resource_tokens:
            raise ValueError(
                f"cannot infer resourceCode from resource {resource!r}; pass --resource-code explicitly"
            )
        resource_code = camelize("_".join(resource_tokens))

    return product, resource_code


def default_host(env: str) -> str:
    return "https://pre-acube.aliyun-inc.com" if env == "pre" else "https://acube.aliyun-inc.com"


def operation_url(host: str, path: str, query: dict[str, Any] | None = None) -> str:
    url = "/".join([host.rstrip("/"), path.strip("/")])
    if query:
        url = f"{url}?{urllib.parse.urlencode(query)}"
    return url


def read_json(path: Path) -> Any:
    with path.open(encoding="utf-8") as fh:
        return json.load(fh)


def write_json(path: Path, value: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8") as fh:
        json.dump(value, fh, indent=2, ensure_ascii=False, sort_keys=True)
        fh.write("\n")


def request_json(url: str, timeout: float, insecure: bool) -> Any:
    request = urllib.request.Request(url, headers={"Accept": "*/*"}, method="GET")
    context = ssl._create_unverified_context() if insecure else None
    try:
        with urllib.request.urlopen(request, timeout=timeout, context=context) as response:
            raw = response.read().decode("utf-8", errors="replace")
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"GET {url} failed with HTTP {exc.code}: {raw[:500]}") from exc
    except urllib.error.URLError as exc:
        raise RuntimeError(f"GET {url} failed: {exc}") from exc

    if not raw.strip():
        return {}
    try:
        return json.loads(raw)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"GET {url} returned non-JSON response: {raw[:500]}") from exc


def post_form_json(url: str, form: dict[str, Any], timeout: float, insecure: bool) -> Any:
    data = urllib.parse.urlencode(form).encode("utf-8")
    request = urllib.request.Request(
        url,
        data=data,
        headers={"Content-Type": "application/x-www-form-urlencoded", "Accept": "*/*"},
        method="POST",
    )
    context = ssl._create_unverified_context() if insecure else None
    try:
        with urllib.request.urlopen(request, timeout=timeout, context=context) as response:
            raw = response.read().decode("utf-8", errors="replace")
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"POST {url} failed with HTTP {exc.code}: {raw[:500]}") from exc
    except urllib.error.URLError as exc:
        raise RuntimeError(f"POST {url} failed: {exc}") from exc

    if not raw.strip():
        return {}
    try:
        return json.loads(raw)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"POST {url} returned non-JSON response: {raw[:500]}") from exc


def response_ok(value: Any) -> bool:
    if not isinstance(value, dict):
        return True
    if "success" in value:
        return bool(value["success"])
    if "successful" in value:
        return bool(value["successful"])
    if "status" in value:
        status = value["status"]
        if isinstance(status, str):
            return status.lower() in {"ok", "success", "succeeded", "done"}
        if isinstance(status, int):
            return status in {0, 200}
    if "code" in value:
        code = value["code"]
        if isinstance(code, int):
            return code in {0, 200}
        if isinstance(code, str):
            return code.lower() in {"0", "200", "ok", "success"}
    return True


def walk(value: Any) -> list[Any]:
    values = [value]
    if isinstance(value, dict):
        for item in value.values():
            values.extend(walk(item))
    elif isinstance(value, list):
        for item in value:
            values.extend(walk(item))
    return values


def find_first_key(value: Any, keys: set[str]) -> Any:
    if isinstance(value, dict):
        for key, item in value.items():
            if key in keys and item not in ("", None):
                return item
        for item in value.values():
            found = find_first_key(item, keys)
            if found not in ("", None):
                return found
    elif isinstance(value, list):
        for item in value:
            found = find_first_key(item, keys)
            if found not in ("", None):
                return found
    return None


def extract_message(value: Any) -> str:
    found = find_first_key(value, {"message", "msg", "errorMessage", "error"})
    if found is None:
        return ""
    if isinstance(found, str):
        return found
    return json.dumps(found, ensure_ascii=False, sort_keys=True)


def extract_resource_type_code(value: Any, product: str, resource_code: str) -> str:
    found = find_first_key(value, {"resourceTypeCode", "resource_type_code"})
    if found is not None:
        return str(found)
    return f"ALIYUN::{product}::{resource_code}"


def extract_task_id(value: Any) -> str | None:
    found = find_first_key(value, {"taskId", "taskID", "buildTaskId", "buildTaskID"})
    return None if found is None else str(found)


def file_content(value: Any) -> str | None:
    if isinstance(value, str):
        return value
    if isinstance(value, dict):
        for key in CONTENT_KEYS:
            item = value.get(key)
            if isinstance(item, str):
                return item
    return None


def is_files_map(value: Any) -> bool:
    if not isinstance(value, dict) or not value:
        return False
    return all(isinstance(key, str) and file_content(item) is not None for key, item in value.items())


def find_files_map(value: Any) -> dict[str, Any]:
    if isinstance(value, dict):
        for key, item in value.items():
            if key in FILE_MAP_KEYS and is_files_map(item):
                return item
        for item in value.values():
            found = find_files_map(item)
            if found:
                return found
    elif isinstance(value, list):
        for item in value:
            found = find_files_map(item)
            if found:
                return found
    return {}


def normalize_repository_path(raw_path: str) -> str:
    path = raw_path.replace("\\", "/").strip()
    if "/repository/" in path:
        path = path.split("/repository/", 1)[1]
    elif path.startswith("repository/"):
        path = path[len("repository/") :]
    path = path.lstrip("/")
    while path.startswith("./"):
        path = path[2:]

    parts = PurePosixPath(path).parts
    if not parts or any(part in {"", ".", ".."} for part in parts):
        raise ValueError(f"unsafe generated file path from Acube files map: {raw_path!r}")
    return "/".join(parts)


def append_log_value(lines: list[str], value: Any) -> None:
    if isinstance(value, list):
        for item in value:
            append_log_value(lines, item)
    elif isinstance(value, dict):
        lines.append(json.dumps(value, ensure_ascii=False, sort_keys=True))
    elif value not in (None, ""):
        lines.extend(str(value).splitlines())


def extract_logs(value: Any) -> list[str]:
    lines: list[str] = []
    if isinstance(value, dict):
        for key, item in value.items():
            if key in LOG_KEYS:
                append_log_value(lines, item)
            else:
                lines.extend(extract_logs(item))
    elif isinstance(value, list):
        for item in value:
            lines.extend(extract_logs(item))
    return lines


def write_generated_files(generated_dir: Path, files_map: dict[str, Any]) -> list[str]:
    generated_files: list[str] = []
    for raw_path, raw_content in sorted(files_map.items()):
        rel_path = normalize_repository_path(raw_path)
        content = file_content(raw_content)
        if content is None:
            raise ValueError(f"Acube file map entry has no text content: {raw_path!r}")
        destination = generated_dir / rel_path
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_text(content, encoding="utf-8")
        generated_files.append(rel_path)
    return generated_files


class OperationRunner:
    def __init__(
        self,
        *,
        host: str,
        api_prefix: str,
        timeout: float,
        output_dir: Path,
        offline: bool,
        insecure: bool,
    ) -> None:
        self.host = host
        self.api_prefix = api_prefix
        self.timeout = timeout
        self.output_dir = output_dir
        self.offline = offline
        self.insecure = insecure
        self.operations: list[dict[str, Any]] = []
        self.logs: list[str] = []

    def run(
        self,
        *,
        name: str,
        raw_filename: str,
        method: str,
        path: str,
        query: dict[str, Any] | None = None,
        form: dict[str, Any] | None = None,
        fixture: Path | None,
        optional: bool = False,
    ) -> Any:
        raw_path = self.output_dir / raw_filename
        url = operation_url(self.host, path, query)
        record: dict[str, Any] = {"name": name, "method": method, "url": url, "rawJson": raw_filename}
        try:
            if fixture:
                response = read_json(fixture)
                record["source"] = "fixture"
                record["fixture"] = str(fixture)
            elif self.offline:
                if optional:
                    response = {"skipped": True, "message": f"offline fixture not provided for {name}"}
                    record.update({"source": "offline", "skipped": True})
                    self.operations.append(record)
                    self.logs.append(f"[{name}] skipped: offline fixture not provided")
                    return response
                raise RuntimeError(f"offline mode requires a fixture for {name}")
            else:
                if method == "GET":
                    response = request_json(url, self.timeout, self.insecure)
                elif method == "POST_FORM":
                    response = post_form_json(url, form or {}, self.timeout, self.insecure)
                else:
                    raise RuntimeError(f"unsupported HTTP method for {name}: {method}")
                record["source"] = "network"

            write_json(raw_path, response)
            ok = response_ok(response)
            record["ok"] = ok
            message = extract_message(response)
            if message:
                record["message"] = message
            self.operations.append(record)
            self.logs.append(f"[{name}] source={record['source']} ok={ok}")
            self.logs.extend(extract_logs(response))
            if not ok and not optional:
                raise RuntimeError(f"{name} failed: {message or json.dumps(response, ensure_ascii=False)[:500]}")
            return response
        except Exception as exc:
            if optional:
                error_response = {"success": False, "message": str(exc), "nonBlocking": True}
                write_json(raw_path, error_response)
                record.update({"ok": False, "source": record.get("source", "network"), "message": str(exc)})
                self.operations.append(record)
                self.logs.append(f"[{name}] non-blocking failure: {exc}")
                return error_response
            raise

    def skip(self, *, name: str, reason: str) -> None:
        self.operations.append({"name": name, "skipped": True, "reason": reason})
        self.logs.append(f"[{name}] skipped: {reason}")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Generate Terraform provider files from Acube.")
    parser.add_argument("--resource", help="Terraform resource name, e.g. alicloud_oss_bucket")
    parser.add_argument("--product", help="Explicit Acube product, overriding resource-name inference")
    parser.add_argument("--resource-code", help="Explicit Acube resourceCode, overriding resource-name inference")
    parser.add_argument("--env", default="pre", help="Acube environment. Defaults to pre.")
    parser.add_argument("--host", help="Acube host. Defaults to pre-acube for env=pre, otherwise acube.")
    parser.add_argument(
        "--api-prefix",
        default=DEFAULT_GENERATOR_PREFIX,
        help=f"Acube Terraform generator API prefix. Defaults to {DEFAULT_GENERATOR_PREFIX}.",
    )
    parser.add_argument("--output-dir", "--out-dir", required=True, type=Path, help="Directory for raw JSON and generated files")
    parser.add_argument("--skip-mapping", action="store_true", help="Skip createMapping and go directly to local build")
    parser.add_argument("--timeout", default=60.0, type=float, help="Network timeout in seconds")
    parser.add_argument(
        "--insecure",
        action="store_true",
        help="Disable TLS certificate verification for internal Acube hosts with local CA issues.",
    )
    parser.add_argument("--fixture-resource-type-json", type=Path, help="Offline fixture for resourceTypeCode/get")
    parser.add_argument("--fixture-mapping-json", type=Path, help="Offline fixture for createMapping")
    parser.add_argument("--fixture-build-json", type=Path, help="Offline fixture for createLocalBuildTask")
    parser.add_argument("--fixture-document-json", type=Path, help="Offline fixture for queryLatestDocument")
    return parser


def ensure_fixture(path: Path | None, parser: argparse.ArgumentParser, option: str) -> None:
    if path is not None and not path.is_file():
        parser.error(f"{option} does not exist or is not a file: {path}")


def validate_args(args: argparse.Namespace, parser: argparse.ArgumentParser) -> None:
    if not args.resource and not (args.product and args.resource_code):
        parser.error("--resource is required unless both --product and --resource-code are provided")
    ensure_fixture(args.fixture_resource_type_json, parser, "--fixture-resource-type-json")
    ensure_fixture(args.fixture_mapping_json, parser, "--fixture-mapping-json")
    ensure_fixture(args.fixture_build_json, parser, "--fixture-build-json")
    ensure_fixture(args.fixture_document_json, parser, "--fixture-document-json")


def run(args: argparse.Namespace) -> dict[str, Any]:
    product, resource_code = infer_product_and_resource_code(args.resource, args.product, args.resource_code)
    resource = args.resource or f"alicloud_{snake_case(product)}_{snake_case(resource_code)}"
    host = args.host or default_host(args.env)
    output_dir: Path = args.output_dir
    generated_dir = output_dir / "generated"
    output_dir.mkdir(parents=True, exist_ok=True)
    generated_dir.mkdir(parents=True, exist_ok=True)

    offline = any(
        (
            args.fixture_resource_type_json,
            args.fixture_mapping_json,
            args.fixture_build_json,
            args.fixture_document_json,
        )
    )
    runner = OperationRunner(
        host=host,
        api_prefix=args.api_prefix,
        timeout=args.timeout,
        output_dir=output_dir,
        offline=offline,
        insecure=args.insecure,
    )

    base_payload = {
        "terraformResource": resource,
        "product": product,
        "resourceCode": resource_code,
    }

    resource_type_response = runner.run(
        name="resourceTypeCode/get",
        raw_filename="resource_type_code_get.json",
        method="GET",
        path=f"{args.api_prefix}/cloudspec/resourceTypeCode/get",
        query={
            "env": args.env,
            "isShowChangeLog": "false",
            "product": product,
            "resourceCode": resource_code,
        },
        fixture=args.fixture_resource_type_json,
    )
    resource_type_code = extract_resource_type_code(resource_type_response, product, resource_code)

    mapping_response: Any = None
    if args.skip_mapping:
        runner.skip(name="createMapping", reason="--skip-mapping")
    else:
        mapping_response = runner.run(
            name="createMapping",
            raw_filename="create_mapping.json",
            method="GET",
            path=MAPPING_PATH,
            query={
                "name": resource,
                "namespace": product,
                "resourceCode": resource_code,
                "isSpec": "true",
                "ignoreVerify": "true",
            },
            fixture=args.fixture_mapping_json,
        )

    build_payload = {"namespace": product, "resourceTypeCode": resource_code, "env": args.env}
    build_response = runner.run(
        name="createLocalBuildTask",
        raw_filename="create_local_build_task.json",
        method="POST_FORM",
        path=BUILD_PATH,
        form=build_payload,
        fixture=args.fixture_build_json,
    )

    files_map = find_files_map(build_response)
    generated_files = write_generated_files(generated_dir, files_map)
    task_id = extract_task_id(build_response)

    document_payload = {
        "namespace": product,
        "resourceTypeCode": resource_code,
        "terraformResourceType": resource,
    }
    document_response = runner.run(
        name="queryLatestDocument",
        raw_filename="query_latest_document.json",
        method="GET",
        path=f"{args.api_prefix}/queryLatestDocument",
        query=document_payload,
        fixture=args.fixture_document_json,
        optional=True,
    )
    document_ok = response_ok(document_response)
    document_message = extract_message(document_response)

    logs_path = output_dir / "logs.txt"
    files_path = output_dir / "files.txt"
    runner.logs.insert(0, f"resourceTypeCode: {resource_type_code}")
    runner.logs.insert(0, f"resourceCode: {resource_code}")
    runner.logs.insert(0, f"product: {product}")
    runner.logs.insert(0, f"resource: {resource}")
    logs_path.write_text("\n".join(line for line in runner.logs if line) + "\n", encoding="utf-8")
    files_path.write_text("".join(f"{path}\n" for path in generated_files), encoding="utf-8")

    summary = {
        "resource": resource,
        "product": product,
        "resourceCode": resource_code,
        "resourceTypeCode": resource_type_code,
        "env": args.env,
        "host": host,
        "apiPrefix": args.api_prefix,
        "insecure": args.insecure,
        "offline": offline,
        "outputDir": str(output_dir),
        "generatedDir": str(generated_dir),
        "logsPath": "logs.txt",
        "filesPath": "files.txt",
        "rawJson": {
            "resourceTypeCode/get": "resource_type_code_get.json",
            "createLocalBuildTask": "create_local_build_task.json",
            "queryLatestDocument": "query_latest_document.json",
            **({} if args.skip_mapping else {"createMapping": "create_mapping.json"}),
        },
        "files": generated_files,
        "taskId": task_id,
        "document": {
            "ok": document_ok,
            "nonBlocking": True,
            "message": document_message,
        },
        "operations": runner.operations,
    }
    write_json(output_dir / "summary.json", summary)
    return summary


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    validate_args(args, parser)
    try:
        summary = run(args)
    except Exception as exc:
        print(f"acube_terraform_generate: ERROR: {exc}", file=sys.stderr)
        return 1

    print(f"Resource: {summary['resource']}")
    print(f"Product: {summary['product']}")
    print(f"ResourceCode: {summary['resourceCode']}")
    print(f"ResourceTypeCode: {summary['resourceTypeCode']}")
    print(f"Generated files: {len(summary['files'])}")
    print(f"Output dir: {summary['outputDir']}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
