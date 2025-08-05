---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_attack_path_sensitive_asset_config"
description: |-
  Provides a Alicloud Threat Detection Attack Path Sensitive Asset Config resource.
---

# alicloud_threat_detection_attack_path_sensitive_asset_config

Provides a Threat Detection Attack Path Sensitive Asset Config resource.

Attack Path Sensitive Asset Configuration.

For information about Threat Detection Attack Path Sensitive Asset Config and how to use it, see [What is Attack Path Sensitive Asset Config](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-createattackpathsensitiveassetconfig).

-> **NOTE:** Available since v1.257.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_slb_load_balancers" "default" {
}

resource "alicloud_threat_detection_attack_path_sensitive_asset_config" "default" {
  attack_path_asset_list {
    instance_id    = data.alicloud_slb_load_balancers.default.balancers.0.id
    vendor         = "0"
    asset_type     = "1"
    asset_sub_type = "0"
    region_id      = "cn-hangzhou"
  }
}
```

## Argument Reference

The following arguments are supported:
* `attack_path_asset_list` - (Required, Set) The attack path sensitive asset configuration list. See [`attack_path_asset_list`](#attack_path_asset_list) below.

### `attack_path_asset_list`

The attack_path_asset_list supports the following:
* `asset_sub_type` - (Required, Int) Cloud product asset subtype.
* `asset_type` - (Required, Int) The asset type of the cloud product asset.
* `instance_id` - (Required) The ID of the cloud product instance.
* `region_id` - (Required) The region ID of the cloud product.
* `vendor` - (Required, Int) Cloud product asset vendor. Valid values: `0`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Attack Path Sensitive Asset Config.
* `delete` - (Defaults to 5 mins) Used when delete the Attack Path Sensitive Asset Config.
* `update` - (Defaults to 5 mins) Used when update the Attack Path Sensitive Asset Config.

## Import

Threat Detection Attack Path Sensitive Asset Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_attack_path_sensitive_asset_config.example <id>
```
