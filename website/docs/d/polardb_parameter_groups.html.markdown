---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_parameter_groups"
sidebar_current: "docs-alicloud-datasource-polardb-parameter-groups"
description: |-
  Provides a list of PolarDB Parameter Groups to the user.
---

# alicloud\_polardb\_parameter\_groups

This data source provides the PolarDB Parameter Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.183.0+.

## Example Usage

Basic Usage

```terraform

data "alicloud_polardb_parameter_groups" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
}

data "alicloud_polardb_parameter_groups" "ids" {
  ids = [data.alicloud_polardb_parameter_groups.default.groups.0.id]
}

output "polardb_parameter_group_id_1" {
  value = data.alicloud_polardb_parameter_groups.ids.groups.0.id
}

data "alicloud_polardb_parameter_groups" "nameRegex" {
  name_regex = data.alicloud_polardb_parameter_groups.default.groups.0.parameter_group_name
}
output "polardb_parameter_group_id_2" {
  value = data.alicloud_polardb_parameter_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Parameter Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Parameter Group name.
* `db_type` - (Optional, ForceNew) The type of the database engine. Only `MySQL` is supported.
* `db_version` - (Optional, ForceNew) The version number of the database engine. Valid values: `5.6`, `5.7`, `8.0`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Parameter Group names.
* `groups` - A list of PolarDB Parameter Groups. Each element contains the following attributes:
	* `id` - The ID of the Parameter Group.
	* `parameter_group_id` - The ID of the Parameter Group.
	* `parameter_group_name` - The name of the parameter template.
	* `db_type` - The type of the database engine.
	* `db_version` - The version number of the database engine.
	* `parameter_group_desc` - The description of the parameter template.
	* `parameter_group_type` - The type of the parameter template.  
	* `parameter_counts` - The number of parameters in the parameter template.
	* `force_restart` - Indicates whether to restart the cluster when this parameter template is applied.
	* `create_time` - The time when the parameter template was created. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.