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

-> **NOTE:** Available since v2.2.0.

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
* `all` - (Optional, Bool, Available since v2.2.0) The flag that indicates whether to query all records.
* `ids` - (Optional, Computed, List, Available since v2.2.0) A list of Dedicated Ip Pool IDs.
* `keyword` - (Optional, String, Available since v2.2.0) The keyword used to search by name.
* `pool_id` - (Optional, String, Available since v2.2.0) The ID of the IP pool.
* `output_file` - (Optional, String, Available since v2.2.0) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Computed, List, Available since v2.2.0) A list of Dedicated Ip Pool IDs.
* `pools` - (Computed, List, Available since v2.2.0) A list of Dedicated Ip Pool Entries. Each element contains the following attributes:
  * `create_time` - (Computed, String, Available since v2.2.0) The creation time.
  * `id` - (Computed, String, Available since v2.2.0) The ID of the IP pool.
  * `ip_count` - (Computed, Int, Available since v2.2.0) The number of source IP addresses.
  * `name` - (Computed, String, Available since v2.2.0) The name of the IP pool.
  * `ips` - (Computed, List, Available since v2.2.0) The list of IP addresses in the pool:
    * `id` - (Computed, String, Available since v2.2.0) The ID of the purchased instance.
    * `ip` - (Computed, String, Available since v2.2.0) The IP address.
    * `zone_id` - (Computed, String, Available since v2.2.0) The zone ID.
