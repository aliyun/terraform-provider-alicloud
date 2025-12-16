---
subcategory: "Message Center (MscSub)"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_subscription"
sidebar_current: "docs-alicloud-resource-msc-sub-subscription"
description: |-
  Provides a Alicloud Msc Sub Subscription resource.
---

# alicloud_msc_sub_subscription

Provides a Msc Sub Subscription resource.

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_msc_sub_subscription&exampleId=71036df4-2b30-a554-1074-ce2145544083caf6e4b0&activeTab=example&spm=docs.r.msc_sub_subscription.0.71036df42b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_msc_sub_subscription" "example" {
  item_name      = "Notifications of Product Expiration"
  sms_status     = "1"
  email_status   = "1"
  pmsg_status    = "1"
  tts_status     = "1"
  webhook_status = "0"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_msc_sub_subscription&spm=docs.r.msc_sub_subscription.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `contact_ids` - (Optional) The ids of subscribed contacts.
  **NOTE:** There is a potential diff error because of the order of `contact_ids` values indefinite.
  So, from version 1.161.0, `contact_ids` type has been updated as `set` from `list`,
  and you can use [tolist](https://www.terraform.io/language/functions/tolist) to convert it to a list.
* `email_status` - (Optional) The status of email subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
* `item_name` - (Required, ForceNew) The name of the Subscription. **NOTE:**  You should use the `alicloud_msc_sub_subscriptions` to query the available subscription item name.
* `pmsg_status` - (Optional) The status of pmsg subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
* `sms_status` - (Optional) The status of sms subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
* `tts_status` - (Optional) The status of tts subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.
* `webhook_ids` - (Optional) The ids of subscribed webhooks.
* `webhook_status` - (Optional) The status of webhook subscription. Valid values: `-1`, `-2`, `0`, `1`. `-1` means required, `-2` means banned; `1` means subscribed; `0` means not subscribed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Subscription.
* `channel` - The channel the Subscription.
* `description` - The description of the Subscription.

## Import

Msc Sub Subscription can be imported using the id, e.g.

```shell
$ terraform import alicloud_msc_sub_subscription.example <id>
```
