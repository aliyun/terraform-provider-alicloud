---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_member"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-member"
description: |-
  Provides a Alicloud Cen Transit Router Multicast Domain Member resource.
---

# alicloud_cen_transit_router_multicast_domain_member

Provides a Cen Transit Router Multicast Domain Member resource.

For information about Cen Transit Router Multicast Domain Member and how to use it, see [What is Transit Router Multicast Domain Member](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-registertransitroutermulticastgroupmembers).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_multicast_domain_member&exampleId=ec201f74-c4c0-0019-d791-cc9fc3743e497ba0dc64&activeTab=example&spm=docs.r.cen_transit_router_multicast_domain_member.0.ec201f74c4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

locals {
  zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.zone
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_ecs_network_interface" "example" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.example.id
  primary_ip_address     = cidrhost(alicloud_vswitch.example.cidr_block, 100)
  security_group_ids     = [alicloud_security_group.example.id]
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
  support_multicast   = true
}

resource "alicloud_cen_transit_router_multicast_domain" "example" {
  transit_router_id                    = alicloud_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name = var.name
}

resource "alicloud_cen_transit_router_vpc_attachment" "example" {
  cen_id            = alicloud_cen_transit_router.example.cen_id
  transit_router_id = alicloud_cen_transit_router_multicast_domain.example.transit_router_id
  vpc_id            = alicloud_vpc.example.id
  zone_mappings {
    zone_id    = local.zone
    vswitch_id = alicloud_vswitch.example.id
  }
}

resource "alicloud_cen_transit_router_multicast_domain_association" "example" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.example.id
  transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.example.transit_router_attachment_id
  vswitch_id                         = alicloud_vswitch.example.id
}

resource "alicloud_cen_transit_router_multicast_domain_member" "example" {
  vpc_id                             = alicloud_vpc.example.id
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain_association.example.transit_router_multicast_domain_id
  network_interface_id               = alicloud_ecs_network_interface.example.id
  group_ip_address                   = "239.1.1.1"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_transit_router_multicast_domain_member&spm=docs.r.cen_transit_router_multicast_domain_member.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast member belongs.
* `group_ip_address` - (Required, ForceNew) The IP address of the multicast group to which the multicast member belongs. If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you in the current multicast domain.
* `network_interface_id` - (Required, ForceNew) The ID of the ENI.
* `vpc_id` - (Optional, ForceNew) The VPC to which the ENI of the multicast member belongs. This field is mandatory for VPCs owned by another accounts.
* `dry_run` - (Optional, Bool) Specifies whether only to precheck the request.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. It formats as `<transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>`.
* `status` - The status of the Transit Router Multicast Domain Member.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Transit Router Multicast Domain Member.
* `delete` - (Defaults to 10 mins) Used when delete the Transit Router Multicast Domain Member.

## Import

Cen Transit Router Multicast Domain Member can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_member.example <transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>
```
