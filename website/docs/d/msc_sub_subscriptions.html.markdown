---
subcategory: "Message Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_subscriptions"
sidebar_current: "docs-alicloud-datasource-msc-sub-subscriptions"
description: |- 
    Provides a list of Message Center Subscriptions to the user.
---

# alicloud\_msc\_sub\_subscriptions

This data source provides the Message Center Subscriptions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_msc_sub_subscriptions" "default" {}
output "msc_sub_subscription_id_1" {
  value = data.alicloud_msc_sub_subscriptions.default.subscriptions.0.id
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `subscriptions` - A list of Msc Sub Subscriptions. Each element contains the following attributes:
    * `channel` - The channel the Subscription.
    * `contact_ids` - The ids of subscribed contacts.
    * `description` - The description of the Subscription.
    * `email_status` - The status of email subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
    * `id` - The ID of the Subscription.
    * `item_id` - The ID of the Subscription.
    * `item_name` - The name of the Subscription.
    * `pmsg_status` - The status of pmsg subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
    * `sms_status` - The status of sms subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
    * `tts_status` - The status of tts subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
    * `webhook_ids` - The ids of subscribed webhooks.
    * `webhook_status` - The status of webhook subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
