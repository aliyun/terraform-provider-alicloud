---
subcategory: "ApsaraDB for MyBase (CDDC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_dedicated_propre_host"
description: |-
  Provides a Alicloud CDDC Dedicated Propre Host resource.
---

# alicloud_cddc_dedicated_propre_host

Provides a CDDC Dedicated Propre Host resource. MyBase proprietary cluster host resources, you need to add a whitelist to purchase a proprietary version of the cluster.

For information about CDDC Dedicated Propre Host and how to use it, see [What is Dedicated Propre Host](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/api-cddc-2020-03-20-creatededicatedhostgroup).

-> **NOTE:** Available since v1.210.0.

-> **DEPRECATED:**  This resource has been [deprecated](https://www.alibabacloud.com/help/en/apsaradb-for-mybase/latest/notice-stop-selling-mybase-hosted-instances-from-august-31-2023) from version `1.225.1`. 

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6e"
  network_type         = "Vpc"
}

data "alicloud_images" "default" {
  name_regex = "^aliyun_3_x64_20G_scc*"
  owners     = "system"
}

data "alicloud_instance_types" "essd" {
  cpu_core_count       = 2
  memory_size          = 4
  system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}

data "alicloud_security_groups" "default" {
  name_regex = "tf-exampleacc-cddc-dedicated_propre_host"
}

resource "alicloud_security_group" "default" {
  count  = length(data.alicloud_security_groups.default.ids) > 0 ? 0 : 1
  vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
  name   = "tf-exampleacc-cddc-dedicated_propre_host"
}

data "alicloud_ecs_deployment_sets" "default" {
  name_regex = "tf-exampleacc-cddc-dedicated_propre_host"
}

resource "alicloud_ecs_deployment_set" "default" {
  count               = length(data.alicloud_ecs_deployment_sets.default.ids) > 0 ? 0 : 1
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "tf-exampleacc-cddc-dedicated_propre_host"
  description         = "tf-exampleacc-cddc-dedicated_propre_host"
}

data "alicloud_key_pairs" "default" {
  name_regex = "tf-exampleacc-cddc-dedicated_propre_host"
}

resource "alicloud_key_pair" "default" {
  count         = length(data.alicloud_key_pairs.default.ids) > 0 ? 0 : 1
  key_pair_name = "tf-exampleacc-cddc-dedicated_propre_host"
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  engine     = "MySQL"
  name_regex = "tf-exampleacc-cddc-dedicated_propre_host"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  count                     = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
  engine                    = "MySQL"
  vpc_id                    = data.alicloud_vpcs.default.ids.0
  cpu_allocation_ratio      = 101
  mem_allocation_ratio      = 50
  disk_allocation_ratio     = 200
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = "tf-exampleacc-cddc-dedicated_propre_host"
  open_permission           = true
}

locals {
  alicloud_security_group_id     = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.ids.0 : concat(alicloud_security_group.default[*].id, [""])[0]
  alicloud_ecs_deployment_set_id = length(data.alicloud_ecs_deployment_sets.default.ids) > 0 ? data.alicloud_ecs_deployment_sets.default.sets.0.deployment_set_id : concat(alicloud_ecs_deployment_set.default[*].id, [""])[0]
  alicloud_key_pair_id           = length(data.alicloud_key_pairs.default.ids) > 0 ? data.alicloud_key_pairs.default.ids.0 : concat(alicloud_key_pair.default[*].id, [""])[0]
  dedicated_host_group_id        = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : concat(alicloud_cddc_dedicated_host_group.default[*].id, [""])[0]
}

resource "alicloud_cddc_dedicated_propre_host" "default" {
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  ecs_instance_name       = "exampleTf"
  ecs_deployment_set_id   = local.alicloud_ecs_deployment_set_id
  auto_renew              = "false"
  security_group_id       = local.alicloud_security_group_id
  dedicated_host_group_id = local.dedicated_host_group_id
  ecs_host_name           = "exampleTf"
  vpc_id                  = data.alicloud_vpcs.default.ids.0
  ecs_unique_suffix       = "false"
  password_inherit        = "false"
  engine                  = "mysql"
  period                  = "1"
  os_password             = "YourPassword123!"
  ecs_zone_id             = "cn-hangzhou-i"
  ecs_class_list {
    disk_type                     = "cloud_essd"
    sys_disk_type                 = "cloud_essd"
    disk_count                    = "1"
    system_disk_performance_level = "PL1"
    data_disk_performance_level   = "PL1"
    disk_capacity                 = "40"
    instance_type                 = "ecs.c6a.large"
    sys_disk_capacity             = "40"
  }

  payment_type = "Subscription"
  image_id     = "m-bp1d13fxs1ymbvw1dk5g"
  period_type  = "Monthly"
}
```

### Deleting `alicloud_cddc_dedicated_propre_host` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cddc_dedicated_propre_host`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cddc_dedicated_propre_host&spm=docs.r.cddc_dedicated_propre_host.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional, Available since v1.215.0) Whether to pay automatically when the host is created.
* `auto_renew` - (Optional) Whether to enable automatic renewal. Valid values:
  - **true**: On
  - **false** (default): Off
