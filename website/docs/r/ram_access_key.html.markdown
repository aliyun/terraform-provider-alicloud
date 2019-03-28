---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_access_key"
sidebar_current: "docs-alicloud-resource-ram-access-key"
description: |-
  Provides a RAM User access key resource.
---

# alicloud\_ram\_access\_key

Provides a RAM User access key resource.

-> **NOTE:**  You should set the `secret_file` if you want to get the access key.  

## Example Usage

```
# Create a new RAM access key for user.
resource "alicloud_ram_user" "user" {
  name = "user_test"
  display_name = "user_display_name"
  mobile = "86-18688888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
  force = true
}

resource "alicloud_ram_access_key" "ak" {
  user_name = "${alicloud_ram_user.user.name}"
  secret_file = "/xxx/xxx/xxx.txt"
}
```
## Argument Reference

The following arguments are supported:

* `user_name` - (Optional, ForceNew) Name of the RAM user. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen.
* `secret_file` - (Optional, ForceNew) The name of file that can save access key id and access key secret. Strongly suggest you to specified it when you creating access key, otherwise, you wouldn't get its secret ever.
* `status` - (Optional) Status of access key. It must be `Active` or `Inactive`. Default value is `Active`.

## Attributes Reference

The following attributes are exported:

* `id` - The access key ID.
* `status` - The access key status.