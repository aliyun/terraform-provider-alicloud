---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_global_security_ip_group"
description: |-
  Provides a Alicloud Polardb Global Security Ip Group resource.
---

# alicloud_polardb_global_security_ip_group

Provides a Polardb Global Security Ip Group resource.

Global Security IP Group.

For information about Polardb Global Security Ip Group and how to use it, see [What is Global Security Ip Group](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateGlobalSecurityIPGroup).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_polardb_global_security_ip_group" "default" {
  global_ip_list       = "192.168.0.1"
  global_ip_group_name = "example_template"
}
```

## Argument Reference

The following arguments are supported:
* `global_ip_group_name` - (Required) The name of the IP whitelist template. The name of the IP whitelist template must meet the following requirements:

  - The name can contain lowercase letters, digits, and underscores (\_).
  - The name must start with a letter and end with a letter or digit.
  - The name must be 2 to 120 characters in length.
* `global_ip_list` - (Required) The IP address in the whitelist template.

-> **NOTE:**   Multiple IP addresses are separated by commas (,). You can create up to 1,000 IP addresses or CIDR blocks for all IP whitelists.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID of the IP whitelist template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Global Security Ip Group.
* `delete` - (Defaults to 5 mins) Used when delete the Global Security Ip Group.
* `update` - (Defaults to 5 mins) Used when update the Global Security Ip Group.

## Import

Polardb Global Security Ip Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_global_security_ip_group.example <id>
```