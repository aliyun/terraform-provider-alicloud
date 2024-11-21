---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerate_ip_endpoint_relation"
sidebar_current: "docs-alicloud-resource-ga-basic-accelerate-ip-endpoint-relation"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation resource.
---

# alicloud_ga_basic_accelerate_ip_endpoint_relation

Provides a Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation resource.

For information about Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation and how to use it, see [What is Basic Accelerate Ip Endpoint Relation](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createbasicaccelerateipendpointrelation).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_basic_accelerate_ip_endpoint_relation&exampleId=0cf2c574-ffe0-b728-603f-7a4558b608ec3d0b96a2&activeTab=example&spm=docs.r.ga_basic_accelerate_ip_endpoint_relation.0.0cf2c574ff&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_ga_basic_ip_set" "default" {
  accelerator_id       = alicloud_ga_basic_accelerator.default.id
  accelerate_region_id = var.endpoint_region
  isp_type             = "BGP"
  bandwidth            = "5"
}

resource "alicloud_ga_basic_accelerate_ip" "default" {
  accelerator_id = alicloud_ga_basic_accelerator.default.id
  ip_set_id      = alicloud_ga_basic_ip_set.default.id
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
  endpoint_sub_address_type = "primary"
  endpoint_sub_address      = "192.168.0.1"
  basic_endpoint_name       = "terraform-example"
}

resource "alicloud_ga_basic_accelerate_ip_endpoint_relation" "default" {
  accelerator_id   = alicloud_ga_basic_accelerate_ip.default.accelerator_id
  accelerate_ip_id = alicloud_ga_basic_accelerate_ip.default.id
  endpoint_id      = alicloud_ga_basic_endpoint.default.endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Basic GA instance.
* `accelerate_ip_id` - (Required, ForceNew) The ID of the Basic Accelerate IP.
* `endpoint_id` - (Required, ForceNew) The ID of the Basic Endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Accelerate Ip Endpoint Relation. It formats as `<accelerator_id>:<accelerate_ip_id>:<endpoint_id>`.
* `status` - The status of the Basic Accelerate Ip Endpoint Relation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Basic Accelerate Ip Endpoint Relation.
* `delete` - (Defaults to 5 mins) Used when delete the Basic Accelerate Ip Endpoint Relation.

## Import

Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_accelerate_ip_endpoint_relation.example <accelerator_id>:<accelerate_ip_id>:<endpoint_id>
```
