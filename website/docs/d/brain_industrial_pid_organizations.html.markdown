---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_organizations"
sidebar_current: "docs-alicloud-datasource-brain-industrial-pid-organizations"
description: |-
  Provides a list of Brain Industrial Pid Organizations to the user.
---

# alicloud\_brain\_industrial\_pid\_organizations

This data source provides the Brain Industrial Pid Organizations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_brain_industrial_pid_organizations" "example" {
  ids        = ["3e74e684-cbb5-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_brain_industrial_pid_organization_id" {
  value = data.alicloud_brain_industrial_pid_organizations.example.organizations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Pid Organization IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Pid Organization name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `parent_organization_id` - (Optional, ForceNew) The parent organization id.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Pid Organization names.
* `organizations` - A list of Brain Industrial Pid Organizations. Each element contains the following attributes:
	* `id` - The ID of the Pid Organization.
	* `parent_pid_organization_id` - The parent organization id.
	* `pid_organization_id` - The organization id.
	* `pid_organization_level` - The organization level.
	* `pid_organization_name` - The organization name.
