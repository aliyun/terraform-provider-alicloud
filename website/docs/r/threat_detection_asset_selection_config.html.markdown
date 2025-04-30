---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_asset_selection_config"
description: |-
  Provides a Alicloud Threat Detection Asset Selection Config resource.
---

# alicloud_threat_detection_asset_selection_config

Provides a Threat Detection Asset Selection Config resource.

Asset selection configuration.

For information about Threat Detection Asset Selection Config and how to use it, see [What is Asset Selection Config](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateAssetSelectionConfig).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_asset_selection_config&exampleId=73c06543-fe96-80fa-0987-84a880aa8bd43cfc3322&activeTab=example&spm=docs.r.threat_detection_asset_selection_config.0.73c06543fe&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}


resource "alicloud_threat_detection_asset_selection_config" "default" {
  business_type = "agentlesss_vul_white_1"
  target_type   = "instance"
  platform      = "all"
}
```

### Deleting `alicloud_threat_detection_asset_selection_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_asset_selection_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `business_type` - (Required, ForceNew) The first ID of the resource
* `platform` - (Optional, ForceNew) The operating system type.
* `target_type` - (Required, ForceNew) Target object type.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Asset Selection Config.

## Import

Threat Detection Asset Selection Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_asset_selection_config.example <id>
```