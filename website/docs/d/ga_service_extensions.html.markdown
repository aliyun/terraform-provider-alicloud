---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_service_extensions"
sidebar_current: "docs-alicloud-datasource-ga-service-extensions"
description: |-
  Provides a list of Global Accelerator (GA) Service Extensions to the user.
---

# alicloud_ga_service_extensions

This data source provides the list of Global Accelerator (GA) Service Extensions and their basic information.

-> **NOTE:** Available since v1.235.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_service_extensions" "default" {
  output_file = "extensions.txt"
}

output "first_extension_id" {
  value = data.alicloud_ga_service_extensions.default.service_extensions.0.id
}
```

Filter by name and ids

```terraform
data "alicloud_ga_service_extensions" "default" {
  ids            = ["se-xxx"]
  name_regex     = "^tf-test.*"
  enable_details = true
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of service extension IDs.
* `name_regex` - (Optional) A regex string to filter results by service extension name.
* `service_extension_name` - (Optional) The name of the service extension.
* `status` - (Optional) The status of the service extension.
* `output_file` - (Optional) File name where to save the data source results.
* `enable_details` - (Optional, Available since v1.235.0+) Default to false. Set it to true to get the details of each service extension, including components, resources and tags.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of service extension IDs.
* `names` - A list of service extension names.
* `service_extensions` - A list of service extensions. Each element contains the following attributes:
  * `id` - The ID of the service extension, which equals to `service_extension_id`.
  * `service_extension_id` - The ID of the service extension.
  * `name` - The name of the service extension.
  * `description` - The description of the service extension.
  * `type` - The type of the service extension.
  * `state` - The state of the service extension.
  * `resource_group_id` - The ID of the resource group.
  * `tags` - The tags of the service extension.
  * `create_time` - The time when the service extension was created.
  * `update_time` - The time when the service extension was last updated.
* `components.0.service_component_id` - The ID of the service component.
* `components.0.priority` - The priority of the service component.
* `components.0.timeout` - The timeout of the service component.
* `components.0.config` - The configuration of the service component.
* `components.0.fail_policy` - The fail policy of the service component.
* `components.0.component_name` - The name of the service component.
* `resources.0.accelerator_id` - The ID of the accelerator.
* `resources.0.resource_type` - The type of the resource.
* `resources.0.resource_id` - The ID of the resource.
* `resources.0.associate_id` - The ID of the association record.
