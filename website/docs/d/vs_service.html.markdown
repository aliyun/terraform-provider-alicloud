---
subcategory: "Video Surveillance (VS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vs_service"
sidebar_current: "docs-alicloud-datasource-vs-service"
description: |-
    Provides a datasource to open the Video Surveillance (VS) service automatically.
---

# alicloud\_vs\_service

Using this data source can open Video Surveillance (VS) service automatically. If the service has been opened, it will return opened.

For information about Video Surveillance (VS) and how to use it, see [What is VS](https://help.aliyun.com/product/108765.html).

-> **NOTE:** Available in v1.116.0+

## Example Usage

```terraform
data "alicloud_vs_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Video Surveillance (VS) service that means you have read and agreed the [VS Terms of Service](https://help.aliyun.com/document_detail/109213.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 