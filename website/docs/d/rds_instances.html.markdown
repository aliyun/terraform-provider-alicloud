---
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_instances"
sidebar_current: "docs-alicloud-datasource-rds-instances"
description: |-
    Provides a collection of RDS instances according to the specified filters.
---

# alicloud\_rds\_instances

The `alicloud_rds_instances` data source provides a collection of RDS instances available in Alicloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```
data "alicloud_rds_instances" "rds" {
  name_regex = "data-\\d+"
  status     = "Running"
  tags       = <<EOF
{
  "type": "database",
  "size": "tiny"
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name.
* `engine` - (Optional) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
* `status` - (Optional) Status of the instance.
* `instance_type` - (Optional) `Primary` for primary instance, `ReadOnly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
* `instance_network_type` - (Optional) Either `Classic` or `VPC` network.
* `vpc_id` - (Optional) Used to retrieve instances belong to specified VPC.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `connection_mode` - (Optional) `Standard` for standard access mode and `Safe` for high security access mode.
* `tags` - (Optional) Query the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.
* `owner_account` - (Optional) Owner of Alicloud account.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of RDS instances. Its every element contains the following attributes:
  * `id` - The ID of the RDS instance.
  * `name` - The name of the RDS instance.
  * `pay_type` - Billing method. Value options: `Postpaid` for  Pay-As-You-Go and `Prepaid` for subscription.
  * `instance_type` - `Primary` for primary instance, `ReadOnly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - Creation time of the instance.
  * `expire_time` - Expiration time. Pay-As-You-Go instances are never expire.
  * `status` - Status of the instance.
  * `engine` - Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
  * `engine_version` - Database version.
  * `db_instance_net_type` - `Internet` for public network or `Intranet` for private network.
  * `connection_mode` - `Standard` for standard access mode and `Safe` for high security access mode.
  * `lock_mode` - `Unlock` normal operation, `ManualLock` locked when manually triggered, `LockByExpiration` automatically locked upon expiration, `LockByRestoration` automatically locked before instance rollback, `LockByDiskQuota` automatically locked when the instance space is full.
  * `lock_reason` - Reason why the instance is locked.
  * `db_instance_class` - Sizing of the RDS instance.
  * `instance_network_type` - Either `Classic` or `VPC` network.
  * `vpc_cloud_instance_id` - VPC cloud instance ID.
  * `zone_id` - Availability zone.
  * `multi_or_single` - `Multi` or `single` instance.
  * `master_instance_id` - ID of the primary instance. If this parameter is not returned, the current instance is a primary instance.
  * `guard_db_instance_id` - If a disaster recovery instance is attached to the current instance, the ID of the disaster recovery instance applies.
  * `temp_db_instance_id` - If a temporary instance is attached to the current instance, the ID of the temporary instance applies.
  * `readonly_db_instance_ids` - A list of the ID's of read-only instances attached to the primary instance.
  * `vpc_id` - VPC ID the instance belongs to.
  * `vswitch_id` - VSwitch ID the instance belongs to.
  * `replicate_id` - Replica ID.
  * `resource_group_id` - Resource group the instance belongs to.
