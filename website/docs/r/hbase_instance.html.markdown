

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
You can see detail product introduction [here](https://help.aliyun.com/product/49055.html?spm=5176.124785.766162.1.67273094TcUuke)

-> **NOTE:**  Available in 1.67.0+

-> **NOTE:**  The following regions don't support create Classic network HBase instance.
[`cn-hangzhou`,`cn-shanghai`,`cn-qingdao`,`cn-beijing`,`cn-shenzhen`,.....]
the official website mark  more regions. or you can call [DescribeRegions](https://help.aliyun.com/document_detail/144489.html?spm=a2c4g.11186623.2.31.704e6a7dxQ6kXN)

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
  pay_type = "Postpaid"
  is_cold_storage = "false"
  security_ip_list = ["127.0.0.1", "127.0.0.2"]
}
```

this is a example for class netType instance. you can find more detail with the examples/hbase dir.

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) The Zone to launch the HBase instance. if vswitch_id is not empty, this zone_id can be "" or consistent.
* `engine` - (Optional) default value = hbase. examples: hbase,hbaseue,serverlesshbase,spark,bds.
* `engine_version` - (Required) hbase major version. hbase:1.1, 2.0; serverless:2.0, bds:1.0. Value options can refer to the latest docs [CreateInstance](https://help.aliyun.com/document_detail/144607.html?spm=a2c4g.11186623.6.768.5d3d2767CyQfIS).
* `master_instance_type`、`core_instance_type` - (Required) Instance specification. see [Instance specifications](https://help.aliyun.com/document_detail/53532.html?spm=a2c4g.11186623.6.547.56cb6233fyWh0Q). or you can call describeInstanceType api.
* `core_instance_quantity`- (Optional. default=2) if core_instance_quantity > 1,this is a cluster.  if core_instance_quantity > 1,this is a single. 
* `core_disk_type`-  (Required) Valid values are `cloud_ssd`, `cloud_efficiency`, `local_hdd_pro`, `local_ssd_pro`. local_disk size is fixed.
* `core_disk_size` -  (Optional) User-defined HBase Cluster one core node's storage space.Unit: GB. Value range:
  - Custom storage space; value range: [100,2000]
  - 10-GB increments. 
* `pay_type` - (Optional) Valid values are `Prepaid`, `Postpaid`,System default to `Postpaid`.
* `pricing_cycle` - (Optional) `year` or `month`.  valid when pay_type = month.
* `duration` - (Optional) use with pricing_cycle, valid when pay_type = Prepaid. unit: pricing_cycle.
* `auto_renew` - (Optional) `true`, `false`, System default to `false`, valid when pay_type = Prepaid.
* `vswitch_id` - (Optional) if vswitch_id is not empty, that mean net_type = vpc and has a same region_id. if vswitch_id is empty, net_type_classic
* `is_cold_storage` - (Optional) `true`, `false`, System default to `false`.
* `security_ip_list` - (Optional) default to '127.0.0.1'.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HBase.

### Timeouts

-> **NOTE:** Available in 1.67.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the HBase instance (until it reaches the initial `ACTIVATION` status). 
* `update` - (Defaults to 30 mins) Used when updating the HBase instance (until it reaches the initial `ACTIVATION` status). 
* `delete` - (Defaults to 30 mins) Used when terminating the HBase instance. 

## Import

HBase can be imported using the id, e.g.

```
$ terraform import alicloud_hbase_instance.example hb-wz96815u13k659fvd
```