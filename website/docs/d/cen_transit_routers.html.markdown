---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_routers"
sidebar_current: "docs-alicloud-datasource-cen-transit-routers"
description: |-
Provides a list of CEN Transit Routers owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_routers

This data source provides CEN Transit Routers available to the user.

## Example Usage

```
data "alicloud_cen_transit_routers" "default" {
  cen_id    = "cen-id1"
  region    = "cn-*****"
}

output "first_transit_routers_type" {
  value = "${data.alicloud_cen_transit_routers.default.transit_routers.0.type}"
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) ID of the CEN instance.
* `region_id` - (Optional) Region ID of the VBR.
* `transit_router_id` - (Optional) ID of the transit router.
* `transit_router_ids` - (Optional) A list of ID of the transit router.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `transit_routers` - A list of CEN Transit Routers. Each element contains the following attributes:
    * `cen_id` - ID of the CEN instance.
    * `ali_uid` - UID of the Aliyun.
    * `transit_router_id` - ID of the transit router.
    * `transit_router_name` - The name of the transit router.
    * `transit_router_description` - The description of the transit router.
    * `status` - The status of the transit router attachment.
    * `type` - Type of the transit router.
    * `region_id` - Region ID of the transit router.
    * `xgw_vip` - The vip of the XGW.
    * `creation_time` - The time when it's created.
