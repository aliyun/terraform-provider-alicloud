---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host"
sidebar_current: "docs-alicloud-resource-cddc-dedicated-host"
description: |-
  Provides a Alicloud ApsaraDB for MyBase Dedicated Host resource.
---

# alicloud\_cddc\_dedicated\_host

Provides a ApsaraDB for MyBase Dedicated Host resource.

For information about ApsaraDB for MyBase Dedicated Host and how to use it, see [What is Dedicated Host](https://www.alibabacloud.com/help/doc-detail/210864.html).

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  engine                    = "MySQL"
  vpc_id                    = data.alicloud_vpcs.default.ids.0
  cpu_allocation_ratio      = 101
  mem_allocation_ratio      = 50
  disk_allocation_ratio     = 200
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = "example_value"
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = "example_value"
  dedicated_host_group_id = alicloud_cddc_dedicated_host_group.default.id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  payment_type            = "Subscription"
  tags = {
    Created = "TF"
    For     = "CDDC_DEDICATED"
  }
}

```

## Argument Reference

The following arguments are supported:

* `allocation_status` - (Optional) Specifies whether instances can be created on the host. Valid values: `Allocatable` or `Suspended`. `Allocatable`: Instances can be created on the host. `Suspended`: Instances cannot be created on the host.
* `auto_renew` - (Optional) Specifies whether to enable the auto-renewal feature.
* `dedicated_host_group_id` - (Required, ForceNew) The ID of the dedicated cluster.
* `host_class` - (Required) The instance type of the host. For more information about the supported instance types of hosts, see [Host specification details](https://www.alibabacloud.com/help/doc-detail/206343.htm).
* `host_name` - (Optional) The name of the host. The name must be `1` to `64` characters in length and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter.
* `image_category` - (Optional) Host Image Category. Valid values: `WindowsWithMssqlEntAlwaysonLicense`, `WindowsWithMssqlStdLicense`, `WindowsWithMssqlEntLicense`, `WindowsWithMssqlWebLicense`, `AliLinux`.
* `os_password` - (Optional) Host password. **NOTE:** The creation of a host password is supported only when the database type is `Tair-PMem`.
* `payment_type` - (Required) The payment type of the resource. Valid values: `Subscription`.
* `period` - (Optional) The unit of the subscription duration. Valid values: `Year`, `Month`, `Week`.
* `used_time` - (Optional) The subscription duration of the host. Valid values: 
  * If the Period parameter is set to `Year`, the value of the UsedTime parameter ranges from `1` to `5`. 
  * If the Period parameter is set to `Month`, the value of the UsedTime parameter ranges from `1` to `9`.
  * If the Period parameter is set to `Week`, the value of the UsedTime parameter ranges from `1`, `2` and `3`.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch to which the host is connected.
* `zone_id` - (Required, ForceNew) The ID of the zone.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dedicated Host. The value formats as `<dedicated_host_group_id>:<dedicated_host_id>`.
* `dedicated_host_id` - The ID of the host.
* `status` - The state of the host. Valid values: `0:` The host is being created. `1`: The host is running. `2`: The host is faulty. `3`: The host is ready for deactivation. `4`: The host is being maintained. `5`: The host is deactivated. `6`: The host is restarting. `7`: The host is locked.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Dedicated Host.
* `update` - (Defaults to 20 mins) Used when update the Dedicated Host.
* `delete` - (Defaults to 20 mins) Used when delete the Dedicated Host.

## Import

ApsaraDB for MyBase Dedicated Host can be imported using the id, e.g.

```
$ terraform import alicloud_cddc_dedicated_host.example <dedicated_host_group_id>:<dedicated_host_id>
```