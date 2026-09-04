# Estimate Cost for Terraform Provider for Alibaba Cloud

> Bring **plan-time price estimation** to your `terraform plan` workflow, with
> zero changes to the provider itself.
>
> ```
> $ terraform plan --estimate-cost -out=tfplan
> [normal terraform plan output ...]
>
> ─────────────────────────────────────────────────────────────
>                        Cost Estimate
> ─────────────────────────────────────────────────────────────
> RESOURCE                       AMOUNT     CURRENCY  MODE        NOTE
> alicloud_instance.demo         235.70     CNY       create      create
> ── TOTAL ──                    235.70     CNY
> ```

## What this is

A pair of small Go binaries that **wrap Terraform** to add a `--estimate-cost`
flag to `terraform plan`. When enabled, the tool reads the plan JSON, translates
each resource change into an Alibaba Cloud OpenAPI call via [CloudControl's
`GetApiPrice`](https://help.aliyun.com/document_detail/cloudcontrol-getapiprice),
and renders a cost table at the end of the plan output.

Provider code is **not modified** in any way. The wrapper transparently
forwards every other Terraform subcommand to the real `terraform` binary
(`init`, `apply`, `state`, `destroy`, etc.).

## How it ships

Built and released **alongside** the provider via `goreleaser`. Every
provider tag (`v1.282.0`, etc.) automatically produces a matching estimate-cost
archive on the same GitHub Release:

```
alicloud-tf-estimate-cost_1.282.0_darwin_arm64.tar.gz
alicloud-tf-estimate-cost_1.282.0_darwin_amd64.tar.gz
alicloud-tf-estimate-cost_1.282.0_linux_amd64.tar.gz
alicloud-tf-estimate-cost_1.282.0_linux_arm64.tar.gz
alicloud-tf-estimate-cost_1.282.0_windows_amd64.zip
```

The mapping files (schema ↔ OpenAPI translation rules) are **embedded inside
the `estimate-cost` binary** via `go:embed`. Users do **not** need to download
or configure any additional files.

## Install

### Manual

1. Download the archive matching your platform from the GitHub Release page.
2. Extract `terraform-wrapper` and `estimate-cost` into a directory on your PATH
   that comes **before** the real `terraform`'s directory (e.g. `~/bin/` if
   that is first in your PATH).
3. Rename `terraform-wrapper` → `terraform`. Keep `estimate-cost` as a sibling
   file (the wrapper looks for it next to itself).

```bash
tar -xzf alicloud-tf-estimate-cost_1.282.0_darwin_arm64.tar.gz
mv terraform-wrapper ~/bin/terraform
mv estimate-cost     ~/bin/estimate-cost
chmod +x ~/bin/terraform ~/bin/estimate-cost
```

Verify:
```bash
which terraform                 # should be ~/bin/terraform
terraform plan --estimate-cost  # if this triggers cost output, install succeeded
```

### Uninstall

```bash
rm ~/bin/terraform ~/bin/estimate-cost
# Real terraform (e.g. /opt/homebrew/bin/terraform) becomes active again.
```

## Usage

The wrapper is **fully transparent** for every Terraform command:

```bash
terraform init                     # forwarded to real terraform
terraform plan                     # forwarded to real terraform
terraform apply tfplan             # forwarded to real terraform
terraform plan --estimate-cost     # ← intercepted, adds cost table
```

`--estimate-cost` works on both create and update plans:

```bash
# Cost of creating a new resource
terraform plan --estimate-cost -out=tfplan

# Cost delta for a modification (refresh skipped for speed)
terraform plan --estimate-cost -refresh=false -out=tfplan
```

For `update` actions the tool calls the appropriate Modify-style OpenAPI
(`ModifyPrepayInstanceSpec` for PrePaid spec changes, `ResizeDisk` for disk
expansion, etc.) and CC API returns the **real delta amount** — the same number
the actual upgrade order will show.

## Authentication

The tool uses the same Alibaba Cloud credentials as your Terraform provider.
It signs requests with POP v3 (`ACS3-HMAC-SHA256`) **locally**; your AK/SK
never leaves your machine.

Set one of these env var pairs (checked in order):

| Var pair | Notes |
|---|---|
| `ALIBABA_CLOUD_ACCESS_KEY_ID` / `ALIBABA_CLOUD_ACCESS_KEY_SECRET` | preferred (matches `aliyun` CLI) |
| `ALICLOUD_ACCESS_KEY` / `ALICLOUD_SECRET_KEY` | fallback (matches Terraform provider) |

### Endpoint

By default, requests go to the CloudControl endpoint
`cloudcontrol.aliyuncs.com`. Two env vars let you redirect them to a
different endpoint when needed (e.g. for testing):

| Var | Default | Notes |
|---|---|---|
| `ALICLOUD_CC_API_ENDPOINT` | `cloudcontrol.aliyuncs.com` | Override the request hostname. |
| `ALICLOUD_CC_API_HOST` | `cloudcontrol.aliyuncs.com` | Override the HTTP Host header. Rarely needed — only when a custom gateway requires it. |

## Architecture

