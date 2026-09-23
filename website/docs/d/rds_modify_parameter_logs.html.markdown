---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_modify_parameter_logs"
sidebar_current: "docs-alicloud-datasource-rds-modify_parameter_logs"
description: |-
  Provides a list of Rds Modify Parameter Logs to the user.
---

# alicloud\_rds\_modify\_parameter\_logs

This data source provides the Rds Modify Parameter Logs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.174.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_modify_parameter_logs" "example" {
  db_instance_id = "example_value"
  start_time     = "2022-06-04T13:56Z"
  end_time       = "2022-06-08T13:56Z"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The db instance id.
* `end_time` - (Required, ForceNew) The end time.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `start_time` - (Required, ForceNew) The start time.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `logs` - A list of Rds Modify Parameter Logs. Each element contains the following attributes:
  * `modify_time` - The time when the parameter was reconfigured. This value is a UNIX timestamp. Unit: milliseconds.
  * `new_parameter_value` - The new value of the parameter.
  * `old_parameter_value` - The original value of the parameter.
  * `parameter_name` - The name of the parameter.
  * `status` - The status of the new value specified for the parameter. Valid values:
    * **Applied**: The new value has taken effect.
    * **Syncing**: The new value is being applied and has not taken effect.
  
  