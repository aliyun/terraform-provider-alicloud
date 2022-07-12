---
subcategory: "EIP Bandwidth Plan (CBWP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_common_bandwidth_package"
sidebar_current: "docs-alicloud-resource-common-bandwidth-package"
description: |-
  Provides a Alicloud Common Bandwidth Package resource.
---

# alicloud\_common_bandwidth_package

Provides a common bandwidth package resource.

-> **NOTE:** Terraform will auto build common bandwidth package instance while it uses `alicloud_common_bandwidth_package` to build a common bandwidth package resource.

For information about common bandwidth package and how to use it, see [What is Common Bandwidth Package](https://www.alibabacloud.com/help/product/55092.htm).

For information about common bandwidth package billing methods, see [Common Bandwidth Package Billing Methods](https://www.alibabacloud.com/help/doc-detail/67459.html).

## Example Usage

Basic Usage

```terraform
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth              = "1000"
  internet_charge_type   = "PayByBandwidth"
  bandwidth_package_name = "test-common-bandwidth-package"
  description            = "test-common-bandwidth-package"
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth of the common bandwidth package. Unit: Mbps.
* `internet_charge_type` - (Optional, ForceNew) The billing method of the common bandwidth package. Valid values are `PayByBandwidth` and `PayBy95` and `PayByTraffic`, `PayByDominantTraffic`. `PayBy95` is pay by classic 95th percentile pricing. International Account doesn't supports `PayByBandwidth` and `PayBy95`. Default to `PayByTraffic`. **NOTE:** From 1.176.0+, `PayByDominantTraffic` is available. 
* `ratio` - (Optional, ForceNew, Available in 1.55.3+) Ratio of the common bandwidth package. It is valid when `internet_charge_type` is `PayBy95`. Default to `100`. Valid values: [10-100].
* `name` - (Optional, Deprecated form v1.120.0) Field `name` has been deprecated from provider version 1.120.0. New field `bandwidth_package_name` instead.
* `bandwidth_package_name` - (Optional, Available in 1.120.0+) The name of the common bandwidth package.
* `description` - (Optional) The description of the common bandwidth package instance.
* `resource_group_id` - (Optional, Available in 1.58.0+, Modifiable in 1.115.0+) The Id of resource group which the common bandwidth package belongs.
* `isp` - (Optional, Available in 1.90.1+) The type of the Internet Service Provider. Valid values: `BGP` and `BGP_PRO`. Default to `BGP`.
* `zone` - (Optional, ForceNew, Available in 1.120.0+) The zone of bandwidth package.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.
* `deletion_protection` - (Optional, Available in v1.124.4+) Whether enable the deletion protection or not. Default value: `false`.
  - true: Enable deletion protection.
  - false: Disable deletion protection.
  
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package instance id.
* `status` - (Available in 1.120.0+) The status of bandwidth package.

### Timeouts

-> **NOTE:** Available in 1.120.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the bandwidth package.
* `delete` - (Defaults to 10 mins) Used when delete the bandwidth package.

## Import

The common bandwidth package can be imported using the id, e.g.

```
$ terraform import alicloud_common_bandwidth_package.foo cbwp-abc123456
```


