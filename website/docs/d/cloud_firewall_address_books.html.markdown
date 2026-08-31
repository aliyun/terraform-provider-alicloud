---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_address_books"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-address-books"
description: |-
  Provides a list of Cloud Firewall Address Books to the user.
---

# alicloud_cloud_firewall_address_books

This data source provides the Cloud Firewall Address Books of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.178.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_cloud_firewall_address_book" "default" {
  group_name       = var.name
  group_type       = "ip"
  description      = "tf-description"
  auto_add_tag_ecs = 0
  address_list     = ["10.21.0.0/16", "10.168.0.0/16"]
}

data "alicloud_cloud_firewall_address_books" "ids" {
  ids = [alicloud_cloud_firewall_address_book.default.id]
}

output "cloud_firewall_address_book_id_1" {
  value = data.alicloud_cloud_firewall_address_books.ids.books.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Address Book IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results Address Book name.
* `group_type` - (Optional, ForceNew) The type of the Address Book. Valid values: `ip`, `ipv6`, `domain`, `port`, `tag`, `asset`, `assetIpv6`.
  **NOTE:** From version 1.213.1, `group_type` can be set to `ipv6`, `domain`, `port`. From version 1.286.0, `group_type` can be set to `asset`, `assetIpv6`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Address Book names.
* `books` - A list of Cloud Firewall Address Books. Each element contains the following attributes:
  * `id` - The ID of the Address Book.
  * `group_uuid` - The ID of the Address Book.  
  * `group_name` - The name of the Address Book.
  * `group_type` - The type of the Address Book.
  * `description` - The description of the Address Book.
  * `auto_add_tag_ecs` - Whether you want to automatically add new matching tags of the ECS IP address to the Address Book.
  * `tag_relation` - One or more tags for the relationship between.
  * `address_list` - The addresses in the Address Book.
  * `address_list_count` - (Available since v1.286.0) The number of addresses in the Address Book.
  * `reference_count` - (Available since v1.286.0) The number of times that the Address Book is referenced.
  * `asset_member_uids` - (Available since v1.286.0) The list of member account UIDs of the asset Address Book.
  * `asset_region_resource_types` - (Available since v1.286.0) The list of regions and asset types of the asset Address Book.
      * `asset_region_id` - The region ID of the assets.
      * `resource_type` - The types of the assets.
          * `ipv4` - The IPv4 asset types.
              * `eip` - Whether the assets of the type EIP are included.
              * `ecs_eip` - Whether the assets of the type EcsEIP are included.
              * `ecs_public_ip` - Whether the assets of the type EcsPublicIP are included.
              * `slb_eip` - Whether the assets of the type SlbEIP are included.
              * `slb_public_ip` - Whether the assets of the type SlbPublicIP are included.
              * `nlb_eip` - Whether the assets of the type NlbEIP are included.
              * `alb_eip` - Whether the assets of the type AlbEIP are included.
              * `nat_eip` - Whether the assets of the type NatEIP are included.
              * `nat_public_ip` - Whether the assets of the type NatPublicIP are included.
              * `eni_eip` - Whether the assets of the type EniEIP are included.
              * `ga_eip` - Whether the assets of the type GaEIP are included.
              * `api_gateway_eip` - Whether the assets of the type ApiGatewayEIP are included.
              * `ai_gateway_eip` - Whether the assets of the type AiGatewayEIP are included.
              * `bastion_host_ip` - Whether the assets of the type BastionHostIP are included.
              * `bastion_host_ingress_ip` - Whether the assets of the type BastionHostIngressIP are included.
              * `bastion_host_egress_ip` - Whether the assets of the type BastionHostEgressIP are included.
              * `havip` - Whether the assets of the type HAVIP are included.
          * `ipv6` - The IPv6 asset types.
              * `ecs_ipv6` - Whether the assets of the type EcsIPv6 are included.
              * `slb_ipv6` - Whether the assets of the type SlbIPv6 are included.
              * `nlb_ipv6` - Whether the assets of the type NlbIPv6 are included.
              * `alb_ipv6` - Whether the assets of the type AlbIPv6 are included.
              * `eni_eipv6` - Whether the assets of the type EniEIPv6 are included.
              * `ga_eipv6` - Whether the assets of the type GaEIPv6 are included.
              * `api_gateway_eipv6` - Whether the assets of the type ApiGatewayEIPv6 are included.
              * `ai_gateway_eipv6` - Whether the assets of the type AiGatewayEIPv6 are included.
  * `ecs_tags` - The logical relation among the ECS tags that to be matched.
    * `tag_key` - The key of ECS tag that to be matched.
    * `tag_value` - The value of ECS tag that to be matched.
