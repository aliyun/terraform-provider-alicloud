#!/usr/bin/env python3
"""
OpenAPI Explorer — GetApiDefinition

Calls the GetApiDefinition API (OpenAPIExplorer/2024-11-30) to retrieve the
full open metadata for a specific cloud API action, including:
  - Request parameters (name, type, required, description)
  - Response schema
  - Error codes
  - System tags, deprecation status, etc.

This is the same data that powers the Alibaba Cloud API Explorer (OpenAPI门户).

Usage:
    python3 get_api_definition.py --product Ecs --api-version 2014-05-26 --api DescribeInstances
    python3 get_api_definition.py --product Vpc --api-version 2016-04-28 --api CreateVpc --output json

Environment variables:
    ALIBABA_CLOUD_ACCESS_KEY_ID
    ALIBABA_CLOUD_ACCESS_KEY_SECRET

Author: AMP Tools Team
"""

import os
import sys
import json
import argparse
from typing import Dict, Any

from alibabacloud_tea_openapi import models as open_api_models
from alibabacloud_tea_openapi.client import Client as OpenApiClient
from alibabacloud_tea_util import models as util_models


# ---------------------------------------------------------------------------
# Constants — from OpenAPIExplorer/2024-11-30::GetApiDefinition
# ---------------------------------------------------------------------------
API_VERSION = "2024-11-30"
API_ACTION = "GetApiDefinition"
API_METHOD = "GET"
API_STYLE = "ROA"
API_PATHNAME = "/api/definition"
DEFAULT_ENDPOINT = "openapiexplorer.aliyuncs.com"


# ---------------------------------------------------------------------------
# Client helper
# ---------------------------------------------------------------------------

def create_client(endpoint: str = DEFAULT_ENDPOINT) -> OpenApiClient:
    """Create an OpenAPI client using AK credentials from environment variables."""
    access_key_id = os.environ.get("AMP_ACCESS_KEY_ID") or os.environ.get(
        "ALIBABA_CLOUD_ACCESS_KEY_ID")
    access_key_secret = os.environ.get("AMP_ACCESS_KEY_SECRET") or os.environ.get(
        "ALIBABA_CLOUD_ACCESS_KEY_SECRET")

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
    config.read_timeout = 60000
    config.connect_timeout = 10000

    return OpenApiClient(config)


# ---------------------------------------------------------------------------
# API call
# ---------------------------------------------------------------------------

