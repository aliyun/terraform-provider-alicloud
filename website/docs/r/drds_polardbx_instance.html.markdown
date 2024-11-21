---
subcategory: "Distributed Relational Database Service (DRDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_drds_polardbx_instance"
description: |-
  Provides a Alicloud DRDS Polardbx Instance resource.
---

# alicloud_drds_polardbx_instance

Provides a DRDS Polardb X Instance resource.

For information about DRDS Polardb X Instance and how to use it, see [What is Polardb X Instance](https://www.alibabacloud.com/help/en/polardb/polardb-for-xscale/api-createdbinstance-1).

-> **NOTE:** Available since v1.211.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_drds_polardbx_instance&exampleId=fbf375da-c462-25f1-a343-5a1971bc805bcf3c5f9b&activeTab=example&spm=docs.r.drds_polardbx_instance.0.fbf375dac4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
provider "alicloud" {
  region = "ap-southeast-1"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "example" {
  vpc_name = var.name
}
resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}
resource "alicloud_drds_polardbx_instance" "default" {
  topology_type  = "3azones"
  vswitch_id     = alicloud_vswitch.example.id
  primary_zone   = "ap-southeast-1a"
  cn_node_count  = "2"
  dn_class       = "mysql.n4.medium.25"
  cn_class       = "polarx.x4.medium.2e"
  dn_node_count  = "2"
  secondary_zone = "ap-southeast-1b"
  tertiary_zone  = "ap-southeast-1c"
  vpc_id         = alicloud_vpc.example.id
}
```

## Argument Reference

The following arguments are supported:
* `cn_class` - (Required, ForceNew) Compute node specifications.
* `cn_node_count` - (Required) Number of computing nodes.
* `dn_class` - (Required, ForceNew) Storage node specifications.
* `dn_node_count` - (Required) The number of storage nodes.
* `primary_zone` - (Required, ForceNew) Primary Availability Zone.
* `resource_group_id` - (Optional, Computed) The resource group ID can be empty. This parameter is not supported for the time being.
* `secondary_zone` - (Optional, ForceNew) Secondary availability zone.
* `tertiary_zone` - (Optional, ForceNew) Third Availability Zone.
* `topology_type` - (Required, ForceNew) Topology type:
  - **3azones**: three available areas;
  - **1azone**: Single zone.
* `vswitch_id` - (Required, ForceNew) The ID of the virtual switch.
* `vpc_id` - (Required, ForceNew) The VPC ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Polardbx Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Polardbx Instance.
* `update` - (Defaults to 5 mins) Used when update the Polardbx Instance.

## Import

DRDS Polardb X Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_drds_polardb_x_instance.example <id>
```