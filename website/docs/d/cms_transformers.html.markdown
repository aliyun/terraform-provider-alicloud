---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_transformers"
description: |-
  Provides a list of Cms Transformer owned by an Alibaba Cloud account.
---

# alicloud_cms_transformers

This data source provides Cms Transformer available to the user.[What is Transformer](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateTransformer)

-> **NOTE:** Available since v1.289.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cms_transformer" "default" {
  transformer_name = var.name
  description      = "Transformer managed by Terraform"
  workspace        = "default-workspace-cn-hangzhou"
}

data "alicloud_cms_transformers" "default" {
  ids              = ["${alicloud_cms_transformer.default.id}"]
  transformer_name = var.name
  workspace        = "default-workspace-cn-hangzhou"
}

output "alicloud_cms_transformer_example_id" {
  value = data.alicloud_cms_transformers.default.transformers.0.id
}
```

## Argument Reference

The following arguments are supported:
* `transformer_name` - (Optional) Filters results by the transformer name.
* `transformer_id` - (Optional) Filters results by the transformer ID.
* `enable` - (Optional) Filters results by whether the transformer is enabled.
* `workspace` - (Required) The workspace ID, which is used to isolate transformer resources for different business workspaces.
* `ids` - (Optional, Computed) A list of Transformer IDs. The value is formulated as `<transformer_id>:<workspace>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Transformer IDs.
* `transformers` - A list of Transformer Entries. Each element contains the following attributes:
  * `id` - The ID of the resource supplied above. The value is formulated as `<transformer_id>:<workspace>`.
  * `transformer_id` - The unique identifier of the transformer.
  * `transformer_name` - The name of the transformer.
  * `workspace` - The workspace ID.
  * `description` - The description of the transformer.
  * `quit_after_match` - Specifies whether to stop processing after the first match.
  * `sort_id` - The sort order of the transformer.
  * `enable` - Indicates whether the transformer is enabled.
  * `create_time` - The creation time.
  * `update_time` - The update time.
  * `user_id` - The user ID.
  * `region_id` - The region ID.
