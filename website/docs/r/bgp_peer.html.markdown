---
subcategory: "BGP"
layout: "alicloud"
page_title: "Alicloud: alicloud_bgp_peer"
sidebar_current: "docs-alicloud-resource-bgp-peer"
description: |-
  Provides a Alicloud BGP Peer resource.
---

# alicloud\_bgp_peer

Provides a BGP Peer resource.

For information about BGP Peer and how to use it, see [alicloud_bgp_peer](https://www.alibabacloud.com/help/doc-detail/144682.html)

-> **NOTE:** Terraform will auto build bgp peer instance while it uses `alicloud_bgp_peer` to build a bgp peer resource.

-> **NOTE:** Available in v1.91.0.

## Example Usage

Basic Usage

```
resource "alicloud_bgp_group" "foo" {
    peer_asn = 2
    router_id = "vbr-xxxxxxx"
    description = "test-description11"
    name = "test-name"
    is_fake_asn = true
    auth_key= "dasdasda"
}

resource "alicloud_bgp_peer" "foo" {
  bgp_group_id    = "${alicloud_bgp_group.foo.id}"
  peer_ip_address   = "192.168.2.11"
}
```
## Argument Reference

The following arguments are supported:

* `bgp_group_id` - (Required, ForceNew) The Id of the BGP Group.
* `peer_ip_address` - (Required, ForceNew) The IP address of the BGP neighbor.
* `enable_bfd` - (Optional) Whether enable BFD function or not. Default to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the BGP Peer id.

## Import

BGP Peer can be imported using the id, e.g.

```
$ terraform import alicloud_bgp_peer.example bgp-awjhxxxxxxxx
```


