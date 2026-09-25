---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_service_extension"
sidebar_current: "docs-alicloud-resource-ga-service-extension"
description: |-
  Provides a Global Accelerator (GA) Service Extension resource.
---

# alicloud_ga_service_extension

Provides a Global Accelerator (GA) Service Extension resource.

For information about Global Accelerator Service Extension and how to use it, see [What is Service Extension](https://www.alibabacloud.com/help/en/global-accelerator/developer-reference/api-ga-2019-11-20-createserviceextension).

-> **NOTE:** Available since v1.235.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_service_extension" "default" {
  name        = "example-value"
  description = "example description"
  tags = {
    Created = "tfTest"
    For     = "example"
  }
}
```

Associate resources and add service components

```terraform
resource "alicloud_ga_service_extension" "default" {
  name        = "example-value"
  description = "example description"
  resources {
    accelerator_id = alicloud_ga_accelerator.default.id
    resource_type  = "accelerator"
    resource_id    = alicloud_ga_accelerator.default.id
  }
  components {
    service_component_id = "sc-xxx"
    priority             = 1
    timeout              = 60
    fail_policy          = "stop"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service extension. The name must be `2` to `128` characters in length, and can contain letters, digits, hyphens (-) and underscores (_). It must start with a letter.
* `description` - (Optional) The description of the service extension. The description must be `1` to `256` characters in length.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `components` - (Optional) The list of service components added to the service extension. See [`components`](#components) below.
* `resources` - (Optional) The list of resources associated with the service extension. See [`resources`](#resources) below.

### `components`

The components block supports the following:

* `service_component_id` - (Required) The ID of the service component.
* `priority` - (Optional) The priority of the service component.
* `timeout` - (Optional) The timeout of the service component. Unit: seconds.
* `config` - (Optional) The configuration of the service component.
* `fail_policy` - (Optional) The fail policy of the service component.

### `resources`

The resources block supports the following:

* `accelerator_id` - (Required) The ID of the accelerator to associate.
* `resource_type` - (Required) The type of the resource to associate.
* `resource_id` - (Required) The ID of the resource to associate.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service extension, which equals to `service_extension_id`.
* `service_extension_id` - The ID of the service extension.
* `type` - The type of the service extension.
* `state` - The state of the service extension.
* `create_time` - The time when the service extension was created.
* `update_time` - The time when the service extension was last updated.
* `components.0.component_name` - The name of the service component.
* `resources.0.associate_id` - The ID of the association record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#custom-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the service extension.
* `update` - (Defaults to 5 mins) Used when updating the service extension.
* `delete` - (Defaults to 5 mins) Used when deleting the service extension.

## Import

Global Accelerator Service Extension can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_service_extension.example <service_extension_id>
```
