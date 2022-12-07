---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_endpoint_group"
sidebar_current: "docs-alicloud-resource-ga-basic-endpoint-group"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Endpoint Group resource.
---

# alicloud\_ga\_basic\_endpoint\_group

Provides a Global Accelerator (GA) Basic Endpoint Group resource.

For information about Global Accelerator (GA) Basic Endpoint Group and how to use it, see [What is Basic Endpoint Group](https://www.alibabacloud.com/help/en/global-accelerator/latest/createbasicendpointgroup).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  pricing_cycle          = "Month"
  bandwidth_billing_type = "CDT"
  auto_pay               = true
  auto_use_coupon        = "true"
  auto_renew             = false
  auto_renew_duration    = 1
}

resource "alicloud_ga_basic_endpoint_group" "default" {
  accelerator_id            = alicloud_ga_basic_accelerator.default.id
  endpoint_group_region     = "cn-beijing"
  endpoint_type             = "SLB"
  endpoint_address          = alicloud_slb_load_balancer.default.id
  endpoint_sub_address      = "192.168.0.1"
  basic_endpoint_group_name = "example_value"
  description               = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the basic GA instance.
* `endpoint_group_region` - (Required, ForceNew) The ID of the region where you want to create the endpoint group.
* `endpoint_type` - (Optional, ForceNew, Computed) The type of the endpoint. Valid values: `ENI`, `SLB` and `ECS`.
* `endpoint_address` - (Optional, ForceNew, Computed) The address of the endpoint.
* `endpoint_sub_address` - (Optional, ForceNew, Computed) The sub address of the endpoint.
* `basic_endpoint_group_name` - (Optional) The name of the endpoint group. The `basic_endpoint_group_name` must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter.
* `description` - (Optional) The description of the endpoint group. The `description` cannot exceed 256 characters in length and cannot contain http:// or https://.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Endpoint Group.
* `status` - The status of the Basic Endpoint Group.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Endpoint Group.
* `update` - (Defaults to 3 mins) Used when update the Basic Endpoint Group.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Endpoint Group.

## Import

Global Accelerator (GA) Basic Endpoint Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_endpoint_group.example <id>
```
