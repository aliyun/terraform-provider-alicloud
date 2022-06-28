---
subcategory: "Message Notification Service (MNS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_service"
sidebar_current: "docs-alicloud-datasource-mns-service"
description: |-
    Provides a datasource to open the MNS service automatically.
---

# alicloud\_mns\_service

Using this data source can open MNS service automatically. If the service has been opened, it will return opened.

For information about MNS and how to use it, see [What is MNS](https://www.alibabacloud.com/help/en/product/27412.htm).

-> **NOTE:** Available in v1.118.0+

-> **NOTE:** The MNS service is not support in the international site.

## Example Usage

```terraform
data "alicloud_mns_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the MNS service that means you have read and agreed the [MNS Terms of Service](https://help.aliyun.com/document_detail/27418.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
