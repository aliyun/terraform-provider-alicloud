---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_service"
sidebar_current: "docs-alicloud-datasource-vpc-traffic-mirror-service"
description: |-
    Provides a datasource to open the VPC Traffic Mirror service automatically.
---

# alicloud\_event\_bridge\_service

Using this data source can open VPC Traffic Mirror service automatically. If the service has been opened, it will return opened.

For information about VPC Traffic Mirror and how to use it, see [What is VPC Traffic Mirror](https://www.alibabacloud.com/help/en/doc-detail/207513.htm).

-> **NOTE:** Available in v1.141.0+


## Example Usage

```
data "alicloud_vpc_traffic_mirror_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the VPC Traffic Mirror service that means you have read and agreed the [VPC Traffic Mirror Terms of Service](https://help.aliyun.com/document_detail/325573.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
