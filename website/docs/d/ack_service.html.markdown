---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ack_service"
sidebar_current: "docs-alicloud-datasource-ack-service"
description: |-
    Provides a datasource to open the Container Service (CS) service automatically.
---

# alicloud\_ack\_service

Using this data source can open Container Service (CS) service automatically. If the service has been opened, it will return opened.

For information about Container Service (CS) and how to use it, see [What is Container Service (CS)](https://www.alibabacloud.com/help/en/product/85222.htm).

-> **NOTE:** Available in v1.113.0+

## Example Usage

```
data "alicloud_ack_service" "open" {
	enable = "On"
    type   = "propayasgo"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.
* `type` - (Required, Available in 1.114.0+) Types of services opened. Valid values: `propayasgo`: Container service ack Pro managed version, `edgepayasgo`: Edge container service, `gspayasgo`: Gene computing services.

-> **NOTE:** Setting `enable = "On"` to open the Container Service (CS) service that means you have read and agreed the [Container Service (CS) Terms of Service](https://help.aliyun.com/document_detail/157971.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
