---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_service"
sidebar_current: "docs-alicloud-datasource-log-service"
description: |-
    Provides a datasource to open the Log service automatically.
---

# alicloud\_log\_service

Using this data source can enable Log service automatically. If the service has been enabled, it will return `Opened`.

For information about Log service and how to use it, see [What is Log Service](https://www.alibabacloud.com/help/product/28958.htm).

-> **NOTE:** Available in v1.96.0+

## Example Usage

```
data "alicloud_log_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the Log service that means you have read and agreed the [Log Terms of Service](https://help.aliyun.com/document_detail/53476.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
