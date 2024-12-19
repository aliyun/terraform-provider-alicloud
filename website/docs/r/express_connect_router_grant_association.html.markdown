---
subcategory: "Express Connect Router"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_grant_association"
description: |-
  Provides a Alicloud Express Connect Router Grant Association resource.
---

# alicloud_express_connect_router_grant_association

Provides a Express Connect Router Grant Association resource.

Network instances authorized to the leased line Gateway.

For information about Express Connect Router Grant Association and how to use it, see [What is Grant Association](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.239.0.

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
  default = "vpc-7qbx5bpq0axxbt9pqnzel"
}

variable "ecr_owner_uid" {
  default = <<EOF
1891593620094065
EOF
}

variable "ecr_id" {
  default = "ecr-0a6p1fk05ji1whedvj"
}

variable "region" {
  default = "cn-wulanchabu-example-5"
}


resource "alicloud_express_connect_router_grant_association" "default" {
  ecr_id             = var.ecr_id
  instance_region_id = var.region
  instance_id        = var.vpc_id
  ecr_owner_ali_uid  = var.ecr_owner_uid
  instance_type      = "VPC"
}
```

## Argument Reference

The following arguments are supported:
* `ecr_id` - (Required, ForceNew) The ID of the associated Leased Line Gateway instance.
* `ecr_owner_ali_uid` - (Required, ForceNew, Int) The ID of the Alibaba Cloud account (primary account) to which the leased line gateway instance is authorized.
* `instance_id` - (Required, ForceNew) The ID of the network instance.
* `instance_region_id` - (Required, ForceNew) The ID of the region where the authorized network instance is located.
* `instance_type` - (Required, ForceNew) The type of the network instance. Value:
  - `VBR`: the VBR instance.
  - `VPC`: VPC instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ecr_id>:<instance_id>:<instance_region_id>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Grant Association.
* `delete` - (Defaults to 5 mins) Used when delete the Grant Association.

## Import

Express Connect Router Grant Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_grant_association.example <ecr_id>:<instance_id>:<instance_region_id>
```