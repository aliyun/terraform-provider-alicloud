---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_vpc_associations"
description: |-
  Provides a list of Express Connect Router Vpc Association to the user.
---

# alicloud_express_connect_router_vpc_associations

This data source provides the  Express Connect Router Vpc Association of the current Alibaba Cloud user.

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

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_express_connect_router_vpc_association" "default" {
  ecr_id                = alicloud_express_connect_router_express_connect_router.default.id
  vpc_id                = alicloud_vpc.default.id
  association_region_id = data.alicloud_regions.default.regions.0.id
  vpc_owner_id          = data.alicloud_account.default.id
  allowed_prefixes      = ["172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"]
}

data "alicloud_express_connect_router_vpc_associations" "ids" {
  ids    = [alicloud_express_connect_router_vpc_association.default.id]
  ecr_id = alicloud_express_connect_router_vpc_association.default.ecr_id
}

output "express_connect_router_vpc_associations_id_0" {
  value = data.alicloud_express_connect_router_vpc_associations.ids.associations.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Vpc Association IDs.
* `ecr_id` - (Required) The ID of the Express Connect Router instance.
* `association_id` - (Optional) The ID of the association between the Express Connect Router and the VPC.
* `vpc_id` - (Optional) The ID of the VPC instance.
* `association_region_id` - (Optional) The region ID of the associated VPC.
* `status` - (Optional) The status of the association. Valid values: `CREATING`, `ACTIVE`, `INACTIVE`, `ASSOCIATING`, `DISSOCIATING`, `UPDATING`, `DELETING`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `associations` - A list of Vpc Associations. Each element contains the following attributes:
  * `id` - The ID of the Vpc Association.
  * `ecr_id` - The ID of the Express Connect Router instance.
  * `association_id` - The ID of the association between the Express Connect Router and the VPC.
  * `vpc_id` - The ID of the VPC instance.
  * `association_node_type` - The type of the associated resource.
  * `vpc_owner_id` - The ID of the Alibaba Cloud account that owns the VPC.
  * `allowed_prefixes_mode` - The prefix-based routing mode.
  * `status` - The status of the association.
  * `create_time` - The time when the association was created.
  * `modify_time` - The time when the association was modified.
  * `allowed_prefixes` - The prefix-based routing mode.
