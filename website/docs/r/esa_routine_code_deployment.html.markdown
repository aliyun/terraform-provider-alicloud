---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine_code_deployment"
description: |-
  Provides a Alicloud ESA Routine Code Deployment resource.
---

# alicloud_esa_routine_code_deployment

Provides a ESA Routine Code Deployment resource.

A code deployment releases one or two committed code versions of an ESA Routine to a target
environment using a percentage strategy, supporting both a single-version full release and a
two-version canary release.

For information about ESA Routine Code Deployment and how to use it, see [What is Routine Code Deployment](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutineCodeDeployment).

-> **NOTE:** A code deployment is an immutable release record. Changing any argument creates a new deployment. Because ESA does not provide an API to delete a deployment record, destroying this resource only removes it from Terraform state; deleting the parent `alicloud_esa_routine` removes its deployment records on the server side.

-> **NOTE:** Available since v1.251.0.

## Example Usage

Full release of a single code version

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

resource "alicloud_esa_routine_code_deployment" "default" {
  routine_name = alicloud_esa_routine.default.name
  env          = "staging"
  strategy     = "percentage"
  code_versions {
    code_version = alicloud_esa_routine.default.latest_code_version
    percentage   = 100
  }
}
```

## Argument Reference

The following arguments are supported:
* `routine_name` - (Required, ForceNew) The name of the routine.
* `env` - (Required, ForceNew) The target environment. Valid values: `staging`, `production`.
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `percentage`. Default value: `percentage`.
* `code_versions` - (Required, ForceNew) The list of code version percentage configurations. At most two entries are supported, and the sum of the percentages must equal 100. See [`code_versions`](#code_versions) below.

### `code_versions`

The `code_versions` supports the following:
* `code_version` - (Required, ForceNew) The code version number to deploy.
* `percentage` - (Required, ForceNew) The percentage of traffic routed to this code version. Valid values: 1 to 100.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is formatted as `<routine_name>:<env>`.
* `deployment_id` - The ID of the deployment record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Routine Code Deployment.
* `delete` - (Defaults to 5 mins) Used when delete the Routine Code Deployment.

## Import

ESA Routine Code Deployment can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_routine_code_deployment.example <routine_name>:<env>
```
