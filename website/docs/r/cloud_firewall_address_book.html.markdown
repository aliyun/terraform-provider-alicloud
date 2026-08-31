---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_address_book"
sidebar_current: "docs-alicloud-resource-cloud-firewall-address-book"
description: |-
  Provides a Alicloud Cloud Firewall Address Book resource.
---

# alicloud_cloud_firewall_address_book

Provides a Cloud Firewall Address Book resource.

For information about Cloud Firewall Address Book and how to use it, see [What is Address Book](https://www.alibabacloud.com/help/en/cloud-firewall/developer-reference/api-cloudfw-2017-12-07-addaddressbook).

-> **NOTE:** Available since v1.178.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_address_book&exampleId=25f27cdd-69ee-ab3a-9836-9ee11efea44fac907b7d&activeTab=example&spm=docs.r.cloud_firewall_address_book.0.25f27cdd69&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cloud_firewall_address_book" "example" {
  description      = "example_value"
  group_name       = "example_value"
  group_type       = "tag"
  tag_relation     = "and"
  auto_add_tag_ecs = 0
  ecs_tags {
    tag_key   = "created"
    tag_value = "tfTestAcc0"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_address_book&spm=docs.r.cloud_firewall_address_book.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `group_name` - (Required) The name of the Address Book.
* `group_type` - (Required, ForceNew) The type of the Address Book. Valid values: `ip`, `ipv6`, `domain`, `port`, `tag`, `asset`, `assetIpv6`.
**NOTE:** From version 1.213.1, `group_type` can be set to `ipv6`, `domain`, `port`. From version 1.286.0, `group_type` can be set to `asset`, `assetIpv6`.
* `description` - (Required) The description of the Address Book.
* `auto_add_tag_ecs` - (Optional, Int) Whether you want to automatically add new matching tags of the ECS IP address to the Address Book. Valid values: `0`, `1`.
* `tag_relation` - (Optional) The logical relation among the ECS tags that to be matched. Default value: `and`. Valid values:
  - `and`: Only the public IP addresses of ECS instances that match all the specified tags can be added to the Address Book.
  - `or`: The public IP addresses of ECS instances that match one of the specified tags can be added to the Address Book.
* `lang` - (Optional) The language of the content within the request and response. Valid values: `zh`, `en`.
* `address_list` - (Optional, List) The list of addresses.
* `ecs_tags` - (Optional, Set) A list of ECS tags. See [`ecs_tags`](#ecs_tags) below.
* `asset_member_uids` - (Optional, List, Available since v1.286.0) The list of member account UIDs of the asset Address Book.
* `asset_region_resource_types` - (Optional, List, Available since v1.286.0) The list of regions and asset types of the asset Address Book. See [`asset_region_resource_types`](#asset_region_resource_types) below.

### `ecs_tags`

The ecs_tags supports the following:

* `tag_key` - (Optional) The key of ECS tag that to be matched.
* `tag_value` - (Optional) The value of ECS tag that to be matched.

### `asset_region_resource_types`

The asset_region_resource_types supports the following:

* `asset_region_id` - (Optional) The region ID of the assets. Set the value to `all` to specify all the regions. **NOTE:** `asset_region_id` cannot be modified after the Address Book is created.
* `resource_type` - (Optional, List) The types of the assets. See [`resource_type`](#asset_region_resource_types-resource_type) below.

### `asset_region_resource_types-resource_type`

The resource_type supports the following:

* `ipv4` - (Optional, List) The IPv4 asset types. See [`ipv4`](#asset_region_resource_types-resource_type-ipv4) below.
* `ipv6` - (Optional, List) The IPv6 asset types. See [`ipv6`](#asset_region_resource_types-resource_type-ipv6) below.

### `asset_region_resource_types-resource_type-ipv4`

The ipv4 supports the following:

* `eip` - (Optional, Bool) Whether to include the assets of the type EIP.
* `ecs_eip` - (Optional, Bool) Whether to include the assets of the type EcsEIP.
* `ecs_public_ip` - (Optional, Bool) Whether to include the assets of the type EcsPublicIP.
* `slb_eip` - (Optional, Bool) Whether to include the assets of the type SlbEIP.
* `slb_public_ip` - (Optional, Bool) Whether to include the assets of the type SlbPublicIP.
* `nlb_eip` - (Optional, Bool) Whether to include the assets of the type NlbEIP.
* `alb_eip` - (Optional, Bool) Whether to include the assets of the type AlbEIP.
* `nat_eip` - (Optional, Bool) Whether to include the assets of the type NatEIP.
* `nat_public_ip` - (Optional, Bool) Whether to include the assets of the type NatPublicIP.
* `eni_eip` - (Optional, Bool) Whether to include the assets of the type EniEIP.
* `ga_eip` - (Optional, Bool) Whether to include the assets of the type GaEIP.
* `api_gateway_eip` - (Optional, Bool) Whether to include the assets of the type ApiGatewayEIP.
* `ai_gateway_eip` - (Optional, Bool) Whether to include the assets of the type AiGatewayEIP.
* `bastion_host_ip` - (Optional, Bool) Whether to include the assets of the type BastionHostIP.
* `bastion_host_ingress_ip` - (Optional, Bool) Whether to include the assets of the type BastionHostIngressIP.
* `bastion_host_egress_ip` - (Optional, Bool) Whether to include the assets of the type BastionHostEgressIP.
* `havip` - (Optional, Bool) Whether to include the assets of the type HAVIP.

### `asset_region_resource_types-resource_type-ipv6`

The ipv6 supports the following:

* `ecs_ipv6` - (Optional, Bool) Whether to include the assets of the type EcsIPv6.
* `slb_ipv6` - (Optional, Bool) Whether to include the assets of the type SlbIPv6.
* `nlb_ipv6` - (Optional, Bool) Whether to include the assets of the type NlbIPv6.
* `alb_ipv6` - (Optional, Bool) Whether to include the assets of the type AlbIPv6.
* `eni_eipv6` - (Optional, Bool) Whether to include the assets of the type EniEIPv6.
* `ga_eipv6` - (Optional, Bool) Whether to include the assets of the type GaEIPv6.
* `api_gateway_eipv6` - (Optional, Bool) Whether to include the assets of the type ApiGatewayEIPv6.
* `ai_gateway_eipv6` - (Optional, Bool) Whether to include the assets of the type AiGatewayEIPv6.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Address Book.
* `address_list_count` - (Available since v1.286.0) The number of addresses in the Address Book.
* `reference_count` - (Available since v1.286.0) The number of times that the Address Book is referenced.

## Import

Cloud Firewall Address Book can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_address_book.example <id>
```
