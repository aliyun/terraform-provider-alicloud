---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_service"
sidebar_current: "docs-alicloud-datasource-pvtz-service"
description: |-
    Provides a datasource to open the Private Zone service automatically.
---

# alicloud\_pvtz\_service

Using this data source can open Private Zone service automatically. If the service has been opened, it will return opened.

For information about Priavte Zone and how to use it, see [What is Private Zone](https://www.alibabacloud.com/help/en/product/64583.htm).

-> **NOTE:** Available in v1.114.0+

## Example Usage

```terraform
data "alicloud_pvtz_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the Private Zone service that means you have read and agreed the [Private Zone Terms of Service](https://help.aliyun.com/document_detail/65657.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
