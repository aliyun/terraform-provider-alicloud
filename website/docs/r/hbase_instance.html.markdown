---
subcategory: "HBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbase_instance"
sidebar_current: "docs-alicloud-resource-hbase-instance"
description: |-
  Provides a HBase instance resource.
---

# alicloud\_hbase\_instance

Provides a HBase instance resource supports replica set instances only. the HBase provides stable, reliable, and automatic scalable database services. 
It offers a full range of database solutions, such as disaster recovery, backup, recovery, monitoring, and alarms.
You can see detail product introduction [here](https://help.aliyun.com/product/49055.html)

-> **NOTE:**  Available in 1.67.0+

-> **NOTE:**  The following regions don't support create Classic network HBase instance.
[`cn-hangzhou`,`cn-shanghai`,`cn-qingdao`,`cn-beijing`,`cn-shenzhen`,.....]
the official website mark  more regions. or you can call [DescribeRegions](https://help.aliyun.com/document_detail/144489.html)

-> **NOTE:**  Create HBase instance or change instance type and storage would cost 15 minutes. Please make full preparation

## Example Usage

### Create a hbase instance

```
resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic"
  zone_id = "cn-shenzhen-b"
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
}
```

this is a example for class netType instance. you can find more detail with the examples/hbase dir.

## Argument Reference

The following arguments are supported:

* `name` - (Required) HBase instance name. Length must be 2-128 characters long. Only Chinese characters, English letters, numbers, period (.), underline (_), or dash (-) are permitted. 
* `zone_id` - (Optional, ForceNew) The Zone to launch the HBase instance. if vswitch_id is not empty, this zone_id can be "" or consistent.
* `engine` - (Optional, ForceNew) default value = hbase. examples: hbase,hbaseue,serverlesshbase,spark,bds.
* `engine_version` - (Required, ForceNew) hbase major version. hbase:1.1, 2.0; serverless:2.0, bds:1.0. Value options can refer to the latest docs [CreateInstance](https://help.aliyun.com/document_detail/144607.html).
* `master_instance_type`、`core_instance_type` - (Required, ForceNew) Instance specification. see [Instance specifications](https://help.aliyun.com/document_detail/53532.html). or you can call describeInstanceType api.
* `core_instance_quantity`- (Optional. ForceNew) default=2. if core_instance_quantity > 1,this is cluster's instance.  if core_instance_quantity = 1,this is a single instance. 
* `core_disk_type`-  (Required, ForceNew) Valid values are `cloud_ssd`, `cloud_efficiency`, `local_hdd_pro`, `local_ssd_pro`. local_disk size is fixed.
* `core_disk_size` -  (Optional, ForceNew) User-defined HBase instance one core node's storage space.Unit: GB. Value range:
  - Custom storage space; value range: [100,2000]
  - 10-GB increments. 
* `pay_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`.
* `duration` - (Optional, ForceNew) 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60, valid when pay_type = PrePaid. unit: month.
* `auto_renew` - (Optional, ForceNew) `true`, `false`, System default to `false`, valid when pay_type = PrePaid.
* `vswitch_id` - (Optional, ForceNew) if vswitch_id is not empty, that mean net_type = vpc and has a same region. if vswitch_id is empty, net_type_classic
* `cold_storage_size` - (Optional, ForceNew) 0 or 0+. 0 means is_cold_storage = false. 0+ means is_cold_storage = true

-> **NOTE:** now only instance name can be change. the others(instance_type, disk_size, core_instance_quantity and so on) will be supported in the furture.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HBase.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the HBase instance (until it reaches the initial `ACTIVATION` status). 
* `delete` - (Defaults to 30 mins) Used when terminating the HBase instance. 

## Import

HBase can be imported using the id, e.g.

```
$ terraform import alicloud_hbase_instance.example hb-wz96815u13k659fvd
```