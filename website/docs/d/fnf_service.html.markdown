---
subcategory: "Serverless Workflow"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_service"
sidebar_current: "docs-alicloud-datasource-fnf-service"
description: |-
    Provides a datasource to open the Fnf service automatically.
---

# alicloud\_fnf\_service

Using this data source can open Fnf service automatically. If the service has been opened, it will return opened.

For information about Fnf and how to use it, see [What is Fnf](https://www.alibabacloud.com/help/en/product/113549.htm).

-> **NOTE:** Available in v1.114.0+

## Example Usage

```terraform
data "alicloud_fnf_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the Fnf service that means you have read and agreed the [Fnf Terms of Service](https://help.aliyun.com/document_detail/117831.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
