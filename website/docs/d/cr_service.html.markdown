---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_service"
sidebar_current: "docs-alicloud-datasource-cr-service"
description: |-
    Provides a datasource to open the Container Registry (CR) service automatically.
---

# alicloud\_cr\_service

Using this data source can open Container Registry (CR) service automatically. If the service has been opened, it will return opened.

For information about Container Registry (CR) and how to use it, see [What is Container Registry (CR)](https://www.alibabacloud.com/help/en/doc-detail/142759.htm).

-> **NOTE:** Available in v1.116.0+

## Example Usage

```terraform
data "alicloud_cr_service" "open" {
  enable   = "On"
  password = "1111aaaa"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.
* `password` - (Required) The user password. The password must be 8 to 32 characters in length, and must contain at least two of the following character types: letters, special characters, and digits.

-> **NOTE:** Setting `enable = "On"` to open the Container Registry (CR) service that means you have read and agreed the [Container Registry (CR) Terms of Service](https://help.aliyun.com/document_detail/190602.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
