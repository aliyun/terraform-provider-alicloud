---
name: amp-resource-metadata
description: Fetch resource and API metadata from AMP (镇元平台) and OpenAPI Explorer. Use when querying resource type definitions (resource schema/Meta), retrieving POP product mapping, fetching API-level OpenAPI specifications (request params, response schema, error codes), or getting API documentation with parameter descriptions. Three scripts — get_resource_type.py for resource schema, get_runtime_api.py for raw API metadata, get_api_definition.py for API documentation with descriptions and error codes.
---

# AMP Resource & API Metadata

Three SDK-based scripts for querying AMP (镇元平台) and OpenAPI Explorer. All use AK/SK auth.

| Script | Action | Returns |
|--------|--------|---------|
| `get_resource_type.py` | `GetResourceType` | Resource schema (Meta), PopCode, GatewayType |
| `get_runtime_api.py` | `GetRuntimeApiForInner` | API OpenAPI spec (raw, includes backend/gateway info) |
| `get_api_definition.py` | `GetApiDefinition` | API doc with parameter descriptions, error codes, response schema |

**All paths use `{SKILL_DIR}` = absolute path to this skill's directory.**

## Setup

```bash
bash {SKILL_DIR}/scripts/setup.sh
source {SKILL_DIR}/scripts/.venv/bin/activate
```

Required environment variables:
```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID=xxx
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=xxx
```

**AK 白名单注意**：APISpecData（apispecdata-share 端点）按账号白名单开放。若运行账号（如数字机器人的 terraform_integration）无权限，调用会报 `InvalidApi.NotFound` 404（POP 网关对无权限内部 API 隐藏为 404，不是 403）。此时用专用 AK 覆盖，脚本优先读 `AMP_*`：
```bash
export AMP_ACCESS_KEY_ID=xxx        # 优先级高于 ALIBABA_CLOUD_*
export AMP_ACCESS_KEY_SECRET=xxx
```
数字机器人上该 AK 存放于 `bootstrap/.env`（gitignored，chmod 600），使用前 `set -a; source bootstrap/.env; set +a`。长期方案是为运行账号申请 APISpecData 白名单。

## 1. Resource Schema — `get_resource_type.py`

Get resource-level metadata: PopCode, PopVersion, GatewayType, and full resource schema (Meta with CRUD mappings, attribute definitions).

```bash
python3 {SKILL_DIR}/scripts/get_resource_type.py \
  --service-code VPC --resource-code Vpc
```

| Flag | Required | Description |
|------|----------|-------------|
| `--service-code` | Yes | Namespace, e.g. `VPC`, `ECS` |
| `--resource-code` | Yes | Resource code, e.g. `Vpc`, `Instance` |
| `--env` | No | `pre` or `daily` (default: online) |
| `--output` | No | `pretty` (default) or `json` |

## 2. API Raw Metadata — `get_runtime_api.py`

Get the full OpenAPI specification for a specific API action: internal details including backend service config, gateway options, RAM policy, etc.

```bash
python3 {SKILL_DIR}/scripts/get_runtime_api.py \
  --pop-code Ecs --pop-version 2014-05-26 --api-name DescribeInstances
```

| Flag | Required | Description |
|------|----------|-------------|
| `--pop-code` | Yes | POP product code, e.g. `Ecs`, `Vpc`, `Slb` |
| `--pop-version` | Yes | POP version, e.g. `2014-05-26`, `2016-04-28` |
| `--api-name` | Yes | API action name, e.g. `DescribeInstances`, `CreateVpc` |
| `--runtime-type` | No | Gateway type (default: `pop`) |
| `--env` | No | `online` (default), `pre`, `daily` |
| `--output` | No | `pretty` (default) or `json` |

## 3. API Documentation — `get_api_definition.py`

Get the published API definition from OpenAPI Explorer — the same data shown on the Alibaba Cloud API documentation portal. Includes **parameter descriptions (中文)**, **error codes**, **response field descriptions**.

```bash
python3 {SKILL_DIR}/scripts/get_api_definition.py \
  --product Ecs --api-version 2014-05-26 --api DescribeInstances
```

| Flag | Required | Description |
|------|----------|-------------|
| `--product` | Yes | Product code, e.g. `Ecs`, `Vpc`, `Slb` |
| `--api-version` | Yes | API version, e.g. `2014-05-26`, `2016-04-28` |
| `--api` | Yes | API action name, e.g. `DescribeInstances`, `CreateVpc` |
| `--output` | No | `pretty` (default) or `json` |

**pretty** mode shows: parameter table (name, type, required, description), response fields, error codes.
**json** mode outputs the complete definition for programmatic use.

### Examples

```bash
# ECS DescribeInstances — see all parameters and descriptions
python3 {SKILL_DIR}/scripts/get_api_definition.py \
  --product Ecs --api-version 2014-05-26 --api DescribeInstances

# VPC CreateVpc — check which fields are required
python3 {SKILL_DIR}/scripts/get_api_definition.py \
  --product Vpc --api-version 2016-04-28 --api CreateVpc

# Full JSON for programmatic use
python3 {SKILL_DIR}/scripts/get_api_definition.py \
  --product Ecs --api-version 2014-05-26 --api RunInstances --output json
```

## get_runtime_api vs get_api_definition

| | `get_runtime_api.py` | `get_api_definition.py` |
|---|---|---|
| **Data source** | AMP internal (GetRuntimeApiForInner) | OpenAPI Explorer (GetApiDefinition) |
| **Parameter descriptions** | No | Yes (Chinese) |
| **Error codes** | No | Yes |
| **Response field descriptions** | No | Yes |
| **Backend service info** | Yes (dubbo, appName, timeout) | No |
| **Gateway config** | Yes (RAM auth, akProven, etc.) | No |
| **Use when** | Need internal/infra details | Need API documentation for users/developers |

## Typical Workflow

1. **Know namespace, need PopCode?** → `get_resource_type.py --service-code ECS --resource-code Instance` → read PopCode/PopVersion from output
2. **Need resource schema?** → same command, read Meta field
3. **Need API documentation?** → `get_api_definition.py --product Ecs --api-version 2014-05-26 --api CreateInstance`
4. **Need internal API details?** → `get_runtime_api.py --pop-code Ecs --pop-version 2014-05-26 --api-name CreateInstance`
5. **Check if a field is required?** → use `get_api_definition.py`, look at the Required column
