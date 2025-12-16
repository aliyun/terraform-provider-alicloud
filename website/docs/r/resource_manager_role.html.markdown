---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_role"
sidebar_current: "docs-alicloud-resource-resource-manager-role"
description: |-
  Provides a Resource Manager role resource.
---

# alicloud_resource_manager_role

Provides a Resource Manager role resource. Members are resource containers in the resource directory, which can physically isolate resources to form an independent resource grouping unit. You can create members in the resource folder to manage them in a unified manner.
For information about Resource Manager role and how to use it, see [What is Resource Manager role](https://www.alibabacloud.com/help/en/doc-detail/111231.htm).

-> **NOTE:** Available since v1.82.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_role&exampleId=8a705046-004e-a81d-e9ef-5b4264155fefe6cec622&activeTab=example&spm=docs.r.resource_manager_role.0.8a70504600&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
data "alicloud_account" "default" {}

resource "alicloud_resource_manager_role" "example" {
  role_name                   = var.name
  assume_role_policy_document = <<EOF
     {
          "Statement": [
               {
                    "Action": "sts:AssumeRole",
                    "Effect": "Allow",
                    "Principal": {
                        "RAM":[
                                "acs:ram::${data.alicloud_account.default.id}:root"
                        ]
                    }
                }
          ],
          "Version": "1"
     }
	 EOF
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_role&spm=docs.r.resource_manager_role.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `assume_role_policy_document` - (Required) The content of the permissions strategy that plays a role.
* `description` - (Optional, ForceNew) The description of the Resource Manager role.
* `max_session_duration` - (Optional) Role maximum session time. Valid values: [3600-43200]. Default to `3600`.
* `role_name` - (Required, ForceNew) Role Name. The length is 1 ~ 64 characters, which can include English letters, numbers, dots "." and dashes "-".
* `create_date` (Removed form v1.114.0) - Role creation time.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Resource Manager role. The value is set to `role_name`.
* `role_id` - This ID of Resource Manager role. The value is set to `role_name`.
* `arn` - The resource descriptor of the role.
* `update_date` - Role update time.

## Import

Resource Manager can be imported using the id or role_name, e.g.

```shell
$ terraform import alicloud_resource_manager_role.example testrd
```
