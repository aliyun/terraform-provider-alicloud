---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_tenant_role_user_attachment"
description: |-
  Provides a Alicloud Max Compute Tenant Role User Attachment resource.
---

# alicloud_max_compute_tenant_role_user_attachment

Provides a Max Compute Tenant Role User Attachment resource.



For information about Max Compute Tenant Role User Attachment and how to use it, see [What is Tenant Role User Attachment](https://next.api.alibabacloud.com/document/MaxCompute/2022-01-04/UpdateTenantUserRoles).

-> **NOTE:** Available since v1.270.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_max_compute_tenant_role_user_attachment" "default0" {
  account_id  = "p4_200053869413670560"
  tenant_role = "admin"
}
```

## Argument Reference

The following arguments are supported:
* `account_id` - (Optional, ForceNew, Computed) Account UID

1. If the user is a primary account, the AccountId format is UID.  
   Example: 200231703336555555

2. If the user is a RAM user, the AccountId format is p4_UID.  
   Example: p4_200531704446555555

3. If the user is a RAM role, the AccountId format is v4_UID.  
   Example: v4_300007628597555555

* `tenant_role` - (Optional, ForceNew, Computed) Tenant role. By default, admin and super_administrator are available. You can add more roles in the console.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<account_id>:<tenant_role>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Tenant Role User Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Tenant Role User Attachment.

## Import

Max Compute Tenant Role User Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_tenant_role_user_attachment.example <account_id>:<tenant_role>
```