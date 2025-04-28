---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_peer_member"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-peer-member"
description: |-
  Provides a Alicloud Cen Transit Router Multicast Domain Peer Member resource.
---

# alicloud_cen_transit_router_multicast_domain_peer_member

Provides a Cen Transit Router Multicast Domain Peer Member resource.

For information about Cen Transit Router Multicast Domain Peer Member and how to use it, see [What is Transit Router Multicast Domain Peer Member](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-cbn-2017-09-12-deregistertransitroutermulticastgroupmembers).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_multicast_domain_peer_member&exampleId=1b46812f-14ac-adb4-aa3f-4d6bcab51713fe9da6d7&activeTab=example&spm=docs.r.cen_transit_router_multicast_domain_peer_member.0.1b46812f14&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "qd"
  region = "cn-qingdao"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_bandwidth_package" "default" {
  bandwidth                  = 5
  cen_bandwidth_package_name = var.name
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "default" {
  provider          = alicloud.hz
  cen_id            = alicloud_cen_bandwidth_package_attachment.default.instance_id
  support_multicast = true
}

resource "alicloud_cen_transit_router" "peer" {
  provider          = alicloud.qd
  cen_id            = alicloud_cen_bandwidth_package_attachment.default.instance_id
  support_multicast = true
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  provider                              = alicloud.hz
  cen_id                                = alicloud_cen_bandwidth_package_attachment.default.instance_id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  peer_transit_router_id                = alicloud_cen_transit_router.peer.transit_router_id
  peer_transit_router_region_id         = "cn-qingdao"
  cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  bandwidth                             = 5
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  provider                                    = alicloud.hz
  transit_router_id                           = alicloud_cen_transit_router_peer_attachment.default.transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
}

resource "alicloud_cen_transit_router_multicast_domain" "peer" {
  provider                                    = alicloud.qd
  transit_router_id                           = alicloud_cen_transit_router_peer_attachment.default.peer_transit_router_id
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = var.name
}

resource "alicloud_cen_transit_router_multicast_domain_peer_member" "default" {
  provider                                = alicloud.hz
  transit_router_multicast_domain_id      = alicloud_cen_transit_router_multicast_domain.default.id
  peer_transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.peer.id
  group_ip_address                        = "224.1.0.0"
}
```

## Argument Reference

The following arguments are supported:
* `group_ip_address` - (Required, ForceNew) The IP address of the multicast group to which the multicast member belongs. Value range: **224.0.0.1** to **239.255.255.254**.If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you in the current multicast domain.
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast member belongs.
* `peer_transit_router_multicast_domain_id` - (Required, ForceNew) The IDs of the inter-region multicast domains.
* `dry_run` - (Optional) Specifies whether only to precheck the request.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<transit_router_multicast_domain_id>:<group_ip_address>:<peer_transit_router_multicast_domain_id>`.
* `status` - The status of the multicast resource. Valid values:
  - Registering: being created
  - Registered: available
  - Deregistering: being deleted

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Transit Router Multicast Domain Peer Member.
* `delete` - (Defaults to 10 mins) Used when delete the Transit Router Multicast Domain Peer Member.

## Import

Cen Transit Router Multicast Domain Peer Member can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_peer_member.example <transit_router_multicast_domain_id>:<group_ip_address>:<peer_transit_router_multicast_domain_id>
```