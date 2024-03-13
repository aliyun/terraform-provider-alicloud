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

-> **NOTE:** This operation overwrites the permissions that have been granted to the specified RAM user. When you call this operation, make sure that the required permissions are included.

-> **NOTE:** Available since v1.122.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

# leave it to empty would create a new one
variable "vpc_id" {
  description = "Existing vpc id used to create several vswitches and other resources."
  default     = ""
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

# leave it to empty then terraform will create several vswitches
variable "vswitch_ids" {
  description = "List of existing vswitch id."
  type        = list(string)
  default     = []
}


variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type        = list(string)
  default     = ["10.1.0.0/16", "10.2.0.0/16"]
}

# options: between 24-28
variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 24
}

# options: ipvs|iptables
variable "proxy_mode" {
  description = "Proxy mode is option of kube-proxy."
  default     = "ipvs"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "192.168.0.0/16"
}

variable "terway_vswitch_ids" {
  description = "List of existing vswitch ids for terway."
  type        = list(string)
  default     = []
}

variable "terway_vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'terway_vswitch_cidrs' is not specified."
  type        = list(string)
  default     = ["10.4.0.0/16", "10.5.0.0/16"]
}

variable "cluster_addons" {
  type = list(object({
    name   = string
    config = string
  }))

  default = [
    {
      "name"   = "terway-eniip",
      "config" = "",
    },
    {
      "name"   = "csi-plugin",
      "config" = "",
    },
    {
      "name"   = "csi-provisioner",
      "config" = "",
    },
    {
      "name"   = "logtail-ds",
      "config" = "{'IngressDashboardEnabled':'true'}",
    },
    {
      "name"   = "nginx-ingress-controller",
      "config" = "{'IngressSlbNetworkType':'internet'}",
    },
    {
      "name"   = "arms-prometheus",
      "config" = "",
    },
    {
      "name"   = "ack-node-problem-detector",
      "config" = "{'sls_project_name':''}",
    }
  ]
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

# If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count      = length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "terway_vswitches" {
  count      = length(var.terway_vswitch_ids) > 0 ? 0 : length(var.terway_vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.terway_vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name         = var.name
  cluster_spec = "ack.pro.small"
  # version can not be defined in variables.tf.
  version            = "1.26.3-aliyun.1"
  worker_vswitch_ids = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)) : length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  pod_vswitch_ids    = length(var.terway_vswitch_ids) > 0 ? split(",", join(",", var.terway_vswitch_ids)) : length(var.terway_vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
  new_nat_gateway    = true
  node_cidr_mask     = var.node_cidr_mask
  proxy_mode         = var.proxy_mode
  service_cidr       = var.service_cidr

  dynamic "addons" {
    for_each = var.cluster_addons
    content {
      name   = lookup(addons.value, "name", var.cluster_addons)
      config = lookup(addons.value, "config", var.cluster_addons)
    }
  }
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
    role_type   = "cluster"
    role_name   = "dev"
    namespace   = ""
    is_custom   = false
    is_ram_role = false
  }
}


#If you already have users and clusters, to complete RBAC authorization, you only need to run the following code to use Terraform. 
locals {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
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

resource "alicloud_cs_kubernetes_permissions" "already_attach" {
  uid = alicloud_ram_user.default.id
  permissions {
    cluster     = local.cluster_id
    role_type   = "cluster"
    role_name   = "ops"
    is_custom   = false
    is_ram_role = false
    namespace   = ""
  }
}

# Remove user permissions,Remove the current user's permissions on the cluster
resource "alicloud_cs_kubernetes_permissions" "remove_permissions" {
  uid = alicloud_cs_kubernetes_permissions.already_attach.uid
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

