---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_login_profile"
description: |-
  Provides a Alicloud RAM Login Profile resource.
---

# alicloud_ram_login_profile

Provides a RAM Login Profile resource.



For information about RAM Login Profile and how to use it, see [What is Login Profile](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createloginprofile).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_login_profile&exampleId=442ec83a-6fe5-90e2-17a2-351fa3abb6a2738d91e8&activeTab=example&spm=docs.r.ram_login_profile.0.442ec83a6f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ram_user" "user" {
  name         = "terraform_example"
  display_name = "terraform_example"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "terraform_example"
  force        = true
}

resource "alicloud_ram_login_profile" "profile" {
  user_name = alicloud_ram_user.user.name
  password  = "Example_1234"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ram_login_profile&spm=docs.r.ram_login_profile.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `mfa_bind_required` - (Optional, Computed) Specifies whether to forcefully enable multi-factor authentication (MFA) for the RAM user. Valid values:
  - true: forcefully enables MFA for the RAM user. The RAM user must bind an MFA device upon the next logon.
  - false (default): does not forcefully enable MFA for the RAM user.
* `password` - (Required, Sensitive) The password must meet the Password strength requirements. For more information about password strength setting requirements, see [GetPasswordPolicy](https://help.aliyun.com/document_detail/2337691.html).
* `password_reset_required` - (Optional) Whether the user must reset the password at the next logon. Value:
  - true
  - false (default)
* `user_name` - (Required, ForceNew) The user name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Login Profile.
* `delete` - (Defaults to 5 mins) Used when delete the Login Profile.
* `update` - (Defaults to 5 mins) Used when update the Login Profile.

## Import

RAM Login Profile can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_login_profile.example <id>
```