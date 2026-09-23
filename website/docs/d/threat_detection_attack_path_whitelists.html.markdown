---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_attack_path_whitelists"
sidebar_current: "docs-alicloud-datasource-threat-detection-attack-path-whitelists"
description: |-
  Provides a list of Threat Detection Attack Path Whitelist owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_attack_path_whitelists

This data source provides Threat Detection Attack Path Whitelist available to the user.[What is Attack Path Whitelist](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateAttackPathWhitelist)

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_threat_detection_attack_path_whitelist" "default" {
  path_type      = "role_escalation"
  whitelist_type = "PART_ASSET"
  whitelist_name = var.name
  path_name      = "ecs_get_credential_by_create_login_profile"
  remark         = var.name
  attack_path_asset_list {
    instance_id    = "AliyunYundunSASReadOnlyAccess::System"
    region_id      = "cn-hangzhou"
    vendor         = 0
    asset_type     = 15
    asset_sub_type = 2
    node_type      = "end"
  }
}

data "alicloud_threat_detection_attack_path_whitelists" "default" {
  ids            = ["${alicloud_threat_detection_attack_path_whitelist.default.id}"]
  path_type      = "role_escalation"
  whitelist_name = var.name
}

output "alicloud_threat_detection_attack_path_whitelist_example_id" {
  value = data.alicloud_threat_detection_attack_path_whitelists.default.whitelists.0.id
}
```

## Argument Reference

The following arguments are supported:
* `lang` - (Optional) The language of the request and response. Valid values: `zh` (Chinese, default), `en` (English).
* `path_name_desc` - (Optional) The description of the path name. You can call [ListAvailableAttackPath](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListAvailableAttackPath) to query the path name descriptions.
* `path_type` - (Optional) The path type of the whitelist. You can call [ListAvailableAttackPath](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListAvailableAttackPath) to query the available path types.
* `whitelist_name` - (Optional) The name of the whitelist.
* `ids` - (Optional, Computed) A list of Attack Path Whitelist IDs.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Attack Path Whitelist IDs.
* `whitelists` - A list of Attack Path Whitelist Entries. Each element contains the following attributes:
  * `attack_path_asset_list` - **NOTE:** This field is only available when `enable_details` is `true`. The list of attack path cloud product assets.
    * `asset_sub_type` - The subtype of the cloud product asset.
    * `asset_type` - The type of the cloud product asset.
    * `instance_id` - The instance ID of the cloud product asset.
    * `node_type` - The type of the whitelist node.
    * `region_id` - The region ID of the cloud product asset instance.
    * `vendor` - The vendor of the cloud product asset.
  * `attack_path_whitelist_id` - The ID of the attack path whitelist.
  * `path_name` - The path name of the whitelist.
  * `path_type` - The path type of the whitelist.
  * `remark` - The remarks of the whitelist.
  * `whitelist_name` - The name of the whitelist.
  * `whitelist_type` - The type of the whitelist.
  * `id` - The ID of the resource supplied above.
