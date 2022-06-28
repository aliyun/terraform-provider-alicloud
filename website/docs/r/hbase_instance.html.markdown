---
subcategory: "HBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbase_instance"
sidebar_current: "docs-alicloud-resource-hbase-instance"
description: |-
  Provides a HBase instance resource.
---

# alicloud\_hbase\_instance

Provides a HBase instance resource supports replica set instances only. The HBase provides stable, reliable, and automatic scalable database services.
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://help.aliyun.com/product/49055.html)

-> **NOTE:**  Available in 1.67.0+

-> **NOTE:**  The following regions don't support create Classic network HBase instance.
[`cn-hangzhou`,`cn-shanghai`,`cn-qingdao`,`cn-beijing`,`cn-shenzhen`,`ap-southeast-1a`,.....]
The official website mark  more regions. or you can call [DescribeRegions](https://help.aliyun.com/document_detail/144489.html)

-> **NOTE:**  Create HBase instance or change instance type and storage would cost 15 minutes. Please make full preparation

## Example Usage

### Create a hbase instance

```terraform
resource "alicloud_hbase_instance" "default" {
  name                   = "tf_testAccHBase_vpc"
  zone_id                = "cn-shenzhen-b"
  vswitch_id             = "vsw-123456"
  engine                 = "hbaseue"
  engine_version         = "2.0"
  master_instance_type   = "hbase.sn1.large"
  core_instance_type     = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type         = "cloud_efficiency"
  core_disk_size         = 400
  pay_type               = "PostPaid"
  cold_storage_size      = 0
}
```

this is a example for class netType instance. you can find more detail with the examples/hbase dir.

## Argument Reference

The following arguments are supported:

* `name` - (Required) HBase instance name. Length must be 2-128 characters long. Only Chinese characters, English letters, numbers, period (.), underline (_), or dash (-) are permitted. 
* `zone_id` - (Optional, ForceNew) The Zone to launch the HBase instance. If vswitch_id is not empty, this zone_id can be "" or consistent.
* `engine` - (Optional, ForceNew) Valid values are "hbase/hbaseue/bds". The following types are supported after v1.73.0: `hbaseue` and `bds`. Single hbase instance need to set engine=hbase, core_instance_quantity=1.
* `engine_version` - (Required, ForceNew) HBase major version. hbase:1.1/2.0, hbaseue:2.0, bds:1.0, unsupport other engine temporarily. Value options can refer to the latest docs [CreateInstance](https://help.aliyun.com/document_detail/144607.html).
* `master_instance_type` - (Required) Instance specification. See [Instance specifications](https://help.aliyun.com/document_detail/53532.html), or you can call describeInstanceType api.
* `core_instance_type` - (Required) Instance specification. See [Instance specifications](https://help.aliyun.com/document_detail/53532.html), or you can call describeInstanceType api.
* `core_instance_quantity`- (Optional) Default=2, [1-200]. If core_instance_quantity > 1, this is cluster's instance. If core_instance_quantity = 1, this is a single instance.
* `core_disk_type`- (Optional, ForceNew) Valid values are `cloud_ssd`, `cloud_essd_pl1`, `cloud_efficiency`, `local_hdd_pro`, `local_ssd_pro`，``, local_disk size is fixed. When engine=bds, no need to set disk type(or empty string).
* `core_disk_size` - (Optional) User-defined HBase instance one core node's storage. Valid when engine=hbase/hbaseue. Bds engine no need core_disk_size, space.Unit: GB. Value range:
  - Custom storage space, value range: [20, 64000].
  - Cluster [400, 64000], step:40-GB increments.
  - Single [20-500GB], step:1-GB increments.
* `pay_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, System default to `PostPaid`. You can also convert PostPaid to PrePaid. And support convert PrePaid to PostPaid from 1.115.0+.
* `duration` - (Optional, ForceNew) 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, valid when pay_type = PrePaid,  unit: month. 12, 24, 36 mean 1, 2, 3 years.
* `auto_renew` - (Optional, ForceNew) Valid values are `true`, `false`, system default to `false`, valid when pay_type = PrePaid.
* `vswitch_id` - (Optional, ForceNew) If vswitch_id is not empty, that mean net_type = vpc and has a same region. If vswitch_id is empty, net_type=classic. Intl site not support classic network.
* `cold_storage_size` - (Optional) 0 or [800, 1000000], step:10-GB increments. 0 means is_cold_storage = false. [800, 1000000] means is_cold_storage = true.
* `maintain_start_time` - (Optional, Available in 1.73.0) The start time of the operation and maintenance time period of the instance, in the format of HH:mmZ (UTC time), for example 02:00Z.
* `maintain_end_time` - (Optional, Available in 1.73.0) The end time of the operation and maintenance time period of the instance, in the format of HH:mmZ (UTC time), for example 04:00Z.
* `deletion_protection` - (Optional, Available in 1.73.0) The switch of delete protection. True: delete protect, False: no delete protect. You must set false when you want to delete cluster.
* `immediate_delete_flag` - (Optional, Available in 1.109.0) The switch of delete immediate. True: delete immediate, False: delete delay. You will not found the cluster no matter set true or false.
* `tags` - (Optional, Available in 1.73.0) A mapping of tags to assign to the resource.
* `account` - (Optional, Available in 1.105.0+) The account of the cluster web ui. Size [0-128].
* `password` - (Optional, Available in 1.105.0+) The password of the cluster web ui account. Size [0-128].
* `ip_white` - (Optional, Available in 1.105.0+) The white ip list of the cluster.
* `security_groups` - (Optional, Available in 1.105.0+) The security group resource of the cluster.


-> **NOTE:** Now only instance name can be change. The others(instance_type, disk_size, core_instance_quantity and so on) will be supported in the furture.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HBase.
* `master_instance_quantity` - Count nodes of the master node.
* `ui_proxy_conn_addrs` - (Available in 1.105.0+) The Web UI proxy addresses of the cluster.
* `zk_conn_addrs` - (Available in 1.105.0+) The zookeeper addresses of the cluster.
* `slb_conn_addrs` - (Available in 1.105.0+) The slb service addresses of the cluster.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the HBase instance (until it reaches the initial `ACTIVATION` status).
* `update` - (Defaults to 60 mins) Used when updating the HBase instance (until it reaches the initial `ACTIVATION` status).
* `delete` - (Defaults to 30 mins) Used when terminating the HBase instance. 

## Import

HBase can be imported using the id, e.g.

```
$ terraform import alicloud_hbase_instance.example hb-wz96815u13k659fvd
```
