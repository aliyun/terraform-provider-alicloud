#!/usr/bin/env python3
"""
AMP Runtime API Metadata Fetcher

Calls the GetRuntimeApiForInner API on AMP (镇元平台) to retrieve the full
OpenAPI specification (schema) for a specific cloud API action.

This returns the complete API metadata: request parameters, response schema,
error codes, etc. — everything needed to understand how an API works.

Usage:
    python3 get_runtime_api.py \
        --pop-code Ecs \
        --pop-version 2014-05-26 \
        --api-name DescribeInstances \
        --runtime-type pop

Environment variables:
    ALIBABA_CLOUD_ACCESS_KEY_ID
    ALIBABA_CLOUD_ACCESS_KEY_SECRET

Author: AMP Tools Team
"""

import os
import sys
import json
import hmac
import hashlib
import base64
import argparse
from typing import Dict, Any, Optional

from alibabacloud_tea_openapi import models as open_api_models
from alibabacloud_tea_openapi.client import Client as OpenApiClient
from alibabacloud_tea_util import models as util_models

try:
    import yaml  # optional — provides better JSON conversion
    HAS_YAML = True
except ImportError:
    HAS_YAML = False


# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------
API_VERSION = "2021-06-01"
API_ACTION = "GetRuntimeApiForInner"
API_METHOD = "GET"
API_STYLE = "RPC"
DEFAULT_ENDPOINT = "apispecdata-share.aliyuncs.com"

# App-level sign (same as Java: ProductMetaUtil.getSign)
APP_NAME = "a-cube-aliyun-com"
APP_SIGN_SECRET = f"{APP_NAME}:137e92e05e524e6c9e3d117d24a0a68f"


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def get_app_sign() -> str:
    """
    Generate HMAC-SHA1 signature for the RequestContext.Sign parameter.
    Mirrors com.aliyun.openplatform.acube.autogenerate.specification.util.ProductMetaUtil#getSign
    """
    signature = hmac.new(
        APP_SIGN_SECRET.encode("utf-8"),
        APP_NAME.encode("utf-8"),
        hashlib.sha1,
    ).digest()
    return base64.b64encode(signature).decode("utf-8")


def create_client(endpoint: str = DEFAULT_ENDPOINT) -> OpenApiClient:
    """Create an OpenAPI client using AK credentials from environment variables."""
    access_key_id = os.environ.get("ALIBABA_CLOUD_ACCESS_KEY_ID")
    access_key_secret = os.environ.get("ALIBABA_CLOUD_ACCESS_KEY_SECRET")

    if not access_key_id or not access_key_secret:
        print(
            "Error: Please set ALIBABA_CLOUD_ACCESS_KEY_ID and "
            "ALIBABA_CLOUD_ACCESS_KEY_SECRET environment variables.",
            file=sys.stderr,
        )
        sys.exit(1)

    config = open_api_models.Config(
        access_key_id=access_key_id,
        access_key_secret=access_key_secret,
    )
    config.endpoint = endpoint
    config.protocol = "HTTPS"
    config.read_timeout = 30000
    config.connect_timeout = 10000

    return OpenApiClient(config)


# ---------------------------------------------------------------------------
# API call
# ---------------------------------------------------------------------------

def get_runtime_api(
    client: OpenApiClient,
    runtime_type: str,
    pop_code: str,
    pop_version: str,
    api_name: str,
    env: str = "online",
) -> Dict[str, Any]:
    """
    Call GetRuntimeApiForInner on AMP to retrieve the full API schema.

    Args:
        client:       Initialised OpenApiClient.
        runtime_type: Gateway type, e.g. "pop", "inner-api".
        pop_code:     POP product code, e.g. "Ecs", "Vpc".
        pop_version:  POP version, e.g. "2014-05-26".
        api_name:     API action name, e.g. "DescribeInstances".
        env:          Environment: "online", "pre", "daily".

    Returns:
        The full API response dict.
    """
    params = open_api_models.Params(
        style=API_STYLE,
        version=API_VERSION,
        action=API_ACTION,
        method=API_METHOD,
        pathname="/",
        protocol="HTTPS",
        auth_type="AK",
        req_body_type="json",
        body_type="json",
    )

    query: Dict[str, Any] = {
        "RuntimeType": runtime_type,
        "PopCode": pop_code,
        "PopVersion": pop_version,
        "ApiName": api_name,
        "Env": env,
        "RequestContext.AppName": APP_NAME,
        "RequestContext.Sign": get_app_sign(),
    }

    request = open_api_models.OpenApiRequest(query=query)

    runtime = util_models.RuntimeOptions(
        autoretry=True,
        max_attempts=3,
    )

    response = client.call_api(params, request, runtime)
    return response


# ---------------------------------------------------------------------------
# Output helpers
# ---------------------------------------------------------------------------

def extract_schema(response: Dict[str, Any]) -> Optional[str]:
    """Extract the Schema string from the API response."""
    body = response.get("body", {})
    data = body.get("Data", body.get("data", {}))
    if isinstance(data, dict):
        return data.get("Schema", data.get("schema"))
    return None


def _parse_schema(schema_str: str) -> Optional[Dict[str, Any]]:
    """Try to parse Schema string as JSON, then YAML."""
    # Try JSON first
    try:
        return json.loads(schema_str)
    except (json.JSONDecodeError, TypeError):
        pass
    # Try YAML if available
    if HAS_YAML:
        try:
            return yaml.safe_load(schema_str)
        except Exception:
            pass
    return None


