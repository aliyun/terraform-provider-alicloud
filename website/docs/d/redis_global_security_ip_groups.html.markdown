---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_redis_global_security_ip_groups"
sidebar_current: "docs-alicloud-datasource-redis-global-security-ip-groups"
description: |-
  Provides a list of Tair (Redis OSS-Compatible) And Memcache (KVStore) Global Security Ip Group owned by an Alibaba Cloud account.
---

# alicloud_redis_global_security_ip_groups

This data source provides Tair (Redis OSS-Compatible) And Memcache (KVStore) Global Security Ip Group available to the user.[What is Global Security Ip Group](https://next.api.alibabacloud.com/document/R-kvstore/2015-01-01/CreateGlobalSecurityIPGroup)

-> **NOTE:** Available since v1.287.0.

## Example Usage

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

data "alicloud_redis_global_security_ip_groups" "default" {
  ids = ["${alicloud_redis_global_security_ip_group.default.id}"]
}

output "alicloud_redis_global_security_ip_group_example_id" {
  value = data.alicloud_redis_global_security_ip_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `engine` - (ForceNew, Optional) The engine type of the resource, such as Redis or Tair.
* `global_security_group_id` - (ForceNew, Optional) The ID of the IP whitelist template.
* `ids` - (Optional, Computed) A list of Global Security Ip Group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Global Security Ip Group IDs.
* `groups` - A list of Global Security Ip Group Entries. Each element contains the following attributes:
  * `global_ip_list` - The IP address in the whitelist template.
  * `global_ip_group_name` - The name of the IP whitelist template.
  * `global_security_group_id` - The ID of the IP whitelist template.
  * `region_id` - The region ID of the resource.
  * `id` - The ID of the resource supplied above.
