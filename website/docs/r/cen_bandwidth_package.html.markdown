---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_package"
sidebar_current: "docs-alicloud-resource-cen-bandwidth-package"
description: |-
  Provides a Alicloud CEN bandwidth package resource.
---

# alicloud\_cen_bandwidth_package

Provides a CEN bandwidth package resource. The CEN bandwidth package is an abstracted object that includes an interconnection bandwidth and interconnection areas. To buy a bandwidth package, you must specify the areas to connect. An area consists of one or more Alibaba Cloud regions. The areas in CEN include Mainland China, Asia Pacific, North America, and Europe.

For information about CEN and how to use it, see [Manage bandwidth packages](https://www.alibabacloud.com/help/doc-detail/65982.htm).

## Example Usage

Basic Usage

```
resource "alicloud_cen_bandwidth_package" "foo" {
    name = "tf-testAccCenBandwidthPackageConfig"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth in Mbps of the bandwidth package. Cannot be less than 1Mbps.
* `geographic_region_ids` - (Required) List of the two areas to connect. Valid value: China | North-America | Asia-Pacific | Europe | Middle-East.
* `name` - (Optional) The name of the bandwidth package. Defaults to null.
* `description` - (Optional) The description of the bandwidth package. Default to null.
* `charge_type` - (Optional) The billing method. Valid value: PostPaid | PrePaid. Default to PostPaid. If set to PrePaid, the bandwidth package can't be deleted before expired time.
* `period` - (Optional) The purchase period in month. Valid value: 1, 2, 3, 6, 12. Default to 1.

~>**NOTE:** PrePaid mode will deduct fees from the account directly. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the bandwidth package.
* `expired_time` - The time of the bandwidth package to expire.
* `status` - The status of the bandwidth, including "InUse" and "Idle".

## Import

CEN bandwidth package can be imported using the id, e.g.

```
$ terraform import alicloud_cen_bandwidth_package.example cenbwp-abc123456
```

