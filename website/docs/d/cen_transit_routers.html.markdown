---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_routers"
sidebar_current: "docs-alicloud-datasource-cen-transit-routers"
description: |-
  Provides a list of CEN Transit Routers owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_routers

This data source provides CEN Transit Routers available to the user.[What is Cen Transit Routers](https://help.aliyun.com/document_detail/261219.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

```
data "alicloud_cen_transit_routers" "default" {
  cen_id    = "cen-id1"
}

output "first_transit_routers_type" {
  value = data.alicloud_cen_transit_routers.default.transit_routers.0.type
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `ids` - (Optional, ForceNew, Available in 1.151.0+) A list of resource id. The element value is same as <cen_id>:<transit_router_id>`.
* `name_regex` - (Optional, ForceNew, Available in 1.151.0+) A regex string to filter CEN Transit Routers by name.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `status` - (Optional, ForceNew) The status of the resource. Valid values `Active`, `Creating`, `Deleting` and `Updating`.  
* `transit_router_ids` - (Optional, ForceNew) A list of ID of the transit router.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of  CEN Transit Routers names.
* `transit_routers` - A list of CEN Transit Routers. Each element contains the following attributes:
    * `cen_id` - The ID of the CEN instance.
    * `ali_uid` - The UID of the Aliyun.
    * `transit_router_id` - The ID of the transit router.
    * `transit_router_name` - The name of the transit router.
    * `transit_router_description` - The description of the transit router.
    * `status` - The status of the transit router attachment.
    * `type` - The Type of the transit router.
    * `region_id` - The Region ID of the transit router.
    * `xgw_vip` - The vip of the XGW.
    * `id` - The ID of the resource, It is formatted to `<cen_id>:<transit_router_id>`. **NOTE:** Before 1.151.0, It is formatted to `<transit_router_id>`.
