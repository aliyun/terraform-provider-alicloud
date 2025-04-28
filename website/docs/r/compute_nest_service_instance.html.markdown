---
subcategory: "Compute Nest"
layout: "alicloud"
page_title: "Alicloud: alicloud_compute_nest_service_instance"
sidebar_current: "docs-alicloud-resource-compute-nest-service-instance"
description: |-
  Provides a Alicloud Compute Nest Service Instance resource.
---

# alicloud_compute_nest_service_instance

Provides a Compute Nest Service Instance resource.

For information about Compute Nest Service Instance and how to use it, see [What is Service Instance](https://www.alibabacloud.com/help/zh/compute-nest/developer-reference/api-computenest-2021-06-01-createserviceinstance).

-> **NOTE:** Available since v1.205.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_compute_nest_service_instance&exampleId=e3737a4e-d172-5b12-fb9e-94f87d27c85d91285e7f&activeTab=example&spm=docs.r.compute_nest_service_instance.0.e3737a4ed1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}

provider "alicloud" {
  region = "cn-hangzhou"
}
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
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
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
  vswitch_id                 = alicloud_vswitch.default.id
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
    resources            = <<EOF
    {
      "Type": "ResourceIds",
      "RegionId": "cn-hangzhou",
      "ResourceIds": {
      "ALIYUN::ECS::INSTANCE": [
        "${alicloud_instance.default.id}"
        ]
      } 
    }
    EOF
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
* `service_instance_name` - (Optional, ForceNew) The name of the Service Instance.
* `parameters` - (Optional) The parameters entered by the deployment service instance.
* `enable_instance_ops` - (Optional, ForceNew) Whether the service instance has the O&M function. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `template_name` - (Optional, ForceNew) The name of the template.
* `resource_group_id` - (Optional) The ID of the resource group.
* `specification_name` - (Optional, ForceNew) The name of the specification.
* `payment_type` - (Optional, ForceNew) The type of payment. Valid values: `Permanent`, `Subscription`, `PayAsYouGo`, `CustomFixTime`.
* `enable_user_prometheus` - (Optional, ForceNew) Whether Prometheus monitoring is enabled. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `operation_metadata` - (Optional, ForceNew) The configuration of O&M. See [`operation_metadata`](#operation_metadata) below.
* `commodity` - (Optional) The order information of cloud market. See [`commodity`](#commodity) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `operation_metadata`

The operation_metadata supports the following:

* `operation_start_time` - (Optional, ForceNew) The start time of O&M.
* `operation_end_time` - (Optional, ForceNew) The end time of O&M.
* `resources` - (Optional, ForceNew) The list of imported resources.
* `operated_service_instance_id` - (Optional, ForceNew) The ID of the imported service instance.

### `commodity`

The commodity supports the following:

* `pay_period` - (Optional) Length of purchase.
* `pay_period_unit` - (Optional) Duration unit. Valid values: `Year`, `Month`, `Day`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Service Instance.
* `status` - The status of the Service Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Service Instance.
* `update` - (Defaults to 5 mins) Used when update the Service Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Service Instance.

## Import

Compute Nest Service Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_compute_nest_service_instance.example <id>
```
