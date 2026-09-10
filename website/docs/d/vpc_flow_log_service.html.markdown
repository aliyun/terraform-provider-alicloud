---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_log_service"
sidebar_current: "docs-alicloud-datasource-vpc-flow-log-service"
description: |-
  Provides a datasource to open the Vpc Flow Log service automatically.
---

# alicloud_vpc_flow_log_service

Using this data source can open Vpc Flow Log service automatically. If the service has been opened, it will return opened.

For information about Vpc Flow Log and how to use it, see [What is Vpc Flow Log](https://www.alibabacloud.com/help/en/vpc/developer-reference/api-openflowlog).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_flow_log_service" "default" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Default value: `Off`. Valid values: `On` and `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Vpc Flow Log service that means you have read and agreed the [Vpc Flow Log Terms of Service](https://help.aliyun.com/zh/vpc/support/vpc-terms-of-service). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
