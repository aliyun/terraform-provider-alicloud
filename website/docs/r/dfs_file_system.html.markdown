---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_file_system"
description: |-
  Provides a Alicloud DFS File System resource.
---

# alicloud_dfs_file_system

Provides a DFS File System resource. 

For information about DFS File System and how to use it, see [What is File System](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_dfs_zones" "default" {}

resource "alicloud_dfs_file_system" "default" {
  storage_type                     = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
  zone_id                          = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type                    = "HDFS"
  description                      = var.name
  file_system_name                 = var.name
  throughput_mode                  = "Provisioned"
  space_capacity                   = "1024"
  provisioned_throughput_in_mi_bps = "512"
}
```

## Argument Reference

The following arguments are supported:
* `data_redundancy_type` - (Optional) Redundancy mode of the file system. Value:
  - LRS (default): Local redundancy.
  - ZRS: Same-City redundancy. When ZRS is selected, zoneId is a string consisting of multiple zones that are expected to be redundant in the same city, for example,  'zoneId1,zoneId2 '.
* `description` - (Optional) The description of the file system resource. No more than 32 characters in length.
* `file_system_name` - (Required) The file system name. The naming rules are as follows: The length is 6~64 characters. Globally unique and cannot be an empty string. English letters are supported and can contain numbers, underscores (_), and dashes (-).
* `partition_number` - (Optional, Available since v1.218.0) Save set sequence number, the user selects the content of the specified sequence number in the Save set.
* `protocol_type` - (Required, ForceNew) The protocol type.  Only HDFS(Hadoop Distributed File System) is supported.
* `provisioned_throughput_in_mi_bps` - (Optional) Provisioned throughput. This parameter is required when ThroughputMode is set to Provisioned. Unit: MB/s Value range: 1~5120.
* `space_capacity` - (Required) File system capacity.  When the actual amount of data stored reaches the capacity of the file system, data cannot be written.  Unit: GiB.
* `storage_set_name` - (Optional, Available since v1.218.0) Save set identity, used to select a user-specified save set.
* `storage_type` - (Required, ForceNew) The storage media type. Value: STANDARD (default): STANDARD PERFORMANCE: PERFORMANCE type.
* `throughput_mode` - (Optional) The throughput mode. Value: Standard (default): Standard throughput Provisioned: preset throughput.
* `zone_id` - (Required, ForceNew) Zone Id, which is used to create file system resources to the specified zone.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the file system instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the File System.
* `delete` - (Defaults to 5 mins) Used when delete the File System.
* `update` - (Defaults to 5 mins) Used when update the File System.

## Import

DFS File System can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_file_system.example <id>
```