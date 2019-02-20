---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_accessrule"
sidebar_current: "docs-alicloud-resource-nas_accessrule"
description: |-
  Provides a Alicloud NAS_AccessRule resource.
---

# alicloud\_nas_accessrule

Provides a NAS_AccessRule resource.

## Example Usage

Basic Usage

```
resource "alicloud_nas_accessrule" "foo" {
		accessgroup_name = "CreateAccessGroup"
		sourcecidr_ip = "168.1.1.0/16"
		rwaccess_type = "RDWR"
		useraccess_type = "no_squash"
}
```
## Argument Reference

The following arguments are supported:

* `accessgroup_name` - (Required, Forces new resource) The AccessGroupName block for the AccessRule.
* `sourcecidr_ip`    - (Required, Forces new resource) The SourceCidrIp block for the AccessRule
* `rwaccess_type`    - (Optional) The AccessRule RWAccessType. Defaults to "RDWR".
* `useraccess_type`  - (Optional) The AccessRule UserAccessType. Defaults to "no_squash".
* `Priority`         - (Optional) The AccessRule Priority. Defaults to 1.

## Attributes Reference

The following attributes are exported:


	d.Set("sourcecidr_ip", resp.SourceCidrIp)
	d.Set("accessrule_id", resp.AccessRuleId)
	d.Set("priority", resp.Priority)
	if resp.RWAccess != "" {
		d.Set("rwaccess_type", resp.RWAccess)
	}
	if resp.UserAccess != "" {
		d.Set("useraccess_type", resp.UserAccess)
	}


* `sourcecidr_ip`    - The SourceCidrIp of the AccessRule.
* `accessrule_id`    - The AccessRuleId block for the AccessRule.
* `priority`         - The Priority block for the AccessRule.
* `rwaccess_type`    - The RWAccess of the AccessRule.
* `useraccess_type`  - The UserAccess of the AccessRule.


