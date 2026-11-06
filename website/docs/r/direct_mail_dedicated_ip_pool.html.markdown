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

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_direct_mail_dedicated_ip_pool&exampleId=72bfb390-8270-6462-1519-f54492ee10acbbd57023&activeTab=example&spm=docs.r.direct_mail_dedicated_ip_pool.0.72bfb39082&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_direct_mail_dedicated_ip_pool" "default" {
  name = var.name
}
```

Create an IP pool from purchased dedicated IP instances:

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_direct_mail_dedicated_ip_pool&exampleId=c580fef4-67bc-a565-d2fa-b58de8534e79ffe3becf&activeTab=example&spm=docs.r.direct_mail_dedicated_ip_pool.1.c580fef467&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

# Query the IDs of purchased dedicated IP instances that are not assigned to
# any pool via the DedicatedIpNonePoolList operation.
resource "alicloud_direct_mail_dedicated_ip_pool" "with_ips" {
  name             = var.name
  buy_resource_ids = "your-purchased-ip-instance-id-1,your-purchased-ip-instance-id-2"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_direct_mail_dedicated_ip_pool&spm=docs.r.direct_mail_dedicated_ip_pool.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `buy_resource_ids` - (Optional, String, Available since v2.2.0) The IDs of the purchased IP instances. Separate multiple IDs with commas (,). You can obtain the instance IDs from the response of the DedicatedIpNonePoolList operation. Member order is insignificant: the provider compares configurations and state by the trimmed, sorted ID set, so listing the IDs in a different order does not trigger an update. The configured IDs are the complete member set: an update replaces the pool members with the configured IDs, and removing `buy_resource_ids` removes all members. If members are changed outside Terraform, the next refresh records the actual (possibly empty) membership and the next apply restores the configured members.
* `name` - (Required, ForceNew, String, Available since v2.2.0) The name of the IP pool. The name must be 1 to 50 characters in length. It can contain letters, digits, underscores (_), and hyphens (-). The name cannot be changed after the IP pool is created; changing `name` replaces the pool with a new one that has a different ID.

## Attributes Reference

The following attributes are exported:
* `id` - (Computed, String, Available since v2.2.0) The ID of the resource supplied above.
* `create_time` - (Computed, String, Available since v2.2.0) The creation time.
* `ip_count` - (Computed, Int, Available since v2.2.0) The number of source IP addresses.
* `ips` - (Computed, List, Available since v2.2.0) The list of IP addresses in the pool.
  * `id` - (Computed, String, Available since v2.2.0) The ID of the purchased instance.
  * `ip` - (Computed, String, Available since v2.2.0) The IP address.
  * `zone_id` - (Computed, String, Available since v2.2.0) The zone ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dedicated Ip Pool.
* `update` - (Defaults to 5 mins) Used when update the Dedicated Ip Pool.
* `delete` - (Defaults to 5 mins) Used when delete the Dedicated Ip Pool.

## Import

Direct Mail Dedicated Ip Pool can be imported using the id, e.g.

```shell
$ terraform import alicloud_direct_mail_dedicated_ip_pool.example <id>
```
