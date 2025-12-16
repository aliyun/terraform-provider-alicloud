---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_resource_group"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) Resource Group resource.
---

# alicloud_adb_resource_group

Provides a AnalyticDB for MySQL (ADB) Resource Group resource.

For information about AnalyticDB for MySQL (ADB) Resource Group and how to use it, see [What is Resource Group](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/api-doc-adb-2019-03-15-api-doc-createdbresourcegroup).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_resource_group&exampleId=b585e2dc-7099-5d61-5285-d0cca6665da8440cfdc2&activeTab=example&spm=docs.r.adb_resource_group.0.b585e2dc70&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_adb_zones" "default" {
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.4.0.0/24"
  zone_id      = data.alicloud_adb_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_adb_db_cluster" "default" {
  compute_resource    = "48Core192GB"
  db_cluster_category = "MixedStorage"
  db_cluster_version  = "3.0"
  db_node_class       = "E32"
  db_node_storage     = 100
  description         = var.name
  elastic_io_resource = 1
  maintain_time       = "04:00Z-05:00Z"
  mode                = "flexible"
  payment_type        = "PayAsYouGo"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_ips        = ["10.168.1.12", "10.168.1.11"]
  vpc_id              = alicloud_vpc.default.id
  vswitch_id          = alicloud_vswitch.default.id
  zone_id             = data.alicloud_adb_zones.default.zones[0].id
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_adb_resource_group" "default" {
  group_name    = "TF_EXAMPLE"
  group_type    = "batch"
  node_num      = 0
  db_cluster_id = alicloud_adb_db_cluster.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_adb_resource_group&spm=docs.r.adb_resource_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cluster_mode` - (Optional, Available since v1.261.0) The working mode of the resource group. Default value: `Disable`. Valid values: `Disable`, `AutoScale`.
* `cluster_size_resource` - (Optional, Available since v1.261.0) The resource specifications of a single compute cluster. Unit: ACU.
* `db_cluster_id` - (Required, ForceNew) The ID of the DBCluster.
* `engine` - (Optional, ForceNew, Available since v1.261.0) The engine of the resource group. Default value: `AnalyticDB`. Valid values: `AnalyticDB`, `SparkWarehouse`.
* `engine_params` - (Optional, Map, Available since v1.261.0) The Spark application configuration parameters that can be applied to all Spark jobs executed in the resource group.
* `group_name` - (Required, ForceNew) The name of the resource group. The `group_name` can be up to 255 characters in length and can contain digits, uppercase letters, hyphens (-), and underscores (_). It must start with a digit or uppercase letter.
* `group_type` - (Optional) The query execution mode. Default value: `interactive`. Valid values: `interactive`, `batch`.
* `max_cluster_count` - (Optional, Int, Available since v1.261.0) The maximum number of compute clusters that are allowed in the resource group.
* `max_compute_resource` - (Optional, Available since v1.261.0) The maximum amount of reserved computing resources, which refers to the amount of resources that are not allocated in the cluster.
* `min_cluster_count` - (Optional, Int, Available since v1.261.0) The minimum number of compute clusters that are required in the resource group.
* `min_compute_resource` - (Optional, Available since v1.261.0) The minimum amount of reserved computing resources. Unit: AnalyticDB compute unit (ACU).
* `node_num` - (Optional, Int) The number of nodes.
* `users` - (Optional, List, Available since v1.227.0) The database accounts with which to associate the resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Resource Group. It formats as `<db_cluster_id>:<group_name>`.
* `user` - The database accounts that are associated with the resource group.
* `port` - (Available since v1.261.0) The port number of the resource group.
* `connection_string` - (Available since v1.261.0) The endpoint of the resource group.
* `create_time` - The time when the resource group was created.
* `update_time` - The time when the resource group was updated.
* `status` - (Available since v1.261.0) The status of the resource group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 31 mins) Used when create the Resource Group.
* `delete` - (Defaults to 31 mins) Used when delete the Resource Group.
* `update` - (Defaults to 31 mins) Used when update the Resource Group.

## Import

Adb Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_resource_group.example <db_cluster_id>:<group_name>
```
