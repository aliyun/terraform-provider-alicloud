---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_service"
sidebar_current: "docs-alicloud-datasource-privatelink-service"
description: |-
    Provides a datasource to open the Privatelink service automatically.
---

# alicloud\_privatelink\_service

Using this data source can open Privatelink service automatically. If the service has been opened, it will return opened.

For information about Privatelink and how to use it, see [What is Privatelink](https://www.alibabacloud.com/help/en/product/120462.htm).

-> **NOTE:** Available in v1.113.0+

## Example Usage

```terraform
data "alicloud_privatelink_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Privatelink service that means you have read and agreed the [Privatelink Terms of Service](https://help.aliyun.com/document_detail/197619.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
