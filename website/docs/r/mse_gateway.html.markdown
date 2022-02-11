---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_gateway"
sidebar_current: "docs-alicloud-resource-mse-gateway"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Gateway resource.
---

# alicloud\_mse\_gateway

Provides a Microservice Engine (MSE) Gateway resource.

For information about Microservice Engine (MSE) Gateway and how to use it, see [What is Gateway](https://help.aliyun.com/document_detail/347638.html).

-> **NOTE:** Available in v1.157.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_mse_gateway" "example" {
  gateway_name      = "example_value"
  replica           = 2
  spec              = "MSE_GTW_2_4_200_c"
  vswitch_id        = data.alicloud_vswitches.default.ids.0
  backup_vswitch_id = data.alicloud_vswitches.default.ids.1
  vpc_id            = data.alicloud_vpcs.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `backup_vswitch_id` - (Optional, ForceNew) The backup vswitch id.
* `enterprise_security_group` - (Optional) Whether the enterprise security group type.
* `gateway_name` - (Optional) The name of the Gateway .
* `internet_slb_spec` - (Optional) Public network SLB specifications.
* `replica` - (Required, ForceNew) Number of Gateway Nodes.
* `slb_spec` - (Optional) Private network SLB specifications.
* `spec` - (Required, ForceNew) Gateway Node Specifications. Valid values: `MSE_GTW_2_4_200_c`, `MSE_GTW_4_8_200_c`, `MSE_GTW_8_16_200_c`, `MSE_GTW_16_32_200_c`.
* `vswitch_id` - (Required, ForceNew) The ID of the vswitch.
* `vpc_id` - (Required, ForceNew) The ID of the vpc.
* `delete_slb` - (Optional) Whether to delete the SLB purchased on behalf of the gateway at the same time.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway.
* `status` - The status of the gateway.
* `slb_list` - A list of gateway Slb.
  * `associate_id` - The associate id.
  * `slb_id` - The ID of the gateway slb.
  * `slb_ip` - The ip of the gateway slb.
  * `slb_port` - The port of the gateway slb.
  * `type` - The type of the gateway slb.
  * `gmt_create` - The creation time of the gateway slb.
  * `gateway_slb_mode` - The Mode of the gateway slb.
  * `gateway_slb_status` - The Status of the gateway slb.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway.

## Import

Microservice Engine (MSE) Gateway can be imported using the id, e.g.

```
$ terraform import alicloud_mse_gateway.example <id>
```