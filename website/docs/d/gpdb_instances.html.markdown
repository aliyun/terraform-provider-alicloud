---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_instances"
sidebar_current: "docs-alicloud-datasource-gpdb-instances"
description: |-
    Provides a collection of AnalyticDB for PostgreSQL instances according to the specified filters.
---

# alicloud\_gpdb\_instances

The `alicloud_gpdb_instances` data source provides a collection of AnalyticDB for PostgreSQL instances available in Alicloud account.
Filters support regular expression for the instance name or availability_zone.

-> **NOTE:**  Available in 1.47.0+

## Example Usage

```
data "alicloud_gpdb_instances" "gpdb" {
  availability_zone = "cn-beijing-c"
  name_regex        = "gp-.+\\d+"
  output_file       = "instances.txt"
}

output "instance_id" {
  value = data.alicloud_gpdb_instances.gpdb.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `name_regex` - (Optional) A regex string to apply to the instance name.
* `availability_zone` - (Optional) Instance availability zone.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - The ids list of AnalyticDB for PostgreSQL instances.
* `names` - The names list of AnalyticDB for PostgreSQL instance.
* `instances` - A list of AnalyticDB for PostgreSQL instances. Its every element contains the following attributes:
  * `id` - The instance id.
  * `description` - The description of an instance.
  * `charge_type` - Billing method. Value options are `PostPaid` for  Pay-As-You-Go and `PrePaid` for yearly or monthly subscription.
  * `region_id` - Region ID the instance belongs to.
  * `availability_zone` - Instance availability zone.
  * `creation_time` - The time when you create an instance. The format is YYYY-MM-DDThh:mm:ssZ, such as 2011-05-30T12:11:4Z.
  * `status` - Status of the instance.
  * `engine` - Database engine type. Supported option is `gpdb`.
  * `engine_version` - Database engine version.
  * `network_type` - Classic network or VPC.
  * `instance_class` - The group type.
  * `instance_group_count` - The number of groups.

