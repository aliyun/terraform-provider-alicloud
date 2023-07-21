---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_endpoint"
sidebar_current: "docs-alicloud-resource-ga-basic-endpoint"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Endpoint resource.
---

# alicloud_ga_basic_endpoint

Provides a Global Accelerator (GA) Basic Endpoint resource.

For information about Global Accelerator (GA) Basic Endpoint and how to use it, see [What is Basic Endpoint](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createbasicendpoint).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-shenzhen"
}

variable "endpoint_region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region = var.region
  alias  = "sz"
}

provider "alicloud" {
  region = var.endpoint_region
  alias  = "hz"
}

data "alicloud_zones" "default" {
  provider                    = alicloud.sz
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  provider   = alicloud.sz
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  provider     = alicloud.sz
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  provider = alicloud.sz
  vpc_id   = alicloud_vpc.default.id
  name     = "terraform-example"
}

resource "alicloud_ecs_network_interface" "default" {
  provider           = alicloud.sz
  vswitch_id         = alicloud_vswitch.default.id
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  basic_accelerator_name = "terraform-example"
  description            = "terraform-example"
  bandwidth_billing_type = "CDT"
  auto_use_coupon        = "true"
  auto_pay               = true
}

resource "alicloud_ga_basic_endpoint_group" "default" {
  accelerator_id            = alicloud_ga_basic_accelerator.default.id
  endpoint_group_region     = var.region
  basic_endpoint_group_name = "terraform-example"
  description               = "terraform-example"
}

resource "alicloud_ga_basic_endpoint" "default" {
  provider                  = alicloud.hz
  accelerator_id            = alicloud_ga_basic_accelerator.default.id
  endpoint_group_id         = alicloud_ga_basic_endpoint_group.default.id
  endpoint_type             = "ENI"
  endpoint_address          = alicloud_ecs_network_interface.default.id
  endpoint_sub_address_type = "secondary"
  endpoint_sub_address      = "192.168.0.1"
  basic_endpoint_name       = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Basic GA instance.
* `endpoint_group_id` - (Required, ForceNew) The ID of the Basic Endpoint Group.
* `endpoint_type` - (Required, ForceNew) The type of the Basic Endpoint. Valid values: `ENI`, `SLB`, `ECS` and `NLB`.
* `endpoint_address` - (Required, ForceNew) The address of the Basic Endpoint.
* `endpoint_sub_address_type` - (Optional, ForceNew) The sub address type of the Basic Endpoint. Valid values: `primary`, `secondary`.
* `endpoint_sub_address` - (Optional, ForceNew) The sub address of the Basic Endpoint.
* `endpoint_zone_id` - (Optional, ForceNew) The zone id of the Basic Endpoint.
* `basic_endpoint_name` - (Optional) The name of the Basic Endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Endpoint. It formats as `<endpoint_group_id>:<endpoint_id>`.
* `endpoint_id` - The ID of the Basic Endpoint.
* `status` - The status of the Basic Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Endpoint.
* `update` - (Defaults to 3 mins) Used when update the Basic Endpoint.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Endpoint.

## Import

Global Accelerator (GA) Basic Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_endpoint.example <endpoint_group_id>:<endpoint_id>
```
