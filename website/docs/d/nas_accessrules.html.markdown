---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_accessrules"
sidebar_current: "docs-alicloud-datasource-nas-accessrules"
description: |-
    Provides a list of AccessRules owned by an Alibaba Cloud account.
---

# alicloud\_nas_accessrules

This data source provides AccessRule available to the user.

## Example Usage

```
data "alicloud_nas_accessrules" "ar" {
  accessgroup_name = "classic"
}

output "first_nas_accessrules_id" {
  value = "${data.alicloud_nas_accessrules.nas_accessrules_ds.accessrules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `accessgroup_name` - (Required) Filter results by a specific AccessGroupName block. 

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `accessrules` - A list of VPCs. Each element contains the following attributes:
  * `sourcecidr_ip`    - SourceCidrIp of the AccessRule.
  * `priority`         - Priority of the AccessRule.
  * `accessrule_id`    - AccessRuleId of the AccessRule.
  * `user_access`      - UserAccess block of the AccessRule
  * `rw_access`        - RWAccess of the AccessRule.
