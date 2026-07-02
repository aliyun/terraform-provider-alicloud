---
subcategory: "DAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_das_sql_log_config"
description: |-
  Provides a Alicloud DAS Sql Log Config resource.
---

# alicloud_das_sql_log_config

Provides a DAS Sql Log Config resource.

SQL audit log configuration for database instances.

For information about DAS Sql Log Config and how to use it, see [What is Sql Log Config](https://next.api.alibabacloud.com/document/DAS/2020-01-16/DescribeSqlLogConfig).

-> **NOTE:** Available since v1.284.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_db_instances" "default" {
  status = "Running"
}

resource "alicloud_das_sql_log_config" "default" {
  instance_id    = data.alicloud_db_instances.default.instances.0.id
  enable         = true
  request_enable = true
  retention      = 30
  hot_retention  = 7
}
```

### Deleting `alicloud_das_sql_log_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_das_sql_log_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the database instance.
* `enable` - (Optional, Computed) Specifies whether SQL Explorer is enabled.
* `request_enable` - (Optional, Computed) The requested state of SQL Explorer.
* `retention` - (Optional, Computed, Int) The retention period of SQL audit logs. Unit: days.
* `hot_retention` - (Optional, Computed, Int) The retention period of hot SQL audit logs. Unit: days.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource, formatted as `<instance_id>`.
* `cold_retention` - The retention period of cold SQL audit logs. Calculated as `retention - hot_retention`.
* `log_filter` - The configuration of log filters.
* `sql_log_visible_time` - The visible start time of SQL audit logs.
* `version` - The current version of SQL audit logs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when creating the Sql Log Config.
* `update` - (Defaults to 5 mins) Used when updating the Sql Log Config.
* `delete` - (Defaults to 5 mins) Used when deleting the Sql Log Config.

## Import

DAS Sql Log Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_das_sql_log_config.example <instance_id>
```