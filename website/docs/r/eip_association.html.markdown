---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_association"
description: |-
  Provides a Alicloud EIP Association resource.
---

# alicloud_eip_association

Provides a EIP Association resource.

-> **NOTE:** `alicloud_eip_association` is useful in scenarios where EIPs are either
 pre-existing or distributed to customers or users and therefore cannot be changed.

-> **NOTE:** From version 1.7.1, the resource support to associate EIP to SLB Instance or Nat Gateway.

-> **NOTE:** One EIP can only be associated with ECS or SLB instance which in the VPC.

For information about EIP Association and how to use it, see [What is Association](https://www.alibabacloud.com/help/en/vpc/developer-reference/api-vpc-2016-04-28-associateeipaddress).

-> **NOTE:** Available since v1.117.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eip_association&exampleId=b09889bb-255e-2b33-a830-8d7e1c81b76bcebf48a3&activeTab=example&spm=docs.r.eip_association.0.b09889bb25&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.example.images.0.id
  instance_type     = data.alicloud_instance_types.example.instance_types.0.id
  security_groups   = [alicloud_security_group.example.id]
  vswitch_id        = alicloud_vswitch.example.id
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_eip_address" "example" {
  address_name = var.name
}

resource "alicloud_eip_association" "example" {
  allocation_id = alicloud_eip_address.example.id
  instance_id   = alicloud_instance.example.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eip_association&spm=docs.r.eip_association.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `allocation_id` - (Required, ForceNew) The ID of the EIP instance.
* `force` - (Optional) Specifies whether to disassociate the EIP from a NAT gateway if a DNAT or SNAT entry is added to the NAT gateway. Valid values:
  - `false` (default)
  - `true`

* `instance_id` - (Required, ForceNew) The ID of the instance with which you want to associate the EIP. You can enter the ID of a NAT gateway, CLB instance, ECS instance, secondary ENI, HAVIP, or IP address. 
* `instance_type` - (Optional, ForceNew, Computed) The type of the instance with which you want to associate the EIP. Valid values:
  - `Nat`: NAT gateway
  - `SlbInstance`: CLB instance
  - `EcsInstance` (default): ECS instance
  - `NetworkInterface`: secondary ENI
  - `HaVip`: HAVIP
  - `IpAddress`: IP address

-> **NOTE:**   The default value is `EcsInstance`. If the instance with which you want to associate the EIP is not an ECS instance, this parameter is required.

* `mode` - (Optional, Computed, Available since v1.216.0) The association mode. Valid values:
  - `NAT` (default): NAT mode
  - `MULTI_BINDED`: multi-EIP-to-ENI mode
  - `BINDED`: cut-network interface controller mode

-> **NOTE:**   This parameter is required only when `instance_type` is set to `NetworkInterface`.

* `private_ip_address` - (Optional, ForceNew) The IP address in the CIDR block of the vSwitch.

  If you leave this parameter empty, the system allocates a private IP address based on the VPC ID and vSwitch ID.

-> **NOTE:**   This parameter is required if `instance_type` is set to `IpAddress`, which indicates that the EIP is to be associated with an IP address.

* `vpc_id` - (Optional, ForceNew) The ID of the VPC in which an IPv4 gateway is created. The VPC and the EIP must be in the same region.

  When you associate an EIP with an IP address, the system can enable the IP address to access the Internet based on VPC route configurations.

-> **NOTE:**   This parameter is required if `instance_type` is set to `IpAddress`, which indicates that the EIP is to be associated with an IP address.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<allocation_id>:<instance_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Association.
* `delete` - (Defaults to 5 mins) Used when delete the Association.
* `update` - (Defaults to 5 mins) Used when update the Association.

## Import

EIP Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_eip_association.example <allocation_id>:<instance_id>
```