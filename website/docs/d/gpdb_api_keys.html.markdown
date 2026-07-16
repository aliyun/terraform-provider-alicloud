---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_api_keys"
sidebar_current: "docs-alicloud-datasource-gpdb-api-keys"
description: |-
  Provides a list of Gpdb Api Key owned by an Alibaba Cloud account.
---

# alicloud_gpdb_api_keys

This data source provides Gpdb Api Key available to the user.[What is Api Key](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateApiKey)

-> **NOTE:** Available since v1.286.0.

## Example Usage

```terraform
variable "workspace_id" {
  default = "ws-xxxxxxx"
}

data "alicloud_gpdb_api_keys" "default" {
  workspace_id = var.workspace_id
}

output "gpdb_api_key_id" {
  value = data.alicloud_gpdb_api_keys.default.keys[0].key_id
}
```

## Argument Reference

The following arguments are supported:
* `workspace_id` - (Required) The ID of the workspace.
* `ids` - (Optional, Computed) A list of Api Key IDs. The value is formulated as `<workspace_id>:<key_id>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Api Key IDs.
* `keys` - A list of Api Key Entries. Each element contains the following attributes:
    * `create_time` - The creation time of the resource.
    * `description` - The description of the API key.
    * `key_id` - The ID of the API key.
    * `key_name` - The name of the API key.
    * `key_prefix` - The prefix of the API key.
    * `id` - The ID of the resource supplied above.
