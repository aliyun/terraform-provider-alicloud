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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_attack_path_sensitive_asset_config&exampleId=5d019d14-a7d8-7d93-b741-bd7f268b262df34b7e58&activeTab=example&spm=docs.r.threat_detection_attack_path_sensitive_asset_config.0.5d019d14a7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_threat_detection_attack_path_sensitive_asset_config&spm=docs.r.threat_detection_attack_path_sensitive_asset_config.example&intl_lang=EN_US)

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
