---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_dedicated_ip_pool"
description: |-
  Provides a Alicloud Direct Mail Dedicated Ip Pool resource.
---

# alicloud_direct_mail_dedicated_ip_pool

Provides a Direct Mail Dedicated Ip Pool resource.

A dedicated IP pool groups purchased dedicated IP addresses. Sender addresses can be associated with an IP pool through a configuration set, so that emails are sent from the dedicated IPs in the pool.

For information about Direct Mail Dedicated Ip Pool and how to use it, see [What is Dedicated Ip Pool](https://next.api.alibabacloud.com/document/Dm/2015-11-23/DedicatedIpPoolCreate).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_direct_mail_dedicated_ip_pool" "default" {
  name = var.name
}
```

Create an IP pool from purchased dedicated IP instances:

```terraform
variable "name" {
  default = "terraform-example"
}

# Query the IDs of purchased dedicated IP instances that are not assigned to any pool.
data "alicloud_direct_mail_dedicated_ip_pools" "default" {}

resource "alicloud_direct_mail_dedicated_ip_pool" "with_ips" {
  name             = var.name
  buy_resource_ids = "your-purchased-ip-instance-id-1,your-purchased-ip-instance-id-2"
}
```

## Argument Reference

The following arguments are supported:
* `buy_resource_ids` - (Optional, ForceNew) The IDs of the purchased IP instances. Separate multiple IDs with commas (,). You can obtain the instance IDs from the response of the DedicatedIpNonePoolList operation.
* `name` - (Required, ForceNew) The name of the IP pool. The name must be 1 to 50 characters in length. It can contain letters, digits, underscores (_), and hyphens (-). The name cannot be changed after the IP pool is created.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time.
* `ip_count` - The number of source IP addresses.
* `ips` - The list of IP addresses in the pool.
  * `id` - The ID of the purchased instance.
  * `ip` - The IP address.
  * `zone_id` - The zone ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dedicated Ip Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Dedicated Ip Pool.

## Import

Direct Mail Dedicated Ip Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_direct_mail_dedicated_ip_pool.example <id>
```
