---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_permissions"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-permissions"
description: |-
  Provides a Alicloud resource to grant RBAC permissions for ACK cluster.
---

# alicloud_cs_kubernetes_permissions

This resource will help you implement RBAC authorization for the kubernetes cluster, see [What is kubernetes permissions](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-grantpermissions).

-> **NOTE:** Please make sure that the target RAM user has been granted a RAM policy with at least read-only permission of the target cluster in the RAM console. Otherwise, the `ErrorRamPolicyConfig` error will be returned. 
For more information about how to authorize a RAM user by attaching RAM policies, see [Create a custom RAM policy](https://www.alibabacloud.com/help/doc-detail/86485.htm).

-> **NOTE:** If you call this operation as a RAM user, make sure that this RAM user has the permissions to grant other RAM users the permissions to manage ACK clusters. Otherwise, the `StatusForbidden` or `ForbiddenGrantPermissions` errors will be returned. For more information, see [Use a RAM user to grant RBAC permissions to other RAM users](https://www.alibabacloud.com/help/faq-detail/119035.htm).

-> **NOTE:** This operation **overwrites** the permissions that have been granted to the specified RAM user. When you call this operation, make sure that the required permissions are included.

-> **NOTE:** Available since v1.122.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_ram_user" "default" {
  name         = var.name
  display_name = var.name
  mobile       = "86-18688888888"    # replace to your tel
  email        = "hello.uuu@aaa.com" # replace to your email
  comments     = "yoyoyo"
  force        = true
}

resource "alicloud_ram_policy" "default" {
  policy_name     = var.name
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "cs:Get*",
          "cs:List*",
          "cs:Describe*"
        ],
        "Effect": "Allow",
        "Resource": [
          "acs:cs:*:*:cluster/${alicloud_cs_managed_kubernetes.default.id}"
        ]
      }
    ],
    "Version": "1"
  }
  EOF
  description     = "this is a policy test by tf"
  force           = true
}

resource "alicloud_ram_user_policy_attachment" "default" {
  policy_name = alicloud_ram_policy.default.name
  policy_type = alicloud_ram_policy.default.type
  user_name   = alicloud_ram_user.default.name
}

resource "alicloud_cs_kubernetes_permissions" "default" {
  uid = alicloud_ram_user.default.id
  permissions {
    cluster     = alicloud_cs_managed_kubernetes.default.id
    role_type   = "namespace"
    role_name   = "dev"
    is_custom   = false
    is_ram_role = false
    namespace   = "kube-system"
  }
  permissions {
    cluster     = alicloud_cs_managed_kubernetes.default.id
    role_type   = "namespace"
    role_name   = "dev"
    is_custom   = false
    is_ram_role = false
    namespace   = "default"
  }
}

# If you already have users and clusters, to complete RBAC authorization, you only need to run the following code to use Terraform.
locals {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  user_id    = alicloud_ram_user.default.id
  user_name  = alicloud_ram_user.default.name
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "AckClusterReadOnlyAccess"
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "cs:Get*",
          "cs:List*",
          "cs:Describe*"
        ],
        "Effect": "Allow",
        "Resource": [
          "acs:cs:*:*:cluster/${local.cluster_id}"
        ]
      }
    ],
    "Version": "1"
  }
  EOF
  description     = "this is a policy test by tf"
  force           = true
}

resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  user_name   = local.user_name
}

# You can import resource by uid and add all required permissions to resource firstly.
# Make sure that the required permissions are included because this resource will overwrite the permissions that have been granted to the specified RAM user.
resource "alicloud_cs_kubernetes_permissions" "already_attach" {
  uid = local.user_id
  # Define all required permissions in one resource block for one user using list permissions, otherwise they will overwrite each other.
  permissions {
    cluster     = local.cluster_id
    role_type   = "namespace"
    role_name   = "dev"
    is_custom   = false
    is_ram_role = false
    namespace   = "kube-system"
  }
  permissions {
    cluster     = local.cluster_id
    role_type   = "namespace"
    role_name   = "dev"
    is_custom   = false
    is_ram_role = false
    namespace   = "default"
  }
}
```
## Argument Reference

The following arguments are supported.

* `uid` - (Required, ForceNew) The ID of the Ram user, and it can also be the id of the Ram Role. If you use Ram Role id, you need to set `is_ram_role` to `true` during authorization.
* `permissions` - (Optional) A list of user permission. See [`permissions`](#permissions) below.

### `permissions`

The permissions mapping supports the following:

* `cluster` - (Required) The ID of the cluster that you want to manage.
* `role_name` - (Required) Specifies the predefined role that you want to assign. Valid values `admin`, `ops`, `dev`, `restricted` and the custom cluster roles.
* `role_type` - (Required) The authorization type. Valid values `cluster`, `namespace`.
* `namespace` - (Optional) The namespace to which the permissions are scoped. This parameter is required only if you set role_type to namespace.
* `is_ram_role` - (Optional) Specifies whether the permissions are granted to a RAM role. When `uid` is ram role id, the value of `is_ram_role` must be `true`.
* `is_custom` - (Optional) Specifies whether to perform a custom authorization. To perform a custom authorization, set `role_name` to a custom cluster role.

## Attributes Reference

The following attributes are exported:
* `id` - Resource id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status).
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster.

## Import

alicloud_cs_kubernetes_permissions can be imported using the RAM user id or Ram Role id, e.g.

```shell
$ terraform import alicloud_cs_kubernetes_permissions.user <uid>
```
