---
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

For information about common bandwidth package billing methods, see [Common Bandwidth Package Billing Methods](https://www.alibabacloud.com/help/doc-detail/67459.html?spm=a2c5t.11065259.1996646101.searchclickresult.7ec93235Vfkwhy).

## Example Usage

Basic Usage

```
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth            = "200"
  internet_charge_type = "PayByBandwidth"
  name                 = "test-common-bandwidth-package"
  description          = "test-common-bandwidth-package"
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth of the common bandwidth package, in Mbps.
* `internet_charge_type` - (Optional, ForceNew) The billing method of the common bandwidth package. Valid values are "PayByBandwidth" and "PayBy95" and "PayByTraffic". "PayBy95" is pay by classic 95th percentile pricing. International Account doesn't supports "PayByBandwidth" and "PayBy95". Default to "PayByTraffic".
* `ratio` - (Optional) Ratio of the common bandwidth package. It is valid when `internet_charge_type` is `PayBy95`. Default to 100. Valid values: [10-100].
* `name` - (Optional) The name of the common bandwidth package.
* `description` - (Optional) The description of the common bandwidth package instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package instance id.

## Import

The common bandwidth package can be imported using the id, e.g.

```
$ terraform import alicloud_common_bandwidth_package.foo cbwp-abc123456
```


