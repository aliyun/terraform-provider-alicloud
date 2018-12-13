---
layout: "alicloud"
page_title: "Alicloud: alicloud_drds_instances"
sidebar_current: "docs-alicloud-drds-instances"
description: |-
  Provides a collection of DRDS instances according to the specified filters.
---
 # alicloud_drds_instance
 The `alicloud_drds_instance` data source provides a collection of DRDS instances available in Alibaba Cloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.
 ## Example Usage
 ```
data "alicloud_drds_instances" "drds_instances_ds" {
  name_regex = "drds-\\d+"
  regionId     = "cn-hangzhou"
}
 output "first_db_instance_id" {
  value = "${data.alicloud_drds_instances.drds_instances_ds.instances.0.drdsInstanceId}"
}
```
 ## Argument Reference
 The following arguments are supported:
 * `name_regex` - A regex string to filter results by instance name.
* `regionId` - Region ID the DRDS instance belongs to.
 ## Attributes Reference
 The following attributes are exported in addition to the arguments listed above:
 * `instances` - A list of RDS instances. Each element contains the following attributes:
  * `drdsInstanceId` - The ID of the DRDS instance.
  * `name` - The name of the RDS instance.
  * `status` - Status of the instance.
  * `type` - The DRDS Instance type.
  * `createTime` - Creation time of the instance.
  * `networkType` - `Classic` for public classic network or `VPC` for private network.
  * `zoneId` - Zone ID the instance belongs to.
  * `version` - The DRDS Instance version.
  * `vips` - The DRDS instance vips info