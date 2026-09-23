---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_private_zones"
sidebar_current: "docs-alicloud-datasource-cen-private-zones"
description: |-
    Provides a list of CEN(Cloud Enterprise Network) Private Zones owned by an Alibaba Cloud account.
---

# alicloud\_cen\_private\_zones

This data source provides CEN Private Zones available to the user.

-> **NOTE:** Available in v1.88.0+.

## Example Usage

```
data "alicloud_cen_private_zones" "this" {
  cen_id         = "cen-o40h17ll9w********"
  ids            = ["cn-hangzhou"]
  status         = "Active"
}

output "first_cen_private_zones_id" {
  value = "${data.alicloud_cen_private_zones.this.zones.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN instance.
* `ids` - (Optional) A list of CEN private zone IDs. Each element format as `<cen_id>:<access_region_id>`. 
  **NOTE:** Before 1.162.0, each element same as `access_region_id`.
* `host_region_id ` - (Optional) The service region is the target region of the PrivateZone service accessed through CEN.
* `status` - (Optional) The status of the PrivateZone service, including `Creating`, `Active` and `Deleting`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN private zone IDs. Each element format as `<cen_id>:<access_region_id>`.
  **NOTE:** Before 1.162.0, each element same as `access_region_id`.
* `zones` - A list of CEN private zones. Each element contains the following attributes:
  * `id` - The ID of the private zone. It formats as `<cen_id>:<access_region_id>`.
  * `cen_id` - The ID of the CEN instance.
  * `private_zone_dns_servers` - The DNS IP addresses of the PrivateZone service.
  * `access_region_id` - The access region. The access region is the region of the cloud resource that accesses the PrivateZone service through CEN.
  * `host_region_id` - The service region. The service region is the target region of the PrivateZone service accessed through CEN.
  * `host_vpc_id` - The VPC that belongs to the service region.
  * `status` - The status of the PrivateZone service.
