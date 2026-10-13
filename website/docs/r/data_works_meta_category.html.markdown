---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_meta_category"
sidebar_current: "docs-alicloud-resource-data-works-meta-category"
description: |-
  Provides a Alicloud Data Works Meta Category resource.
---

# alicloud_data_works_meta_category

Provides a Data Works Meta Category resource. DataWorks Meta Category groups tables into a hierarchical taxonomy for governance and discovery.

For information about Data Works Meta Category and how to use it, see [What is Meta Category](https://www.alibabacloud.com/help/en/dataworks/user-guide/manage-categories).

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_data_works_meta_category" "example" {
  name               = "tf-meta-category-example"
  comment            = "example root category"
  parent_category_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the meta category. The name must be unique among siblings under the same parent.
* `comment` - (Optional) The comment of the meta category.
* `parent_category_id` - (Optional, ForceNew) The id of the parent meta category. Defaults to `0` which means a root-level category. The upstream API does not support changing the parent of an existing category, so any change to this field forces a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the meta category in the format `<category_id>:<parent_category_id>`.
* `category_id` - The unique identifier of the meta category assigned by the server.
* `create_time` - The creation timestamp of the meta category.

## Import

Data Works Meta Category can be imported using the id, e.g.

```shell
terraform import alicloud_data_works_meta_category.example 123456:0
```

-> **NOTE:** The id is `<category_id>:<parent_category_id>`. The parent segment is required because the upstream `GetMetaCategory` API lists categories by parent, so the provider needs the parent context to locate a single category.
