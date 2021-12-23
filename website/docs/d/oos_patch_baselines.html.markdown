---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_patch_baselines"
sidebar_current: "docs-alicloud-datasource-oos-patch-baselines"
description: |-
  Provides a list of Oos Patch Baselines to the user.
---

# alicloud\_oos\_patch\_baselines

This data source provides the Oos Patch Baselines of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_patch_baselines" "ids" {
  ids = ["my-PatchBaseline-1", "my-PatchBaseline-2"]
}
output "oos_patch_baseline_id_1" {
  value = data.alicloud_oos_patch_baselines.ids.baselines.0.id
}

data "alicloud_oos_patch_baselines" "nameRegex" {
  name_regex = "^my-PatchBaseline"
}
output "oos_patch_baseline_id_2" {
  value = data.alicloud_oos_patch_baselines.nameRegex.baselines.0.id
}

data "alicloud_oos_patch_baselines" "shareType" {
  ids        = ["my-PatchBaseline-1"]
  share_type = "Private"
}
output "oos_patch_baseline_id_3" {
  value = data.alicloud_oos_patch_baselines.shareType.baselines.0.id
}

data "alicloud_oos_patch_baselines" "shareType" {
  ids              = ["my-PatchBaseline-1"]
  operation_system = "AliyunLinux"
}
output "oos_patch_baseline_id_4" {
  value = data.alicloud_oos_patch_baselines.shareType.baselines.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Patch Baseline IDs. Its element value is same as Patch Baseline Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Patch Baseline name.
* `operation_system` - (Optional, ForceNew) Operating system type. Valid values: `AliyunLinux`, `Anolis`, `Centos`, `Debian`, `RedhatEnterpriseLinux`, `Ubuntu`, `Windows`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `share_type` - (Optional, ForceNew) Patch baseline sharing type. Valid values: `Private`, `Public`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Patch Baseline names.
* `baselines` - A list of Oos Patch Baselines. Each element contains the following attributes:
	* `approval_rules` - Accept the rules.
	* `create_time` - The create time of patch baselines.
	* `created_by` - The user who created the patch baselines.
	* `description` - Patches baseline description information.
	* `id` - The ID of the Patch Baseline. Its value is same as `patch_baseline_name`.
	* `is_default` - Whether it is the default patch baseline.
	* `operation_system` - Operating system type.
	* `patch_baseline_id` - Patch baseline ID.
	* `patch_baseline_name` - The name of the patch baseline.
	* `share_type` - Patch baseline sharing type.
	* `updated_by` - The user who updated the patch baselines.
	* `updated_date` - The update time of patch baselines.