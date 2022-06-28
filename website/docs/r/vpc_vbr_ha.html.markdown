---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_vbr_ha"
sidebar_current: "docs-alicloud-resource-vpc-vbr-ha"
description: |-
  Provides a Alicloud VPC Vbr Ha resource.
---

# alicloud\_vpc\_vbr\_ha

Provides a VPC Vbr Ha resource.

For information about VPC Vbr Ha and how to use it, see [What is Vbr Ha](https://www.alibabacloud.com/help/doc-detail/212629.html).

-> **NOTE:** Available in v1.151.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "example" {
  count                      = 2
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections[count.index].id
  virtual_border_router_name = var.name
  vlan_id                    = 100
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "example_value"
  description       = "example_value"
}

resource "alicloud_cen_instance_attachment" "example" {
  count                    = 2
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_express_connect_virtual_border_router.example[count.index].id
  child_instance_type      = "VBR"
  child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_vbr_ha" "example" {
  depends_on  = ["alicloud_cen_instance_attachment.example"]
  vbr_id      = alicloud_cen_instance_attachment.example[0].id
  peer_vbr_id = alicloud_cen_instance_attachment.example[1].id
  vbr_ha_name = "example_value"
  description = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the VBR switching group. It must be `2` to `256` characters in length and must start with a letter or Chinese, but cannot start with `https://` or `https://`.
* `dry_run` - (Optional) The dry run.
* `peer_vbr_id` - (Required, ForceNew) The ID of the other VBR in the VBR failover group.
* `vbr_ha_name` - (Optional, ForceNew) The name of the VBR failover group.
* `vbr_id` - (Required, ForceNew) The ID of the VBR instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vbr Ha.
* `status` - The state of the VBR failover group.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Vbr Ha.
* `delete` - (Defaults to 10 mins) Used when delete the Vbr Ha.

## Import

VPC Vbr Ha can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_vbr_ha.example <id>
```