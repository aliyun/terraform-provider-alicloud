---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policies"
sidebar_current: "docs-alicloud-datasource-ram-policies"
description: |-
    Provides a list of ram policies available to the user.
---

# alicloud_ram_policies

This data source provides a list of RAM policies in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.0.0+.

## Example Usage

```terraform
resource "alicloud_ram_group" "group" {
  name     = "groupName-${random_integer.default.result}"
  comments = "this is a group comments."
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

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.policy_name
  policy_type = alicloud_ram_policy.policy.type
  group_name  = alicloud_ram_group.group.name
}
data "alicloud_ram_policies" "policies_ds" {
  group_name = alicloud_ram_group_policy_attachment.attach.group_name
  type       = "Custom"
}

output "first_policy_name" {
  value = data.alicloud_ram_policies.policies_ds.policies.0.name
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter resulting policies by name.
* `ids` - (Optional, ForceNew) A list of ram group IDs. 
* `type` - (Optional, ForceNew) Filter results by a specific policy type. Valid values are `Custom` and `System`.
* `user_name` - (Optional, ForceNew) Filter results by a specific user name. Returned policies are attached to the specified user.
* `group_name` - (Optional, ForceNew) Filter results by a specific group name. Returned policies are attached to the specified group.
* `role_name` - (Optional, ForceNew) Filter results by a specific role name. Returned policies are attached to the specified role.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional, Available since 1.114.0+) Default to `true`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ram group names.
* `policies` - A list of policies. Each element contains the following attributes:
  * `id` - ID of the policy.
  * `name` - Name of the policy.
  * `policy_name` - Name of the policy.
  * `type` - Type of the policy.
  * `description` - Description of the policy.
  * `default_version` - Default version of the policy.
  * `create_date` - Creation date of the policy.
  * `update_date` - Update date of the policy.
  * `attachment_count` - Attachment count of the policy.
  * `document` - Policy document of the policy.
  * `policy_document` - Policy document of the policy.
  * `version_id` - The ID of default policy.
  * `user_name` - The user name of  policy.
