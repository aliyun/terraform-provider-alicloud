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
  approval_rules      = "example_value"
  operation_system    = "Windows"
  patch_baseline_name = "my-PatchBaseline"
}

```

## Argument Reference

The following arguments are supported:

* `approval_rules` - (Required) Accept the rules. This value follows the json format. For more details, see the [description of ApprovalRules in the Request parameters table for details](https://www.alibabacloud.com/help/zh/doc-detail/311002.html).
* `description` - (Optional) Patches baseline description information.
* `operation_system` - (Required, ForceNew) Operating system type. Valid values: `AliyunLinux`, `Anolis`, `Centos`, `Debian`, `RedhatEnterpriseLinux`, `Ubuntu`, `Windows`.
* `patch_baseline_name` - (Required, ForceNew) The name of the patch baseline.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Patch Baseline. Its value is same as `patch_baseline_name`.

## Import

OOS Patch Baseline can be imported using the id, e.g.

```
$ terraform import alicloud_oos_patch_baseline.example <patch_baseline_name>
```