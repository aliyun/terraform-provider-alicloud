---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_nat_ip"
description: |-
  Provides a Alicloud Nat Gateway Nat Ip resource.
---

# alicloud_vpc_nat_ip

Provides a Nat Gateway Nat Ip resource.

NAT IP address instance.

For information about Nat Gateway Nat Ip and how to use it, see [What is Nat Ip](https://www.alibabacloud.com/help/en/nat-gateway/developer-reference/api-vpc-2016-04-28-createnatip-natgws).

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_nat_ip&exampleId=4ada155c-b87d-bb11-dc6b-9edd006ac02a6b01dd02&activeTab=example&spm=docs.r.vpc_nat_ip.0.4ada155cb8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.example.zones.0.id
  vswitch_name = "terraform-example"
}

resource "alicloud_nat_gateway" "example" {
  vpc_id               = alicloud_vpc.example.id
  internet_charge_type = "PayByLcu"
  nat_gateway_name     = "terraform-example"
  description          = "terraform-example"
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.example.id
  network_type         = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "example" {
  nat_ip_cidr             = "192.168.0.0/16"
  nat_gateway_id          = alicloud_nat_gateway.example.id
  nat_ip_cidr_description = "terraform-example"
  nat_ip_cidr_name        = "terraform-example"
}

resource "alicloud_vpc_nat_ip" "example" {
  nat_ip             = "192.168.0.37"
  nat_gateway_id     = alicloud_nat_gateway.example.id
  nat_ip_description = "example_value"
  nat_ip_name        = "example_value"
  nat_ip_cidr        = alicloud_vpc_nat_ip_cidr.example.nat_ip_cidr
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_nat_ip&spm=docs.r.vpc_nat_ip.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional, Computed) Specifies whether to only precheck the request. Valid values:
  - `true`: Sends a dry-run request without creating the NAT IP address. The check includes verifying the validity of the AccessKey, RAM user permissions, and whether all required parameters are provided. If the check fails, an error is returned. If the check passes, the error code `DryRunOperation` is returned.
  - `false` (default): Sends a normal request. If the check passes, a 2xx HTTP status code is returned and the NAT IP address is created.

-> **NOTE:** This parameter is only evaluated during resource creation, update and deletion. Modifying it in isolation will not trigger any action.

* `nat_gateway_id` - (Required, ForceNew) The ID of the Virtual Private Cloud (VPC) NAT gateway for which you want to create the NAT IP address.
* `nat_ip` - (Optional, ForceNew) The NAT IP address to be created.
* `nat_ip_cidr` - (Required, ForceNew) The CIDR block to which the NAT IP address belongs.
* `nat_ip_description` - (Optional) The description of the NAT IP address. The description must be `2` to `256` characters in length and start with a letter. The description cannot start with `http://` or `https://`.
* `nat_ip_name` - (Optional) The name of the NAT IP address. The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter. The name must start with a letter and cannot start with `http://` or `https://`.
* `nat_ip_cidr_id` - (Removed since v1.273.0) Field `nat_ip_cidr_id` has been removed from provider version 1.273.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<nat_gateway_id>:<nat_ip_id>`.
* `nat_ip_id` - Ihe ID of the Nat Ip.
* `status` - The status of the NAT IP address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Ip.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Ip.
* `update` - (Defaults to 5 mins) Used when update the Nat Ip.

## Import

Nat Gateway Nat Ip can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_nat_ip.example <nat_gateway_id>:<nat_ip_id>
```
