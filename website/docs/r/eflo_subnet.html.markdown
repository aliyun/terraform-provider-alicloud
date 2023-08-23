---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_subnet"
sidebar_current: "docs-alicloud-resource-eflo-subnet"
description: |-
  Provides a Alicloud Eflo Subnet resource.
---

# alicloud_eflo_subnet

Provides a Eflo Subnet resource.

For information about Eflo Subnet and how to use it, see [What is Subnet](https://www.alibabacloud.com/help/en/pai/user-guide/overview-of-intelligent-computing-lingjun).

-> **NOTE:** Available since v1.204.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}
data "alicloud_zones" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eflo_vpd" "default" {
  cidr              = "10.0.0.0/8"
  vpd_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

resource "alicloud_eflo_subnet" "default" {
  subnet_name = var.name
  zone_id     = data.alicloud_zones.default.zones.0.id
  cidr        = "10.0.0.0/16"
  vpd_id      = alicloud_eflo_vpd.default.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr` - (Required, ForceNew) CIDR network segment.
* `subnet_name` - (Required) The Subnet name.
* `type` - (Optional, ForceNew) Eflo subnet usage type. optional value:
  - General type is not filled in
  - OOB:OOB type
  - LB: LB type
* `zone_id` - (Required, ForceNew) The zone ID  of the resource.
* `vpd_id` - (Required, ForceNew) The Eflo VPD ID.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<vpd_id>:<subnet_id>`.
* `create_time` - The creation time of the resource.
* `gmt_modified` - Modification time.
* `message` - Error message.
* `resource_group_id` - Resource Group ID.
* `status` - The status of the resource.
* `subnet_id` - The id of the subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Subnet.
* `delete` - (Defaults to 5 mins) Used when delete the Subnet.
* `update` - (Defaults to 5 mins) Used when update the Subnet.

## Import

Eflo Subnet can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_subnet.example <vpd_id>:<subnet_id>
```