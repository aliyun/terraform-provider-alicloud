---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_commit_container_task"
description: |-
  Provides a Alicloud ECI Commit Container Task resource.
---

# alicloud_eci_commit_container_task

Provides a ECI Commit Container Task resource.

Build a mirror in the specified container and push it to the Mirror repository.

For information about ECI Commit Container Task and how to use it, see [What is Commit Container Task](https://next.api.alibabacloud.com/document/Eci/2018-08-08/CommitContainer).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_eci_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_eci_zones.default.zones.0.zone_ids.0
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
}

resource "alicloud_eci_container_group" "default" {
  container_group_name = var.name
  security_group_id    = alicloud_security_group.default.id
  vswitch_id           = alicloud_vswitch.default.id

  containers {
    name  = var.name
    image = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/nginx:alpine"
  }
}

# CommitContainer authorizes the ECI service to push to your ACR repository
# through this RAM role.
resource "alicloud_ram_role" "default" {
  role_name                   = var.name
  force                       = true
  assume_role_policy_document = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["eci.aliyuncs.com"]
      }
    }
  ],
  "Version": "1"
}
EOF
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.role_name
  policy_type = "System"
  policy_name = "AliyunContainerRegistryFullAccess"
}

resource "alicloud_eci_commit_container_task" "example" {
  container_group_id = alicloud_eci_container_group.default.id
  container_name     = var.name

  arn {
    role_type = "service"
    role_arn  = alicloud_ram_role.default.arn
  }

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

* `acr_registry_info` - (Optional, ForceNew, List) Access credential configuration information for ACR Enterprise Edition instance. See [`acr_registry_info`](#acr_registry_info) below.
* `arn` - (Optional, ForceNew, List) ARN information required for authorization. See [`arn`](#arn) below.

  -> **NOTE:** Although `arn` is marked optional by the OpenAPI definition, the live `CommitContainer` API returns `MissingParameter` when it is absent. Prepare the target ACR repository and a RAM role trusted by the ECI service before creating this resource.
* `container_group_id` - (Required, ForceNew, String) The ID of the ECI container group that contains the source container.
* `container_name` - (Required, ForceNew, String) The name of the source container in the container group.
* `image` - (Required, ForceNew, List) Container image information, including repository, tag, message and author. See [`image`](#image) below.

### `acr_registry_info`

The acr_registry_info supports the following:

* `arn_service` - (Optional, ForceNew, String) The RAM role ARN of the user (authorized party) in cross-account authorization scenarios.
* `arn_user` - (Optional, ForceNew, String) The RAM role ARN of the authorizer in cross-account authorization scenarios.
* `instance_id` - (Optional, ForceNew, String) ACR Enterprise Edition instance ID.
* `region_id` - (Optional, ForceNew, String) The region of the ACR Enterprise Edition instance.

### `arn`

The arn supports the following:

* `role_arn` - (Optional, ForceNew, String) The ARN of the authorization role.
* `role_type` - (Optional, ForceNew, String) Authorization type. Valid values: `service`.

### `image`

The image supports the following:

* `author` - (Optional, ForceNew, String) Image author information.
* `message` - (Optional, ForceNew, String) Image description message.
* `repository` - (Required, ForceNew, String) ACR image repository address.
* `tag` - (Required, ForceNew, String) Image tag. It is required; an empty string is accepted by the API and means no tag modification.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource supplied above.
* `create_time` - (String) The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Commit Container Task.

## Import

ECI Commit Container Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_eci_commit_container_task.example <commit_container_task_id>
```
