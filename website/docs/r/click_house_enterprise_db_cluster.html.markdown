---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster"
description: |-
  Provides a Alicloud Click House Enterprise D B Cluster resource.
---

# alicloud_click_house_enterprise_db_cluster

Provides a Click House Enterprise D B Cluster resource.

Enterprise Edition Cluster Resources.

For information about Click House Enterprise D B Cluster and how to use it, see [What is Enterprise D B Cluster](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/CreateDBInstance).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "vsw_ip_range_k" {
  default = "172.16.3.0/24"
}

variable "vsw_ip_range_l" {
  default = "172.16.2.0/24"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

variable "zone_id_l" {
  default = "cn-beijing-l"
}

variable "zone_id_k" {
  default = "cn-beijing-k"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_vswitch" "defaultylyLu8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_l
  cidr_block = var.vsw_ip_range_l
}

resource "alicloud_vswitch" "defaultRNbPh8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_k
  cidr_block = var.vsw_ip_range_k
}


resource "alicloud_click_house_enterprise_db_cluster" "default" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
  multi_zones {
    vswitch_ids = ["${alicloud_vswitch.defaultTQWN3k.id}"]
    zone_id     = var.zone_id_i
  }
  multi_zones {
    vswitch_ids = ["${alicloud_vswitch.defaultylyLu8.id}"]
    zone_id     = var.zone_id_l
  }
  multi_zones {
    vswitch_ids = ["${alicloud_vswitch.defaultRNbPh8.id}"]
    zone_id     = var.zone_id_k
  }
}
```

## Argument Reference

The following arguments are supported:
* `multi_zones` - (Optional, ForceNew, Computed, Set) The list of multi-zone information. See [`multi_zones`](#multi_zones) below.
* `scale_max` - (Optional) The maximum value of serverless auto scaling.
* `scale_min` - (Optional) The minimum value of serverless auto scaling.
* `vpc_id` - (Optional, ForceNew) The VPC ID.
* `vswitch_id` - (Optional, ForceNew) The vSwitch ID.
* `zone_id` - (Optional, ForceNew) The zone ID.

### `multi_zones`

The multi_zones supports the following:
* `vswitch_ids` - (Optional, ForceNew, Set) The vSwtichID list.
* `zone_id` - (Optional, ForceNew) The zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Enterprise D B Cluster.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise D B Cluster.
* `update` - (Defaults to 60 mins) Used when update the Enterprise D B Cluster.

## Import

Click House Enterprise D B Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster.example <id>
```