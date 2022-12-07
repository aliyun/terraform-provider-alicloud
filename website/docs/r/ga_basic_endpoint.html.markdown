---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_endpoint"
sidebar_current: "docs-alicloud-resource-ga-basic-endpoint"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Endpoint resource.
---

# alicloud\_ga\_basic\_endpoint

Provides a Global Accelerator (GA) Basic Endpoint resource.

For information about Global Accelerator (GA) Basic Endpoint and how to use it, see [What is Basic Endpoint](https://help.aliyun.com/document_detail/466839.html).

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

resource "alicloud_ecs_network_interface" "default" {
  provider           = "alicloud.sz"
  vswitch_id         = "your_vswitch_id"
  security_group_ids = ["your_security_group_id"]
}

resource "alicloud_ga_basic_endpoint" "default" {
  provider                  = "alicloud.hz"
  accelerator_id            = "your_accelerator_id"
  endpoint_group_id         = "your_endpoint_group_id"
  endpoint_type             = "ENI"
  endpoint_address          = alicloud_ecs_network_interface.default.id
  endpoint_sub_address_type = "secondary"
  endpoint_sub_address      = "192.168.0.1"
  basic_endpoint_name       = "example_value"
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

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Endpoint.
* `update` - (Defaults to 3 mins) Used when update the Basic Endpoint.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Endpoint.

## Import

Global Accelerator (GA) Basic Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_endpoint.example <endpoint_group_id>:<endpoint_id>
```
