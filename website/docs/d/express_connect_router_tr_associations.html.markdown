---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_tr_associations"
description: |-
  Provides a list of Express Connect Router Tr Association to the user.
---

# alicloud_express_connect_router_tr_associations

This data source provides the  Express Connect Router Tr Association of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn = "65532"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_express_connect_router_tr_association" "default" {
  ecr_id                  = alicloud_express_connect_router_express_connect_router.default.id
  transit_router_id       = alicloud_cen_transit_router.default.transit_router_id
  cen_id                  = alicloud_cen_transit_router.default.cen_id
  transit_router_owner_id = data.alicloud_account.default.id
  association_region_id   = data.alicloud_regions.default.regions.0.id
  allowed_prefixes        = ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
}

data "alicloud_express_connect_router_tr_associations" "ids" {
  ids    = [alicloud_express_connect_router_tr_association.default.id]
  ecr_id = alicloud_express_connect_router_tr_association.default.ecr_id
}

output "express_connect_router_tr_associations_id_0" {
  value = data.alicloud_express_connect_router_tr_associations.ids.associations.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Tr Association IDs.
* `ecr_id` - (Required) The ID of the Express Connect Router instance.
* `association_id` - (Optional) The ID of the association between the Express Connect Router and the TR.
* `transit_router_id` - (Optional) The ID of the transit router instance.
* `association_region_id` - (Optional) The region ID of the associated TR.
* `cen_id` - (Optional) The ID of the Cloud Enterprise Network instance.
* `status` - (Optional) The status of the association. Valid values: `CREATING`, `ACTIVE`, `INACTIVE`, `ASSOCIATING`, `DISSOCIATING`, `UPDATING`, `DELETING`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `associations` - A list of Tr Associations. Each element contains the following attributes:
  * `id` - The ID of the Tr Association.
  * `ecr_id` - The ID of the Express Connect Router instance.
  * `association_id` - The ID of the association between the Express Connect Router and the TR.
  * `transit_router_id` - The ID of the TR instance.
  * `association_node_type` - The type of the associated resource.
  * `transit_router_owner_id` - The ID of the Alibaba Cloud account that owns the TR.
  * `cen_id` - The ID of the CEN instance.
  * `allowed_prefixes_mode` - The prefix-based routing mode.
  * `status` - The status of the association.
  * `create_time` - The time when the association was created.
  * `modify_time` - The time when the association was modified.
  * `allowed_prefixes` - The prefix-based routing mode.
