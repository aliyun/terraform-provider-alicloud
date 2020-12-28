---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_service"
sidebar_current: "docs-alicloud-datasource-ons-service"
description: |-
    Provides a datasource to open the ONS service automatically.
---

# alicloud\_ons\_service

Using this data source can open ONS service automatically. If the service has been opened, it will return opened.

For information about ONS and how to use it, see [What is ONS](https://help.aliyun.com/product/29530.html).

-> **NOTE:** Available in v1.111.0+

## Example Usage

```terraform
data "alicloud_ons_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 