---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_instance_refresh"
sidebar_current: "docs-alicloud-resource-ess-instance-refresh"
description: |-
  Provides a ESS instance refresh resource.
---

# alicloud_ess_instance_refresh

Provides a ESS instance refresh resource.

For information about ess instance refresh, see [StartInstanceRefresh](https://www.alibabacloud.com/help/en/auto-scaling/developer-reference/api-startinstancerefresh).

-> **NOTE:** Available since v1.261.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ess_instance_refresh&exampleId=36c0c7b1-51ca-80c7-4692-88ba02e441a3e4680173&activeTab=example&spm=docs.r.ess_instance_refresh.0.36c0c7b151&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  name = "${var.name}-${random_integer.default.result}"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.name
  cidr_block = "172.16.0.0/16"
}
data "alicloud_instance_types" "default1" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = local.name
}

resource "alicloud_security_group" "default" {
  security_group_name = local.name
  vpc_id              = alicloud_vpc.default.id
}
data "alicloud_images" "default1" {
  name_regex  = "^ubu"
  most_recent = true
  owners      = "system"
}
data "alicloud_images" "default2" {
  name_regex  = "^aliyun"
  most_recent = true
  owners      = "system"
}
resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 10
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
  desired_capacity   = 1
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default1.images[0].id
  instance_type     = data.alicloud_instance_types.default1.instance_types.0.id
  security_group_id = alicloud_security_group.default.id
  force_delete      = true
  active            = true
  enable            = true
}

resource "alicloud_ess_instance_refresh" "default" {
  scaling_group_id               = alicloud_ess_scaling_configuration.default.scaling_group_id
  desired_configuration_image_id = data.alicloud_images.default2.images.0.id
  min_healthy_percentage         = 90
  max_healthy_percentage         = 150
  checkpoint_pause_time          = 60
  skip_matching                  = false
  checkpoints {
    percentage = 100
  }
}
```
### Deleting `alicloud_ess_instance_refresh` or removing it from your configuration

The `alicloud_ess_instance_refresh` resource allows you to manage  `status = "RollbackInProgress"`  instance refresh, but Terraform cannot destroy it.
Deleting will remove it from your state file and management, but will not destroy the Instance Refresh.
You can resume managing the instance refresh via the AlibabaCloud Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ess_instance_refresh&spm=docs.r.ess_instance_refresh.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the scaling group.
* `min_healthy_percentage` - (Optional, ForceNew) The percentage of instances that must be healthy in the scaling group during the instance refresh. The value is a percentage of the scaling group's capacity.
* `max_healthy_percentage` - (Optional, ForceNew) The percentage by which the number of instances in the scaling group can exceed the group's capacity during the instance refresh.
* `desired_configuration_image_id` - (Optional, ForceNew) The ID of the image file. This is the image resource used for automatic instance creation.
* `desired_configuration_launch_template_id` - (Optional, ForceNew) The ID of the launch template. The scaling group uses this template to obtain launch configuration information.
* `desired_configuration_launch_template_version` - (Optional, ForceNew) The version of the launch template.
* `desired_configuration_launch_template_overrides` - (Optional, ForceNew) The instance type information in the launch template overrides. See [`desired_configuration_launch_template_overrides`](#desired_configuration_launch_template_overrides) below for details.
* `desired_configuration_containers` - (Optional, ForceNew) The list of containers in the instance. See [`desired_configuration_containers`](#desired_configuration_containers) below for details.
* `skip_matching` - (Optional, ForceNew) Indicates whether to skip instances that match the desired configuration.
* `status` - (Optional) The current status of the instance refresh task. Possible values:
    - Pending: The instance refresh task is created and waiting to be scheduled.
    - InProgress: The instance refresh task is in progress. 
    - Paused: The instance refresh task is paused. 
    - CheckpointPause: The task is paused because it has reached a checkpoint (Checkpoint.Percentage).
    - Failed: The instance refresh task failed.
    - Successful: The instance refresh task was successful.
    - Cancelling: The instance refresh task is being canceled. 
    - RollbackInProgress: The instance refresh task is being rolled back. 
    - RollbackSuccessful: The instance refresh task was rolled back successfully. Set RollbackSuccessful to rollback the instance refresh task.
    - RollbackFailed: The rollback of the instance refresh task failed.
    - Cancelled:  The instance refresh task is canceled. Set Cancelled to cancel the instance refresh task.
* `checkpoints` - (Optional, ForceNew) The checkpoints for the refresh task. The task automatically pauses for the duration specified by CheckpointPauseTime when the percentage of new instances reaches a specified value. See [`checkpoints`](#checkpoints) below for details.
* `checkpoint_pause_time` - (Optional, ForceNew) The duration of the pause when the task reaches a checkpoint. Unit: minutes.

### `desired_configuration_launch_template_overrides`

The desired_configuration_launch_template_overrides supports the following:

* `instance_type` - (Optional) The specified instance type, which overwrites the instance type in the launch template.

### `desired_configuration_containers`

The desired_configuration_containers supports the following:

* `name` - (Optional, ForceNew) The custom name of the container.
* `image` - (Optional, ForceNew) The container image.
* `commands` - (Optional, ForceNew) The container startup command.
* `args` - (Optional, ForceNew) The arguments for the container startup command.
* `environment_vars` - (Optional, ForceNew) Information about the environment variables. See [`environment_vars`](#desired_configuration_containers-environment_vars) below for details.

### `checkpoints`

The checkpoints supports the following:

* `percentage` - (Optional) The percentage of new instances out of the total instances in the scaling group. The task automatically pauses when this percentage is reached.

### `desired_configuration_containers-environment_vars`

The environment_vars supports the following:

* `key` - (Optional) The name of the environment variable.
* `value` - (Optional) The value of the environment variable.
* `field_ref_field_path` - (Optional) This parameter is not available for use.

## Attributes Reference

The following attributes are exported:

* `id` - The instance refresh ID.

## Import

ESS instance refresh can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_instance_refresh.example ir-abc123456
```

