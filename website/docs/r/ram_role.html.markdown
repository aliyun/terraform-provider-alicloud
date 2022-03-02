---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role"
sidebar_current: "docs-alicloud-resource-ram-role"
description: |-
  Provides a RAM Role resource.
---

# alicloud\_ram\_role

Provides a RAM Role resource.

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

## Example Usage

```terraform
# Create a new RAM Role.
resource "alicloud_ram_role" "role" {
  name        = "testrole"
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
  force       = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the RAM role. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-", "_", and must not begin with a hyphen.
* `services` - (Deprecated, Optional, Type: list, Conflicts with `document`) (It has been deprecated from version 1.49.0, and use field 'document' to replace.) List of services which can assume the RAM role. The format of each item in this list is `${service}.aliyuncs.com` or `${account_id}@${service}.aliyuncs.com`, such as `ecs.aliyuncs.com` and `1234567890000@ots.aliyuncs.com`. The `${service}` can be `ecs`, `log`, `apigateway` and so on, the `${account_id}` refers to someone's Alicloud account id.
* `ram_users` - (Deprecated, Optional, Type: list, Conflicts with `document`) (It has been deprecated from version 1.49.0, and use field 'document' to replace.) List of ram users who can assume the RAM role. The format of each item in this list is `acs:ram::${account_id}:root` or `acs:ram::${account_id}:user/${user_name}`, such as `acs:ram::1234567890000:root` and `acs:ram::1234567890001:user/Mary`. The `${user_name}` is the name of a RAM user which must exists in the Alicloud account indicated by the `${account_id}`.
* `version` - (Deprecated, Optional, Conflicts with `document`) (It has been deprecated from version 1.49.0, and use field 'document' to replace.) Version of the RAM role policy document. Valid value is `1`. Default value is `1`.
* `document` - (Optional, Conflicts with `services`, `ram_users` and `version`) Authorization strategy of the RAM role. It is required when the `services` and `ram_users` are not specified.
* `description` - (Optional) Description of the RAM role. This name can have a string of 1 to 1024 characters. **NOTE:** The `description` supports modification since V1.144.0.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.
* `max_session_duration` - (Optional, Available in 1.105.0+) The maximum session duration of the RAM role. Valid values: 3600 to 43200. Unit: seconds. Default value: 3600. The default value is used if the parameter is not specified.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is set to `role_name`.
* `role_id` - The role ID.
* `name` - The role name.
* `arn` - The role arn.
* `description` - The role description.
* `version` - The role policy document version.
* `document` - Authorization strategy of the role.
* `ram_users` - List of services which can assume the RAM role. 
* `services` - List of services which can assume the RAM role.

### Timeouts

-> **NOTE:** Available in v1.159.0+

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the ram role.
* `update` - (Defaults to 10 mins) Used when updating the ram role.
* `delete` - (Defaults to 10 mins) Used when deleting the ram role.

## Import

RAM role can be imported using the id or name, e.g.

```
$ terraform import alicloud_ram_role.example my-role
```
