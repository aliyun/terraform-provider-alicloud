---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_threat_intelligence_switch"
description: |-
  Provides a Alicloud Cloud Firewall Threat Intelligence Switch resource.
---

# alicloud_cloud_firewall_threat_intelligence_switch

Provides a Cloud Firewall Threat Intelligence Switch resource.

Cloud Firewall Switch Threat Intelligence.

For information about Cloud Firewall Threat Intelligence Switch and how to use it, see [What is Threat Intelligence Switch](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyThreatIntelligenceSwitch).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_firewall_threat_intelligence_switch" "default" {
  action        = "alert"
  enable_status = "0"
  category_id   = "IpOutThreatTorExit"
}
```

### Deleting `alicloud_cloud_firewall_threat_intelligence_switch` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_threat_intelligence_switch`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `action` - (Optional) Rule action. Value:
  - `alert`: Watch
  - `drop`: Intercept
* `category_id` - (Optional, ForceNew, Computed) The threat intelligence classification ID.
* `enable_status` - (Optional, Int) Switch status. Value:
  - `1`: On
  - `0`: Closed

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Threat Intelligence Switch.
* `update` - (Defaults to 5 mins) Used when update the Threat Intelligence Switch.

## Import

Cloud Firewall Threat Intelligence Switch can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_threat_intelligence_switch.example <id>
```