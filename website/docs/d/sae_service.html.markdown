---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_service"
sidebar_current: "docs-alicloud-datasource-sae-service"
description: |-
    Provides a datasource to open the SAE service automatically.
---

# alicloud_sae_service

Using this data source can open SAE service automatically. If the service has been opened, it will return opened.

For information about SAE and how to use it, see [What is SAE](https://www.alibabacloud.com/help/doc-detail/125720.htm).

-> **NOTE:** The SAE service is not support in the international site.

-> **NOTE:** Available since v1.120.0+

## Example Usage

```terraform
data "alicloud_sae_service" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

  -> **NOTE:** Setting `enable = "On"` to open the SAE service that means you have read and agreed the [SAE Terms of Service](https://help.aliyun.com/document_detail/123775.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status.
