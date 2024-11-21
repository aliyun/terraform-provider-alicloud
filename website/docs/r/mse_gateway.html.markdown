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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mse_gateway&exampleId=4623b44d-e20a-2904-4bdb-291791729664050253bc&activeTab=example&spm=docs.r.mse_gateway.0.4623b44de2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count        = 2
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = format("172.16.%d.0/21", (count.index + 1) * 16)
  zone_id      = data.alicloud_zones.example.zones[count.index].id
  vswitch_name = format("terraform_example_%d", count.index + 1)
}

resource "alicloud_mse_gateway" "example" {
  gateway_name      = "terraform-example"
  replica           = 2
  spec              = "MSE_GTW_2_4_200_c"
  vswitch_id        = alicloud_vswitch.example.0.id
  backup_vswitch_id = alicloud_vswitch.example.1.id
  vpc_id            = alicloud_vpc.example.id
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

```shell
$ terraform import alicloud_mse_gateway.example <id>
```