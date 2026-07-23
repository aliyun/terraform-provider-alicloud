---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine"
description: |-
  Provides a Alicloud ESA Routine resource.
---

# alicloud_esa_routine

Provides a ESA Routine resource.

An ESA Routine is an edge function. Besides managing the routine itself, this resource can
optionally upload the routine code from a local file: on create or when the local file
content changes, the code is uploaded to the staging environment and committed as a new
immutable code version.

For information about ESA Routine and how to use it, see [What is Routine](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutine).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_routine" "default" {
  description = var.name
  name        = var.name
}
```

Upload code from a local file

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
```

## Argument Reference

The following arguments are supported:
* `name` - (Required, ForceNew) Routine Name. It must be unique in the same account.
* `description` - (Optional, ForceNew) The description of the routine.
* `filename` - (Optional) The path to the local routine code file. When set, the file is uploaded to the staging environment and committed as a new code version on create; changing the file content triggers a new upload and commit on the next apply.
* `code_description` - (Optional) The description recorded for the committed code version.

-> **NOTE:** `name` and `description` cannot be modified in place. Changing either of them forces a new routine to be created.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is formatted as the routine `name`.
* `create_time` - The time when the routine was created.
* `code_checksum` - The SHA-256 checksum (base64 encoded) of the local code file. It is used to detect code content changes.
* `latest_code_version` - The latest code version committed from the local code file.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Routine.
* `update` - (Defaults to 5 mins) Used when update the Routine code.
* `delete` - (Defaults to 5 mins) Used when delete the Routine.

## Import

ESA Routine can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_routine.example <id>
```

-> **NOTE:** `filename`, `code_description`, `code_checksum` and `latest_code_version` are not restored on import because the routine code is uploaded from a local file that is not readable from the server.
