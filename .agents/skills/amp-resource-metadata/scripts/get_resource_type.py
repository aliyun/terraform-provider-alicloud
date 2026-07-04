#!/usr/bin/env python3
"""
AMP Resource Metadata Fetcher

Calls the GetResourceType API on AMP (镇元平台) to retrieve resource metadata
using Alibaba Cloud Python SDK V2.0 generalized call.

Usage:
    python3 get_resource_type.py --service-code <ServiceCode> --resource-code <ResourceCode> [--env <Env>]

Author: AMP Tools Team
"""

import os
import sys
import json
import argparse
from typing import Dict, Any, Optional

from alibabacloud_tea_openapi import models as open_api_models
from alibabacloud_tea_openapi.client import Client as OpenApiClient
from alibabacloud_tea_util import models as util_models


# ---------------------------------------------------------------------------
# Constants derived from apischema.json
# ---------------------------------------------------------------------------
API_VERSION = "2025-02-27"          # from namespace APISpecData::2025-02-27
API_ACTION = "GetResourceType"
API_METHOD = "POST"
API_STYLE = "RPC"
DEFAULT_ENDPOINT = "apispecdata-share.cn-zhangjiakou.aliyuncs.com"


# ---------------------------------------------------------------------------
# Client helper
# ---------------------------------------------------------------------------

def create_client(endpoint: str = DEFAULT_ENDPOINT) -> OpenApiClient:
    """
    Create an OpenAPI client using AK credentials from environment variables.

    Required environment variables (AMP_* takes precedence; APISpecData 按账号
    白名单开放，运行账号无权限时用 AMP_* 指定专用 AK):
        AMP_ACCESS_KEY_ID / ALIBABA_CLOUD_ACCESS_KEY_ID
        AMP_ACCESS_KEY_SECRET / ALIBABA_CLOUD_ACCESS_KEY_SECRET

    Args:
        endpoint: The service endpoint (default: apispecdata.cn-hangzhou.aliyuncs.com)

    Returns:
        An initialised OpenApiClient instance.
    """
    access_key_id = os.environ.get("AMP_ACCESS_KEY_ID") or os.environ.get(
        "ALIBABA_CLOUD_ACCESS_KEY_ID")
    access_key_secret = os.environ.get("AMP_ACCESS_KEY_SECRET") or os.environ.get(
        "ALIBABA_CLOUD_ACCESS_KEY_SECRET")

    if not access_key_id or not access_key_secret:
        print(
            "Error: Please set AMP_ACCESS_KEY_ID/AMP_ACCESS_KEY_SECRET (or "
            "ALIBABA_CLOUD_ACCESS_KEY_ID/ALIBABA_CLOUD_ACCESS_KEY_SECRET) "
            "environment variables.",
            file=sys.stderr,
        )
        sys.exit(1)

    config = open_api_models.Config(
        access_key_id=access_key_id,
        access_key_secret=access_key_secret,
    )
    config.endpoint = endpoint
    config.protocol = "HTTPS"
    config.read_timeout = 10000
    config.connect_timeout = 5000

    return OpenApiClient(config)


# ---------------------------------------------------------------------------
# API call
# ---------------------------------------------------------------------------

