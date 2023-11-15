---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_endpoint_group_ip_address_cidr_blocks"
sidebar_current: "docs-alicloud-datasource-ga-endpoint-group-ip-address-cidr-blocks"
description: |-
  Provides a list of Global Accelerator (GA) Endpoint Group Ip Address Cidr Blocks to the user.
---

# alicloud_ga_endpoint_group_ip_address_cidr_blocks

This data source provides the Global Accelerator (GA) Endpoint Group Ip Address Cidr Blocks of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_endpoint_group_ip_address_cidr_blocks" "default" {
  endpoint_group_region = "cn-hangzhou"
}

output "ga_endpoint_group_ip_address_cidr_blocks_endpoint_group_region" {
  value = data.alicloud_ga_endpoint_group_ip_address_cidr_blocks.default.endpoint_group_ip_address_cidr_blocks.0.endpoint_group_region
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_region` - (Required, ForceNew) The region ID of the endpoint group.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `endpoint_group_ip_address_cidr_blocks` - A list of Endpoint Group Ip Address Cidr Blocks. Each element contains the following attributes:
  * `endpoint_group_region` - The region ID of the endpoint group.
  * `ip_address_cidr_blocks` - The CIDR blocks.
  * `status` - The status of the endpoint group.
  