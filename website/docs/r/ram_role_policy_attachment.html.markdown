---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role_policy_attachment"
description: |-
  Provides a Alicloud RAM Role Policy Attachment resource.
---

# alicloud_ram_role_policy_attachment

Provides a RAM Role Policy Attachment resource.


For information about RAM Role Policy Attachment and how to use it, see [What is Role Policy Attachment](https://next.api.alibabacloud.com/document/Ram/2015-05-01/AttachPolicyToRole).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_role_policy_attachment&exampleId=d1fa1b9b-e7df-bcf0-7fdb-5ab39b581203ff08373e&activeTab=example&spm=docs.r.ram_role_policy_attachment.0.d1fa1b9be7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a RAM Role Policy attachment.
resource "alicloud_ram_role" "role" {
  name        = "roleName"
  document    = <<EOF
    {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "apigateway.aliyuncs.com", 
              "ecs.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
    }
    EOF
  description = "this is a role test."
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

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.policy_name
  policy_type = alicloud_ram_policy.policy.type
  role_name   = alicloud_ram_role.role.name
}
```

## Argument Reference

The following arguments are supported:
* `policy_name` - (Required, ForceNew) The name of the policy.
* `policy_type` - (Required, ForceNew) The type of the policy. Valid values: `System` and `Custom`.
* `role_name` - (Required, ForceNew) The name of the RAM role.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Role Policy Attachment. It formats as `role:<policy_name>:<policy_type>:<role_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Role Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Role Policy Attachment.

## Import

RAM Role Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_role_policy_attachment.example role:<policy_name>:<policy_type>:<role_name>
```
