---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_saved_query"
description: |-
  Provides a Alicloud Resource Manager Saved Query resource.
---

# alicloud_resource_manager_saved_query

Provides a Resource Manager Saved Query resource. ResourceCenter Saved Query.

For information about Resource Manager Saved Query and how to use it, see [What is Saved Query](https://www.alibabacloud.com/help/zh/resource-management/developer-reference/api-resourcecenter-2022-12-01-createsavedquery).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_saved_query&exampleId=9c8f80b6-fdab-bc34-92a1-d53411b726b65485f90c&activeTab=example&spm=docs.r.resource_manager_saved_query.0.9c8f80b6fd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_resource_manager_saved_query" "default" {
  description      = var.name
  expression       = "select * from resources limit 1;"
  saved_query_name = var.name

}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Query Description.
* `expression` - (Required) Query Expression.
* `saved_query_name` - (Required) The name of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Saved Query.
* `delete` - (Defaults to 5 mins) Used when delete the Saved Query.
* `update` - (Defaults to 5 mins) Used when update the Saved Query.

## Import

Resource Manager Saved Query can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_saved_query.example <id>
```