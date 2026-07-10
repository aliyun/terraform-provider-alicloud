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

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_das_sql_log_config&exampleId=deab1288-8dbb-d089-77c0-43e05a731eae1bdd4f5a&activeTab=example&spm=docs.r.das_sql_log_config.0.deab12888d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_das_sql_log_config&spm=docs.r.das_sql_log_config.example&intl_lang=EN_US)


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