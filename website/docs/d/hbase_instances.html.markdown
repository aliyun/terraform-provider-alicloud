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

```terraform
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
* `tags` - (Optional, Available in 1.73.0) A mapping of tags to assign to the resource.
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
  * `engine` - The engine of the instance.
  * `engine_version` - The engine_version of the instance.
  * `network_type` - Classic network or VPC.
  * `master_instance_type` - Like hbase.sn2.2xlarge, hbase.sn2.4xlarge, hbase.sn2.8xlarge and so on.
  * `master_node_count` - The node count of master
  * `core_instance_type` - Like hbase.sn2.2xlarge, hbase.sn2.4xlarge, hbase.sn2.8xlarge and so on.
  * `core_node_count` - Same with "core_instance_quantity"
  * `core_disk_type` - Cloud_ssd or cloud_efficiency
  * `core_disk_size` - Core node disk size, unit:GB.
  * `vpc_id` - VPC ID the instance belongs to.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `pay_type` - Billing method. Value options are `PostPaid` for  Pay-As-You-Go and `PrePaid` for yearly or monthly subscription.
  * `status` - Status of the instance.
  * `backup_status` - The Backup Status of the instance.
  * `created_time` - The created time of the instance.
  * `expire_time` - The expire time of the instance.
  * `deletion_protection` - The switch of delete protection.
  * `tags` - A mapping of tags to assign to the resource.
