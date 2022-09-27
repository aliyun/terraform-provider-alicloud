---
subcategory: "Cloud Architect Design Tools"
layout: "alicloud"
page_title: "Alicloud: alicloud_bp_studio_application"
sidebar_current: "docs-alicloud-resource-bp-studio-application"
description: |-
  Provides a Alicloud Cloud Architect Design Tools Application resource.
---

# alicloud\_bp\_studio\_application

Provides a Cloud Architect Design Tools Application resource.

For information about Cloud Architect Design Tools Application and how to use it, see [What is Application](https://help.aliyun.com/document_detail/428263.html).

-> **NOTE:** Available in v1.192.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bp_studio_application" "default" {
  application_name  = "example_value"
  template_id       = "example_value"
  resource_group_id = "example_value"
  area_id           = "example_value"
  instances {
    id        = "example_value"
    node_name = "example_value"
    node_type = "ecs"
  }
  configuration = {
    enableMonitor = "1"
  }
  variables = {
    test = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, ForceNew) The name of the application.
* `template_id` - (Required, ForceNew) The id of the template.
* `resource_group_id` - (Optional, ForceNew, Computed) The id of the resource group.
* `area_id` - (Optional, ForceNew) The id of the area.
* `instances` - (Optional, ForceNew) The instance list. Support the creation of instances in the existing vpc under the application. See the following `Block instances`.
* `configuration` - (Optional, ForceNew) The configuration of the application.
* `variables` - (Optional, ForceNew) The variables of the application.

#### Block instances

The instances supports the following:

* `id` - (Optional, ForceNew) The id of the instance.
* `node_name` - (Optional, ForceNew) The name of the instance.
* `node_type` - (Optional, ForceNew) The type of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Application.
* `status` - The status of the Application.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 120 mins) Used when create the Application.
* `delete` - (Defaults to 120 mins) Used when delete the Application.

## Import

Cloud Architect Design Tools Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_bp_studio_application.example <id>
```
