---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_supabase_project"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) Supabase Project resource.
---

# alicloud_gpdb_supabase_project

Provides a AnalyticDB for PostgreSQL (GPDB) Supabase Project resource.



For information about AnalyticDB for PostgreSQL (GPDB) Supabase Project and how to use it, see [What is Supabase Project](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateSupabaseProject).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-j"
}

resource "alicloud_gpdb_supabase_project" "default" {
  project_spec           = "1C2G"
  zone_id                = "cn-hangzhou-j"
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  project_name           = "supabase_example"
  security_ip_list       = ["127.0.0.1"]
  vswitch_id             = data.alicloud_vswitches.default.ids.0
  disk_performance_level = "PL0"
  storage_size           = "1"
  account_password       = "YourPassword123!"
}
```

## Argument Reference

The following arguments are supported:
* `account_password` - (Required) The password for the initial account.
  - Consists of three or more of uppercase letters, lowercase letters, numbers, and special characters.
  - Support for special characters:! @#$%^& *()_+-=
  - Length is 8~32 characters.
* `disk_performance_level` - (Optional, ForceNew) cloud disk performance level
* `project_name` - (Required, ForceNew) The project name. The naming rules are as follows:
  - 1~128 characters in length.
  - Can only contain English letters, numbers, dashes (-) and underscores (_).
  - Must begin with an English letter or an underscore (_).
* `project_spec` - (Required, ForceNew) The performance level of the Supabase instance.
* `security_ip_list` - (Required, List) The IP address whitelist.
* `storage_size` - (Optional, ForceNew, Int) The storage capacity of the instance. Unit: GB.
* `vswitch_id` - (Required, ForceNew) The vSwitch ID.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `zone_id` - (Required, ForceNew) The Zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID.
* `status` - The status of the Supabase instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 15 mins) Used when create the Supabase Project.
* `delete` - (Defaults to 5 mins) Used when delete the Supabase Project.
* `update` - (Defaults to 5 mins) Used when update the Supabase Project.

## Import

AnalyticDB for PostgreSQL (GPDB) Supabase Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_supabase_project.example <id>
```
