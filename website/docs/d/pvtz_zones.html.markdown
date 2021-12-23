---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zones"
sidebar_current: "docs-alicloud-datasource-pvtz-zones"
description: |-
    Provides a list of Private Zones which owned by an Alibaba Cloud account.
---

# alicloud\_pvtz\_zones

This data source lists a number of Private Zones resource information owned by an Alibaba Cloud account.

## Example Usage

```terraform
data "alicloud_pvtz_zones" "pvtz_zones_ds" {
  keyword = "${alicloud_pvtz_zone.basic.zone_name}"
}

output "first_zone_id" {
  value = "${data.alicloud_pvtz_zones.pvtz_zones_ds.zones.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `keyword` - (Optional) keyword for zone name.
* `ids` - (Optional, Available 1.53.0+) A list of zone IDs. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `lang` - (Optional, Available 1.107.0+) User language.
* `query_region_id` - (Optional, Available 1.107.0+) query_region_id for zone regionId.
* `query_vpc_id` - (Optional, Available 1.107.0+) query_vpc_id for zone vpcId.
* `resource_group_id` - (Optional, Available 1.107.0+) resource_group_id for zone resourceGroupId.
* `search_mode` - (Optional, Available 1.107.0+) Search mode. Value: 
    - LIKE: fuzzy search.
    - EXACT: precise search. It is not filled in by default.
* `enable_details` -(Optional, Available 1.107.0+) Default to `false`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs. 
* `names` - A list of zone names. 
* `zones` - A list of zones. Each element contains the following attributes:
  * `id` - ID of the Private Zone.
  * `remark` - Remark of the Private Zone.
  * `record_count` - Count of the Private Zone Record.
  * `name` - Name of the Private Zone.
  * `zone_name` - ZoneName of the Private Zone.
  * `is_ptr` - Whether the Private Zone is ptr.
  * `create_timestamp` - Time of create of the Private Zone.
  * `update_timestamp` - Time of update of the Private Zone.
  * `proxy_pattern` - The recursive DNS proxy.
  * `resource_group_id` - The Id of resource group which the Private Zone belongs.
  * `slave_dns` - Whether to turn on secondary DNS.
  * `zone_id` - ZoneId of the Private Zone.
  * `bind_vpcs` - List of the VPCs is bound to the Private Zone:
    * `region_id` - Binding the regionId of VPC.
    * `region_name` - Binding the regionName of VPC.
    * `vpc_id` - Binding the vpcId of VPC.
