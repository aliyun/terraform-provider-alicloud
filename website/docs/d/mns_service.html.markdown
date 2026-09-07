---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_service"
description: |-
  Provides a datasource to open the MNS service automatically.
---

# alicloud_mns_service

Using this data source can open MNS service automatically. If the service has been opened, it will return opened.

For information about MNS and how to use it, see [What is MNS](https://www.alibabacloud.com/help/en/product/27412.htm).

-> **NOTE:** Deprecated since v1.252.0.

-> **DEPRECATED:**  This datasource has been deprecated from version `1.252.0`. Please use new resource [alicloud_message_service_service](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/message_service_service).

-> **NOTE:** The MNS service is not support in the international site.

## Example Usage

```terraform
data "alicloud_mns_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Default value: `Off`. Valid values: `On` and `Off`.

-> **NOTE:** Setting `enable = "On"` to open the MNS service that means you have read and agreed the [MNS Terms of Service](https://help.aliyun.com/document_detail/27418.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
