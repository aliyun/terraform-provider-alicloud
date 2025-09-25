---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_role_policy_document"
description: |-
  Generates a MaxCompute role policy document to use with MaxCompute role.
---

# alicloud_max_compute_role_policy_document

This data source Generates a MaxCompute Role policy document to use in a MaxCompute role.

-> **NOTE:** Available since v1.261.0.

## Example Usage

### Basic Example

```terraform
data "alicloud_max_compute_role_policy_document" "default" {
  version = "1"
  statement {
    effect   = "Allow"
    action   = ["odps:*"]
    resource = ["acs:odps:*:projects/my_project/schemas/*"]
  }
}

resource "alicloud_maxcompute_project" "default" {
  project_name  = "my_project"
  product_type  = "PayAsYouGo"
  default_quota = "os_PayAsYouGoQuota"
}

resource "alicloud_max_compute_role" "dpu_dbt_role" {
  type         = "resource"
  project_name = alicloud_maxcompute_project.default.project_name
  role_name    = "my_mc_role"
  policy       = data.alicloud_max_compute_role_policy_document.default.document
}
```


## Argument Reference

The following arguments are supported:

* `version` - (Optional) Version of the RAM policy document. Valid value is `1`. Default value is `1`.
* `statement` - (Optional) Statement of the RAM policy document. See the following `Block statement`. See [`statement`](#statement) below.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `statement`

The statement supports the following:

* `effect` - (Optional) This parameter indicates whether or not the `action` is allowed. Valid values are `Allow` and `Deny`. Default value is `Allow`.
* `action` - (Required) Action of the MaxCompute role policy document.
* `resource` - (Required) List of specific objects which will be authorized.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `document` - Standard policy document rendered based on the arguments above.
