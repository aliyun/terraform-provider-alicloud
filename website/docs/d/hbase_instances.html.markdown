---
subcategory: "HBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbase_instances"
sidebar_current: "docs-alicloud-datasource-hbase-instances"
description: |-
    Provides a collection of HBase instances according to the specified filters.
---

# alicloud\_hbase\_instances

The `alicloud_hbase_instances` data source provides a collection of HBase instances available in Alicloud account.
Filters support regular expression for the instance name, ids or availability_zone.

-> **NOTE:**  Available in 1.67.0+

## Example Usage

```
data "alicloud_hbase_instances" "hbase" {
  name_regex        = "tf_testAccHBase"
  availability_zone = "cn-shenzhen-b"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name.
* `ids` - (Optional) The ids list of HBase instances
* `availability_zone` - (Optional) Instance availability zone.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - The ids list of HBase instances
* `names` - The names list of HBase instances
* `instances` - A list of HBase instances. Its every element contains the following attributes:
  * `id` - The ID of the HBase instance.
  * `name` - The name of the HBase instance.
  * `region_id` - Region ID the instance belongs to. 
  * `zone_id` - Zone ID the instance belongs to.
  * `engine` - the engine of the instance.
  * `engine_version` - the engine_version of the instance.
  * `network_type` - Classic network or VPC.
  * `master_instance_type` - hbase.n1.medium, hbase.n1.large, hbase.n1.2xlarge and so on.
  * `master_node_count` - the node count of master
  * `core_instance_type` - hbase.n1.medium, hbase.n1.large, hbase.n1.2xlarge and so on.
  * `core_node_count` - same with "core_instance_quantity"
  * `core_disk_type` - cloud_ssd or cloud_efficiency
  * `core_disk_size` - core node disk size, unit:GB.
  * `vpc_id` - VPC ID the instance belongs to.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `pay_type` - Billing method. Value options are `PostPaid` for  Pay-As-You-Go and `PrePaid` for yearly or monthly subscription.
  * `status` - Status of the instance.
  * `backup_status` - the Backup Status of the instance.
  * `create_time` - the create time of the instance.
  * `expire_time` - the expire time of the instance.