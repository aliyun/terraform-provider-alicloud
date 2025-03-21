---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_user_policy_attachment"
description: |-
  Provides a Alicloud RAM User Policy Attachment resource.
---

# alicloud_ram_user_policy_attachment

Provides a RAM User Policy Attachment resource.


For information about RAM User Policy Attachment and how to use it, see [What is User Policy Attachment](https://next.api.alibabacloud.com/document/Ram/2015-05-01/AttachPolicyToUser).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_user_policy_attachment&exampleId=bb51ee62-921c-07ac-96eb-459c30b73ef77c7c5674&activeTab=example&spm=docs.r.ram_user_policy_attachment.0.bb51ee6292&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a RAM User Policy attachment.
resource "alicloud_ram_user" "user" {
  name         = "userName"
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "tf-example-${random_integer.default.result}"
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "oss:ListObjects",
          "oss:GetObject"
        ],
        "Effect": "Allow",
        "Resource": [
          "acs:oss:*:*:mybucket",
          "acs:oss:*:*:mybucket/*"
        ]
      }
    ],
      "Version": "1"
  }
  EOF
  description     = "this is a policy test"
}

resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.policy_name
  policy_type = alicloud_ram_policy.policy.type
  user_name   = alicloud_ram_user.user.name
}
```

## Argument Reference

The following arguments are supported:
* `policy_name` - (Required, ForceNew) The permission policy name.
* `policy_type` - (Required, ForceNew) Permission policy type. Value:
  - System: System policy.
  - Custom: Custom policy.
* `user_name` - (Required, ForceNew) The RAM user name.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of User Policy Attachment. The value is formulated as `user:<policy_name>:<policy_type>:<user_name>`.

## Timeouts

-> **NOTE:** Available since v1.246.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the User Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the User Policy Attachment.

## Import

RAM User Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_user_policy_attachment.example user:<policy_name>:<policy_type>:<user_name>
```