def print_result(response: Dict[str, Any], output: str, verbose: bool) -> None:
    """Print the API metadata."""
    if verbose:
        print(json.dumps(response, indent=2, ensure_ascii=False))
        return

    body = response.get("body", {})

    # Check for errors
    success = body.get("Success", body.get("success"))
    if success is not None and not success:
        code = body.get("Code", body.get("code", "Unknown"))
        message = body.get("Message", body.get("message", "Unknown error"))
        print(f"API call failed — Code: {code}, Message: {message}", file=sys.stderr)
        sys.exit(1)

    schema_str = extract_schema(response)
    if not schema_str:
        print("No Schema data returned.", file=sys.stderr)
        print("Full response body:", file=sys.stderr)
        print(json.dumps(body, indent=2, ensure_ascii=False), file=sys.stderr)
        sys.exit(1)

    if output == "json":
        schema = _parse_schema(schema_str)
        if schema:
            print(json.dumps(schema, indent=2, ensure_ascii=False))
        else:
            # Cannot parse — output raw text (likely YAML)
            print(schema_str)
    else:
        # Pretty mode: try to parse and summarize, fallback to raw output
        schema = _parse_schema(schema_str)
        if schema:
            _pretty_print_schema(schema)
        else:
            # Raw YAML output (still very readable)
            print(schema_str)


def _pretty_print_schema(schema: Dict[str, Any]) -> None:
    """Human-readable summary of the API schema."""
    info = schema.get("info", {})
    apis = schema.get("apis", {})
    components = schema.get("components", {}).get("schemas", {})

    print("=" * 70)
    print("  AMP API Metadata")
    print("=" * 70)

    if info:
        desc = info.get("description", "N/A")
        style = info.get("apiStyle", info.get("style", "N/A"))
        ns = info.get("namespace", "N/A")
        print(f"  Description: {desc}")
        print(f"  Namespace  : {ns}")
        print(f"  Style      : {style}")

    if apis:
        print("-" * 70)
        for api_name, api_detail in apis.items():
            print(f"\n  API: {api_name}")
            methods = api_detail.get("methods", [])
            print(f"  Methods : {', '.join(methods)}")
            print(f"  Summary : {api_detail.get('summary', 'N/A')}")
            print(f"  Deprecated: {api_detail.get('deprecated', False)}")

            # Parameters
            params = api_detail.get("parameters", [])
            if params:
                print(f"\n  Parameters ({len(params)}):")
                for p in params:
                    name = p.get("name", "?")
                    required = "required" if p.get("required") else "optional"
                    p_schema = p.get("schema", {})
                    p_type = p_schema.get("type", "string")
                    desc = p.get("description", "")
                    in_loc = p.get("in", "query")
                    print(f"    - {name:<30s} {p_type:<10s} {required:<10s} (in: {in_loc})")
                    if desc:
                        # Truncate long descriptions
                        short = desc.replace("\n", " ")[:120]
                        print(f"      {short}")

            # Responses
            responses = api_detail.get("responses", {})
            if responses:
                print(f"\n  Responses:")
                for code, resp in responses.items():
                    ref = resp.get("schema", {}).get("$ref", "")
                    print(f"    {code}: {ref or '(inline schema)'}")

    if components:
        print("-" * 70)
        print(f"\n  Schema Components ({len(components)}):")
        for comp_name in sorted(components.keys()):
            comp = components[comp_name]
            props = comp.get("properties", {})
            comp_type = comp.get("type", "object")
            print(f"    {comp_name} ({comp_type}, {len(props)} fields)")

    print("\n" + "=" * 70)
    print("Tip: Use --output json for the complete machine-readable schema.")


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(
        description="Fetch API metadata (OpenAPI schema) from AMP via GetRuntimeApiForInner.",
        epilog=(
            "Required environment variables:\n"
            "  ALIBABA_CLOUD_ACCESS_KEY_ID       Alibaba Cloud Access Key ID\n"
            "  ALIBABA_CLOUD_ACCESS_KEY_SECRET   Alibaba Cloud Access Key Secret\n"
            "\n"
            "Examples:\n"
            "  %(prog)s --pop-code Ecs --pop-version 2014-05-26 --api-name DescribeInstances\n"
            "  %(prog)s --pop-code Vpc --pop-version 2016-04-28 --api-name CreateVpc --output json\n"
            "  %(prog)s --pop-code Ecs --pop-version 2014-05-26 --api-name RunInstances --runtime-type pop\n"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--pop-code", required=True, help="POP product code, e.g. Ecs, Vpc, Slb")
    parser.add_argument("--pop-version", required=True, help="POP version, e.g. 2014-05-26")
    parser.add_argument("--api-name", required=True, help="API action name, e.g. DescribeInstances")
    parser.add_argument("--runtime-type", default="pop", help="Gateway type (default: pop)")
    parser.add_argument("--env", default="online", help="Environment: online (default), pre, daily")
    parser.add_argument("--endpoint", default=DEFAULT_ENDPOINT, help=f"Endpoint (default: {DEFAULT_ENDPOINT})")
    parser.add_argument("--output", choices=["pretty", "json"], default="pretty", help="Output format")
    parser.add_argument("--verbose", "-v", action="store_true", help="Print full raw response")

    args = parser.parse_args()

    try:
        client = create_client(endpoint=args.endpoint)
        response = get_runtime_api(
            client=client,
            runtime_type=args.runtime_type,
            pop_code=args.pop_code,
            pop_version=args.pop_version,
            api_name=args.api_name,
            env=args.env,
        )
        print_result(response, args.output, args.verbose)
    except Exception as exc:
        print(f"Error: {exc}", file=sys.stderr)
        if args.verbose:
            import traceback
            traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
