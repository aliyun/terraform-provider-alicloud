---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_service"
sidebar_current: "docs-alicloud-datasource-ots-service"
description: |-
    Provides a datasource to open the Table Staore service automatically.
---

# alicloud\_ots\_service

Using this data source can enable Table Staore service automatically. If the service has been enabled, it will return `Opened`.

For information about Table Staore and how to use it, see [What is Table Staore](https://www.alibabacloud.com/help/product/27278.htm).

-> **NOTE:** Available in v1.97.0+

## Example Usage

```
data "alicloud_ots_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off".

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
