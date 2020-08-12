---
subcategory: "BGP"
layout: "alicloud"
page_title: "Alicloud: alicloud_bgp_group"
sidebar_current: "docs-alicloud-resource-bgp-group"
description: |-
  Provides a Alicloud BGP Group resource.
---

# alicloud\_bgp_group

Provides a BGP Group resource.

For information about BGP Group and how to use it, see [alicloud_bgp_group](https://www.alibabacloud.com/help/doc-detail/144675.html)

-> **NOTE:** Terraform will auto build bgp group instance while it uses `alicloud_bgp_group` to build a bgp group resource.

-> **NOTE:** Available in v1.91.0.

## Example Usage

Basic Usage

```
resource "alicloud_bgp_group" "foo" {
  peer_asn    = 2
  router_id   = "vbr-xxxxxxxxxxxxxx"
  description = "test_description"
  name        = "test_name"
  is_fake_asn = false
  auth_key    = "dasdasdasd"
}
```
## Argument Reference

The following arguments are supported:

* `peer_asn` - (Required, ForceNew) The ASN of the side device.
* `router_id` - (Required, ForceNew) The Id of the VBR.
* `description` - (Optional) The Description of the BGP group.
* `name` - (Optional) The Name of the BGP group.
* `is_fake_asn` - (Optional) Generally, a router running BGP can belong to only one AS. In some cases, for example, an AS needs to be migrated or merged with other ASs, and the new AS is used to replace the original AS number.
* `auth_key` - (Optional) The authentication key of the BGP group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the BGP Group id.

## Import

BGP Group can be imported using the id, e.g.

```
$ terraform import alicloud_bgp_group.example bgpg-bp123456
```


