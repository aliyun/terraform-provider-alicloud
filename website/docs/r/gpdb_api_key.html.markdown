---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_api_key"
description: |-
  Provides a Alicloud AnalyticDB for PostgreSQL (GPDB) Api Key resource.
---

# alicloud_gpdb_api_key

Provides a AnalyticDB for PostgreSQL (GPDB) Api Key resource.

The API key under a GPDB SaaS workspace.

For information about AnalyticDB for PostgreSQL (GPDB) Api Key and how to use it, see [What is Api Key](https://next.api.alibabacloud.com/document/gpdb/2016-05-03/CreateApiKey).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "workspace_id" {
  default = "ws-xxxxxxx"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_gpdb_api_key" "default" {
  workspace_id = var.workspace_id
  key_name     = var.name
  description  = "terraform example"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the API key.
* `key_name` - (Required, ForceNew) The name of the API key.
* `service_ids` - (Optional, ForceNew) The list of SaaS service IDs that the API key is authorized to access.
* `workspace_id` - (Required, ForceNew) The ID of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<workspace_id>:<key_id>`.
* `key_id` - The ID of the API key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Api Key.
* `delete` - (Defaults to 5 mins) Used when delete the Api Key.

## Import

AnalyticDB for PostgreSQL (GPDB) Api Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_api_key.example <workspace_id>:<key_id>
```
