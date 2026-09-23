---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_service"
sidebar_current: "docs-alicloud-datasource-kms-service"
description: |-
    Provides a datasource to open the KMS service automatically.
---

# alicloud\_kms\_service

Using this data source can open KMS service automatically. If the service has been opened, it will return opened.

For information about KMS and how to use it, see [What is KMS](https://help.aliyun.com/document_detail/186020.html).

-> **NOTE:** Available in v1.108.0+

## Example Usage

```
data "alicloud_kms_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

-> **NOTE:** Setting `enable = "On"` to open the KMS service that means you have read and agreed the [KMS Terms of Service](https://help.aliyun.com/document_detail/125937.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
