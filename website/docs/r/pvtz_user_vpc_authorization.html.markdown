---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_user_vpc_authorization"
sidebar_current: "docs-alicloud-resource-pvtz-user-vpc-authorization"
description: |-
  Provides a Alicloud Private Zone User Vpc Authorization resource.
---

# alicloud_pvtz_user_vpc_authorization

Provides a Private Zone User Vpc Authorization resource.

-> **NOTE:** Available since v1.138.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_user_vpc_authorization&exampleId=6944948e-7ede-d082-3de8-c2510d1e784483cb5417&activeTab=example&spm=docs.r.pvtz_user_vpc_authorization.0.6944948e7e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "authorized_user_id" {
  default = 123456789
}
resource "alicloud_pvtz_user_vpc_authorization" "example" {
  authorized_user_id = var.authorized_user_id
  auth_channel       = "RESOURCE_DIRECTORY"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_pvtz_user_vpc_authorization&spm=docs.r.pvtz_user_vpc_authorization.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `auth_channel` - (Optional) The auth channel. Valid values: `RESOURCE_DIRECTORY`.
* `authorized_user_id` - (Required, ForceNew) The primary account ID of the user who authorizes the resource.
* `auth_type` - (Optional, ForceNew) The type of Authorization. Valid values: `NORMAL` and `CLOUD_PRODUCT`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User Vpc Authorization. The value formats as `<authorized_user_id>:<auth_type>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the User Vpc Authorization.
* `delete` - (Defaults to 2 mins) Used when delete the User Vpc Authorization.

## Import

Private Zone User Vpc Authorization can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_user_vpc_authorization.example <authorized_user_id>:<auth_type>
```
