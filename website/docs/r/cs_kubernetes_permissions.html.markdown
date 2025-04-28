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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_kubernetes_permissions&exampleId=0f035c3d-0e4f-431c-6b61-5d1ff75a19923cce09ce&activeTab=example&spm=docs.r.cs_kubernetes_permissions.0.0f035c3d0e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "terraform-example"
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type        = list(string)
  default     = ["10.1.0.0/16", "10.2.0.0/16"]
}

variable "pod_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or service's and cannot be in them."
  default     = "172.16.0.0/16"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "192.168.0.0/16"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "ManagedKubernetes"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "default" {
  count      = length(var.vswitch_cidrs)
  vpc_id     = alicloud_vpc.vpc.id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

# Create a new RAM cluster.
resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = "${var.name}-${random_integer.default.result}"
  cluster_spec         = "ack.pro.small"
  version              = data.alicloud_cs_kubernetes_version.default.metadata.0.version
  worker_vswitch_ids   = split(",", join(",", alicloud_vswitch.default.*.id))
  new_nat_gateway      = false
  pod_cidr             = var.pod_cidr
  service_cidr         = var.service_cidr
  slb_internet_enabled = false
}

# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name = "${var.name}-${random_integer.default.result}"
}

# Create a cluster permission for user.
resource "alicloud_cs_kubernetes_permissions" "default" {
  uid = alicloud_ram_user.user.id
  permissions {
    cluster     = alicloud_cs_managed_kubernetes.default.id
    role_type   = "cluster"
    role_name   = "admin"
    namespace   = ""
    is_custom   = false
    is_ram_role = false
  }
}

resource "alicloud_cs_kubernetes_permissions" "attach" {
  uid = alicloud_ram_user.user.id
  permissions {
    cluster     = alicloud_cs_managed_kubernetes.default.id
    role_type   = "namespace"
    role_name   = "cs:dev"
    namespace   = "default"
    is_custom   = true
    is_ram_role = false
  }
}
```

## Argument Reference

The following arguments are supported.

* `uid` - (Required, ForceNew) The ID of the Ram user, and it can also be the id of the Ram Role. If you use Ram Role id, you need to set `is_ram_role` to `true` during authorization.
* `permissions` - (Optional) A list of user permission. See [`permissions`](#permissions) below.

### `permissions`

The permissions mapping supports the following:

* `cluster` - (Required) The ID of the cluster that you want to manage, When `role_type` value is `all-clusters`, the value of `cluster` must be `""`.
* `role_name` - (Required) Specifies the predefined role that you want to assign. Valid values `admin`, `ops`, `dev`, `restricted` and the custom cluster roles.
* `role_type` - (Required) The authorization type. Valid values `cluster`, `namespace` and `all-clusters`.
* `namespace` - (Optional) The namespace to which the permissions are scoped. This parameter is required only if you set role_type to namespace.
* `is_ram_role` - (Optional) Specifies whether the permissions are granted to a RAM role. When `uid` is ram role id, the value of `is_ram_role` must be `true`.
* `is_custom` - (Optional) Specifies whether to perform a custom authorization. To perform a custom authorization, the value of `is_custom` must be `true`, and set `role_name` to a custom cluster role.

## Attributes Reference

The following attributes are exported:
* `id` - Resource id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status).
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster.

