---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_grant_associations"
description: |-
  Provides a list of Express Connect Router Grant Associations to the user.
---

# alicloud_express_connect_router_grant_associations

This data source provides the Express Connect Router Grant Association of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "vpc_id" {
  # You need to modify this value to an existing VPC under your account
  default = "vpc-xxx"
}

variable "ecr_owner_uid" {
  # You need to modify this value to ecr owner ali uid
  default = "18xxx"
}

variable "ecr_id" {
  # You need to modify this value to an existing ecr id
  default = "ecr-xxx"
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alicloud_express_connect_router_grant_association" "default" {
  ecr_id             = var.ecr_id
  instance_region_id = var.region
  instance_id        = var.vpc_id
  ecr_owner_ali_uid  = var.ecr_owner_uid
  instance_type      = "VPC"
}

data "alicloud_express_connect_router_grant_associations" "ids" {
  ids    = [alicloud_express_connect_router_grant_association.default.id]
  ecr_id = alicloud_express_connect_router_grant_association.default.ecr_id
}

output "express_connect_router_grant_associations_id_0" {
  value = data.alicloud_express_connect_router_grant_associations.ids.associations.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Tr Association IDs.
* `ecr_id` - (Required) The ID of the Express Connect Router instance.
* `instance_id` - (Optional) The ID of the authorized instance.
* `instance_region_id` - (Optional) The region where the authorized network instance is located.
* `instance_type` - (Optional) The type of the network instance. Valid values:
  - `VBR`: virtual border router (VBR).
  - `VPC`: virtual private cloud (VPC).
* `caller_type` - (Optional) The type of the caller. Default value: `ECR`. Valid values: `ECR`, `OTHER`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `associations` - A list of Grant Associations. Each element contains the following attributes:
  * `id` - The ID of the Grant Association.
  * `ecr_id` - The ID of the Express Connect Router instance.
  * `instance_id` - The ID of the authorized instance.
  * `instance_region_id` - The ID of the region to which the authorized instance belongs.
  * `owner_id` - The ID of the Alibaba Cloud account that owns the Express Connect Router instance.
  * `instance_type` - The type of the network instance.
  * `grant_id` - The authorization ID.
  * `instance_owner_id` - The ID of the Alibaba Cloud account to which the instance belongs.
  * `instance_owner_bid` - The ID of the enterprise account to which the instance belongs.
  * `status` - The status of the authorized network instance.
  * `create_time` - The time when the instance was created.
  * `modify_time` - The time when the instance was modified.
