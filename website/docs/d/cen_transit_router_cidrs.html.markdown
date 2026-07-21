---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_cidrs"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-cidrs"
description: |-
  Provides a list of Cen Transit Router Cidrs to the user.
---

# alicloud\_cen\_transit\_router\_cidrs

This data source provides the Cen Transit Router Cidrs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.193.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_cidrs" "ids" {
  ids               = ["example_id"]
  transit_router_id = "tr-6ehx7q2jze8ch5ji0****"
}

output "cen_transit_router_cidr_id_0" {
  value = data.alicloud_cen_transit_router_cidrs.ids.cidrs.0.id
}

data "alicloud_cen_transit_router_cidrs" "nameRegex" {
  name_regex        = "^my-name"
  transit_router_id = "tr-6ehx7q2jze8ch5ji0****"
}

output "cen_transit_router_cidr_id_1" {
  value = data.alicloud_cen_transit_router_cidrs.nameRegex.cidrs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Cen Transit Router Cidr IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Transit Router Cidr name.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_cidr_id` - (Optional, ForceNew) The ID of the transit router cidr.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Transit Router Cidr names.
* `cidrs` - A list of Cen Transit Router Cidrs. Each element contains the following attributes:
  * `id` - The ID of the Cen Transit Router Cidr. It formats as `<transit_router_id>:<transit_router_cidr_id>`.
  * `transit_router_id` - The ID of the transit router.
  * `transit_router_cidr_id` - The ID of the transit router cidr.
  * `cidr` - The cidr of the transit router.
  * `transit_router_cidr_name` - The name of the transit router.
  * `description` - The description of the transit router.
  * `publish_cidr_route` - Whether to allow automatically adding Transit Router Cidr in Transit Router Route Table.
  * `family` - The type of the transit router cidr.
