---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_mail_addresses"
sidebar_current: "docs-alicloud-datasource-direct-mail-mail-addresses"
description: |-
  Provides a list of Direct Mail Mail Addresses to the user.
---

# alicloud\_direct\_mail\_mail\_addresses

This data source provides the Direct Mail Mail Addresses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_direct_mail_mail_addresses" "ids" {
  ids = ["example_id"]
}
output "direct_mail_mail_address_id_1" {
  value = data.alicloud_direct_mail_mail_addresses.ids.addresses.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Mail Address IDs.
* `key_word` - (Optional, ForceNew) The key word about account email address.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `sendtype` - (Optional, ForceNew) Account type. Valid values: `batch`, `trigger`.
* `status` - (Optional, ForceNew) Account Status. Valid values: `0`, `1`. Freeze: 1, normal: 0.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `addresses` - A list of Direct Mail Mail Addresses. Each element contains the following attributes:
	* `id` - The ID of the Mail Address.
	* `account_name` - The sender address.
	* `create_time` - The creation of the record time.
	* `daily_count` - On the quota limit.
	* `daily_req_count` - On the quota.
	* `domain_status` - Domain name status. Valid values: `0`, `1`.
	* `mail_address_id` - The sender address ID.
	* `month_count` - Monthly quota limit.
	* `month_req_count` - Months amount.
	* `reply_address` - Return address.
	* `reply_status` - If using STMP address status.
	* `sendtype` - Account type.
	* `status` - Account Status. Valid values: `0`, `1`. Freeze: 1, normal: 0.
