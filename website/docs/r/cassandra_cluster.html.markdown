---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_cluster"
sidebar_current: "docs-alicloud-resource-cassandra-cluster"
description: |-
  Provides a Cassandra cluster resource.
---

# alicloud\_cassandra\_cluster

Provides a Cassandra cluster resource supports replica set clusters only. The Cassandra provides stable, reliable, and automatic scalable database services. 
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://www.alibabacloud.com/help/product/49055.htm).

-> **NOTE:**  Available in 1.88.0+.

-> **NOTE:**  The following regions support create Vpc network Cassandra cluster.
The official website mark more regions. Or you can call [DescribeRegions](https://help.aliyun.com/document_detail/157540.html).

-> **NOTE:**  Create Cassandra cluster or change cluster type and storage would cost 30 minutes. Please make full preparation.

## Example Usage

### Create a cassandra cluster

```
resource "alicloud_cassandra_cluster" "default" {
  cluster_name = "cassandra-cluster-name-tf"
  data_center_name = "dc-1"
  auto_renew = "false"
  instance_type = "cassandra.c.large"
  major_version = "3.11"
  node_count = "2"
  pay_type = "PayAsYouGo"
  vswitch_id = "vsw-xxxx"
  disk_size = "160"
  disk_type = "cloud_ssd"
  maintain_start_time = "18:00Z"
  maintain_end_time = "20:00Z"
  ip_white = "127.0.0.1"
}
```

This is a example for class netType cluster. You can find more detail with the examples/cassandra_cluster dir.

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) Cassandra cluster name. Length must be 2~128 characters long. Only Chinese characters, English letters, numbers, period `.`, underline `_`, or dash `-` are permitted. 
* `data_center_name` - (Required) Cassandra dataCenter-1 name. Length must be 2~128 characters long. Only Chinese characters, English letters, numbers, period `.`, underline `_`, or dash `-` are permitted. 
* `zone_id` - (Optional, ForceNew) The Zone to launch the Cassandra cluster. If vswitch_id is not empty, this zone_id can be "" or consistent.
* `major_version` - (Required, ForceNew) Cassandra major version. Now only support version `3.11`.
* `instance_type` - (Required) Instance specification. See [Instance specifications](https://help.aliyun.com/document_detail/157445.html). Or you can call describeInstanceType api.
* `node_count`- (Optional) The node count of Cassandra dataCenter-1 default to 2. 
* `disk_type`-  (Required) The disk type of Cassandra dataCenter-1. Valid values are `cloud_ssd`, `cloud_efficiency`, `local_hdd_pro`, `local_ssd_pro`, local_disk size is fixed.
* `disk_size` -  (Optional) User-defined Cassandra dataCenter-1 one node's storage space.Unit: GB. Value range:
  - Custom storage space; value range: [160, 2000].
  - 80-GB increments. 
* `pay_type` - (Optional, ForceNew) The pay type of Cassandra dataCenter-1. Valid values are `Subscription`, `PayAsYouGo`,System default to `PayAsYouGo`.
* `auto_renew` - (Optional, ForceNew) Auto renew of dataCenter-1,`true` or `false`. System default to `false`, valid when pay_type = PrePaid.
* `auto_renew_period` - (Optional, ForceNew) Period of dataCenter-1 auto renew, if auto renew is `true`, one of `1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60`, valid when pay_type = Subscription. Unit: month.
* `vswitch_id` - (Optional, ForceNew) The vswitch_id of dataCenter-1, can not empty.
* `maintain_start_time` - (Optional) The start time of the operation and maintenance time period of the cluster, in the format of HH:mmZ (UTC time).
* `maintain_end_time` - (Optional) The end time of the operation and maintenance time period of the cluster, in the format of HH:mmZ (UTC time).
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `ip_white` - (Optional) Set the instance's IP whitelist in VPC network.
* `security_groups` - (Optional)  A list of security group ids to associate with.

-> **NOTE:** Now cluster_name,data_center_name,instance_type,node_count,disk_type,disk_size,maintain_start_time,maintain_end_time,tags,ip_white,security_groups can be change. The others(auto_renew, auto_renew_period and so on) will be supported in the furture.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cassandra.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the Cassandra cluster (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the Cassandra cluster (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 5 mins) Used when terminating the Cassandra cluster. 

## Import

Cassandra cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cassandra_cluster.example cds-wz9sr400dd7xxxxx
```
