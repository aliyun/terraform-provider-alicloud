---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_service"
sidebar_current: "docs-alicloud-datasource-cdn-service"
description: |-
    Provides a datasource to open the CDN service automatically.
---

# alicloud\_cdn\_service

Using this data source can enable CDN service automatically. If the service has been enabled, it will return `Opened`.

For information about CDN and how to use it, see [What is CDN](https://www.alibabacloud.com/help/product/27099.htm).

-> **NOTE:** Available in v1.98.0+

## Example Usage

```
data "alicloud_cdn_service" "open" {
	enable               = "On"
	internet_charge_type = "PayByTraffic"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".
* `internet_charge_type` - (Optional) The new billing method. Valid values: `PayByTraffic` and `PayByBandwidth`. Default value: `PayByTraffic`.
It is required when `enable = on`. If the CDN service has been opened and you can update its internet charge type by modifying the filed `internet_charge_type`. 
As a note, the updated internet charge type will be effective in the next day zero time.

-> **NOTE:** Setting `enable = "On"` to open the CDN service that means you have read and agreed the [CDN Terms of Service](https://help.aliyun.com/document_detail/27110.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
* `opening_time` - The time when the CDN service was activated. The time follows the ISO 8601 standard in the yyyy-MM-ddThh:mmZ format.
* `changing_charge_type` -  The billing method to be effective.
* `changing_affect_time` - 	The time when the change of the billing method starts to take effect. The time is displayed in GMT.
