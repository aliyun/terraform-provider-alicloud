---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_ip"
sidebar_current: "docs-alicloud-resource-ddos-bgp-ip"
description: |-
  Provides a Alicloud Ddos Bgp Ip resource.
---

# alicloud\_ddos\_bgp\_ip

Provides a Ddos Bgp Ip resource.

For information about Ddos Bgp Ip and how to use it, see [What is Ip](https://www.alibabacloud.com/help/en/ddos-protection/latest/addip).

-> **NOTE:** Available in v1.180.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eip_address" "default" {
  address_name = "${var.name}"
}

data "alicloud_ddosbgp_instances" default {}

resource "alicloud_ddosbgp_ip" "default" {
  instance_id       = data.alicloud_ddosbgp_instances.default.ids.0
  ip                = alicloud_eip_address.default.ip_address
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the native protection enterprise instance to be operated.
* `ip` - (Required, ForceNew) The IP address.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Ip. The value formats as `<instance_id>:<ip>`.
* `status` - The current state of the IP address. Valid Value: `normal`, `hole_begin`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ddos Bgp Ip.
* `delete` - (Defaults to 1 mins) Used when deleting the Ddos Bgp Ip.

## Import

Ddos Bgp Ip can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddosbgp_ip.example <instance_id>:<ip>
```