```
                          ┌──────────────────────┐
                          │   internal/embed-    │
                          │   mappings/data/     │
                          │   alicloud_*.json    │  ← single source of truth
                          └──────────┬───────────┘
                                     │ go:embed
                                     ▼
$ terraform plan --estimate-cost     ┌──────────────────────────────────┐
                  │                  │  estimate-cost binary             │
                  ▼                  │   - reads embedded mappings       │
        ┌─────────────────────┐ exec │   - parses plan.json              │
        │ terraform-wrapper   │─────►│   - builds OpenAPI params         │
        │   - intercept       │      │   - calls CC API GetApiPrice      │
        │   - run real plan   │      │   - aggregates result             │
        │   - show -json      │      └──────────────────────────────────┘
        │   - call estimate-cost                  │
        │   - render output   │                   │ HTTPS + POP v3 sig
        └─────────────────────┘                   ▼
                                       Alibaba Cloud CloudControl
                                       (GetApiPrice → BSS)
```

### `plan.json` stays local

The wrapper extracts only the **billing-relevant fields** from `plan.json` and
sends them as standard OpenAPI parameters. Sensitive content (`user_data`,
`password`, `tags`, internal endpoints) is never transmitted.

## Updating mappings

The `internal/embedmappings/data/*.json` files describe how each provider
schema attribute maps to an OpenAPI parameter. They are **the single source
of truth** for the price-mapping layer and live in this same repository so
they evolve with the provider's schema in lockstep.

To add cost coverage for a new resource:

1. Add `internal/embedmappings/data/alicloud_<resource>.json` following the
   format of `alicloud_instance.json`.
2. The binary will pick it up on the next build (`go:embed all:data`).
3. Add a brief test plan if possible.

Target-level matching fields:

| Field | Meaning |
|---|---|
| `actions` | Match against the plan's `change.actions` (`create` / `update` / ...). |
| `when` | All `{path: value}` pairs must hold. `$.x` reads `change.after`, `$before.x` reads `change.before` — the latter enables "changed from A to B" dispatch (e.g. a PostPaid → PrePaid switch). |
| `whenChanged` | At least one listed field must differ between before and after. Prevents a target from firing when its driving field didn't change. |

Param-level extras beyond `from`/`const`/`default`:

| Field | Meaning |
|---|---|
| `wrapArray` | Wrap the resolved value in a single-element list (for list-typed inputs like `InstanceIds`). |
| `PricingContext.*` keys | Dotted `PricingContext.` param names are automatically nested into a structured `PricingContext` JSON object — GetApiPrice requires that form for pricing hints such as `CurrentDisk` (unlike POP params like `SystemDisk.Size`, which stay flat). |

## Development

```bash
# Build both binaries to the repo's bin/ for local testing
go build -o bin/estimate-cost ./estimate-cost/cmd/estimate-cost/
go build -o bin/terraform     ./estimate-cost/cmd/terraform-wrapper/

# Override embedded mapping with a local directory (dev / debug only)
TF_COST_MAPPINGS=$(pwd)/estimate-cost/internal/embedmappings/data \
  terraform plan --estimate-cost -out=tfplan
```

## What's covered today

### `alicloud_instance`

| Plan change | OpenAPI called | Pricing mode |
|---|---|---|
| Create (PrePaid or PostPaid) | `RunInstances` | single |
| PrePaid: change `instance_type` | `ModifyPrepayInstanceSpec` | delta |
| PrePaid: increase `system_disk_size` | `ResizeDisk` | delta |
| PostPaid: change `instance_type` | `ModifyInstanceSpec` | delta |
| Switch to `PayByBandwidth` / raise `internet_max_bandwidth_out` | `ModifyInstanceNetworkSpec` | single (new-bandwidth order) |
| PostPaid → PrePaid `instance_charge_type` switch | `ModifyInstanceChargeType` | single (first period) |
| Change `image_id` | `ReplaceSystemDisk` | delta |
| Change `system_disk_performance_level` / `system_disk_provisioned_iops` | `ModifyDiskSpec` | delta |
| Any update not matched above | `RunInstances` × 2 | diff (after − before) |

A single update may match more than one delta target — e.g. PrePaid
changing both `instance_type` and `system_disk_size` returns two orders
(`ModifyPrepayInstanceSpec` + `ResizeDisk`) summed into one total.

### Intentionally not covered

| Case | Why |
|---|---|
| `ResizeDisk` on PostPaid disks | GetApiPrice's rule requires `DiskChargeType='PrePaid'`; pay-as-you-go disks have no upfront order to quote. Falls back to the create-target diff. |
| PrePaid → PostPaid charge-type switch | A refund, not a purchase — GetApiPrice returns `PRICING_PLAN_RESULT_NOT_FOUND`. Falls back to diff. |
| `AllocatePublicIpAddress` | Priceable by GetApiPrice, but deliberately skipped: the bandwidth cost is already captured by the `ModifyInstanceNetworkSpec` target and a separate target would double-count. |
| Bandwidth changes while staying on `PayByTraffic` | Pricing needs an estimated traffic volume (`PricingContext.EstimatedInternetTrafficOutGB`) which a Terraform plan cannot provide. |

More resources will be added in subsequent provider releases.

## Limitations

- **Unknown fields**: if a pricing input depends on a `(known after apply)` value
  (e.g. `region = other_resource.id`), the tool prints `known after apply`
  rather than fabricating a price.
- **No stock check**: a plan that would fail at `apply` time with
  `OperationDenied.NoStock` is still priced. Stock is an orthogonal concern.
- **Endpoint override**: by default the tool targets the CloudControl endpoint
  `cloudcontrol.aliyuncs.com`. To target a different endpoint, set
  `ALICLOUD_CC_API_ENDPOINT` and, if needed, `ALICLOUD_CC_API_HOST`.

## License

Same as the parent provider repository (see `LICENSE` at the repo root).
