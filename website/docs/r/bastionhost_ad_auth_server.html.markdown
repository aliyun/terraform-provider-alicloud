---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_ad_auth_server"
sidebar_current: "docs-alicloud-resource-bastionhost-ad-auth-server"
description: |-
  Provides a Alicloud Bastion Host Ad Auth Server resource.
---

# alicloud\_bastionhost\_ad\_auth\_server

Provides a Bastion Host Ad Auth Server resource.

For information about Bastion Host Ad Auth Server and how to use it, see [What is Ad Auth Server](https://www.alibabacloud.com/help/en/bastion-host/latest/api-modifyinstanceadauthserver-v1).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_instances" "default" {}
resource "alicloud_bastionhost_ad_auth_server" "example" {
  account     = "example_value"
  base_dn     = "example_value"
  domain      = "example_value"
  instance_id = data.alicloud_bastionhost_instances.default.ids.0
  is_ssl      = false
  port        = 80
  password    = "YouPassword123"
  server      = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required) The username of the account that is used for the AD server.
* `base_dn` - (Required) The Base distinguished name (DN).
* `domain` - (Required) The domain on the AD server.
* `email_mapping` - (Optional) The field that is used to indicate the email address of a user on the AD server.
* `filter` - (Optional) The condition that is used to filter users.
* `instance_id` - (Required, ForceNew) The ID of the Bastion machine instance.
* `is_ssl` - (Required) Specifies whether to support SSL.
* `mobile_mapping` - (Optional) The field that is used to indicate the mobile phone number of a user on the AD server.
* `name_mapping` - (Optional) The field that is used to indicate the name of a user on the AD server.
* `password` - (Required, Sensitive) The password of the account that is used for the AD server.
* `port` - (Required) The port that is used to access the AD server.
* `server` - (Required) The address of the AD server.
* `standby_server` - (Optional) The address of the secondary AD server.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ad Auth Server. Its value is same as `instance_id`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ad Auth Server .
* `update` - (Defaults to 1 mins) Used when updating the Ad Auth Server .

## Import

Bastion Host Ad Auth Server can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_ad_auth_server.example <instance_id>
```