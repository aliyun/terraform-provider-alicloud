---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_global_security_ip_group"
description: |-
  Provides a Alicloud Mongodb Global Security IP Group resource.
---

# alicloud_mongodb_global_security_ip_group

Provides a Mongodb Global Security IP Group resource.

Whitelist Template Resources.

For information about Mongodb Global Security IP Group and how to use it, see [What is Global Security IP Group](https://next.api.alibabacloud.com/document/Dds/2015-12-01/CreateGlobalSecurityIPGroup).

-> **NOTE:** Available since v1.257.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraformexample"
}

resource "alicloud_mongodb_global_security_ip_group" "default" {
  global_ig_name          = var.name
  global_security_ip_list = "192.168.1.1,192.168.1.2,192.168.1.3"
}
```

## Argument Reference

The following arguments are supported:
* `global_security_ip_list` - (Required) The IP address in the whitelist template.

-> **NOTE:** Separate multiple IP addresses with commas (,). You can create up to 1000 IP addresses or CIDR blocks for all IP address whitelists.

* `global_ig_name` - (Required) The name of the IP whitelist template.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Global Security IP Group.
* `delete` - (Defaults to 5 mins) Used when delete the Global Security IP Group.
* `update` - (Defaults to 5 mins) Used when update the Global Security IP Group.

## Import

Mongodb Global Security IP Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_global_security_ip_group.example <id>
```
