---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_meta_categories"
sidebar_current: "docs-alicloud-datasource-data-works-meta-categories"
description: |-
  Provides a list of Data Works Meta Categories to the user.
---

# alicloud_data_works_meta_categories

This data source provides the Data Works Meta Categories of the current Alibaba Cloud user by parent id.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_data_works_meta_categories" "default" {
  parent_category_id = 0
  ids                = ["${alicloud_data_works_meta_category.example.category_id}"]
}

output "first_category_id" {
  value = data.alicloud_data_works_meta_categories.default.categories.0.category_id
}
```

## Argument Reference

The following arguments are supported:

* `parent_category_id` - (Optional) The id of the parent meta category. Defaults to `0` which means root-level categories. The data source lists the children of this parent.
* `ids` - (Optional) A list of category ids used to filter the returned categories. Only categories whose `category_id` matches one of the given ids are returned.
* `output_file` - (Optional) File name where to write the result after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of matching meta category ids in the format `<category_id>:<parent_category_id>`.
* `categories` - A list of meta categories. Each element contains the following attributes:
  * `id` - The resource id of the meta category, in the format `<category_id>:<parent_category_id>`.
  * `category_id` - The id of the meta category.
  * `parent_category_id` - The parent meta category id.
  * `name` - The name of the meta category.
  * `comment` - The comment of the meta category.
  * `create_time` - The creation timestamp of the meta category.
