---
subcategory: "Cloud Control"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_control_products"
sidebar_current: "docs-alicloud-datasource-cloud-control-products"
description: |-
  Provides a list of Cloud Control Product owned by an Alibaba Cloud account.
---

# alicloud_cloud_control_products

This data source provides Cloud Control Product available to the user.[What is Product](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.241.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_cloud_control_products" "default" {
  ids = ["VPC"]
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Product IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Product IDs.
* `names` - A list of name of Products.
* `products` - A list of Product Entries. Each element contains the following attributes:
  * `product_code` - The first ID of the resource
  * `product_name` - The name of the resource
  * `id` - The ID of the resource supplied above.
