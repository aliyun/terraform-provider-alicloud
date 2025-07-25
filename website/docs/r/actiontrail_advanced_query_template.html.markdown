---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_advanced_query_template"
description: |-
  Provides a Alicloud Actiontrail Advanced Query Template resource.
---

# alicloud_actiontrail_advanced_query_template

Provides a Actiontrail Advanced Query Template resource.

sql template of advanced query.

For information about Actiontrail Advanced Query Template and how to use it, see [What is Advanced Query Template](https://next.api.alibabacloud.com/document/Actiontrail/2020-07-06/CreateAdvancedQueryTemplate).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_actiontrail_advanced_query_template" "default" {
  simple_query  = true
  template_name = "exampleTemplateName"
  template_sql  = "*"
}
```

## Argument Reference

The following arguments are supported:
* `simple_query` - (Required) Distinguish whether the current template is a simple query
* `template_name` - (Optional) The name of the resource
* `template_sql` - (Required) SQL content saved on behalf of the current template

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Advanced Query Template.
* `delete` - (Defaults to 5 mins) Used when delete the Advanced Query Template.
* `update` - (Defaults to 5 mins) Used when update the Advanced Query Template.

## Import

Actiontrail Advanced Query Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_actiontrail_advanced_query_template.example <id>
```