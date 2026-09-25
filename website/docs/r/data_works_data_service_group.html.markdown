---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_service_group"
description: |-
  Provides a Alicloud Data Works Data Service Group resource.
---

# alicloud_data_works_data_service_group

Provides a Data Works Data Service Group resource.

Data Service Business Process.

For information about Data Works Data Service Group and how to use it, see [What is Data Service Group](https://next.api.alibabacloud.com/document/dataworks-public/2020-05-18/CreateDataServiceGroup).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "project_id" {
  description = "The ID of the DataWorks workspace"
  type        = number
}

variable "api_gateway_group_id" {
  description = "The ID of the API gateway group bound to the business process"
  type        = string
}

resource "alicloud_data_works_data_service_group" "default" {
  project_id              = var.project_id
  api_gateway_group_id    = var.api_gateway_group_id
  data_service_group_name = var.name
  description             = "example business process"
}
```

### Deleting `alicloud_data_works_data_service_group` or removing it from your configuration

Terraform cannot destroy resource `alicloud_data_works_data_service_group`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `api_gateway_group_id` - (Required, ForceNew) The ID of the API gateway group to which the business process belongs.
* `data_service_group_name` - (Required, ForceNew) The name of the business process.
* `description` - (Optional, ForceNew) The description of the business process.
* `project_id` - (Required, ForceNew, Int) The ID of the DataWorks workspace to which the business process belongs.
* `tenant_id` - (Optional, ForceNew, Int, Deprecated) The tenant ID. This field is deprecated.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<project_id>:<data_service_group_id>`.
* `create_time` - The time when the business process was created.
* `creator_id` - The ID of the user who created the business process.
* `data_service_group_id` - The ID of the business process.
* `modified_time` - The time when the business process was last modified.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Data Service Group.
* `delete` - (Defaults to 5 mins) Used when delete the Data Service Group.

## Import

Data Works Data Service Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_data_service_group.example <project_id>:<data_service_group_id>
```
