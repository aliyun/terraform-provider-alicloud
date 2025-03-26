---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_network"
description: |-
  Provides a Alicloud Data Works Network resource.
---

# alicloud_data_works_network

Provides a Data Works Network resource.

Resource Group Network.

For information about Data Works Network and how to use it, see [What is Network](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createnetwork).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_network&exampleId=70cb3a95-6fdb-8c77-c084-609d6e55725652561828&activeTab=example&spm=docs.r.data_works_network.0.70cb3a956f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default5Bia4h" {
  description = var.name
  vpc_name    = var.name
  cidr_block  = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultss7s7F" {
  description  = var.name
  vpc_id       = alicloud_vpc.default5Bia4h.id
  zone_id      = "cn-beijing-g"
  vswitch_name = format("%s1", var.name)
  cidr_block   = "10.0.0.0/24"
}

resource "alicloud_data_works_dw_resource_group" "defaultVJvKvl" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = alicloud_vswitch.defaultss7s7F.id
  remark                = var.name
  resource_group_name   = "network_openapi_example01"
  default_vpc_id        = alicloud_vpc.default5Bia4h.id
}

resource "alicloud_vpc" "defaulte4zhaL" {
  description = var.name
  vpc_name    = format("%s3", var.name)
  cidr_block  = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default675v38" {
  description  = var.name
  vpc_id       = alicloud_vpc.defaulte4zhaL.id
  zone_id      = "cn-beijing-g"
  vswitch_name = format("%s4", var.name)
  cidr_block   = "172.16.0.0/24"
}


resource "alicloud_data_works_network" "default" {
  vpc_id               = alicloud_vpc.defaulte4zhaL.id
  vswitch_id           = alicloud_vswitch.default675v38.id
  dw_resource_group_id = alicloud_data_works_dw_resource_group.defaultVJvKvl.id
}
```

## Argument Reference

The following arguments are supported:
* `dw_resource_group_id` - (Required, ForceNew) The ID of the resource group.
* `vswitch_id` - (Required, ForceNew) The vSwitch ID of the network resource.
* `vpc_id` - (Required, ForceNew) Virtual Private Cloud ID of network resources

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Time when the network resource was created
* `status` - Network Resource Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network.
* `delete` - (Defaults to 5 mins) Used when delete the Network.

## Import

Data Works Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_network.example <id>
```