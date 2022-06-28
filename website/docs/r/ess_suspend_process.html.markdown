---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_suspend_process"
sidebar_current: "docs-alicloud-resource-ess-suspend-process"
description: |-
  Provides a ESS Suspend Process resource to suspend or resume process for scaling group.
---

# alicloud\_ess\_suspend\_process

Suspend/Resume processes to a specified scaling group.

For information about scaling group suspend process, see [SuspendProcesses](https://www.alibabacloud.com/help/en/auto-scaling/latest/suspendprocesses).

-> NOTE: Available in v1.166.0+

## Example Usage

```
variable "name" {
  default = "testAccEssSuspendProcess"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  name              = var.name
}


resource "alicloud_alb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  name             = "test"
}


resource "alicloud_ess_scaling_group" "default" {
  min_size           = "2"
  max_size           = "2"
  scaling_group_name = var.name
  vswitch_ids        = [alicloud_vswitch.default.id]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  image_id = data.alicloud_images.default.images[0].id
  instance_type = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete = true
  active = true
  enable = true
}

resource "alicloud_alb_server_group" "default" {
  server_group_name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  health_check_config {
  health_check_enabled = "false"
  }
  sticky_session_config {
  sticky_session_enabled = true
  cookie                 = "tf-testAcc"
  sticky_session_type    = "Server"
  }
}

resource "alicloud_ess_suspend_process" "default" {
  scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
  process             = "ScaleIn"
  depends_on = ["alicloud_ess_scaling_configuration.default"]
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) ID of the scaling group.
* `process` - (Required, ForceNew) Activity type N that you want to suspend. Valid values are: `SCALE_OUT`,`SCALE_IN`,`HealthCheck`,`AlarmNotification` and `ScheduledAction`.




## Attributes Reference

The following attributes are exported:

* `id` - (Required, ForceNew) The ESS suspend process resource ID，in the follwing format: scaling_group_id:process.

## Import

ESS suspend process can be imported using the id, e.g.

```
$ terraform import alicloud_suspend_process.example asg-xxx:sgp-xxx:5000 
```

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the process.
* `delete` - (Defaults to 1 mins) Used when delete the process.