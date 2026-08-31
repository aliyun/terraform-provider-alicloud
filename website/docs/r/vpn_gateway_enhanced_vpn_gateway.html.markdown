---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_enhanced_vpn_gateway"
description: |-
  Provides a Alicloud Vpn Gateway Enhanced Vpn Gateway resource.
---

# alicloud_vpn_gateway_enhanced_vpn_gateway

Provides a Vpn Gateway Enhanced Vpn Gateway resource.

Enhanced VPN gateway.

For information about Vpn Gateway Enhanced Vpn Gateway and how to use it, see [What is Enhanced Vpn Gateway](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateEnhancedVpnGateway).

-> **NOTE:** Available since v1.280.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpn_gateway_enhanced_vpn_gateway&exampleId=2a52a5c5-ac67-82ee-d9aa-f2378db75354bf2245ff&activeTab=example&spm=docs.r.vpn_gateway_enhanced_vpn_gateway.0.2a52a5c5ac&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "ap-southeast-3"
}

variable "region" {
  default = "ap-southeast-3"
}

variable "zone2" {
  default = "ap-southeast-3a"
}

variable "zone1" {
  default = "ap-southeast-3b"
}

resource "alicloud_vpc" "defaulttYTx5F" {
  cidr_block = "192.168.0.0/16"
  is_default = false
}

resource "alicloud_vswitch" "defaultTRk7k3" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone1
  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "default23kGFr" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone2
  cidr_block = "192.168.20.0/24"
}


resource "alicloud_vpn_gateway_enhanced_vpn_gateway" "default" {
  vpn_type                     = "Normal"
  description                  = "default"
  disaster_recovery_vswitch_id = alicloud_vswitch.default23kGFr.id
  vpc_id                       = alicloud_vpc.defaulttYTx5F.id
  vpn_gateway_name             = "default"
  network_type                 = "public"
  vswitch_id                   = alicloud_vswitch.defaultTRk7k3.id
  gateway_type                 = "Enhanced.SiteToSite"
  auto_propagate               = false
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpn_gateway_enhanced_vpn_gateway&spm=docs.r.vpn_gateway_enhanced_vpn_gateway.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `auto_propagate` - (Optional) Specifies whether to automatically propagate BGP routes to the VPC
* `description` - (Optional) The description of the VPN gateway.
* `disaster_recovery_vswitch_id` - (Optional, ForceNew) The ID of the backup VSwitch to which the VPN gateway is attached.
* `gateway_type` - (Required, ForceNew) VPN gateway type
* `network_type` - (Optional, ForceNew) Type of Gateway
* `tags` - (Optional, Map) The Tag of
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch to which the VPN gateway is attached.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the VPN gateway belongs.
* `vpn_gateway_name` - (Optional) The name of the VPN gateway.
* `vpn_type` - (Optional, ForceNew) The Type of Vpn

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the VPN gateway was created.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enhanced Vpn Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Enhanced Vpn Gateway.
* `update` - (Defaults to 5 mins) Used when update the Enhanced Vpn Gateway.

## Import

Vpn Gateway Enhanced Vpn Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_gateway_enhanced_vpn_gateway.example <vpn_instance_id>
```
