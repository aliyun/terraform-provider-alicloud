---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_zones"
sidebar_current: "docs-alicloud-datasource-slb-zones"
description: |-
    Provides a list of availability zones for SLB that can be used by an Alibaba Cloud account.
---

# alicloud\_slb\_zones

This data source provides availability zones for SLB that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.73.0+.

## Example Usage

```terraform
data "alicloud_slb_zones" "zones_ids" {
  available_slb_address_type       = "vpc"
  available_slb_address_ip_version = "ipv4"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Deprecated from v1.154.0+) Default to false and only output `id` in the `zones` block. Set it to true can output more details.
* `available_slb_address_type` - (Optional) Filter the results by a slb instance network type. Valid values:
  * vpc: an internal SLB instance that is deployed in a virtual private cloud (VPC).
  * classic_internet: a public-facing SLB instance. 
  * classic_intranet: an internal SLB instance that is deployed in a classic network.
    
* `available_slb_address_ip_version` - (Optional) Filter the results by a slb instance address version. Can be either `ipv4`, or `ipv6`.
* `master_zone_id` - (Optional, Available in 1.157.0+) The primary zone.
* `slave_zone_id` - (Optional, Available in 1.157.0+) The secondary zone.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of primary zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone. It is same as `master_zone_id`.
  * `master_zone_id` - (Available in 1.157.0+) The primary zone.
  * `slave_zone_id` - (Available in 1.157.0+) The secondary zone.
  * `slb_slave_zone_ids` - (Deprecated from 1.157.0) A list of slb slave zone ids in which the slb master zone. 
    It has been deprecated from v1.157.0 and use `slave_zone_id` instead.
  * `supported_resources` - (Available in 1.154.0+)A list of available resource which the slb master zone supported.
    * `address_type` - The type of network.
    * `address_ip_version` - The type of IP address.

