---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_access_key_policy"
description: |-
  Provides a Alicloud RAM Access Key Policy resource.
---

# alicloud_ram_access_key_policy

Provides a RAM Access Key Policy resource.

Sets the network access restriction policy for the AccessKey of an Alibaba Cloud account (primary account) or a RAM user.

For information about RAM Access Key Policy and how to use it, see [What is Access Key Policy](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-setaccesskeypolicy).

-> **NOTE:** Available since v1.286.0.

-> **NOTE:** There is no dedicated delete API for the network access restriction policy. Destroying this resource clears all whitelist rules by setting a disabled policy with no statements (`{"Version":1,"Status":"Inactive","Statements":[]}`) via the `SetAccessKeyPolicy` API. A disabled policy that carries no statements is treated as "no policy" (equivalent to never having configured one), so a policy document of this form is not considered a managed resource.

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
```

## Argument Reference

The following arguments are supported:

* `user_access_key_id` - (Required, ForceNew) The ID of the access key that the network access restriction policy applies to.
* `user_principal_name` - (Optional, ForceNew) The logon name of the RAM user. Specify this parameter when managing the access key policy of another RAM user. If it is left empty, the policy is applied to the specified access key of the current user.
* `access_key_policy` - (Required) The network access restriction policy, in JSON format. For the structure of the policy document, see [SetAccessKeyPolicy](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-setaccesskeypolicy).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Access Key Policy. When `user_principal_name` is specified, the value is formatted as `<user_principal_name>:<user_access_key_id>`. Otherwise, it is `<user_access_key_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Access Key Policy.
* `update` - (Defaults to 5 mins) Used when update the Access Key Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Access Key Policy.

## Import

RAM Access Key Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_access_key_policy.example <user_access_key_id>
```

If the policy is set for another RAM user, use the composite id:

```shell
$ terraform import alicloud_ram_access_key_policy.example <user_principal_name>:<user_access_key_id>
```
