---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy_attachment"
sidebar_current: "docs-alicloud-resource-resource-manager-control-policy-attachment"
description: |-
  Provides a Alicloud Resource Manager Control Policy Attachment resource.
---

# alicloud_resource_manager_control_policy_attachment

Provides a Resource Manager Control Policy Attachment resource.

For information about Resource Manager Control Policy Attachment and how to use it, see [What is Control Policy Attachment](https://www.alibabacloud.com/help/en/resource-management/latest/api-resourcedirectorymaster-2022-04-19-attachcontrolpolicy).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_control_policy_attachment&exampleId=b1a1a3c8-4441-5fcb-eeec-fc8448ea447cf90ca039&activeTab=example&spm=docs.r.resource_manager_control_policy_attachment.0.b1a1a3c844&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
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

resource "alicloud_resource_manager_folder" "example" {
  folder_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_resource_manager_control_policy_attachment" "example" {
  policy_id = alicloud_resource_manager_control_policy.example.id
  target_id = alicloud_resource_manager_folder.example.id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The ID of control policy.
* `target_id` - (Required, ForceNew) The ID of target.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Control Policy Attachment. The value is formatted `<policy_id>:<target_id>`.

## Import

Resource Manager Control Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_control_policy_attachment.example <policy_id>:<target_id>
```
