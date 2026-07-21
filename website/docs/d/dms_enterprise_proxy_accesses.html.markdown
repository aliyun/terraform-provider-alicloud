---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_proxy_accesses"
sidebar_current: "docs-alicloud-datasource-dms-enterprise-proxy-accesses"
description: |-
  Provides a list of DMS Enterprise Proxy Access owned by an Alibaba Cloud account.
---

# alicloud_dms_enterprise_proxy_accesses

This data source provides DMS Enterprise Proxy Access available to the user.[What is Proxy Access](https://next.api.alibabacloud.com/document/dms-enterprise/2018-11-01/CreateProxyAccess)

-> **NOTE:** Available since v1.195.0.

## Example Usage

```terraform
resource "alicloud_dms_enterprise_proxy_access" "default" {
  indep_password = "PASSWORD-DEMO"
  proxy_id       = 1881
  indep_account  = "dmstest"
  user_id        = 104442
}

data "alicloud_dms_enterprise_proxy_accesses" "default" {
  ids      = ["${alicloud_dms_enterprise_proxy_access.default.id}"]
  proxy_id = 1881
}

output "alicloud_dms_proxy_acceses_example_id" {
  value = data.alicloud_dms_enterprise_proxy_accesses.default.accesses.0.id
}
```

## Argument Reference

The following arguments are supported:
* `proxy_id` - (Required, ForceNew) The ID of the security agent.
* `ids` - (Optional, ForceNew, Computed) A list of Proxy Access IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Proxy Access IDs.
* `accesses` - A list of Proxy Access Entries. Each element contains the following attributes:
    * `id` - Security Protection authorization ID.
    * `access_id` - The authorized account of the security agent.
    * `create_time` - The authorization time of the security access agent permission.
    * `indep_account` - Database account.
    * `instance_id` - The ID of the instance.
    * `origin_info` - The source information of the security access agent permission is enabled, and the return value is as follows:**Owner Authorization**: The UID of the owner in parentheses.**Work Order Authorization**: The ticket number in parentheses is the number of the user to apply for permission.
    * `proxy_access_id` - Security Protection authorization ID. After the target user is authorized by the security protection agent, the system automatically generates a security protection authorization ID, which is globally unique.
    * `proxy_id` - The ID of the security agent.
    * `user_id` - The user ID.
    * `user_name` - User nickname.
    * `user_uid` - User UID.
