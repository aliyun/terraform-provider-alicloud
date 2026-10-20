---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_password_policy"
description: |-
  Provides a Alicloud RAM Password Policy resource.
---

# alicloud_ram_password_policy

Provides a RAM Password Policy resource.

Password strength information.

For information about RAM Password Policy and how to use it, see [What is Password Policy](https://next.api.alibabacloud.com/document/Ram/2015-05-01/SetPasswordPolicy).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_password_policy&exampleId=6d42d04a-2407-784e-3a2c-5a080b15ac65c6902e96&activeTab=example&spm=docs.r.ram_password_policy.0.6d42d04a24&intl_lang=EN_US" target="_blank">
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


resource "alicloud_ram_password_policy" "default" {
  minimum_password_length              = "8"
  require_lowercase_characters         = false
  require_numbers                      = false
  max_password_age                     = "0"
  password_reuse_prevention            = "1"
  max_login_attemps                    = "1"
  hard_expiry                          = false
  require_uppercase_characters         = false
  require_symbols                      = false
  password_not_contain_user_name       = false
  minimum_password_different_character = "1"
}
```

### Deleting `alicloud_ram_password_policy` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ram_password_policy`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ram_password_policy&spm=docs.r.ram_password_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `hard_expiry` - (Optional) Specifies whether to disable logon after the password expires. Valid values:
  - true: After the password expires, you cannot use the password to log on to the console. You can log on to the console only after you reset the password by using your Alibaba Cloud account or as a RAM user that has administrative rights.
  - false: After the password expires, you can change the password to log on to the console. This is the default value.
* `initial_password_age` - (Optional, Int) The initial password age in days.
Valid values: 0 to 90. Default value: 14.
* `intercept_risk_password_on_api` - (Optional) Whether to intercept risk passwords on API calls. Valid values:
  - true
  - false (default)
* `max_login_attemps` - (Optional, Int) The maximum number of password retries. If you enter the wrong passwords for the specified consecutive times, the account is locked for one hour.
Valid values: 0 to 32.
The default value is 0, which indicates that the password retries are not limited.
* `max_password_age` - (Optional, Int) The validity period of the password.
Valid values: 0 to 1095. Unit: days.
The default value is 0, which indicates that the password never expires.
* `minimum_password_different_character` - (Optional, Int) The minimum number of unique characters in the password.
Valid values: 0 to 8.
The default value is 0, which indicates that no limits are imposed on the number of unique characters in a password.
* `minimum_password_length` - (Optional, Computed, Int) The minimum number of characters in the password.
Valid values: 8 to 32. Default value: 8.
* `password_not_contain_user_name` - (Optional) Specifies whether to exclude the username from the password. Valid values:
  - true: A password cannot contain the username.
  - false: A password can contain the username. This is the default value.
* `password_reuse_prevention` - (Optional, Int) The policy for password history check.
The previous N passwords cannot be reused. Valid values of N: 0 to 24.
The default value is 0, which indicates that RAM users can reuse previous passwords.
* `require_lowercase_characters` - (Optional) Specifies whether the password must contain lowercase letters. Valid values:
  - true
  - false (default)
* `require_numbers` - (Optional) Specifies whether the password must contain digits. Valid values:
  - true
  - false (default)
* `require_symbols` - (Optional) Specifies whether the password must contain special characters. Valid values:
  - true
  - false (default)
* `require_uppercase_characters` - (Optional) Specifies whether the password must contain uppercase letters. Valid values:
  - true
  - false (default)

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. This field is set to your Alibaba Cloud Account ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Password Policy.
* `update` - (Defaults to 5 mins) Used when update the Password Policy.

## Import

RAM Password Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_password_policy.example <id>.
```