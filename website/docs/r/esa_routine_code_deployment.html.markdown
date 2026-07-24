---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine_code_deployment"
description: |-
  Provides a Alicloud ESA Routine Code Deployment resource.
---

# alicloud_esa_routine_code_deployment

Provides a ESA Routine Code Deployment resource. It deploys one or two committed code versions of a routine to an environment, optionally splitting traffic between them for a canary (percentage) rollout.

For information about ESA Routine Code Deployment and how to use it, see [What is Routine Code Deployment](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutineCodeDeployment).

-> **NOTE:** Available since v1.287.0.

-> **NOTE:** A code deployment is a rollout event bound to an environment. The Alibaba Cloud ESA API provides no operation to withdraw a deployment; a deployed version can only be replaced by a new deployment. Destroying this resource only removes it from the Terraform state and does not roll back the deployed code.

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
  name             = var.name
  code             = "addEventListener('fetch', e => e.respondWith(new Response('hello')))"
  code_description = "version 1"
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
* `routine_name` - (Required, ForceNew) The name of the routine to deploy.
* `env` - (Required, ForceNew) The target environment. Valid values: `staging`, `production`.
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `percentage`. Default value: `percentage`.
* `code_versions` - (Required, ForceNew) The list of code versions and their traffic percentages. At most two versions are supported, and the percentages must sum to 100. See [`code_versions`](#code_versions) below.

### `code_versions`

The `code_versions` supports the following:
* `code_version` - (Required, ForceNew) The committed code version to deploy.
* `percentage` - (Required, ForceNew) The traffic percentage of this code version. Valid values: 1 to 100.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. It is formatted as `<routine_name>:<env>`.
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
