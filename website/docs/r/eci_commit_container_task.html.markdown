---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_commit_container_task"
description: |-
  Provides a Alicloud ECI Commit Container Task resource.
---

# alicloud_eci_commit_container_task

Provides a ECI Commit Container Task resource.

Build a mirror in the specified container and push it to the Mirror repository.

For information about ECI Commit Container Task and how to use it, see [What is Commit Container Task](https://next.api.alibabacloud.com/document/Eci/2018-08-08/CommitContainer).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_eci_container_group" "default" {
  container_group_name = "terraform-example"
  security_group_id    = alicloud_security_group.default.id
  vswitch_id           = alicloud_vswitch.default.id

  containers {
    name  = "terraform-example"
    image = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/nginx:alpine"
  }
}

resource "alicloud_eci_commit_container_task" "example" {
  container_group_id = alicloud_eci_container_group.default.id
  container_name     = "terraform-example"

  image {
    repository = "registry.cn-hangzhou.aliyuncs.com/namespace/repository"
    tag        = "v1.0.0"
    author     = "terraform"
    message    = "built by terraform"
  }
}
```

### Deleting `alicloud_eci_commit_container_task` or removing it from your configuration

CommitContainer is an asynchronous task API with no delete operation: the task record is managed by the service and retained for query after the task reaches a terminal state. Destroying `alicloud_eci_commit_container_task` removes the resource from Terraform state without deleting the task record (no delete API is invoked). If the task is no longer visible via DescribeCommitContainerTask, Read automatically removes the resource from state.

## Argument Reference

The following arguments are supported:

* `acr_registry_info` - (Optional, ForceNew, Set) Access credential configuration information for ACR Enterprise Edition instance. See [`acr_registry_info`](#acr_registry_info) below.
* `arn` - (Optional, ForceNew, Set) ARN information required for authorization. See [`arn`](#arn) below.
* `container_group_id` - (Required, ForceNew) The ID of the ECI container group that contains the source container.
* `container_name` - (Required, ForceNew) The name of the source container in the container group.
* `image` - (Required, ForceNew, Set) Container image information, including repository, tag, message and author. See [`image`](#image) below.

### `acr_registry_info`

The acr_registry_info supports the following:

* `arn_service` - (Optional, ForceNew) The RAM role ARN of the user (authorized party) in cross-account authorization scenarios.
* `arn_user` - (Optional, ForceNew) The RAM role ARN of the authorizer in cross-account authorization scenarios.
* `instance_id` - (Optional, ForceNew) ACR Enterprise Edition instance ID.
* `region_id` - (Optional, ForceNew) The region of the ACR Enterprise Edition instance.

### `arn`

The arn supports the following:

* `role_arn` - (Optional, ForceNew) The ARN of the authorization role.
* `role_type` - (Optional, ForceNew) Authorization type. Can be configured as service, indicating using RAM role for authorization.

### `image`

The image supports the following:

* `author` - (Optional, ForceNew) Image author information.
* `message` - (Optional, ForceNew) Image description message.
* `repository` - (Required, ForceNew) ACR image repository address.
* `tag` - (Required, ForceNew) Image tag. Default is empty, meaning no tag modification.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `page_truncated` - Page truncation status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Commit Container Task.

## Import

ECI Commit Container Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_eci_commit_container_task.example <commit_container_task_id>
```
