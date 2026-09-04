---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_attack_path_whitelist"
description: |-
  Provides a Alicloud Threat Detection Attack Path Whitelist resource.
---

# alicloud_threat_detection_attack_path_whitelist

Provides a Threat Detection Attack Path Whitelist resource.

Attack Path Whitelist.

For information about Threat Detection Attack Path Whitelist and how to use it, see [What is Attack Path Whitelist](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateAttackPathWhitelist).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

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
  whitelist_name = "example-1"
  path_name      = "ecs_get_credential_by_create_login_profile"
  remark         = "example-1"
  attack_path_asset_list {
    instance_id    = "AliyunYundunSASReadOnlyAccess::System"
    region_id      = "cn-hangzhou"
    vendor         = 0
    asset_type     = 15
    asset_sub_type = 2
    node_type      = "end"
  }
}
```

## Argument Reference

The following arguments are supported:
* `attack_path_asset_list` - (Optional, List) The list of attack path cloud product assets. See [`attack_path_asset_list`](#attack_path_asset_list) below.
* `path_name` - (Required) The path name of the whitelist. You can call [ListAvailableAttackPath](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListAvailableAttackPath) to query the available path names.
* `path_type` - (Required) The path type of the whitelist. You can call [ListAvailableAttackPath](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListAvailableAttackPath) to query the available path types.
* `remark` - (Optional) The remarks of the whitelist.
* `whitelist_name` - (Required) The name of the whitelist.
* `whitelist_type` - (Required) The type of the whitelist. Valid values: `ALL_ASSET` (all assets), `PART_ASSET` (partial assets).

### `attack_path_asset_list`

The attack_path_asset_list supports the following:
* `asset_sub_type` - (Optional, Int) The subtype of the cloud product asset. You can call [ListCloudAssetInstances](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListCloudAssetInstances) to query the asset subtypes.
* `asset_type` - (Optional, Int) The type of the cloud product asset. You can call [ListCloudAssetInstances](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListCloudAssetInstances) to query the asset types.
* `instance_id` - (Optional) The instance ID of the cloud product asset. You can call [ListCloudAssetInstances](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListCloudAssetInstances) to query the instance IDs.
* `node_type` - (Optional) The type of the whitelist node. Valid values: `start` (starting point), `end` (end point).
* `region_id` - (Optional) The region ID of the cloud product asset instance.
* `vendor` - (Required, Int) The vendor of the cloud product asset. You can call [ListCloudAssetInstances](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListCloudAssetInstances) to query the vendors.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Attack Path Whitelist.
* `delete` - (Defaults to 5 mins) Used when delete the Attack Path Whitelist.
* `update` - (Defaults to 5 mins) Used when update the Attack Path Whitelist.

## Import

Threat Detection Attack Path Whitelist can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_attack_path_whitelist.example <attack_path_whitelist_id>
```