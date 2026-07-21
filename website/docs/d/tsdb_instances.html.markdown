---
subcategory: "Time Series Database (TSDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_tsdb_instances"
sidebar_current: "docs-alicloud-datasource-tsdb-instances"
description: |-
  Provides a list of Time Series Database (TSDB) Instances to the user.
---

# alicloud\_tsdb\_instances

This data source provides the Time Series Database (TSDB) Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.112.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_tsdb_instances" "example" {
  ids = ["example_value"]
}

output "first_tsdb_instance_id" {
  value = data.alicloud_tsdb_instances.example.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `app_key` - (Optional, ForceNew) The app key.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `engine_type` - (Optional, ForceNew) The engine type of instance. Enumerative: `tsdb_tsdb` refers to TSDB, `tsdb_influxdb` refers to TSDB for InfluxDB️.
* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `query_str` - (Optional, ForceNew) The query str.
* `status` - (Optional, ForceNew) Instance status, enumerative: `ACTIVATION`, `DELETED`, `CREATING`, `CLASS_CHANGING`, `LOCKED`.
* `status_list` - (Optional, ForceNew) The status list.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of TSDB Instances. Each element contains the following attributes:
	* `auto_renew` - Auto renew.
	* `cpu_number` - The cpu core number of instance.
	* `disk_category` - The disk type of instance. `cloud_ssd` refers to SSD disk, `cloud_efficiency` refers to efficiency disk,cloud_essd refers to ESSD PL1 disk.
	* `engine_type` - The engine type of instance. Enumerative: `tsdb_tsdb` refers to TSDB, `tsdb_influxdb` refers to TSDB for InfluxDB️.
	* `expired_time` - Instance expiration time.
	* `id` - The ID of the Instance.
	* `instance_alias` - The alias of the instance.
	* `instance_class` - The specification of the instance. 
	* `instance_id` - The ID of the instance.
	* `instance_storage` - The storage capacity of the instance. Unit: GB. For example, the value 50 indicates 50 GB.
	* `memory_size` - The memory size of instance.
	* `network_type` - Instance network type.
	* `payment_type` - The billing method. Valid values: `PayAsYouGo` and `Subscription`. The `PayAsYouGo` value indicates the pay-as-you-go method, and the `Subscription` value indicates the subscription method.
	* `status` - Instance status, enumerative: ACTIVATION,DELETED, CREATING,CLASS_CHANGING,LOCKED.
	* `vpc_connection_address` - The vpc connection address of instance.
	* `vpc_id` - The ID of the virtual private cloud (VPC) that is connected to the instance.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The zone ID of the instance.
