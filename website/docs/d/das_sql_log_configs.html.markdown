---
subcategory: "DAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_das_sql_log_configs"
description: |-
  Provides a list of DAS Sql Log Configs to the user.
---

# alicloud_das_sql_log_configs

This data source provides the DAS Sql Log Config of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.284.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_db_instances" "default" {
  status = "Running"
}

data "alicloud_das_sql_log_configs" "default" {
  instance_id = data.alicloud_db_instances.default.instances.0.id
}

output "das_sql_log_config_enable" {
  value = data.alicloud_das_sql_log_configs.default.configs.0.enable
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required) The ID of the database instance to query the SQL log configuration for.
* `ids` - (Optional, Computed) A list of Sql Log Config IDs. Its element value is same as the instance ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `configs` - A list of Sql Log Config Entries. Each element contains the following attributes:
    * `id` - The ID of the Sql Log Config. Its value is same as `instance_id`.
    * `instance_id` - The ID of the database instance.
    * `enable` - Specifies whether SQL Explorer is enabled.
    * `request_enable` - The requested state of SQL Explorer.
    * `retention` - The retention period of SQL audit logs. Unit: days.
    * `hot_retention` - The retention period of hot SQL audit logs. Unit: days.
    * `cold_retention` - The retention period of cold SQL audit logs. Calculated as `retention - hot_retention`.
    * `version` - The current version of SQL audit logs.
    * `log_filter` - The configuration of log filters.
    * `sql_log_visible_time` - The visible start time of SQL audit logs.
