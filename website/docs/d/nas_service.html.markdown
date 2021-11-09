---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_service"
sidebar_current: "docs-alicloud-datasource-nas-service"
description: |-
    Provides a datasource to open the NAS service automatically.
---

# alicloud\_nas\_service

Using this data source can enable NAS service automatically. If the service has been enabled, it will return `Opened`.

For information about NAS and how to use it, see [What is NAS](https://www.alibabacloud.com/help/product/27516.htm).

-> **NOTE:** Available in v1.97.0+

## Example Usage

```
data "alicloud_nas_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the NAS service that means you have read and agreed the [NAS Terms of Service](https://help.aliyun.com/knowledge_detail/44004.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
