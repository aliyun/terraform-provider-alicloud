---
subcategory: "ApsaraDB for MyBase (CDDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_host"
sidebar_current: "docs-alicloud-resource-cddc-dedicated-host"
description: |-
  Provides a Alicloud ApsaraDB for MyBase Dedicated Host resource.
---

# alicloud_cddc_dedicated_host

Provides a ApsaraDB for MyBase Dedicated Host resource.

For information about ApsaraDB for MyBase Dedicated Host and how to use it, see [What is Dedicated Host](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/creatededicatedhost).

-> **NOTE:** Available since v1.147.0.

-> **DEPRECATED:**  This resource has been [deprecated](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/notice-stop-selling-mybase-hosted-instances-from-august-31-2023) from version `1.225.1`. 

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_cddc_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  engine                    = "MySQL"
  vpc_id                    = alicloud_vpc.default.id
  cpu_allocation_ratio      = 101
  mem_allocation_ratio      = 50
  disk_allocation_ratio     = 200
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = var.name
}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = alicloud_cddc_dedicated_host_group.default.id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = alicloud_vswitch.default.id
  payment_type            = "Subscription"
  tags = {
    Created = "TF"
    For     = "CDDC_DEDICATED"
  }
}
```

### Deleting `alicloud_cddc_dedicated_host` or removing it from your configuration

The `alicloud_cddc_dedicated_host` resource allows you to manage `payment_type = "Subscription"` host instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Host Instance.
You can resume managing the subscription host instance via the AlibabaCloud Console.

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Dedicated Host.
* `update` - (Defaults to 20 mins) Used when update the Dedicated Host.
* `delete` - (Defaults to 20 mins) Used when delete the Dedicated Host.

## Import

ApsaraDB for MyBase Dedicated Host can be imported using the id, e.g.

```shell
$ terraform import alicloud_cddc_dedicated_host.example <dedicated_host_group_id>:<dedicated_host_id>
```