---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerate_ip_endpoint_relation"
sidebar_current: "docs-alicloud-resource-ga-basic-accelerate-ip-endpoint-relation"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation resource.
---

# alicloud\_ga\_basic\_accelerate\_ip\_endpoint\_relation

Provides a Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation resource.

For information about Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation and how to use it, see [What is Basic Accelerate Ip Endpoint Relation](https://help.aliyun.com/document_detail/466842.html).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  alias  = "sz"
  region = "cn-shenzhen"
}

provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  provider   = "alicloud.sz"
  name_regex = "your_vpc_name"
}

data "alicloud_vswitches" "default" {
  provider = "alicloud.sz"
  vpc_id   = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  provider = "alicloud.sz"
  vpc_id   = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecs_network_interface" "default" {
  provider           = "alicloud.sz"
  vswitch_id         = data.alicloud_vswitches.default.ids.0
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  pricing_cycle          = "Month"
  basic_accelerator_name = var.name
  description            = var.name
  bandwidth_billing_type = "CDT"
  auto_pay               = true
  auto_use_coupon        = "true"
  auto_renew             = false
  auto_renew_duration    = 1
}

resource "alicloud_ga_basic_ip_set" "default" {
  accelerator_id       = alicloud_ga_basic_accelerator.default.id
  accelerate_region_id = "cn-hangzhou"
  isp_type             = "BGP"
  bandwidth            = "5"
}

resource "alicloud_ga_basic_accelerate_ip" "default" {
  accelerator_id = alicloud_ga_basic_ip_set.default.accelerator_id
  ip_set_id      = alicloud_ga_basic_ip_set.default.id
}

resource "alicloud_ga_basic_endpoint_group" "default" {
  accelerator_id        = alicloud_ga_basic_accelerator.default.id
  endpoint_group_region = "cn-shenzhen"
}

resource "alicloud_ga_basic_endpoint" "default" {
  accelerator_id            = alicloud_ga_basic_accelerator.default.id
  endpoint_group_id         = alicloud_ga_basic_endpoint_group.default.id
  endpoint_type             = "ENI"
  endpoint_address          = alicloud_ecs_network_interface.default.id
  endpoint_sub_address_type = "primary"
  endpoint_sub_address      = "192.168.0.1"
  basic_endpoint_name       = var.name
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

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Basic Accelerate Ip Endpoint Relation.
* `delete` - (Defaults to 5 mins) Used when delete the Basic Accelerate Ip Endpoint Relation.

## Import

Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_accelerate_ip_endpoint_relation.example <accelerator_id>:<accelerate_ip_id>:<endpoint_id>
```
