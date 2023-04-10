---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_default_patch_baseline"
sidebar_current: "docs-alicloud-resource-oos-default-patch-baseline"
description: |-
  Provides a Alicloud Oos Default Patch Baseline resource.
---

# alicloud_oos_default_patch_baseline

Provides a Oos Default Patch Baseline resource.

For information about Oos Default Patch Baseline and how to use it, see [What is Default Patch Baseline](https://www.alibabacloud.com/help/en/operation-orchestration-service/latest/api-doc-oos-2019-06-01-api-doc-registerdefaultpatchbaseline).

-> **NOTE:** Available in v1.203.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_oos_patch_baseline" "default" {
  operation_system    = "Windows"
  patch_baseline_name = var.name
  description         = var.name
  approval_rules      = "{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}"
}
resource "alicloud_oos_default_patch_baseline" "default" {
  patch_baseline_name = alicloud_oos_patch_baseline.default.id
}
```

## Argument Reference

The following arguments are supported:
* `patch_baseline_name` - (Required,ForceNew) The name of the patch baseline.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `patch_baseline_id` - The ID of the patch baseline.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Default Patch Baseline.
* `delete` - (Defaults to 5 mins) Used when delete the Default Patch Baseline.

## Import

Oos Default Patch Baseline can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_default_patch_baseline.example <id>
```