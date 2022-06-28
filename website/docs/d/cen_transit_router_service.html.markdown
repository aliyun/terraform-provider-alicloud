---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_service"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-service"
description: |-
    Provides a datasource to open the CEN Transit Router Service automatically.
---

# alicloud\_cen\_transit\_router\_service

Using this data source can open CEN Transit Router Service automatically. If the service has been opened, it will return opened.

For information about CEN and how to use it, see [What is CEN](https://www.alibabacloud.com/help/en/doc-detail/59870.htm).

-> **NOTE:** Available in v1.139.0+

## Example Usage

```
data "alicloud_cen_transit_router_service" "open" {
	enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. If has been enabled, return the result. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the CEN Transit Router Service that means you have read and agreed the [CEN Terms of Service](https://help.aliyun.com/document_detail/66667.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
