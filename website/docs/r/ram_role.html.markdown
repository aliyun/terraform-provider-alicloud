---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role"
description: |-
  Provides a Alicloud RAM Role resource.
---

# alicloud_ram_role

Provides a RAM Role resource.



For information about RAM Role and how to use it, see [What is Role](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createrole).

-> **NOTE:** Available since v1.0.0.

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_role&exampleId=5c79f663-631e-b723-f6d9-85e5d8983056326aab52&activeTab=example&spm=docs.r.ram_role.0.5c79f66363&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_role" "default" {
  role_name                   = "terraform-example-${random_integer.default.result}"
  assume_role_policy_document = <<EOF
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
  description                 = "this is a role test."
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ram_role&spm=docs.r.ram_role.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `assume_role_policy_document` - (Optional, Available since v1.252.0) The trust policy that specifies one or more trusted entities to assume the RAM role. The trusted entities can be Alibaba Cloud accounts, Alibaba Cloud services, or identity providers (IdPs).
* `description` - (Optional) The description of the RAM role.
* `max_session_duration` - (Optional, Int, Available since v1.105.0) The maximum session time of the RAM role. Default value: `3600`. Valid values: `3600` to `43200`.
* `role_name` - (Optional, ForceNew, Available since v1.252.0) The name of the RAM role.
* `tags` - (Optional, Map, Available since v1.252.0) The list of tags for the role.
* `force` - (Optional, Bool) Specifies whether to force delete the Role. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `name` - (Optional, ForceNew, Deprecated since v1.252.0) Field `name` has been deprecated from provider version 1.252.0. New field `role_name` instead.
* `document` - (Optional, Deprecated since v1.252.0) Field `document` has been deprecated from provider version 1.252.0. New field `assume_role_policy_document` instead.
* `version` - (Optional, Deprecated since v1.49.0) Field `version` has been deprecated from provider version 1.49.0. New field `document` instead.
* `ram_users` - (Optional, List, Deprecated since v1.49.0) Field `ram_users` has been deprecated from provider version 1.49.0. New field `document` instead.
* `services` - (Optional, List, Deprecated since v1.49.0) Field `services` has been deprecated from provider version 1.49.0. New field `document` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `arn` - The Alibaba Cloud Resource Name (ARN) of the RAM role.
* `create_time` - (Available since v1.252.0) The time when the RAM role was created.
* `role_id` - The ID of the RAM role.

## Timeouts

-> **NOTE:** Available since v1.159.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Role.
* `delete` - (Defaults to 5 mins) Used when delete the Role.
* `update` - (Defaults to 5 mins) Used when update the Role.

## Import

RAM Role can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_role.example <id>
```
