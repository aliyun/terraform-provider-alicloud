---
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instances"
sidebar_current: "docs-alicloud-datasource-db-instances"
description: |-
    Provides a collection of RDS instances according to the specified filters.
---

# alicloud\_db\_instances

The `alicloud_db_instances` data source provides a collection of RDS instances available in Alibaba Cloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```
data "alicloud_db_instances" "db_instances_ds" {
  name_regex = "data-\\d+"
  status     = "Running"
  tags       = <<EOF
{
  "type": "database",
  "size": "tiny"
}
EOF
}

output "first_db_instance_id" {
  value = "${data.alicloud_db_instances.db_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by instance name.
* `engine` - (Optional) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
* `status` - (Optional) Status of the instance.
* `db_type` - (Optional) `Primary` for primary instance, `ReadOnly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
* `vpc_id` - (Optional) Used to retrieve instances belong to specified VPC.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `connection_mode` - (Optional) `Standard` for standard access mode and `Safe` for high security access mode.
* `tags` - (Optional) Query the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of RDS instances. Each element contains the following attributes:
  * `id` - The ID of the RDS instance.
  * `name` - The name of the RDS instance.
  * `charge_type` - Billing method. Value options: `Postpaid` for Pay-As-You-Go and `Prepaid` for subscription.
  * `db_type` - `Primary` for primary instance, `ReadOnly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - Creation time of the instance.
  * `expire_time` - Expiration time. Pay-As-You-Go instances never expire.
  * `status` - Status of the instance.
  * `engine` - Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
  * `engine_version` - Database version.
  * `net_type` - `Internet` for public network or `Intranet` for private network.
  * `connection_mode` - `Standard` for standard access mode and `Safe` for high security access mode.
  * `instance_type` - Sizing of the RDS instance.
  * `availability_zone` - Availability zone.
  * `master_instance_id` - ID of the primary instance. If this parameter is not returned, the current instance is a primary instance.
  * `guard_instance_id` - If a disaster recovery instance is attached to the current instance, the ID of the disaster recovery instance applies.
  * `temp_instance_id` - If a temporary instance is attached to the current instance, the ID of the temporary instance applies.
  * `readonly_instance_ids` - A list of IDs of read-only instances attached to the primary instance.
  * `vpc_id` - ID of the VPC the instance belongs to.
  * `vswitch_id` - ID of the VSwitch the instance belongs to.
