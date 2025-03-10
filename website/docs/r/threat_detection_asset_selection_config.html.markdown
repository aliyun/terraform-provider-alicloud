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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Asset Selection Config.

## Import

Threat Detection Asset Selection Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_asset_selection_config.example <id>
```