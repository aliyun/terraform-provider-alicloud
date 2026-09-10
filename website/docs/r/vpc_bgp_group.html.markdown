---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_group"
description: |-
  Provides a Alicloud Express Connect Bgp Group resource.
---

# alicloud_vpc_bgp_group

Provides a Express Connect Bgp Group resource.



For information about Express Connect Bgp Group and how to use it, see [What is Bgp Group](https://www.alibabacloud.com/help/en/doc-detail/91267.html).

-> **NOTE:** Available since v1.152.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_bgp_group&exampleId=2e8b1fd0-1a30-9625-f33a-b8c5d9f2adc74c87ff15&activeTab=example&spm=docs.r.vpc_bgp_group.0.2e8b1fd01a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}
resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}
resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.vlan_id.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_group" "example" {
  auth_key       = "YourPassword+12345678"
  bgp_group_name = var.name
  description    = var.name
  peer_asn       = 1111
  router_id      = alicloud_express_connect_virtual_border_router.example.id
  is_fake_asn    = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_bgp_group&spm=docs.r.vpc_bgp_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `auth_key` - (Optional) The authentication key of the BGP group. 
* `bgp_group_name` - (Optional) The name of the BGP group. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The name must start with a letter. 
* `clear_auth_key` - (Optional, Available since v1.261.0) Specifies whether to clear the secret key. Valid values: `true`, `false`.
* `description` - (Optional) The description of the BGP group. The description must be 2 to 256 characters in length. It must start with a letter and cannot start with `http://` or `https://`. 
* `ip_version` - (Optional, ForceNew, Computed, Available since v1.261.0) The IP version. Valid values:
  - `IPv4`: This is the default value.
  - `IPv6`: IPv6 is supported only if the VBR for which you want to create the BGP group has IPv6 enabled.
* `is_fake_asn` - (Optional, Computed) Specifies whether to use a fake AS number. Valid values: `true`, `false`.

-> **NOTE:** Note A router that runs BGP typically belongs to only one AS. If you need to replace an AS with a new one, but you cannot immediately modify BGP configurations due to business requirements, you can specify a fake AS number to establish a connection with the local end. This ensures service continuity in scenarios such as AS migration or AS merging.

* `local_asn` - (Optional, Computed) The custom ASN on the Alibaba Cloud side. Valid values:
  - `45104`
  - **64512~65534** 
  - **4200000000~4294967294**.

-> **NOTE:**  `65025` is reserved by Alibaba Cloud. By default, Alibaba Cloud uses `45104` as `LocalAsn`. If you use custom `LocalAsn` in multi-line access scenarios, loops in BGP may occur. 

* `peer_asn` - (Required) The ASN of the gateway device in the data center. 
* `route_limit` - (Optional, Computed, Int, Available since v1.261.0) The maximum number of routes supported by a BGP peer. Default value: 110.
* `router_id` - (Required, ForceNew) The ID of the virtual border router (VBR) that is associated with the BGP group. 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The Region ID of the BGP group.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bgp Group.
* `delete` - (Defaults to 5 mins) Used when delete the Bgp Group.
* `update` - (Defaults to 5 mins) Used when update the Bgp Group.

## Import

Express Connect Bgp Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_bgp_group.example <id>
```