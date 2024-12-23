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

For information about Resource Manager Control Policy Attachment and how to use it, see [What is Control Policy Attachment](https://www.alibabacloud.com/help/en/resource-management/resource-directory/developer-reference/api-resourcemanager-2020-03-31-attachcontrolpolicy).

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
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_resource_manager_control_policy" "default" {
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

resource "alicloud_resource_manager_folder" "default" {
  folder_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_resource_manager_control_policy_attachment" "default" {
  policy_id = alicloud_resource_manager_control_policy.default.id
  target_id = alicloud_resource_manager_folder.default.id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The ID of the access control policy.
* `target_id` - (Required, ForceNew) The ID of the object to which you want to attach the access control policy.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Control Policy Attachment. It formats as `<policy_id>:<target_id>`.

## Timeouts

-> **NOTE:** Available since v1.240.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Control Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Control Policy Attachment.

## Import

Resource Manager Control Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_control_policy_attachment.example <policy_id>:<target_id>
```
