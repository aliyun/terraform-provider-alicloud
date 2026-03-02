---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_patch_baseline"
description: |-
  Provides a Alicloud Operation Orchestration Service (OOS) Patch Baseline resource.
---

# alicloud_oos_patch_baseline

Provides a Operation Orchestration Service (OOS) Patch Baseline resource.



For information about Operation Orchestration Service (OOS) Patch Baseline and how to use it, see [What is Patch Baseline](https://www.alibabacloud.com/help/en/operation-orchestration-service/latest/patch-manager-overview).

-> **NOTE:** Available since v1.146.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_oos_patch_baseline" "default" {
  patch_baseline_name = var.name
  operation_system    = "Windows"
  approval_rules      = "{\"PatchRules\":[{\"EnableNonSecurity\":true,\"PatchFilterGroup\":[{\"Values\":[\"*\"],\"Key\":\"Product\"},{\"Values\":[\"Security\",\"Bugfix\"],\"Key\":\"Classification\"},{\"Values\":[\"Critical\",\"Important\"],\"Key\":\"Severity\"}],\"ApproveAfterDays\":7,\"ComplianceLevel\":\"Unspecified\"}]}"
}
```

## Argument Reference

The following arguments are supported:
* `approval_rules` - (Required) This value follows the json format. For more details, see the description of [ApprovalRules in the Request parameters table for details](https://www.alibabacloud.com/help/zh/operation-orchestration-service/latest/api-oos-2019-06-01-createpatchbaseline).
* `approved_patches` - (Optional, List, Available since v1.219.0) List of approved patches.
* `approved_patches_enable_non_security` - (Optional, Available since v1.219.0) Specifies whether approved patches include updates other than security updates.
* `description` - (Optional) Description of the patch baseline.
* `operation_system` - (Required, ForceNew) Operating system type. Supported values:
  - Windows
  - Ubuntu
  - CentOS
  - Debian
  - AliyunLinux
  - RedhatEnterpriseLinux
  - Anolis
  - Fedora
  - Suse
  - RockyLinux
  - AlmaLinux.
* `patch_baseline_name` - (Required, ForceNew) Patch baseline name.
* `rejected_patches` - (Optional, List, Available since v1.210.0) List of rejected patches.
* `rejected_patches_action` - (Optional, Computed, Available since v1.210.0) Rejected patches action. Valid values: `ALLOW_AS_DEPENDENCY`, `BLOCK`.
* `resource_group_id` - (Optional, Computed, Available since v1.219.0) Resource group ID.
* `sources` - (Optional, List, Available since v1.219.0) List of patch source configurations.
* `tags` - (Optional, Map) A key-value map of tags. The number of key-value pairs must be between 1 and 20.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - Creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Patch Baseline.
* `delete` - (Defaults to 5 mins) Used when delete the Patch Baseline.
* `update` - (Defaults to 5 mins) Used when update the Patch Baseline.

## Import

Operation Orchestration Service (OOS) Patch Baseline can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_patch_baseline.example <patch_baseline_name>
```