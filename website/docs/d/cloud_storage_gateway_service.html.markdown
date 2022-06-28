---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_service"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-service"
description: |-
    Provides a datasource to open the Cloud Storage Gateway service automatically.
---

# alicloud\_cloud\_storage\_gateway\_service

Using this data source can open Cloud Storage Gateway service automatically. If the service has been opened, it will return opened.

For information about Cloud Storage Gateway and how to use it, see [What is Cloud Storage Gateway](https://www.alibabacloud.com/help/en/product/53923.htm).

-> **NOTE:** Available in v1.117.0+

## Example Usage

```terraform
data "alicloud_cloud_storage_gateway_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the Cloud Storage Gateway service that means you have read and agreed the [Cloud Storage Gateway Terms of Service](https://help.aliyun.com/document_detail/117679.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
