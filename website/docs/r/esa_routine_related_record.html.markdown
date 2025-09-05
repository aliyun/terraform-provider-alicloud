---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine_related_record"
description: |-
  Provides a Alicloud ESA Routine Related Record resource.
---

# alicloud_esa_routine_related_record

Provides a ESA Routine Related Record resource.



For information about ESA Routine Related Record and how to use it, see [What is Routine Related Record](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutineRelatedRecord).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_esa_routine" "default" {
  description = "example-routine2"
  name        = "example-routine2"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "example.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_routine_related_record" "default" {
  record_name = "example.com"
  site_id     = alicloud_esa_site.default.id
  name        = alicloud_esa_routine.default.id
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required, ForceNew) The routine name.
* `record_name` - (Required, ForceNew) The record name.
* `site_id` - (Required, ForceNew, Int) The website ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<name>:<record_id>`.
* `record_id` - The record ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Routine Related Record.
* `delete` - (Defaults to 5 mins) Used when delete the Routine Related Record.

## Import

ESA Routine Related Record can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_routine_related_record.example <name>:<record_id>
```