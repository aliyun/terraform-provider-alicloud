---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_security_group"
description: |-
  Provides a Alicloud ENS Security Group resource.
---

# alicloud_ens_security_group

Provides a ENS Security Group resource.



For information about ENS Security Group and how to use it, see [What is Security Group](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createsnapshot).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_security_group" "default" {
  description         = var.name
  security_group_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Security group description information
It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with http:// or https://
* `permissions` - (Optional, List, Available since v1.287.0) A collection of rules for a security group instance See [`permissions`](#permissions) below.
* `security_group_name` - (Optional) Security group name
The security group name. The length is 2~128 English or Chinese characters. It must start with an uppercase or lowcase letter or a Chinese character and cannot start with http:// or https. Can contain digits, colons (:), underscores (_), or hyphens (-)

### `permissions`

The permissions supports the following:
* `creation_time` - (Computed, Available since v1.287.0) Creation time, UTC time.
* `description` - (Optional, Available since v1.287.0) Rule description information
* `dest_cidr_ip` - (Optional, Available since v1.287.0) Destination IP address segment for outbound authorization
Example value: 0.0.0.0/0
* `direction` - (Optional, Available since v1.287.0) Authorized direction
Example value: ingress
* `ip_protocol` - (Optional, Computed, Available since v1.287.0) IP protocol
Example value: TCP
* `ipv6_dest_cidr_ip` - (Optional, Available since v1.287.0) The target IPv6 address segment.
* `ipv6_source_cidr_ip` - (Optional, Available since v1.287.0) The source IPv6 address segment.
* `policy` - (Optional, Computed, Available since v1.287.0) Authorization Policy
Example value: Accept
* `port_range` - (Optional, Computed, Available since v1.287.0) Source end port range.
* `priority` - (Optional, Computed, Int, Available since v1.287.0) Rule Priority
Example value: 1
* `source_cidr_ip` - (Optional, Computed, Available since v1.287.0) Source IP address segment, used for inbound authorization
Example value: 0.0.0.0/0
* `source_port_range` - (Optional, Computed, Available since v1.287.0) The port range of the source security group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Group.
* `delete` - (Defaults to 5 mins) Used when delete the Security Group.
* `update` - (Defaults to 5 mins) Used when update the Security Group.

## Import

ENS Security Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_security_group.example <security_group_id>
```