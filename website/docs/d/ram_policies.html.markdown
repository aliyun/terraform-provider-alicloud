---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policies"
description: |-
  Provides a list of RAM Policies to the user.
---

# alicloud_ram_policies

This data source provides the RAM Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_policy" "default" {
  policy_name     = "${var.name}-${random_integer.default.result}"
  description     = "${var.name}-${random_integer.default.result}"
  force           = true
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Effect": "Allow",
        "Action": "*",
        "Resource": "*"
      }
    ],
    "Version": "1"
  }
  EOF
  tags = {
    Created = "TF"
    For     = "Policy"
  }
}

data "alicloud_ram_policies" "ids" {
  ids = [alicloud_ram_policy.default.id]
}

output "ram_policies_id_0" {
  value = data.alicloud_ram_policies.ids.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` (Optional, ForceNew, List, Available since v1.114.0) - A list of Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Policy name.
* `type` - (Optional, ForceNew) The type of the policy. Valid values: `System` and `Custom`.
* `user_name` - (Optional, ForceNew) The name of the RAM user.
* `group_name` - (Optional, ForceNew) The name of the user group.
* `role_name` - (Optional, ForceNew) The name of the RAM role.
* `tags` - (Optional, ForceNew, Available since v1.262.1) A mapping of tags to assign to the resource.
* `enable_details` -(Optional, Bool, Available since v1.114.0) Whether to query the detailed list of resource attributes. Default value: `true`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - (Available since v1.42.0) A list of Policy names.
* `policies` - A list of Policy. Each element contains the following attributes:
  * `id` - (Available since v1.114.0) The ID of the Policy.
  * `policy_name` - (Available since v1.114.0) The name of the policy.
  * `name` - The name of the policy.
  * `type` - The type of the policy.
  * `description` - The description of the policy.
  * `tags` - (Available since v1.262.1) The tags of the Policy.
  * `default_version` - The default version of the policy.
  * `attachment_count` - The number of references to the policy.
  * `policy_document` - (Available since v1.114.0) The document of the policy. **Note:** `policy_document` takes effect only if `enable_details` is set to `true`.
  * `document` - The document of the policy. **Note:** `document` takes effect only if `enable_details` is set to `true`.
  * `version_id` - (Available since v1.114.0) The ID of the default policy version. **Note:** `version_id` takes effect only if `enable_details` is set to `true`.
  * `create_date` - The time when the policy was created.
  * `update_date` - The time when the policy was modified.
  * `user_name` - (Removed since v1.262.1) Field `user_name` has been removed from provider version 1.262.1.
