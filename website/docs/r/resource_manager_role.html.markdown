---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_role"
sidebar_current: "docs-alicloud-resource-resource-manager-role"
description: |-
  Provides a Resource Manager role resource.
---

# alicloud\_resource\_manager\_role

Provides a Resource Manager role resource. Members are resource containers in the resource directory, which can physically isolate resources to form an independent resource grouping unit. You can create members in the resource folder to manage them in a unified manner.
For information about Resource Manager role and how to use it, see [What is Resource Manager role](https://www.alibabacloud.com/help/en/doc-detail/111231.htm).

-> **NOTE:** Available in v1.82.0+.

## Example Usage

```terraform
# Add a Resource Manager role.
resource "alicloud_resource_manager_role" "example" {
  role_name                   = "testrd"
  assume_role_policy_document = <<EOF
     {
          "Statement": [
               {
                    "Action": "sts:AssumeRole",
                    "Effect": "Allow",
                    "Principal": {
                        "RAM":[
                                "acs:ram::103755469187****:root"，
                                "acs:ram::104408977069****:root"
                        ]
                    }
                }
          ],
          "Version": "1"
     }
	 EOF
}
```
## Argument Reference

The following arguments are supported:

* `assume_role_policy_document` - (Required) The content of the permissions strategy that plays a role.
* `description` - (Optional, ForceNew) The description of the Resource Manager role.
* `max_session_duration` - (Optional) Role maximum session time. Valid values: [3600-43200]. Default to `3600`.
* `role_name` - (Required, ForceNew) Role Name. The length is 1 ~ 64 characters, which can include English letters, numbers, dots "." and dashes "-".

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Resource Manager role. The value is set to `role_name`.
* `role_id` - This ID of Resource Manager role. The value is set to `role_name`.
* `arn` - The resource descriptor of the role.
* `create_date` (Removed form v1.114.0) - Role creation time.
* `update_date` - Role update time.

## Import

Resource Manager can be imported using the id or role_name, e.g.

```
$ terraform import alicloud_resource_manager_role.example testrd
```
