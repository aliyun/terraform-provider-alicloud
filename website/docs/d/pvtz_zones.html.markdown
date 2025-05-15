---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zones"
sidebar_current: "docs-alicloud-datasource-pvtz-zones"
description: |-
  Provides a list of Private Zones to the user.
---

# alicloud_pvtz_zones

This data source provides the Private Zones of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.13.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example.com"
}

resource "alicloud_pvtz_zone" "default" {
  zone_name = var.name
}

data "alicloud_pvtz_zones" "ids" {
  ids = [alicloud_pvtz_zone.default.id]
}

output "pvtz_zones_id_0" {
  value = data.alicloud_pvtz_zones.ids.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List, Available since v1.53.0) A list of Zones IDs.
* `name_regex` - (Optional, ForceNew, Available since v1.107.0) A regex string to filter results by Zone name.
* `keyword` - (Optional, ForceNew) The keyword of the zone name.
* `resource_group_id` - (Optional, ForceNew, Available since v1.107.0) The ID of the resource group to which the zone belongs.
* `query_vpc_id` - (Optional, ForceNew, Available since v1.107.0) The ID of the VPC associated with the zone.
* `query_region_id` - (Optional, ForceNew, Available since v1.107.0) The region ID of the virtual private cloud (VPC) associated with the zone.
* `search_mode` - (Optional, ForceNew, Available since v1.107.0) The search mode. The value of Keyword is the search scope. Default value: `LIKE`. Valid values:
  - `LIKE`: Fuzzy search.
  - `EXACT`: Exact search.
* `lang` - (Optional, ForceNew, Available since v1.107.0) The language of the response. Default value: `en`. Valid values: `en`, `zh`.
* `enable_details` -(Optional, Bool, Available since v1.107.0) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Zone names. 
* `zones` - A list of Zone. Each element contains the following attributes:
  * `id` - The ID of the Private Zone.
  * `zone_id` - The ID of the Zone.
  * `zone_name` - The Name of the Private Zone.
  * `name` - The Name of the Zone.
  * `proxy_pattern` - Indicates whether the recursive resolution proxy for subdomain names is enabled.
  * `record_count` - The number of Domain Name System (DNS) records added in the zone.
  * `resource_group_id` - The ID of the resource group to which the zone belongs.
  * `remark` - The description of the zone.
  * `is_ptr` - Indicates whether the zone is a reverse lookup zone.
  * `create_timestamp` - The time when the zone was created.
  * `update_timestamp` - The time when the DNS record was updated.
  * `slave_dns` - Indicates whether the secondary Domain Name System (DNS) feature is enabled for the zone. **Note:** `slave_dns` takes effect only if `enable_details` is set to `true`.
  * `bind_vpcs` - The VPCs associated with the zone. **Note:** `bind_vpcs` takes effect only if `enable_details` is set to `true`.
    * `vpc_id` - The ID of the VPC.
    * `vpc_name` - The Name of the VPC.
    * `region_id` - The region ID of the VPC.
    * `region_name` - The name of the region where the VPC resides.
