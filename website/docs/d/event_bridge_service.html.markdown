---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_service"
sidebar_current: "docs-alicloud-datasource-event-bridge-service"
description: |-
    Provides a datasource to open the Event Bridge service automatically.
---

# alicloud\_event\_bridge\_service

Using this data source can open Event Bridge service automatically. If the service has been opened, it will return opened.

For information about Event Bridge and how to use it, see [What is Event Bridge](https://www.alibabacloud.com/help/en/doc-detail/163239.htm).

-> **NOTE:** Available in v1.126.0+

-> **NOTE:** This data source supports `cn-shanghai`, `cn-hangzhou` and `ap-southeast-1` regions.

## Example Usage

```
data "alicloud_event_bridge_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Event Bridge service that means you have read and agreed the [Event Bridge Terms of Service](https://help.aliyun.com/document_detail/163911.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
