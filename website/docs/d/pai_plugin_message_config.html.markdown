---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_message_config"
description: |-
  This data source provides the PAIPlugin SMS message config of the current account.
---

# alicloud_pai_plugin_message_config

This data source provides the PAIPlugin SMS message config of the current account
(singleton).

For information about PAIPlugin message config and how to use it, see
[What is GetMessageConfig](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-getmessageconfig).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
data "alicloud_pai_plugin_message_config" "default" {}

output "sms_report_url" {
  value = data.alicloud_pai_plugin_message_config.default.sms_report_url
}
```

## Argument Reference

> **NOTE:** This data source has no arguments; it always reflects the current account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `sms_report_url` - The URL that receives the SMS delivery report callback.
* `sms_up_url` - The URL that receives the uplink SMS message callback.
