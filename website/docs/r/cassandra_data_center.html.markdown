---
subcategory: "Cassandra"
layout: "alicloud"
page_title: "Alicloud: alicloud_cassandra_data_center"
sidebar_current: "docs-alicloud-resource-cassandra-data-center"
description: |-
  Provides a Cassandra dataCenter resource.
---

# alicloud\_cassandra\_data\_center

Provides a Cassandra dataCenter resource supports replica set dataCenters only. The Cassandra provides stable, reliable, and automatic scalable database services. 
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://www.alibabacloud.com/help/product/49055.htm).

-> **NOTE:**  Available in 1.88.0+.

-> **NOTE:**  Create a cassandra dataCenter need a clusterId,so need create a cassandra cluster first.

-> **NOTE:**  The following regions support create Vpc network Cassandra cluster.
The official website mark  more regions. Or you can call [DescribeRegions](https://help.aliyun.com/document_detail/157540.html).

-> **NOTE:**  Create Cassandra dataCenter or change dataCenter type and storage would cost 30 minutes. Please make full preparation.

## Example Usage

### Create a cassandra dataCenter

```
resource "alicloud_cassandra_cluster" "default" {
  cluster_name        = "cassandra-cluster-name-tf"
  data_center_name    = "dc-1"
  auto_renew          = "false"
  instance_type       = "cassandra.c.large"
  major_version       = "3.11"
  node_count          = "2"
  pay_type            = "PayAsYouGo"
  vswitch_id          = "vsw-xxxx1"
  disk_size           = "160"
  disk_type           = "cloud_ssd"
  maintain_start_time = "18:00Z"
  maintain_end_time   = "20:00Z"
  ip_white            = "127.0.0.1"
}

resource "alicloud_cassandra_data_center" "default" {
  cluster_id       = alicloud_cassandra_cluster.default.id
  data_center_name = "dc-2"
  auto_renew       = "false"
  instance_type    = "cassandra.c.large"
  node_count       = "2"
  pay_type         = "PayAsYouGo"
  vswitch_id       = "vsw-xxxx2"
  disk_size        = "160"
  disk_type        = "cloud_ssd"
}
```

This is a example for class netType dataCenter. You can find more detail with the examples/cassandra_data_center dir.

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) Cassandra cluster id of dataCenter-2 belongs to.  
* `data_center_name` - (Required) Cassandra dataCenter-2 name. Length must be 2~128 characters long. Only Chinese characters, English letters, numbers, period `.`, underline `_`, or dash `-` are permitted. 
* `zone_id` - (Optional, ForceNew) The Zone to launch the Cassandra dataCenter-2. If vswitch_id is not empty, this zone_id can be "" or consistent.
* `instance_type` - (Required) Instance specification. See [Instance specifications](https://help.aliyun.com/document_detail/157445.html). Or you can call describeInstanceType api.
* `node_count`- (Optional) The node count of Cassandra dataCenter-2, default to 2. 
* `disk_type`-  (Required) The disk type of Cassandra dataCenter-2. Valid values are `cloud_ssd`, `cloud_efficiency`, `local_hdd_pro`, `local_ssd_pro`, local_disk size is fixed.
* `disk_size` -  (Optional) User-defined Cassandra dataCenter one core node's storage space.Unit: GB. Value range:
  - Custom storage space; value range: [160, 2000].
  - 80-GB increments. 
* `pay_type` - (Optional, ForceNew) The pay type of Cassandra dataCenter-2. Valid values are `Subscription`, `PayAsYouGo`. System default to `PayAsYouGo`.
* `auto_renew` - (Optional, ForceNew) Auto renew of dataCenter-2,`true` or `false`. System default to `false`, valid when pay_type = Subscription.
* `auto_renew_period` - (Optional, ForceNew) Period of dataCenter-2 auto renew, if auto renew is `true`, one of `1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60`, valid when pay_type = Subscription. Unit: month.
* `vswitch_id` - (Optional, ForceNew) The vswitch_id of dataCenter-2, mast different of vswitch_id(dc-1), can not empty.

-> **NOTE:** Now data_center_name,instance_type,node_count,disk_type,disk_size can be change. The others(auto_renew, auto_renew_period and so on) will be supported in the furture.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cassandra.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the Cassandra dataCenter (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the Cassandra dataCenter (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 5 mins) Used when terminating the Cassandra dataCenter. 

## Import

If you need full function, please import Cassandra cluster first.
Cassandra dataCenter can be imported using the dcId:clusterId, e.g.

```
$ terraform import alicloud_cassandra_data_center.dc_2 cn-shenxxxx-x:cds-wz933ryoaurxxxxx
```
