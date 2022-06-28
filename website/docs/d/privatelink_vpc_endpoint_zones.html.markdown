---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_zones"
sidebar_current: "docs-alicloud-datasource-privatelink-vpc-endpoint-zones"
description: |-
  Provides a list of Privatelink Vpc Endpoint Zones to the user.
---

# alicloud\_privatelink\_vpc\_endpoint\_zones

This data source provides the Privatelink Vpc Endpoint Zones of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_privatelink_vpc_endpoint_zones" "example" {
  endpoint_id = "ep-gw8boxxxxx"
}

output "first_privatelink_vpc_endpoint_zone_id" {
  value = data.alicloud_privatelink_vpc_endpoint_zones.example.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Required, ForceNew) The ID of the Vpc Endpoint.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The Status of Vpc Endpoint Zone. Valid Values: `Connected`, `Connecting`, `Creating`, `Deleted`, `Deleting`, `Disconnected`, `Disconnecting` and `Wait`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Vpc Endpoint Zone IDs.
* `zones` - A list of Privatelink Vpc Endpoint Zones. Each element contains the following attributes:
	* `eni_id` - Terminal node network card.
	* `eni_ip` - IP address of the terminal node network card.
	* `id` - The ID of the Vpc Endpoint Zone.
	* `status` - The Status of Vpc Endpoint Zone..
	* `vswitch_id` - The VSwitch id.
	* `zone_domain` - The Zone Domain.
	* `zone_id` - The Zone Id.