* `dedicated_host_group_id` - (Optional, ForceNew, Computed) You have a dedicated cluster ID.
* `ecs_class_list` - (Required, ForceNew) ECS specifications. See [`ecs_class_list`](#ecs_class_list) below.
* `ecs_deployment_set_id` - (Optional, ForceNew) The ID of the cloud server deployment set.
* `ecs_host_name` - (Optional, ForceNew) Windows system: length of 2 to 15 characters, allowing the use of upper and lower case letters, numbers. You cannot use only numbers. Other operating systems (such as Linux): the length of 2 to 64 characters, allowing the use of dot (.) to separate characters into multiple segments, each segment allows the use of upper and lower case letters, numbers, but can not use continuous dot (.). Cannot start or end with a dot (.).
* `ecs_instance_name` - (Optional, ForceNew) The instance name. It must be 2 to 128 characters in length and must start with an uppercase or lowercase letter or a Chinese character. It cannot start with http:// or https. Can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-). The default value is the InstanceId of the instance.
* `ecs_unique_suffix` - (Optional) Whether to automatically add an ordered suffix for HostName and InstanceName when creating multiple instances. The ordered suffix starts from 001 and cannot exceed 999. Value Description:
  - **true**: added.
  - **false** (default): Do not add.
When the HostName or InstanceName is set according to the specified sorting format, and the naming suffix name_suffix is not set, that is, when the naming format is name_prefix[begin_number,bits], the UniqueSuffix does not take effect, and the names are only sorted according to the specified order.
* `ecs_zone_id` - (Required, ForceNew) The ID of the zone.
* `engine` - (Required, ForceNew) Database type, value:
  - **alisql**
  - **tair**
  - **mssql**
Must be consistent with the parent resource cluster engine attributes.
* `image_id` - (Optional, ForceNew) The ID of the custom image.
-> **NOTE:**  If you need to use the default image, you do not need to fill it in.
* `internet_charge_type` - (Optional, Available since v1.215.0) Network billing type. Value range: PayByBandwidth: Billing based on fixed bandwidth. PayByTraffic: charges by using the flow meter.
* `internet_max_bandwidth_out` - (Optional, Available since v1.215.0) The maximum outbound bandwidth of the public network, in Mbit/s. Value range: 0~100.  Default value: 0. When set to greater than 0, a public IP is automatically created.
* `key_pair_name` - (Optional, ForceNew) The key pair name.
* `os_password` - (Optional) Host login password, which can be set later. The password must meet the following requirements:
  - Length is 8~30 characters.
  - Must contain at least three items: uppercase letters, lowercase letters, numbers, and special characters.
  - Special symbol '()\' ~! @#$%^& *-_+ =|{}[]:;',.? /'
-> **NOTE:** - If you need to set the host login password later, fill in an empty string for this parameter. If you need to set a host login password, we recommend that you use the HTTPS protocol to send requests to avoid password leakage.
* `password_inherit` - (Optional) Whether to use the default password of the image.
  - **false**: (default)Do not use
  - **true**: Use
-> **NOTE:**  If the default password of the image is used, the **OSPassword** parameter is not required.
* `payment_type` - (Required, ForceNew) The Payment type. Currently, only **Subscription** is supported.
* `period` - (Optional) Duration of purchase.
* `period_type` - (Optional) The subscription type. Currently, only **Monthly** (subscription) is supported.
* `resource_group_id` - (Optional, ForceNew, Computed, Available since v1.215.0) The ID of the resource group.
* `security_group_id` - (Required, ForceNew) The ID of the security group.
* `tags` - (Optional, ForceNew, Map, Available since v1.215.0) Host tag information.
* `user_data` - (Optional, Available since v1.215.0) User-defined script data. The maximum size of the original data is 16kB.
* `user_data_encoded` - (Optional, Available since v1.215.0) Whether custom data is encoded in Base64 format.
* `vswitch_id` - (Required, ForceNew) The ID of the virtual switch.
* `vpc_id` - (Required, ForceNew) VPCID of the VPC.

### `ecs_class_list`

The ecs_class_list supports the following:
* `data_disk_performance_level` - (Optional, ForceNew) Data disk PL level.
* `disk_capacity` - (Optional, ForceNew) The capacity of the data disk.
* `disk_count` - (Optional, ForceNew) Number of mounted data disks.
* `disk_type` - (Optional, ForceNew) Data disk type, value range:
  - **cloud_essd**: the ESSD cloud disk.
  - **cloud_ssd**: SSD cloud disk.
  - **cloud_efficiency**: The ultra cloud disk.
  - **cloud_auto**: ESSD AutoPL cloud disk.
* `instance_type` - (Required, ForceNew) ECS specifications.
* `sys_disk_capacity` - (Required, ForceNew) System disk capacity.
* `sys_disk_type` - (Required, ForceNew) System disk type, value:
  - **cloud_essd**: the ESSD cloud disk. 
  - **cloud_ssd**: SSD cloud disk.
  - **cloud_efficiency**: The ultra cloud disk.
  - **cloud_auto**: ESSD AutoPL cloud disk.
* `system_disk_performance_level` - (Optional, ForceNew) System disk PL level.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<dedicated_host_group_id>:<ecs_instance_id>`.
* `ecs_instance_id` - ECS instance ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dedicated Propre Host.

## Import

CDDC Dedicated Propre Host can be imported using the id, e.g.

```shell
$ terraform import alicloud_cddc_dedicated_propre_host.example <dedicated_host_group_id>:<ecs_instance_id>
```