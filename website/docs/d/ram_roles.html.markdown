---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_roles"
description: |-
  Provides a list of RAM Roles to the user.
---

# alicloud_ram_roles

This data source provides the RAM Roles of the current Alibaba Cloud user.

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

resource "alicloud_ram_role" "default" {
  role_name                   = "${var.name}-${random_integer.default.result}"
  description                 = "${var.name}-${random_integer.default.result}"
  force                       = true
  assume_role_policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "ecs.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  tags = {
    Created = "TF"
    For     = "Role"
  }
}

data "alicloud_ram_roles" "ids" {
  ids = [alicloud_ram_role.default.role_id]
}

output "ram_roles_id_0" {
  value = data.alicloud_ram_roles.ids.roles.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` (Optional, ForceNew, List, Available since v1.42.0) - A list of Role IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Role name.
* `policy_name` - (Optional, ForceNew) The name of the policy.
* `policy_type` - (Optional, ForceNew) The type of the policy. Default value: `System`. Valid values: `System`, `Custom`. **Note:** `policy_type` takes effect only when `policy_name` is set.
* `tags` - (Optional, ForceNew, Available since v1.262.1) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - (Available since v1.42.0) A list of Role names.
* `roles` - A list of Role. Each element contains the following attributes:
  * `id` - The ID of the RAM role.
  * `name` - The name of the RAM role.
  * `assume_role_policy_document` - The policy that specifies the trusted entity to assume the RAM role.
  * `document` - The policy that specifies the trusted entity to assume the RAM role.
  * `description` - The description of the RAM role.
  * `tags` - (Available since v1.262.1) The tags of the RAM role.
  * `arn` - The Alibaba Cloud Resource Name (ARN) of the RAM role.
  * `create_date` - The creation time.
  * `update_date` - The update time.
