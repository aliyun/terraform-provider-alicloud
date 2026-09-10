---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_service"
sidebar_current: "docs-alicloud-datasource-oss-service"
description: |-
    Provides a datasource to open the OSS service automatically.
---

# alicloud\_oss\_service

Using this data source can enable OSS service automatically. If the service has been enabled, it will return `Opened`.

For information about OSS and how to use it, see [What is OSS](https://www.alibabacloud.com/help/product/31815.htm).

-> **NOTE:** Available in v1.97.0+

## Example Usage

```
data "alicloud_oss_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the OSS service that means you have read and agreed the [OSS Terms of Service](https://help.aliyun.com/document_detail/31821.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
