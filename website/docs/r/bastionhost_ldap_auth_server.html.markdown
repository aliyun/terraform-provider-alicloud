---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_ldap_auth_server"
sidebar_current: "docs-alicloud-resource-bastionhost-ldap-auth-server"
description: |-
  Provides a Alicloud Bastion Host Ldap Auth Server resource.
---

# alicloud\_bastionhost\_ldap\_auth\_server

Provides a Bastion Host Ldap Auth Server resource.

For information about Bastion Host Ldap Auth Server and how to use it, see [What is Ldap Auth Server](https://www.alibabacloud.com/help/en/bastion-host/latest/api-modifyinstanceldapauthserver-v1).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_instances" "default" {}
resource "alicloud_bastionhost_ldap_auth_server" "example" {
  account     = "example_value"
  base_dn     = "example_value"
  instance_id = data.alicloud_bastionhost_instances.default.ids.0
  is_ssl      = false
  port        = 80
  server      = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Required) The username of the account that is used for the LDAP server.
* `base_dn` - (Required) The Base distinguished name (DN).
* `email_mapping` - (Optional) The field that is used to indicate the email address of a user on the LDAP server.
* `filter` - (Optional) The condition that is used to filter users.
* `instance_id` - (Required, ForceNew) The ID of the Bastion machine instance.
* `is_ssl` - (Optional) Specifies whether to support SSL.
* `login_name_mapping` - (Optional) The field that is used to indicate the logon name of a user on the LDAP server.
* `mobile_mapping` - (Optional) The field that is used to indicate the mobile phone number of a user on the LDAP server.
* `name_mapping` - (Optional) The field that is used to indicate the name of a user on the LDAP server.
* `password` - (Required, Sensitive) The password of the account that is used for the LDAP server.
* `port` - (Required) The port that is used to access the LDAP server.
* `server` - (Required) The address of the LDAP server.
* `standby_server` - (Optional) The address of the secondary LDAP server.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ldap Auth Server. Its value is same as `instance_id`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ldap Auth Server.
* `update` - (Defaults to 1 mins) Used when updating the Ldap Auth Server.

## Import

Bastion Host Ldap Auth Server can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_ldap_auth_server.example <instance_id>
```