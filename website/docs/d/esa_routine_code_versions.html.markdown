---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine_code_versions"
description: |-
  Provides a list of ESA Routine Code Versions to the user.
---

# alicloud_esa_routine_code_versions

This data source provides the ESA Routine Code Versions of the current Alibaba Cloud user.

For information about ESA Routine Code Versions and how to use it, see [ListRoutineCodeVersions](https://next.api.alibabacloud.com/document/ESA/2024-09-10/ListRoutineCodeVersions).

-> **NOTE:** Available since v1.287.0.

-> **NOTE:** The underlying `ListRoutineCodeVersions` API caps pagination at 2 pages of up to 20 items each, so at most the 40 most recent code versions are returned.

## Example Usage

```terraform
data "alicloud_esa_routine_code_versions" "default" {
  routine_name = "terraform-example"
}

output "code_versions" {
  value = data.alicloud_esa_routine_code_versions.default.versions
}
```

## Argument Reference

The following arguments are supported:
* `routine_name` - (Required) The name of the routine to query code versions for.
* `search_key_word` - (Optional) A keyword used to fuzzy-match code versions.
* `ids` - (Optional) A list of code version IDs (code version numbers) used to filter the result.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `versions` - A list of Routine Code Versions. Each element contains the following attributes:
  * `id` - The ID of the code version. Same as `code_version`.
  * `code_version` - The code version number.
  * `code_description` - The description of the code version.
  * `create_time` - The time when the code version was created, in RFC 3339 (UTC) format.
  * `status` - The status of the code version.
  * `deploy_env` - The environment bundled with the code version. Valid values: `staging`, `production`.
  * `build_id` - The build ID of the code version.
  * `has_env_vars` - Whether the code version bundles environment variables.
