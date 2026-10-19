---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_extensions"
sidebar_current: "docs-alicloud-datasource-data-works-extensions"
description: |-
  Provides a list of Data Works Extensions to the user.
---

# alicloud_data_works_extensions

This data source provides the Data Works Extensions of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.231.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_extensions" "default" {
  output_file = "extensions.json"
}

output "data_works_extension_id" {
  value = data.alicloud_data_works_extensions.default.extensions.0.id
}
```

Filter by extension code

```terraform
data "alicloud_data_works_extensions" "example" {
  extension_code = "example-extension-code"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Extension IDs. The ID is the same as `extension_code`.
* `extension_code` - (Optional, ForceNew) The unique code of the extension. If specified, the data source queries the details of this single extension via `GetExtension`; otherwise, it lists all extensions via `ListExtensions`.
* `name_regex` - (Optional, ForceNew) A regex string to filter extensions by `extension_name`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Extension IDs.
* `extensions` - A list of Data Works Extensions. Each element contains the following attributes:
  * `id` - The ID of the extension, same as `extension_code`.
  * `extension_code` - The unique code of the extension.
  * `extension_name` - The name of the extension.
  * `extension_desc` - The description of the extension.
  * `status` - The status of the extension.
  * `help_doc_url` - The help document URL of the extension.
  * `project_testing` - Whether the extension is in project testing.
  * `detail_url` - The detail URL of the extension.
  * `parameter_setting` - The parameter setting of the extension, serialized as a JSON string.
  * `option_setting` - The option setting of the extension, serialized as a JSON string.
  * `bind_event_list` - The list of bound events, serialized as a JSON string.
  * `event_category_list` - The list of event categories, serialized as a JSON string.
  * `region_id` - The region ID of the extension.
