---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_tenant_role_user_attachment"
description: |-
  Provides a Alicloud Max Compute Tenant Role User Attachment resource.
---

# alicloud_max_compute_tenant_role_user_attachment

Provides a Max Compute Tenant Role User Attachment resource.

Binding relationship between tenant roles and users

-> **WARNING:** Using TenantRoleUserAttachment restricts the use of the [MaxCompute Console > Tenant Management > Tenant Attributes > Use Account ID] feature. If you use the [Use Account ID] feature in the console, TenantRoleUserAttachment becomes unavailable and requires approximately one hour to recover.

For information about Max Compute Tenant Role User Attachment and how to use it, see [What is Tenant Role User Attachment](https://next.api.alibabacloud.com/document/MaxCompute/2022-01-04/UpdateTenantUserRoles).

-> **NOTE:** Available since v1.270.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_max_compute_tenant_role_user_attachment&exampleId=145d8862-76ff-5a2f-f507-39f2f80120cd49b6fc27&activeTab=example&spm=docs.r.max_compute_tenant_role_user_attachment.0.145d886276&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_max_compute_tenant_role_user_attachment&spm=docs.r.max_compute_tenant_role_user_attachment.example&intl_lang=EN_US)

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