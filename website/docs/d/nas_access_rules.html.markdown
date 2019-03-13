---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_rules"
sidebar_current: "docs-alicloud-datasource-nas-access-rules"
description: |-
  Provides a list of AccessRules owned by an Alibaba Cloud account.
---

# alicloud\_nas_accessrules

This data source provides AccessRule available to the user.

-> NOTE: Available in 1.35.0+

## Example Usage

```
data "alicloud_nas_access_rules" "foo" {
  access_group_name = "tf-testAccAccessGroupsdatasource"
  source_cidr_ip = "168.1.1.0/16"
  rw_access_type = "RDWR"
  user_access_type = "no_squash"
  priority = 2
}

output "first_nas_accessrules_id" {
  value = "${data.alicloud_nas_accessrules.nas_accessrules_ds.accessrules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required ForceNew) Filter results by a specific AccessGroupName.
* `source_cidr_ip` - (Optional) Filter results by a specific SourceCidrIp. 
* `user_access` - (Optional) Filter results by a specific UserAccess. 
* `rw_access` - (Optional) Filter results by a specific RWAccess. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of rule id. Each item formats as `<access_group_name>:<access rule id>`.
* `rules` - A list of VPCs. Each element contains the following attributes:
 * `source_cidr_ip`    - SourceCidrIp of the AccessRule.
 * `priority`         - Priority of the AccessRule.
 * `access_rule_id`    - AccessRuleId of the AccessRule.
 * `user_access`      - UserAccess block of the AccessRule
 * `rw_access`        - RWAccess of the AccessRule.
