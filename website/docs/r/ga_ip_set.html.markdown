---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_ip_set"
sidebar_current: "docs-alicloud-resource-ga-ip-set"
description: |-
  Provides a Alicloud Global Accelerator (GA) Ip Set resource.
---

# alicloud_ga_ip_set

Provides a Global Accelerator (GA) Ip Set resource.

For information about Global Accelerator (GA) Ip Set and how to use it, see [What is Ip Set](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createipsets).

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region  = var.region
  profile = "default"
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_ip_set" "example" {
  accelerate_region_id = var.region
  bandwidth            = "5"
  accelerator_id       = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator (GA) instance.
* `accelerate_region_id` - (Required, ForceNew) The ID of an acceleration region.
* `bandwidth` - (Optional, Int) The bandwidth allocated to the acceleration region.
-> **NOTE:** The minimum bandwidth of each accelerated region is 2Mbps. The total bandwidth of the acceleration region should be less than or equal to the bandwidth of the basic bandwidth package you purchased.
* `ip_version` - (Optional, ForceNew) The IP protocol used by the GA instance. Valid values: `IPv4`, `IPv6`. Default value: `IPv4`.
* `isp_type` - (Optional, ForceNew, Available since v1.207.0) The line type of the elastic IP address (EIP) in the acceleration region. Valid values: `BGP`, `BGP_PRO`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ip Set.
* `ip_address_list` - The list of accelerated IP addresses in the acceleration region.
* `status` -  The status of the acceleration region.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Ip Set.
* `update` - (Defaults to 2 mins) Used when update the Ip Set.
* `delete` - (Defaults to 10 mins) Used when delete the Ip Set.

## Import

Ga Ip Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_ip_set.example <id>
```