def get_resource_type(
    client: OpenApiClient,
    service_code: str,
    resource_code: str,
    env: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Call GetResourceType on AMP to retrieve resource metadata.

    Args:
        client:        Initialised OpenApiClient.
        service_code:  The service code (e.g. "VPC", "ECS").
        resource_code: The resource code (e.g. "Vpc", "Instance").
        env:           Optional environment identifier (e.g. "pre", "daily").

    Returns:
        The full API response dict with keys: body, headers, statusCode.
    """
    # -- 1. Configure API params (RPC style) --------------------------------
    params = open_api_models.Params(
        style=API_STYLE,
        version=API_VERSION,
        action=API_ACTION,
        method=API_METHOD,
        pathname="/",
        protocol="HTTPS",
        auth_type="AK",
        req_body_type="formData",
        body_type="json",
    )

    # -- 2. Build formData body ---------------------------------------------
    body: Dict[str, Any] = {
        "ServiceCode": service_code,
        "ResourceCode": resource_code,
    }
    if env:
        body["Env"] = env

    request = open_api_models.OpenApiRequest(body=body)

    # -- 3. Runtime options -------------------------------------------------
    runtime = util_models.RuntimeOptions(
        autoretry=True,
        max_attempts=3,
    )

    # -- 4. Execute ---------------------------------------------------------
    response = client.call_api(params, request, runtime)
    return response


# ---------------------------------------------------------------------------
# Pretty-print helpers
# ---------------------------------------------------------------------------

def print_metadata(response: Dict[str, Any], verbose: bool = False) -> None:
    """
    Print the resource metadata in a human-readable format.

    Args:
        response: The raw API response.
        verbose:  If True, print the full JSON response.
    """
    if verbose:
        print(json.dumps(response, indent=2, ensure_ascii=False))
        return

    body = response.get("body", {})
    success = body.get("Success", body.get("success"))
    if not success:
        code = body.get("Code", body.get("code", "Unknown"))
        message = body.get("Message", body.get("message", "Unknown error"))
        print(f"API call failed — Code: {code}, Message: {message}", file=sys.stderr)
        sys.exit(1)

    data = body.get("Data", body.get("data", {}))
    if not data:
        print("No data returned.", file=sys.stderr)
        sys.exit(1)

    fields = [
        ("ServiceCode",  "服务 Code"),
        ("PopCode",      "POP 产品"),
        ("PopVersion",   "POP 版本"),
        ("ResourceName", "资源名称"),
        ("GatewayType",  "网关类型"),
        ("Env",          "环境"),
        ("Description",  "描述"),
    ]

    print("=" * 60)
    print("  AMP Resource Metadata")
    print("=" * 60)
    for key, label in fields:
        value = data.get(key, "N/A")
        if value:
            print(f"  {label:<12} : {value}")

    meta_raw = data.get("Meta")
    if meta_raw:
        print("-" * 60)
        print("  Meta (资源 Schema):")
        try:
            meta = json.loads(meta_raw) if isinstance(meta_raw, str) else meta_raw
            print(json.dumps(meta, indent=4, ensure_ascii=False))
        except (json.JSONDecodeError, TypeError):
            print(f"  {meta_raw}")
    print("=" * 60)


# ---------------------------------------------------------------------------
# CLI entry point
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(
        description="Fetch resource metadata from AMP (镇元平台) via GetResourceType API.",
        epilog=(
            "Required environment variables:\n"
            "  ALIBABA_CLOUD_ACCESS_KEY_ID       Alibaba Cloud Access Key ID\n"
            "  ALIBABA_CLOUD_ACCESS_KEY_SECRET   Alibaba Cloud Access Key Secret\n"
            "\n"
            "Examples:\n"
            "  %(prog)s --service-code VPC --resource-code Vpc\n"
            "  %(prog)s --service-code ECS --resource-code Instance --env pre\n"
            "  %(prog)s --service-code VPC --resource-code Vpc --output json\n"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--service-code",
        required=True,
        help="Service code, e.g. VPC, ECS, OOS",
    )
    parser.add_argument(
        "--resource-code",
        required=True,
        help="Resource code, e.g. Vpc, Instance",
    )
    parser.add_argument(
        "--env",
        default=None,
        help="Environment identifier (optional), e.g. pre, daily",
    )
    parser.add_argument(
        "--endpoint",
        default=DEFAULT_ENDPOINT,
        help=f"Service endpoint (default: {DEFAULT_ENDPOINT})",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Print full JSON response",
    )
    parser.add_argument(
        "--output",
        choices=["pretty", "json"],
        default="pretty",
        help="Output format (default: pretty)",
    )

    args = parser.parse_args()

    try:
        client = create_client(endpoint=args.endpoint)
        response = get_resource_type(
            client=client,
            service_code=args.service_code,
            resource_code=args.resource_code,
            env=args.env,
        )

        if args.output == "json":
            print(json.dumps(response, indent=2, ensure_ascii=False))
        else:
            print_metadata(response, verbose=args.verbose)

    except Exception as exc:
        print(f"Error: {exc}", file=sys.stderr)
        if args.verbose:
            import traceback
            traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
