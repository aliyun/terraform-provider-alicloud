---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_access_key_policy"
description: |-
  Provides a Alicloud RAM Access Key Policy data source.
---

# alicloud_ram_access_key_policy

Provides a RAM Access Key Policy data source.

Reads the network access restriction policy for the AccessKey of an Alibaba Cloud account (primary account) or a RAM user.

For information about RAM Access Key Policy and how to use it, see [What is Access Key Policy](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-getaccesskeypolicy).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_ram_access_key" "default" {
  user_name = alicloud_ram_user.default.name
}

resource "alicloud_ram_access_key_policy" "default" {
  user_access_key_id  = alicloud_ram_access_key.default.id
  user_principal_name = "${alicloud_ram_user.default.name}@${data.alicloud_account.default.id}.onaliyun.com"
  access_key_policy = jsonencode({
    Status = "Active"
    Statements = [{
      Type   = "ClassicWhiteList"
      IPList = ["10.0.0.1/32"]
    }]
  })
}

data "alicloud_ram_access_key_policy" "default" {
  user_access_key_id  = alicloud_ram_access_key_policy.default.user_access_key_id
  user_principal_name = alicloud_ram_access_key_policy.default.user_principal_name
}
```

## Argument Reference

The following arguments are supported:

* `user_access_key_id` - (Required) The ID of the access key whose network access restriction policy is to be read.
* `user_principal_name` - (Optional) The logon name of the RAM user. Specify this parameter when reading the access key policy of another RAM user. If it is left empty, the policy of the specified access key of the current user is read.

## Attributes Reference

The following attributes are exported:

* `id` - The data source ID. When `user_principal_name` is specified, the value is formatted as `<user_principal_name>:<user_access_key_id>`. Otherwise, it is `<user_access_key_id>`.
* `access_key_policy` - The network access restriction policy in JSON format.
