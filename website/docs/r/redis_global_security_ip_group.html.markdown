---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_redis_global_security_ip_group"
description: |-
  Provides a Alicloud Tair (Redis OSS-Compatible) And Memcache (KVStore) Global Security Ip Group resource.
---

# alicloud_redis_global_security_ip_group

Provides a Tair (Redis OSS-Compatible) And Memcache (KVStore) Global Security Ip Group resource.



For information about Tair (Redis OSS-Compatible) And Memcache (KVStore) Global Security Ip Group and how to use it, see [What is Global Security Ip Group](https://next.api.alibabacloud.com/document/R-kvstore/2015-01-01/CreateGlobalSecurityIPGroup).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}


resource "alicloud_redis_global_security_ip_group" "default" {
  global_ip_group_name = "ggn_example_create"
  global_ip_list       = "192.168.0.1,10.10.10.10,172.16.0.1"
}
```

## Argument Reference

The following arguments are supported:
* `global_ip_list` - (Required) The IP address in the whitelist template.
* `global_ip_group_name` - (Required) The name of the IP whitelist template.

-> **NOTE:**   Multiple IP addresses are separated by commas (,). You can create up to 1,000 IP addresses or CIDR blocks for all IP whitelists.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Global Security Ip Group.
* `delete` - (Defaults to 5 mins) Used when delete the Global Security Ip Group.
* `update` - (Defaults to 5 mins) Used when update the Global Security Ip Group.

## Import

Tair (Redis OSS-Compatible) And Memcache (KVStore) Global Security Ip Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_redis_global_security_ip_group.example <global_security_group_id>
```