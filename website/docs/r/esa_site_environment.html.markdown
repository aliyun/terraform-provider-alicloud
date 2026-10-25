---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site_environment"
description: |-
  Provides a Alicloud ESA Site Environment resource.
---

# alicloud_esa_site_environment

Provides a ESA Site Environment resource.



For information about ESA Site Environment and how to use it, see [What is Site Environment](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateEnvironment).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}

resource "alicloud_esa_site_environment" "default" {
  site_id               = data.alicloud_esa_sites.default.sites.0.id
  environment_name      = var.name
  rule                  = "(ip.src in {1.1.1.1})"
  next_environment_name = "Production"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The site ID.
* `environment_name` - (Required, ForceNew) The name of the environment.
* `rule` - (Required) The environment rule.
* `next_environment_name` - (Required, ForceNew) The name of the environment with the next priority.
* `read_only` - (Optional, Bool) Specifies whether the environment is read-only. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Site Environment. It formats as `<site_id>:<environment_name>`.
* `create_time` - The time when the environment was created.
* `update_time` - The time when the environment was updated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Site Environment.
* `update` - (Defaults to 5 mins) Used when update the Site Environment.
* `delete` - (Defaults to 5 mins) Used when delete the Site Environment.

## Import

ESA Site Environment can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site_environment.example <site_id>:<environment_name>
```
