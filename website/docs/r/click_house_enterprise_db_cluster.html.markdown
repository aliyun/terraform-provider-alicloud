---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_enterprise_db_cluster"
description: |-
  Provides a Alicloud Click House Enterprise Db Cluster resource.
---

# alicloud_click_house_enterprise_db_cluster

Provides a Click House Enterprise Db Cluster resource.

Enterprise Edition Cluster Resources.

For information about Click House Enterprise Db Cluster and how to use it, see [What is Enterprise Db Cluster](https://next.api.alibabacloud.com/document/clickhouse/2023-05-22/CreateDBInstance).

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
* `description` - (Optional, Computed, Available since v1.270.0) Cluster description.
* `multi_zones` - (Optional, ForceNew, Computed, List) The multi-zone configuration. See [`multi_zones`](#multi_zones) below.
* `node_count` - (Optional, Computed, Int, Available since v1.270.0) The number of nodes. Valid values: 2 to 16. This parameter is required when NodeScaleMin and NodeScaleMax are configured to define the auto-scaling range.
* `node_scale_max` - (Optional, Computed, Int, Available since v1.270.0) Maximum value for serverless node auto scaling. Valid values range from 4 to 32 and must be greater than the minimum value.  
* `node_scale_min` - (Optional, Computed, Int, Available since v1.270.0) The minimum value for serverless node auto-scaling. Valid values: 4–32.
* `resource_group_id` - (Optional, Computed, Available since v1.270.0) Resource group ID of the cluster.
* `scale_max` - (Optional) The maximum value for serverless auto scaling. This parameter is not recommended. We recommend that you use NodeCount, NodeScaleMin, and NodeScaleMax to configure auto scaling capabilities.
* `scale_min` - (Optional) The minimum value for serverless auto scaling. This parameter is not recommended. We recommend that you use NodeCount, NodeScaleMin, and NodeScaleMax to configure auto scaling capabilities.
* `tags` - (Optional, Map, Available since v1.270.0) Tag information.  
* `vpc_id` - (Optional, ForceNew) The VPC ID.
* `vswitch_id` - (Optional, ForceNew) vSwitch ID.
* `zone_id` - (Optional, ForceNew) The zone ID.

### `multi_zones`

The multi_zones supports the following:
* `vswitch_ids` - (Optional, ForceNew, List) List of vSwitch IDs.
* `zone_id` - (Optional, ForceNew) Zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `category` - Instance type.
* `charge_type` - The billing method.
* `computing_group_ids` - List of computing group IDs.
* `create_time` - The cluster creation time, in the format yyyy-MM-ddTHH:mm:ssZ.
* `disabled_ports` - Disable specified database ports.
* `endpoints` - List of endpoint details.
  * `computing_group_id` - The computing group ID.
  * `connection_string` - The instance connection string.
  * `endpoint_name` - The endpoint name.
  * `ip_address` - The IP address.
  * `net_type` - The network type of the connection string.
  * `ports` - A list of port details.
    * `port` - The access port.
    * `protocol` - The protocol type.
  * `status` - Status.
  * `vswitch_id` - The vSwitch ID.
  * `vpc_id` - The VPC ID.
  * `vpc_instance_id` - The VPC instance ID.
* `engine_minor_version` - The minor version number of the cluster engine.
* `engine_version` - The database engine version.
* `instance_network_type` - Network type of the instance.
* `maintain_end_time` - The end time of the maintenance window.
* `maintain_start_time` - The start time of the maintenance window.
* `object_store_size` - The object storage size.
* `region_id` - The region ID.
* `status` - The instance status.
* `storage_quota` - Pre-purchased storage capacity (GB).
* `storage_size` - The storage capacity.
* `storage_type` - The storage type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Enterprise Db Cluster.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Db Cluster.
* `update` - (Defaults to 60 mins) Used when update the Enterprise Db Cluster.

## Import

Click House Enterprise Db Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_enterprise_db_cluster.example <db_instance_id>
```