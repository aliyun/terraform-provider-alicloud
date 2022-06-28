---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host_account"
sidebar_current: "docs-alicloud-resource-cddc-dedicated-host-account"
description: |-
  Provides a Alicloud ApsaraDB for MyBase Dedicated Host Account resource.
---

# alicloud\_cddc\_dedicated\_host\_account

Provides a ApsaraDB for MyBase Dedicated Host Account resource.

For information about ApsaraDB for MyBase Dedicated Host Account and how to use it, see [What is Dedicated Host Account](https://www.alibabacloud.com/help/en/doc-detail/196877.html).

-> **NOTE:** Available in v1.148.0+.

-> **NOTE:** Each Dedicated host can have only one account. Before you create an account for a host, make sure that the existing account is deleted.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftestacc"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type        = "mssql"
  zone_id        = data.alicloud_cddc_zones.default.ids.0
  storage_type   = "cloud_essd"
  image_category = "WindowsWithMssqlStdLicense"

}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "default-NODELETING"
  engine     = "mssql"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  count                     = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
  engine                    = "SQLServer"
  vpc_id                    = data.alicloud_vpcs.default.ids.0
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = var.name
  open_permission           = true
}

data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.groups[0].vpc_id : data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_vswitch" "default" {
  count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = data.alicloud_vpcs.default.vpcs[0].cidr_block
  zone_id    = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : alicloud_cddc_dedicated_host_group.default[0].id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : alicloud_vswitch.default[0].id
  payment_type            = "Subscription"
  image_category          = "WindowsWithMssqlStdLicense"
}
resource "alicloud_cddc_dedicated_host_account" "example" {
  account_name      = var.name
  account_password  = "yourpassword"
  dedicated_host_id = alicloud_cddc_dedicated_host.default.dedicated_host_id
  account_type      = "Normal"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) The name of the Dedicated host account. The account name must be 2 to 16 characters in length, contain lower case letters, digits, and underscore(_). At the same time, the name must start with a letter and end with a letter or number.
* `account_password` - (Required, Sensitive) The password of the Dedicated host account. The account password must be 6 to 32 characters in length, and can contain letters, digits, and special characters `!@#$%^&*()_+-=`.
* `account_type` - (Optional, ForceNew) The type of the Dedicated host account. Valid values: `Admin`, `Normal`.
* `dedicated_host_id` - (Required, ForceNew) The ID of Dedicated the host.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Dedicated Host Account. The value formats as `<dedicated_host_id>:<account_name>`.

## Import

ApsaraDB for MyBase Dedicated Host Account can be imported using the id, e.g.

```
$ terraform import alicloud_cddc_dedicated_host_account.example <dedicated_host_id>:<account_name>
```