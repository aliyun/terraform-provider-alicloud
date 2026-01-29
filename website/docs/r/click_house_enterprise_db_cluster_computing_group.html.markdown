---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster_computing_group"
description: |-
  Provides a Alicloud Click House Enterprise Db Cluster Computing Group resource.
---

# alicloud_click_house_enterprise_db_cluster_computing_group

Provides a Click House Enterprise Db Cluster Computing Group resource.

ClickHouse Enterprise Edition Cluster Compute Group.

For information about Click House Enterprise Db Cluster Computing Group and how to use it, see [What is Enterprise Db Cluster Computing Group](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/CreateComputingGroup).

-> **NOTE:** Available since v1.270.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "vsw__ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc__ip_range" {
  default = "172.16.0.0/12"
}

variable "vsw__ip_range_k" {
  default = "172.16.3.0/24"
}

variable "vsw__ip_range_l" {
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

resource "alicloud_vpc" "defaultp2mwWM" {
  cidr_block = var.vpc__ip_range
}

resource "alicloud_vswitch" "defaultkCZhNu" {
  vpc_id     = alicloud_vpc.defaultp2mwWM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw__ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "defaultQ5vukB" {
  zone_id    = alicloud_vswitch.defaultkCZhNu.zone_id
  vpc_id     = alicloud_vpc.defaultp2mwWM.id
  scale_min  = "8"
  scale_max  = "8"
  vswitch_id = alicloud_vswitch.defaultkCZhNu.id
}


resource "alicloud_click_house_enterprise_db_cluster_computing_group" "default" {
  node_scale_min              = "4"
  computing_group_description = "example"
  node_count                  = "2"
  db_instance_id              = alicloud_click_house_enterprise_db_cluster.defaultQ5vukB.id
  node_scale_max              = "4"
  is_readonly                 = false
}
```

## Argument Reference

The following arguments are supported:
* `computing_group_description` - (Optional) Computing group name
* `db_instance_id` - (Required, ForceNew) ClickHouse Enterprise Edition instance ID
* `is_readonly` - (Required) Computing group read-only
* `node_count` - (Required, Int) Number of computing group nodes
* `node_scale_max` - (Required, Int) Computing group node elasticity upper limit
* `node_scale_min` - (Required, Int) Lower limit of computing group node elasticity

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<db_instance_id>:<computing_group_id>`.
* `computing_group_endpoint_names` - 当前属性没有在镇元上录入属性描述�
* `computing_group_endpoints` - 当前属性没有在镇元上录入属性描述�
* `computing_group_id` - Computing Group ID.
* `computing_group_public_endpoints` - 当前属性没有在镇元上录入属性描述�
* `computing_group_status` - Computing group status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Enterprise Db Cluster Computing Group.
* `delete` - (Defaults to 60 mins) Used when delete the Enterprise Db Cluster Computing Group.
* `update` - (Defaults to 60 mins) Used when update the Enterprise Db Cluster Computing Group.

## Import

Click House Enterprise Db Cluster Computing Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster_computing_group.example <db_instance_id>:<computing_group_id>
```