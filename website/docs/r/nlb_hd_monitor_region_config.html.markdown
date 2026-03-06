---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_hd_monitor_region_config"
description: |-
  Provides a Alicloud Network Load Balancer (NLB) Hd Monitor Region Config resource.
---

# alicloud_nlb_hd_monitor_region_config

Provides a Network Load Balancer (NLB) Hd Monitor Region Config resource.

HD monitor config.

For information about Network Load Balancer (NLB) Hd Monitor Region Config and how to use it, see [What is Hd Monitor Region Config](https://next.api.alibabacloud.com/document/Nlb/2022-04-30/SetHdMonitorRegionConfig).

-> **NOTE:** Available since v1.273.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}


resource "alicloud_nlb_hd_monitor_region_config" "default" {
  metric_store = "example"
  log_project  = "example"
}
```

## Argument Reference

The following arguments are supported:
* `log_project` - (Required) The name of the LogProject.
* `metric_store` - (Required) The name of the MetricStore.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `region_id` - The ID of the region in which the resource resides.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Hd Monitor Region Config.
* `delete` - (Defaults to 5 mins) Used when delete the Hd Monitor Region Config.
* `update` - (Defaults to 5 mins) Used when update the Hd Monitor Region Config.

## Import

Network Load Balancer (NLB) Hd Monitor Region Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_hd_monitor_region_config.example <region_id>
```