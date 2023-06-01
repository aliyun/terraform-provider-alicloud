---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_patch_baseline"
sidebar_current: "docs-alicloud-resource-oos-patch-baseline"
description: |-
  Provides a Alicloud OOS Patch Baseline resource.
---

# alicloud\_oos\_patch\_baseline

Provides a OOS Patch Baseline resource.

For information about OOS Patch Baseline and how to use it, see [What is Patch Baseline](https://www.alibabacloud.com/help/en/doc-detail/268700.html).

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_oos_patch_baseline" "example" {
  approval_rules      = "{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}"
  operation_system    = "Windows"
  patch_baseline_name = "terraform-example"
}

```

## Argument Reference

The following arguments are supported:

* `approval_rules` - (Required) Accept the rules. This value follows the json format. For more details, see the [description of ApprovalRules in the Request parameters table for details](https://www.alibabacloud.com/help/zh/doc-detail/311002.html).
* `description` - (Optional) Patches baseline description information.
* `operation_system` - (Required, ForceNew) Operating system type. Valid values: `AliyunLinux`, `Anolis`, `CentOS`, `Debian`, `RedhatEnterpriseLinux`, `Ubuntu`, `Windows`, `AlmaLinux`.
* `patch_baseline_name` - (Required, ForceNew) The name of the patch baseline.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Patch Baseline. Its value is same as `patch_baseline_name`.

## Import

OOS Patch Baseline can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_patch_baseline.example <patch_baseline_name>
```