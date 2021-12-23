---
subcategory: "Internet of Things (Iot)"
layout: "alicloud"
page_title: "Alicloud: alicloud_iot_service"
sidebar_current: "docs-alicloud-datasource-iot-service"
description: |-
    Provides a datasource to open the IOT service automatically.
---

# alicloud\_iot\_service

Using this data source can open IOT service automatically. If the service has been opened, it will return opened.

For information about IOT and how to use it, see [What is IOT](https://www.alibabacloud.com/help/en/product/30520.htm).

-> **NOTE:** Available in v1.115.0+

## Example Usage

```terraform
data "alicloud_iot_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the IOT service that means you have read and agreed the [IOT Terms of Service](https://help.aliyun.com/document_detail/44548.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
