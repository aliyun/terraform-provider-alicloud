---
subcategory: "App Streaming"
layout: "alicloud"
page_title: "Alicloud: alicloud_app_streaming_apps"
sidebar_current: "docs-alicloud-datasource-app-streaming-apps"
description: |-
  Provides a list of App Streaming App available to the user.
---

# alicloud_app_streaming_apps

This data source provides the App Streaming App of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.289.0.

## Example Usage

```terraform
data "alicloud_app_streaming_apps" "default" {
}

output "app_streaming_app_ids" {
  value = data.alicloud_app_streaming_apps.default.ids
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of App IDs.
* `name_regex` - (Optional) A regex string to filter results by App name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of App IDs.
* `apps` - A list of App Entries. Each element contains the following attributes:
  * `id` - The ID of the App.
  * `app_id` - The ID of the App.
  * `app_name` - The name of the App.
  * `app_version` - The version number of the App.
  * `app_version_name` - The version name of the App.
  * `icon_url` - The URL of the App icon.
