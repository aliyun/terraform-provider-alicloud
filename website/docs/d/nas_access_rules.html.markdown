---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_rules"
sidebar_current: "docs-alicloud-datasource-nas-access-rules"
description: |-
  Provides a list of AccessRules owned by an Alibaba Cloud account.
---

# alicloud\_nas_access_rules

This data source provides AccessRule available to the user.

-> **NOTE**: Available in 1.35.0+

## Example Usage

```terraform
data "alicloud_nas_access_rules" "foo" {
  access_group_name = "tf-testAccAccessGroupsdatasource"
  source_cidr_ip    = "168.1.1.0/16"
  rw_access         = "RDWR"
  user_access       = "no_squash"
}

output "alicloud_nas_access_rules_id" {
  value = "${data.alicloud_nas_access_rules.foo.rules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required ForceNew) Filter results by a specific AccessGroupName.
* `ids` - (Optional, Available in 1.53.0+) A list of rule IDs.
* `source_cidr_ip` - (Optional) Filter results by a specific SourceCidrIp. 
* `user_access` - (Optional) Filter results by a specific UserAccess. 
* `rw_access` - (Optional) Filter results by a specific RWAccess. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of rule IDs, Each element set to `access_rule_id` (Each element formats as `<access_group_name>:<access_rule_id>` before 1.53.0).
* `rules` - A list of AccessRules. Each element contains the following attributes:
 * `source_cidr_ip` - SourceCidrIp of the AccessRule.
 * `priority` - Priority of the AccessRule.
 * `access_rule_id` - AccessRuleId of the AccessRule.
 * `user_access` - UserAccess of the AccessRule
 * `rw_access` - RWAccess of the AccessRule.
