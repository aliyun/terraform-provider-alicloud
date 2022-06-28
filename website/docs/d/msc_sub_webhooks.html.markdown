---
subcategory: "Message Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_webhooks"
sidebar_current: "docs-alicloud-datasource-msc-sub-webhooks"
description: |-
  Provides a list of Msc Sub Webhooks to the user.
---

# alicloud\_msc\_sub\_webhooks

This data source provides the Msc Sub Webhooks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.141.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_msc_sub_webhooks" "ids" {
  ids = ["example_id"]
}
output "msc_sub_webhook_id_1" {
  value = data.alicloud_msc_sub_webhooks.ids.webhooks.0.id
}

data "alicloud_msc_sub_webhooks" "nameRegex" {
  name_regex = "^my-Webhook"
}
output "msc_sub_webhook_id_2" {
  value = data.alicloud_msc_sub_webhooks.nameRegex.webhooks.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Webhook IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Webhook name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Webhook names.
* `webhooks` - A list of Msc Sub Webhooks. Each element contains the following attributes:
	* `id` - The ID of the Webhook.
	* `server_url` - The serverUrl of the Subscription.
	* `webhook_id` - The first ID of the resource.
	* `webhook_name` - The name of the Webhook. **Note:** The name must be `2` to `12` characters in length, and can contain uppercase and lowercase letters.
