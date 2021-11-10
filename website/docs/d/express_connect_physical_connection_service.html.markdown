---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_physical_connection_service"
sidebar_current: "docs-alicloud-datasource-express-connect-physical-connection-service"
description: |-
    Provides a datasource to open the physical connection service automatically.
---

# alicloud_express_connect_physical_connection_service

Using this data source can enable outbound traffic for an Express Connect circuit automatically. If the service has been opened, it will return opened.

For information about Express Connect and how to use it, see [What is Express Connect](https://www.alibabacloud.com/help/doc-detail/275179.htm).

-> **NOTE:** Available in v1.132.0+

## Example Usage

```
data "alicloud_express_connect_physical_connection_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to enable outbound traffic for an Express Connect circuit that means you have read and agreed the [Express Connect Terms of Service](https://terms.aliyun.com/legal-agreement/terms/suit_bu1_ali_cloud/suit_bu1_ali_cloud201803060947_16271.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
