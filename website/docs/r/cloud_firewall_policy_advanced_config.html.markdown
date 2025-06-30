---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_policy_advanced_config"
description: |-
  Provides a Alicloud Cloud Firewall Policy Advanced Config resource.
---

# alicloud_cloud_firewall_policy_advanced_config

Provides a Cloud Firewall Policy Advanced Config resource.

Access Control Advanced Configuration.

For information about Cloud Firewall Policy Advanced Config and how to use it, see [What is Policy Advanced Config](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyPolicyAdvancedConfig).

-> **NOTE:** Available since v1.253.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_firewall_policy_advanced_config" "default" {
  internet_switch = "off"
}
```

### Deleting `alicloud_cloud_firewall_policy_advanced_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_policy_advanced_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `internet_switch` - (Required) Access control policy strict mode of on-state. Valid values:
  - `on`: strict mode enabled.
  - `off`: strict mode is turned off.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as ``.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy Advanced Config.
* `update` - (Defaults to 5 mins) Used when update the Policy Advanced Config.

## Import

Cloud Firewall Policy Advanced Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_policy_advanced_config.example 
```