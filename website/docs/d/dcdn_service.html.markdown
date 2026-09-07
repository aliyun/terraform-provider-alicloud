---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_service"
sidebar_current: "docs-alicloud-datasource-dcdn-service"
description: |-
    Provides a datasource to open the DCDN service automatically.
---

# alicloud\_dcdn\_service

Using this data source can open DCDN service automatically. If the service has been opened, it will return opened.

For information about DCDN and how to use it, see [What is DCDN](https://help.aliyun.com/document_detail/197288.html).

-> **NOTE:** Available in v1.111.0+

## Example Usage

```terraform
data "alicloud_dcdn_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the DCDN service that means you have read and agreed the [DCDN Terms of Service](https://help.aliyun.com/document_detail/169354.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
