---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy"
sidebar_current: "docs-alicloud-resource-resource-manager-control-policy"
description: |-
  Provides a Alicloud Resource Manager Control Policy resource.
---

# alicloud_resource_manager_control_policy

Provides a Resource Manager Control Policy resource.

For information about Resource Manager Control Policy and how to use it, see [What is Control Policy](https://www.alibabacloud.com/help/en/resource-management/latest/api-resourcedirectorymaster-2022-04-19-createcontrolpolicy).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_control_policy&exampleId=4853b74e-7473-f7d9-ae4d-554d6448cab81e064016&activeTab=example&spm=docs.r.resource_manager_control_policy.0.4853b74e74&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_resource_manager_control_policy" "example" {
  control_policy_name = var.name
  description         = var.name
  effect_scope        = "RAM"
  policy_document     = <<EOF
  {
    "Version": "1",
    "Statement": [
      {
        "Effect": "Deny",
        "Action": [
          "ram:UpdateRole",
          "ram:DeleteRole",
          "ram:AttachPolicyToRole",
          "ram:DetachPolicyFromRole"
        ],
        "Resource": "acs:ram:*:*:role/ResourceDirectoryAccountAccessRole"
      }
    ]
  }
  EOF
}

```

## Argument Reference

The following arguments are supported:

* `control_policy_name` - (Required) The name of control policy.
* `description` - (Optional) The description of control policy.
* `effect_scope` - (Required, ForceNew) The effect scope. Valid values `RAM`.
* `policy_document` - (Required) The policy document of control policy.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Control Policy.

## Import

Resource Manager Control Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_control_policy.example <id>
```