def get_api_definition(
    client: OpenApiClient,
    product: str,
    api_version: str,
    api: str,
) -> Dict[str, Any]:
    """
    Call GetApiDefinition to retrieve the full API metadata.

    Args:
        client:      Initialised OpenApiClient.
        product:     Product code, e.g. "Ecs", "Vpc", "Slb".
        api_version: API version, e.g. "2014-05-26".
        api:         API action name, e.g. "DescribeInstances".

    Returns:
        The full API response dict.
    """
    params = open_api_models.Params(
        style=API_STYLE,
        version=API_VERSION,
        action=API_ACTION,
        method=API_METHOD,
        pathname=API_PATHNAME,
        protocol="HTTPS",
        auth_type="AK",
        req_body_type="json",
        body_type="json",
    )

    query = {
        "product": product,
        "apiVersion": api_version,
        "api": api,
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

def print_result(response: Dict[str, Any], output: str, verbose: bool) -> None:
    """Print the API definition."""
    body = response.get("body", {})

    if verbose:
        print(json.dumps(response, indent=2, ensure_ascii=False))
        return

    # Check for error
    if isinstance(body, dict) and body.get("code") and body.get("code") != 0:
        print(f"API error: {body.get('message', 'Unknown')}", file=sys.stderr)
        sys.exit(1)

    if output == "json":
        print(json.dumps(body, indent=2, ensure_ascii=False))
    else:
        _pretty_print(body)


def _pretty_print(body: Dict[str, Any]) -> None:
    """Human-readable summary of the API definition."""
    print("=" * 70)
    print("  API Definition (OpenAPI Explorer)")
    print("=" * 70)

    # Basic info
    methods = body.get("methods", [])
    deprecated = body.get("deprecated", False)
    op_type = body.get("operationType", "N/A")
    sys_tags = body.get("systemTags", {})

    print(f"  Methods    : {', '.join(methods)}")
    print(f"  Operation  : {op_type}")
    print(f"  Deprecated : {deprecated}")
    if sys_tags:
        print(f"  Charge     : {sys_tags.get('chargeType', 'N/A')}")
        print(f"  Risk       : {sys_tags.get('riskType', 'N/A')}")

    # Parameters
    params = body.get("parameters", [])
    # Filter to query/body params (skip system params)
    user_params = [p for p in params if p.get("in") in ("query", "body", "path", "formData")]

    if user_params:
        print(f"\n  Parameters ({len(user_params)} user-facing, {len(params)} total):")
        print(f"  {'Name':<35s} {'Type':<12s} {'Required':<10s} {'In':<10s} Description")
        print(f"  {'-'*35} {'-'*12} {'-'*10} {'-'*10} {'-'*30}")
        for p in user_params:
            name = p.get("name", "?")
            schema = p.get("schema", {})
            p_type = schema.get("type", "string")
            required = "required" if schema.get("required") or p.get("required") else "optional"
            in_loc = p.get("in", "query")
            desc = (p.get("description") or schema.get("description") or "")
            desc = desc.replace("\n", " ")[:60]
            print(f"  {name:<35s} {p_type:<12s} {required:<10s} {in_loc:<10s} {desc}")

    # Responses
    responses = body.get("responses", {})
    if responses:
        print(f"\n  Responses:")
        for code, resp in responses.items():
            schema = resp.get("schema", {})
            props = schema.get("properties", {})
            print(f"    {code}: {len(props)} fields")
            for prop_name, prop_def in list(props.items())[:20]:
                prop_type = prop_def.get("type", "?")
                prop_desc = (prop_def.get("description") or "").replace("\n", " ")[:50]
                print(f"      - {prop_name:<30s} {prop_type:<10s} {prop_desc}")

    # Error codes
    error_codes = body.get("errorCodes", {})
    if error_codes:
        total = sum(len(v) if isinstance(v, list) else 1 for v in error_codes.values())
        print(f"\n  Error Codes ({total}):")
        for code, details in list(error_codes.items())[:10]:
            if isinstance(details, list):
                for d in details[:3]:
                    msg = d.get("errorMessage", "")[:60] if isinstance(d, dict) else str(d)[:60]
                    print(f"    {code}: {msg}")
            else:
                print(f"    {code}: {details}")
        if len(error_codes) > 10:
            print(f"    ... and {len(error_codes) - 10} more")

    print("\n" + "=" * 70)
    print("Tip: Use --output json for the complete definition.")


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(
        description="Fetch API definition from OpenAPI Explorer (GetApiDefinition 2024-11-30).",
        epilog=(
            "Required environment variables:\n"
            "  ALIBABA_CLOUD_ACCESS_KEY_ID       Alibaba Cloud Access Key ID\n"
            "  ALIBABA_CLOUD_ACCESS_KEY_SECRET   Alibaba Cloud Access Key Secret\n"
            "\n"
            "Examples:\n"
            "  %(prog)s --product Ecs --api-version 2014-05-26 --api DescribeInstances\n"
            "  %(prog)s --product Vpc --api-version 2016-04-28 --api CreateVpc --output json\n"
            "  %(prog)s --product Slb --api-version 2014-05-15 --api CreateLoadBalancer\n"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--product", required=True, help="Product code, e.g. Ecs, Vpc, Slb")
    parser.add_argument("--api-version", required=True, help="API version, e.g. 2014-05-26")
    parser.add_argument("--api", required=True, help="API action name, e.g. DescribeInstances")
    parser.add_argument("--endpoint", default=DEFAULT_ENDPOINT, help=f"Endpoint (default: {DEFAULT_ENDPOINT})")
    parser.add_argument("--output", choices=["pretty", "json"], default="pretty", help="Output format")
    parser.add_argument("--verbose", "-v", action="store_true", help="Print full raw response")

    args = parser.parse_args()

    try:
        client = create_client(endpoint=args.endpoint)
        response = get_api_definition(
            client=client,
            product=args.product,
            api_version=args.api_version,
            api=args.api,
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
