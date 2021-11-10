---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_ip_set"
sidebar_current: "docs-alicloud-resource-ga-ip-set"
description: |-
  Provides a Alicloud Global Accelerator (GA) Ip Set resource.
---

# alicloud\_ga\_ip\_set

Provides a Global Accelerator (GA) Ip Set resource.

For information about Global Accelerator (GA) Ip Set and how to use it, see [What is Ip Set](https://www.alibabacloud.com/help/en/doc-detail/153246.htm).

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
  ratio          = 30
}
resource "alicloud_ga_bandwidth_package_attachment" "example" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.example.id
}
resource "alicloud_ga_ip_set" "example" {
  depends_on           = [alicloud_ga_bandwidth_package_attachment.example]
  accelerate_region_id = "cn-hangzhou"
  bandwidth            = "5"
  accelerator_id       = alicloud_ga_accelerator.example.id
}

```

## Argument Reference

The following arguments are supported:

* `accelerate_region_id` - (Required, ForceNew)  The ID of an acceleration region.
* `accelerator_id` - (Required) The ID of the Global Accelerator (GA) instance.
* `bandwidth` - (Optional) The bandwidth allocated to the acceleration region.

-> **NOTE:** The minimum bandwidth of each accelerated region is 2Mbps. The total bandwidth of the acceleration region should be less than or equal to the bandwidth of the basic bandwidth package you purchased.
                                                                        
* `ip_version` - (Optional, ForceNew) The IP protocol used by the GA instance. Valid values: `IPv4`, `IPv6`. Default value is `IPv4`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ip Set.
* `ip_address_list` - The list of accelerated IP addresses in the acceleration region.
* `status` -  The status of the acceleration region.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Ip Set.
* `delete` - (Defaults to 1 mins) Used when delete the Ip Set.
* `update` - (Defaults to 2 mins) Used when update the Ip Set.

## Import

Ga Ip Set can be imported using the id, e.g.

```
$ terraform import alicloud_ga_ip_set.example <id>
```
