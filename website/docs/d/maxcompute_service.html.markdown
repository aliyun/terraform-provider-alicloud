---
subcategory: "MaxCompute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_service"
sidebar_current: "docs-alicloud-datasource-maxcompute-service"
description: |-
    Provides a datasource to open the Maxcompute service automatically.
---

# alicloud\_maxcompute\_service

-> **NOTE:** When you open MaxCompute service, you'd better open [DataWorks service](https://www.alibabacloud.com/help/en/product/72772.htm) as well.

Using this data source can open Maxcompute service automatically. If the service has been opened, it will return opened.

For information about Maxcompute and how to use it, see [What is Maxcompute](https://www.alibabacloud.com/help/en/product/27797.htm).

-> **NOTE:** Available in v1.117.0+

## Example Usage

```terraform
data "alicloud_maxcompute_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Maxcompute service that means you have read and agreed the [Maxcompute Terms of Service](https://help.aliyun.com/document_detail/98605.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
