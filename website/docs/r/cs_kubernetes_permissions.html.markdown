---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_permissions"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-permissions"
description: |-
  Provides a Alicloud resource to grant RBAC permissions for ACK cluster.
---

# alicloud\_cs\_kubernetes\_permissions

This resource will help you implement RBAC authorization for the kubernetes cluster. 

-> **NOTE:** Please make sure that the target RAM user has been granted a RAM policy with at least read-only permission of the target cluster in the RAM console. Otherwise, the `ErrorRamPolicyConfig` error will be returned. 
For more information about how to authorize a RAM user by attaching RAM policies, see [Create a custom RAM policy](https://www.alibabacloud.com/help/doc-detail/86485.htm).

-> **NOTE:** If you call this operation as a RAM user, make sure that this RAM user has the permissions to grant other RAM users the permissions to manage ACK clusters. Otherwise, the `StatusForbidden` or `ForbiddenGrantPermissions` errors will be returned. For more information, see [Use a RAM user to grant RBAC permissions to other RAM users](https://www.alibabacloud.com/help/faq-detail/119035.htm).

-> **NOTE:** This operation overwrites the permissions that have been granted to the specified RAM user. When you call this operation, make sure that the required permissions are included.

-> **NOTE:** Available in v1.122.0+.

## Example Usage
### Grant RBAC permissions
If you don't have users and clusters, to perform RBAC authorization, you need to complete the following steps. 

Step 1, create a cluster using Terraform.
```terraform
variable "name" {
  default = "custom-name"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vswitch_name      = var.name
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
}

# Create a managed cluster
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [alicloud_vswitch.default.id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
}
```
Step 2: Use Teraform to create a RAM user and authorize access to the cluster created in step 1 (at least read-only permission is required). 
```terraform
# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name         = var.name
  display_name = var.name
  mobile       = "86-18688888888"    # replace to your tel
  email        = "hello.uuu@aaa.com" # replace to your email
  comments     = "yoyoyo"
  force        = true
}

# Create a new RAM Policy.
resource "alicloud_ram_policy" "policy" {
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
          "acs:cs:*:*:cluster/${alicloud_cs_managed_kubernetes.default.0.id}"
        ]
      }
    ],
    "Version": "1"
  }
  EOF
  description     = "this is a policy test by tf"
  force           = true
}

# Authorize the RAM user
resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  user_name   = alicloud_ram_user.user.name
}
```
Finally, complete the RBAC authorization.
```terraform
# Grant users developer permissions for the cluster.
resource "alicloud_cs_kubernetes_permissions" "default" {
  # uid
  uid = alicloud_ram_user.user.id
  # permissions
  permissions {
    cluster     = alicloud_cs_managed_kubernetes.default.0.id
    role_type   = "cluster"
    role_name   = "dev"
    namespace   = ""
    is_custom   = false
    is_ram_role = false
  }
  # If you want to grant users multiple cluster permissions, you can define multiple sets of permissions 
  #  permissions {
  #    cluster     = "cluster_id_2"
  #    role_type   = "cluster"
  #    role_name   = "ops"
  #    namespace   =  ""
  #    is_custom   = false
  #    is_ram_role = false
  #  }
  depends_on = [
    alicloud_ram_user_policy_attachment.attach
  ]
}
```
If you already have users and clusters, to complete RBAC authorization, you only need to run the following code to use Terraform. 
```terraform
# Get RAM user ID 
data "alicloud_ram_users" "users_ds" {
  name_regex = "your ram user name"
}

# Create a new RAM Policy.
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
          "acs:cs:*:*:cluster/${target_cluster_ID}"
        ]
      }
    ],
    "Version": "1"
  }
  EOF
  description     = "this is a policy test by tf"
  force           = true
}

# Authorize the RAM user
resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  user_name   = data.alicloud_ram_users.users_ds.users.0.name
}

# RBAC authorization for the cluster
resource "alicloud_cs_kubernetes_permissions" "default" {
  uid = data.alicloud_ram_users.users_ds.users.0.id
  permissions {
    cluster     = "target cluster id1"
    role_type   = "cluster"
    role_name   = "ops"
    is_custom   = false
    is_ram_role = false
    namespace   = ""
  }
  permissions {
    cluster     = "target cluster id2"
    role_type   = "cluster"
    role_name   = "ops"
    is_custom   = false
    is_ram_role = false
    namespace   = ""
  }
  depends_on = [
    alicloud_ram_user_policy_attachment.attach
  ]
}
```
### Remove user permissions
Remove the current user's permissions on the cluster
```terraform
# remove the permissions on the "cluster_id_01", "cluster_id_02".
resource "alicloud_cs_kubernetes_permissions" "default" {
  uid = data.alicloud_ram_users.users_ds.users.0.id
}
```

## Argument Reference

The following arguments are supported.

* `uid` - (Required, ForceNew) The ID of the Ram user, and it can also be the id of the Ram Role. If you use Ram Role id, you need to set `is_ram_role` to `true` during authorization.
* `permissions` - (Optional) A list of user permission.
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

