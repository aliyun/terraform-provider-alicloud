---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_dedicated_ip_pools"
sidebar_current: "docs-alicloud-datasource-direct-mail-dedicated-ip-pools"
description: |-
  Provides a list of Direct Mail Dedicated Ip Pool owned by an Alibaba Cloud account.
---

# alicloud_direct_mail_dedicated_ip_pools

This data source provides Direct Mail Dedicated Ip Pool available to the user. [What is Dedicated Ip Pool](https://next.api.alibabacloud.com/document/Dm/2015-11-23/DedicatedIpPoolCreate)

-> **NOTE:** Available since v1.293.0.

## Example Usage

```terraform
data "alicloud_direct_mail_dedicated_ip_pools" "default" {}

output "first_pool_id" {
  value = data.alicloud_direct_mail_dedicated_ip_pools.default.pools.0.id
}

data "alicloud_direct_mail_dedicated_ip_pools" "by_name" {
  keyword = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:
* `all` - (Optional) The flag that indicates whether to query all records.
* `ids` - (Optional, Computed) A list of Dedicated Ip Pool IDs.
* `keyword` - (Optional) The keyword used to search by name.
* `pool_id` - (Optional) The ID of the IP pool.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Dedicated Ip Pool IDs.
* `pools` - A list of Dedicated Ip Pool Entries. Each element contains the following attributes:
  * `create_time` - The creation time.
  * `id` - The ID of the IP pool.
  * `ip_count` - The number of source IP addresses.
  * `name` - The name of the IP pool.
  * `ips` - The list of IP addresses in the pool:
    * `id` - The ID of the purchased instance.
    * `ip` - The IP address.
    * `zone_id` - The zone ID.
