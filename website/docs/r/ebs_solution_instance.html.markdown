---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_solution_instance"
description: |-
  Provides a Alicloud EBS Solution Instance resource.
---

# alicloud_ebs_solution_instance

Provides a EBS Solution Instance resource. 

For information about EBS Solution Instance and how to use it, see [What is Solution Instance](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

variable "zone_id" {
  default = "cn-shanghai-l"
}

variable "region_id" {
  default = "cn-shanghai"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ebs_solution_instance" "default" {
  solution_instance_name = var.name
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.ids.0
  description            = "description"
  solution_id            = "mysql"
  parameters {
    parameter_key   = "zoneId"
    parameter_value = var.zone_id
  }
  parameters {
    parameter_key   = "ecsType"
    parameter_value = "ecs.c6.large"
  }
  parameters {
    parameter_key   = "ecsImageId"
    parameter_value = "CentOS_7"
  }
  parameters {
    parameter_key   = "internetMaxBandwidthOut"
    parameter_value = "100"
  }
  parameters {
    parameter_key   = "internetChargeType"
    parameter_value = "PayByTraffic"
  }
  parameters {
    parameter_key   = "ecsPassword"
    parameter_value = "Ebs12345"
  }
  parameters {
    parameter_key   = "sysDiskType"
    parameter_value = "cloud_essd"
  }
  parameters {
    parameter_key   = "sysDiskPerformance"
    parameter_value = "PL0"
  }
  parameters {
    parameter_key   = "sysDiskSize"
    parameter_value = "40"
  }
  parameters {
    parameter_key   = "dataDiskType"
    parameter_value = "cloud_essd"
  }
  parameters {
    parameter_key   = "dataDiskPerformance"
    parameter_value = "PL0"
  }
  parameters {
    parameter_key   = "dataDiskSize"
    parameter_value = "40"
  }
  parameters {
    parameter_key   = "mysqlVersion"
    parameter_value = "MySQL80"
  }
  parameters {
    parameter_key   = "mysqlUser"
    parameter_value = "root"
  }
  parameters {
    parameter_key   = "mysqlPassword"
    parameter_value = "Ebs12345"
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Solution Instance Description.
* `parameters` - (Optional) Solution Instance Creation Parameters. See [`parameters`](#parameters) below.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `solution_id` - (Required, ForceNew) Solution ID.
* `solution_instance_name` - (Optional, Computed) Solution Instance Name.

### `parameters`

The parameters supports the following:
* `parameter_key` - (Required) Create parameter Key.
* `parameter_value` - (Required) Create parameter Value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Solution Instance Creation Time.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Solution Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Solution Instance.
* `update` - (Defaults to 5 mins) Used when update the Solution Instance.

## Import

EBS Solution Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_solution_instance.example <id>
```