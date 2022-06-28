---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_service"
sidebar_current: "docs-alicloud-datasource-edas-service"
description: |-
    Provides a datasource to open the EDAS service automatically.
---

# alicloud\_edas\_service

Using this data source can open EDAS service automatically. If the service has been opened, it will return opened.

For information about EDAS and how to use it, see [What is EDAS](https://www.alibabacloud.com/help/product/29500.htm).

-> **NOTE:** Available in v1.98.0+

-> **NOTE:** The EDAS service is not support to be open automatically in the international site.

## Example Usage

```
data "alicloud_edas_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the EDAS service that means you have read and agreed the [EDAS Terms of Service](https://help.aliyun.com/document_detail/44633.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
