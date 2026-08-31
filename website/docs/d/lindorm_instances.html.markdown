---
subcategory: "Lindorm"
layout: "alicloud"
page_title: "Alicloud: alicloud_lindorm_instances"
sidebar_current: "docs-alicloud-datasource-lindorm-instances"
description: |-
  Provides a list of Lindorm Instances to the user.
---

# alicloud\_lindorm\_instances

This data source provides the Lindorm Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_lindorm_instances" "ids" {}
output "lindorm_instance_id_1" {
  value = data.alicloud_lindorm_instances.ids.instances.0.id
}

data "alicloud_lindorm_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "lindorm_instance_id_2" {
  value = data.alicloud_lindorm_instances.nameRegex.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `query_str` - (Optional, ForceNew) The query str, which can use `instance_name` keyword for fuzzy search.
* `status` - (Optional, ForceNew) Instance status, Valid values: `ACTIVATION`, `DELETED`, `CREATING`, `CLASS_CHANGING`, `LOCKED`, `INSTANCE_LEVEL_MODIFY`, `NET_MODIFYING`, `RESIZING`, `RESTARTING`, `MINOR_VERSION_TRANSING`.
* `support_engine` - (Optional, ForceNew) The support engine. Valid values: `1` to `7`. 

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of Lindorm Instances. Each element contains the following attributes:
	* `auto_renew` - AutoRenew.
	* `cold_storage` - The cold storage capacity of the instance. Unit: GB. 
	* `create_time` - The creation date of Instance.
	* `deletion_proection` - The deletion protection of instance.
	* `disk_category` - The disk type of instance. Valid values: `capacity_cloud_storage`, `cloud_efficiency`, `cloud_essd`, `cloud_ssd`.
	* `disk_usage` - The usage of disk.
	* `disk_warning_threshold` - The threshold of disk.
	* `engine_type` -  The type of Instance engine .
	* `expired_time` - The expiration time of Instance.
	* `file_engine_node_count` - The count of file engine.
	* `file_engine_specification` - The specification of file engine. Valid values: `lindorm.c.xlarge`.
	* `id` - The ID of the Instance.
	* `instance_id` - The ID of the instance.
	* `instance_name` - The name of the instance.
	* `instance_storage` - The storage capacity of the instance. Unit: GB. For example, the value 50 indicates 50 GB.
	* `ip_white_list` - The ip white list of instance.
	* `lts_node_count` - The count of lindorm tunnel service.
	* `lts_node_specification` - The specification of lindorm tunnel service. Valid values: `lindorm.g.2xlarge`, `lindorm.g.xlarge`.
	* `network_type` - Instance network type, enumerative.VPC.
	* `payment_type` - The billing method. Valid values: `PayAsYouGo` and `Subscription`.
	* `phoenix_node_count` - The count of phoenix.
	* `phoenix_node_specification` - The specification of phoenix. Valid values: `lindorm.c.2xlarge`, `lindorm.c.4xlarge`, `lindorm.c.8xlarge`, `lindorm.c.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
	* `resource_owner_id` - The owner id of resource.
	* `search_engine_node_count` - The count of search engine.
	* `search_engine_specification` - The specification of search engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
	* `service_type` - The service type of Instance, Valid values:  `lindorm`, `serverless_lindorm`, `lindorm_standalone`.
	* `status` - The status of Instance, enumerative: Valid values: `ACTIVATION`, `DELETED`, `CREATING`, `CLASS_CHANGING`, `LOCKED`, `INSTANCE_LEVEL_MODIFY`, `NET_MODIFYING`, `RESIZING`, `RESTARTING`, `MINOR_VERSION_TRANSING`.
	* `table_engine_node_count` - The count of table engine.
	* `table_engine_specification` - The specification of  table engine. Valid values: `lindorm.c.2xlarge`, `lindorm.c.4xlarge`, `lindorm.c.8xlarge`, `lindorm.c.xlarge`, `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
	* `time_series_engine_node_count` - The count of time series engine.
	* `time_serires_engine_specification` - The specification of time series engine. Valid values: `lindorm.g.2xlarge`, `lindorm.g.4xlarge`, `lindorm.g.8xlarge`, `lindorm.g.xlarge`.
	* `vpc_id` - The ID of the virtual private cloud (VPC) that is connected to the instance.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The zone ID of the instance.
