---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_udp_reflect"
sidebar_current: "docs-alicloud-resource-ddoscoo-udp-reflect"
description: |-
  Provides a Alicloud Ddos Coo Udp Reflect resource.
---

# alicloud\_ddoscoo\_udp\_reflect

Provides a Ddos Coo Udp Reflect resource.

For information about UDP reflect attack protection and how to use it, see [ConfigUdpReflect](https://www.alibabacloud.com/help/en/ddos-protection/latest/api-ddoscoo-2020-01-01-configudpreflect).

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_ddoscoo_instances" "default" {}

resource "alicloud_ddoscoo_udp_reflect" "default" {
  instance_id = data.alicloud_ddoscoo_instances.default.ids.0
  udp_sports  = [17, 19]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Anti-DDoS Pro or Anti-DDoS Premium instance. Each instance owns exactly one UDP reflect protection policy, so the instance ID is also the resource ID; changing it creates a new resource.

  -> **NOTE:** You can call [DescribeInstanceIds](https://www.alibabacloud.com/help/en/doc-detail/157459.html) to query the IDs of all instances.

* `udp_sports` - (Required, List) The full set of UDP reflect source ports to block. The API applies this list as a full replacement, so reducing the list removes the dropped ports from the policy. Recommended common UDP reflect attack source ports: `17` (QOTD), `19` (CharGEN), `69` (TFTP), `111` (Portmap), `123` (NTP), `137` (NetBIOS), `161` (SNMPv2), `389` (CLDAP), `1194` (OpenVPN), `1900` (SSDP), `3389` (RDP), `11211` (Memcached).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which equals the `instance_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when configuring the UDP reflect policy.
* `delete` - (Defaults to 5 mins) Used when clearing the UDP reflect policy.
* `update` - (Defaults to 5 mins) Used when updating the UDP reflect policy.

## Import

Ddos Coo Udp Reflect can be imported using the id, which is the instance id, e.g.

```shell
$ terraform import alicloud_ddoscoo_udp_reflect.example <instance_id>
```
