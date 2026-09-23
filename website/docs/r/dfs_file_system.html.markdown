---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_file_system"
description: |-
  Provides a Alicloud Apsara File Storage for HDFS (DFS) File System resource.
---

# alicloud_dfs_file_system

Provides a Apsara File Storage for HDFS (DFS) File System resource.



For information about Apsara File Storage for HDFS (DFS) File System and how to use it, see [What is File System](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dfs_file_system&exampleId=b26fd6fb-f316-1cc8-a3c2-ff69b2700582a9599cfd&activeTab=example&spm=docs.r.dfs_file_system.0.b26fd6fbf3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

resource "alicloud_dfs_file_system" "default" {
  storage_type                     = "PERFORMANCE"
  zone_id                          = "cn-hangzhou-b"
  protocol_type                    = "PANGU"
  description                      = var.name
  file_system_name                 = var.name
  throughput_mode                  = "Provisioned"
  space_capacity                   = "1024"
  provisioned_throughput_in_mi_bps = "512"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dfs_file_system&spm=docs.r.dfs_file_system.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `data_redundancy_type` - (Optional) Redundancy mode of the file system. Value:
  - LRS (default): Local redundancy.
  - ZRS: Same-City redundancy. When ZRS is selected, zoneId is a string consisting of multiple zones that are expected to be redundant in the same city, for example,  'zoneId1,zoneId2 '.
* `dedicated_cluster_id` - (Optional) Dedicated cluster id, which is used to support scenarios such as group cloud migration.
* `description` - (Optional) The description of the file system resource. No more than 32 characters in length.
* `file_system_name` - (Required) The file system name. The naming rules are as follows: The length is 6~64 characters. Globally unique and cannot be an empty string. English letters are supported and can contain numbers, underscores (_), and dashes (-).
* `partition_number` - (Optional, Available since v1.218.0) Save set sequence number, the user selects the content of the specified sequence number in the Save set.
* `protocol_type` - (Required, ForceNew) The protocol type. Value: `HDFS`, `PANGU`. 
* `provisioned_throughput_in_mi_bps` - (Optional) Provisioned throughput. This parameter is required when ThroughputMode is set to Provisioned. Unit: MB/s Value range: 1~5120.
* `space_capacity` - (Required) File system capacity.  When the actual amount of data stored reaches the capacity of the file system, data cannot be written.  Unit: GiB.
* `storage_set_name` - (Optional, Available since v1.218.0) Save set identity, used to select a user-specified save set.
* `storage_type` - (Required, ForceNew) The storage media type. Value: STANDARD (default): STANDARD PERFORMANCE: PERFORMANCE type.
* `throughput_mode` - (Optional) The throughput mode. Value: Standard (default): Standard throughput Provisioned: preset throughput.
* `zone_id` - (Optional, ForceNew) Zone Id, which is used to create file system resources to the specified zone.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the file system instance.
* `region_id` - (Available since v1.242.0) The region ID of the File System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the File System.
* `delete` - (Defaults to 5 mins) Used when delete the File System.
* `update` - (Defaults to 5 mins) Used when update the File System.

## Import

Apsara File Storage for HDFS (DFS) File System can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_file_system.example <id>
```
