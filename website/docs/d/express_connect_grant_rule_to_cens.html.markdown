---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_grant_rule_to_cens"
sidebar_current: "docs-alicloud-datasource-express-connect-grant-rule-to-cens"
description: |-
  Provides a list of Express Connect Grant Rule To Cens to the user.
---

# alicloud\_express\_connect\_grant\_rule\_to\_cens

This data source provides the Express Connect Grant Rule To Cens of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.196.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_express_connect_grant_rule_to_cens" "ids" {
  ids         = ["example_id"]
  instance_id = "your_vbr_instance_id"
}

output "express_connect_grant_rule_to_cen_id_0" {
  value = data.alicloud_express_connect_grant_rule_to_cens.ids.cens.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Grant Rule To Cen IDs.
* `instance_id` - (Required, ForceNew) The ID of the VBR.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `cens` - A list of Express Connect Grant Rule To Cens. Each element contains the following attributes:
	* `id` - The ID of the Grant Rule To Cen. It formats as `<cen_id>:<cen_owner_id>:<instance_id>`.
	* `cen_id` - The ID of the authorized CEN instance.
	* `cen_owner_id` - The user ID (UID) of the Alibaba Cloud account to which the CEN instance belongs.
	* `create_time` - The time when the instance was created.
	