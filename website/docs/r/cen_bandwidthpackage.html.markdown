---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidthpackage"
sidebar_current: "docs-alicloud-cen-bandwidthpackage"
description: |-
  Provides a Alicloud CEN bandwidth package resource.
---

# alicloud\_cen_bandwidthpackage

Provides a CEN bandwidth package resource.

## Example Usage

Basic Usage

```
resource "alicloud_cen_bandwidthpackage" "foo" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"Asia-Pacific"]
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth in Mbps of the bandwidth package. Cannot be less than 1Mbps.
* `geographic_region_id` - (Required) List of the two areas to connect. Valid value: China | North-America | Asia-Pacific | Europe | Middle-East.
* `name` - (Optional) The name of the bandwidth package. Defaults to null.
* `description` - (Optional) The description of the bandwidth package. Default to null.
* `charge_type` - (Optional) The billing method. Valid value: POSTPAY| PREPAY. Default to POSTPAY. If choose PREPAY, must set auto_pay to true.
* `auto_pay` - (Optional) Enable the automatic payment. Valid value: true | false. Default to false.
* `period` - (Optional) The purchase period. Default to 1.
* `pricing_cycle` - (Optional) The pricing cycle. Valid value: Month | Year. Default to Month.

~>**NOTE:** PREPAY mode while auto_pay enabledwill deduct fees from the account directly.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the bandwidth package.
* `bandwidth` - The bandwidth in Mbps of the bandwidth package.
* `geographic_region_id` - List of the two areas to connect.
* `name` - The name of the bandwidth package.
* `description` - The description of the bandwidth package.
* `expired_time` - The time of the bandwidth package to expire.
* `status` - The status of the bandwidth, including "InUse" and "Idle".
* `charge_type` - The billing method.

## Import

CEN bandwidth package can be imported using the id, e.g.

```
$ terraform import alicloud_cen_bandwidthpackage.example cenbwp-abc123456
```

