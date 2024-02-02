---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_proxy_access"
sidebar_current: "docs-alicloud-resource-dms-enterprise-proxy-access"
description: |-
  Provides a Alicloud DMS Enterprise Proxy Access resource.
---

# alicloud_dms_enterprise_proxy_access

Provides a DMS Enterprise Proxy Access resource.

For information about DMS Enterprise Proxy Access and how to use it, see [What is Proxy Access](https://next.api.alibabacloud.com/document/dms-enterprise/2018-11-01/CreateProxyAccess).

-> **NOTE:** Available since v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dms_enterprise_users" "dms_enterprise_users_ds" {
  role   = "USER"
  status = "NORMAL"
}
data "alicloud_dms_enterprise_proxies" "ids" {}

resource "alicloud_dms_enterprise_proxy_access" "default" {
  proxy_id = data.alicloud_dms_enterprise_proxies.ids.proxies.0.id
  user_id  = data.alicloud_dms_enterprise_users.dms_enterprise_users_ds.users.0.user_id
}
```

## Argument Reference

The following arguments are supported:
* `indep_account` - (ForceNew, Optional) Database account.
* `indep_password` - (ForceNew, Optional) Database password.
* `proxy_id` - (Required, ForceNew) The ID of the security agent. 
* `user_id` - (Required, ForceNew) The user ID.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the proxy access resource.
* `access_id` - The authorized account of the security agent.
* `access_secret` - Secure access agent authorization password.
* `gmt_create` - The authorization time of the security access agent permission.
* `instance_id` - The ID of the instance.
* `origin_info` - The source information of the security access agent permission is enabled, and the return value is as follows:
  * **Owner Authorization**: The UID of the owner in parentheses.
  * **Work Order Authorization**: The ticket number in parentheses is the number of the user to apply for permission.
* `proxy_access_id` - Security Protection authorization ID. After the target user is authorized by the security protection agent, the system automatically generates a security protection authorization ID, which is globally unique.
* `user_name` - User nickname.
* `user_uid` - User UID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Proxy Access.
* `delete` - (Defaults to 5 mins) Used when delete the Proxy Access.

## Import

DMS Enterprise Proxy Access can be imported using the id, e.g.

```shell
$terraform import alicloud_dms_enterprise_proxy_access.example <id>
```