---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_stocks"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-stocks"
description: |- 
   Provides a list of Cloud Storage Gateway Stocks to the user.
---

# alicloud\_cloud\_storage\_gateway\_stocks

This data source provides the Cloud Storage Gateway Stocks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Advanced"
}
output "zone_id" {
  value = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `gateway_class` - (Optional, ForceNew) The gateway class. Valid values: `Basic`, `Standard`,`Enhanced`,`Advanced`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `stocks` - A list of Cloud Storage Gateway Stocks. Each element contains the following attributes:
    * `zone_id` - The Zone ID.
    * `available_gateway_classes` - A list of available gateway class in this Zone ID.
