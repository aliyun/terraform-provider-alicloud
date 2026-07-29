---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine_code_versions"
description: |-
  Provides a list of ESA Routine Code Versions to the user.
---

# alicloud_esa_routine_code_versions

This data source provides the ESA Routine Code Versions of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.251.0.

-> **NOTE:** The underlying `ListRoutineCodeVersions` API supports only page 1 and page 2 with a maximum page size of 20, so this data source returns at most the most recent 40 code versions of the routine.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_routine" "default" {
  name             = var.name
  description      = var.name
  filename         = "${path.module}/index.js"
  code_description = "initial version"
}

data "alicloud_esa_routine_code_versions" "default" {
  name = alicloud_esa_routine.default.name
}

output "esa_routine_code_version_0" {
  value = data.alicloud_esa_routine_code_versions.default.versions.0.code_version
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the routine.
* `search_key_word` - (Optional) The keyword used to search code versions.
* `ids` - (Optional) A list of Routine Code Version IDs (each is a code version number).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Routine Code Version IDs (each is a code version number).
* `versions` - A list of Routine Code Versions. Each element contains the following attributes:
  * `id` - The ID of the code version. It is the same as `code_version`.
  * `code_version` - The code version number.
  * `code_description` - The description of the code version.
  * `create_time` - The time when the code version was created.
  * `status` - The status of the code version.
  * `deploy_env` - The environment the code version is deployed to. Valid values: `staging`, `production`.
  * `build_id` - The build ID of the code version.
  * `has_env_vars` - Whether the code version has environment variables.
