---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_service"
sidebar_current: "docs-alicloud-datasource-fc-service"
description: |-
    Provides a datasource to open the FC service automatically.
---

# alicloud\_fc\_service

Using this data source can open FC service automatically. If the service has been opened, it will return opened.

For information about FC and how to use it, see [What is FC](https://www.alibabacloud.com/help/en/product/50980.htm).

-> **NOTE:** Available in v1.112.0+

## Example Usage

```
data "alicloud_fc_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the FC service that means you have read and agreed the [FC Terms of Service](https://help.aliyun.com/document_detail/52972.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
