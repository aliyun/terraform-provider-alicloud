---
subcategory: "Compute Nest"
layout: "alicloud"
page_title: "Alicloud: alicloud_compute_nest_service_instance"
sidebar_current: "docs-alicloud-resource-compute-nest-service-instance"
description: |-
  Provides a Alicloud Compute Nest Service Instance resource.
---

# alicloud\_compute\_nest\_service\_instance

Provides a Compute Nest Service Instance resource.

For information about Compute Nest Service Instance and how to use it, see [What is Service Instance](https://help.aliyun.com/document_detail/396194.html).

-> **NOTE:** Available in v1.205.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
  name_regex = "your_name_regex"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_compute_nest_service_instance" "default" {
  service_id            = "service-dd475e6e468348799f0f"
  service_version       = "1"
  service_instance_name = var.name
  resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  payment_type          = "Permanent"
  operation_metadata {
    operation_start_time = "1681281179000"
    operation_end_time   = "1681367579000"
    resources            = "{\"Type\":\"ResourceIds\",\"ResourceIds\":{\"ALIYUN::ECS::INSTANCE\":[\"${alicloud_instance.default.id}\"]},\"RegionId\":\"cn-hangzhou\"}"
  }
  tags = {
    Created = "TF"
    For     = "ServiceInstance"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, ForceNew) The ID of the service.
* `service_version` - (Required, ForceNew) The version of the service.
* `service_instance_name` - (Optional, ForceNew, Computed) The name of the Service Instance.
* `parameters` - (Optional) The parameters entered by the deployment service instance.
* `enable_instance_ops` - (Optional, ForceNew, Computed) Whether the service instance has the O&M function. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `template_name` - (Optional, ForceNew, Computed) The name of the template.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `specification_name` - (Optional, ForceNew) The name of the specification.
* `payment_type` - (Optional, ForceNew, Computed) The type of payment. Valid values: `Permanent`, `Subscription`, `PayAsYouGo`, `CustomFixTime`.
* `enable_user_prometheus` - (Optional, ForceNew, Computed) Whether Prometheus monitoring is enabled. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `operation_metadata` - (Optional, ForceNew, Computed) The configuration of O&M. See the following `Block operation_metadata`.
* `commodity` - (Optional) The order information of cloud market. See the following `Block commodity`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

#### Block operation_metadata

The operation_metadata supports the following:

* `operation_start_time` - (Optional, ForceNew) The start time of O&M.
* `operation_end_time` - (Optional, ForceNew) The end time of O&M.
* `resources` - (Optional, ForceNew, Computed) The list of imported resources.
* `operated_service_instance_id` - (Optional, ForceNew) The ID of the imported service instance.

#### Block commodity

The commodity supports the following:

* `pay_period` - (Optional) Length of purchase.
* `pay_period_unit` - (Optional) Duration unit. Valid values: `Year`, `Month`, `Day`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Service Instance.
* `status` - The status of the Service Instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Service Instance.
* `update` - (Defaults to 5 mins) Used when update the Service Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Service Instance.

## Import

Compute Nest Service Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_compute_nest_service_instance.example <id>
```
