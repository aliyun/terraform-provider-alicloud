---
subcategory: "Distributed Relational Database Service (DRDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_drds_instances"
sidebar_current: "docs-alicloud-drds-instances"
description: |-
  Provides a collection of DRDS instances according to the specified filters.
---

# alicloud_drds_instance

 The `alicloud_drds_instance` data source provides a collection of DRDS instances available in Alibaba Cloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

~> **NOTE:** Available in 1.35.0+.

## Example Usage

 ```
data "alicloud_drds_instances" "drds_instances_ds" {
  name_regex = "drds-\\d+"
  ids        = ["drdsabc123456"]
}
output "first_db_instance_id" {
  value = "${data.alicloud_drds_instances.drds_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, Deprecated) A regex string to filter results by instance description. It is deprecated since v1.91.0 and will be removed in a future release, please use 'description_regex' instead.
* `description_regex` - (Optional, Available in 1.91.0+) A regex string to filter results by instance description.
* `ids` - (Optional) A list of DRDS instance IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

 * `ids` - A list of DRDS instance IDs.
 * `descriptions` - A list of DRDS descriptions. 
 * `instances` - A list of DRDS instances.
   * `id` - The ID of the DRDS instance.
   * `description` - The DRDS instance description.
   * `name` - The name of the RDS instance.
   * `status` - Status of the instance.
   * `type` - The DRDS Instance type.
   * `create_time` - Creation time of the instance.
   * `network_type` - `Classic` for public classic network or `VPC` for private network.
   * `zone_id` - Zone ID the instance belongs to.
   * `version` - The DRDS Instance version.
   * `connection_string` - (Available in 1.196.0+) The connection string of the DRDS instance.
   * `port` - (Available in 1.196.0+) The connection port of the DRDS instance.
