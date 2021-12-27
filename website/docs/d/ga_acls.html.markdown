---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acls"
sidebar_current: "docs-alicloud-datasource-ga-acls"
description: |-
  Provides a list of Ga Acls to the user.
---

# alicloud\_ga\_acls

This data source provides the Ga Acls of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.150.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_acls" "ids" {}
output "ga_acl_id_1" {
  value = data.alicloud_ga_acls.ids.acls.0.id
}

data "alicloud_ga_acls" "nameRegex" {
  name_regex = "^my-Acl"
}
output "ga_acl_id_2" {
  value = data.alicloud_ga_acls.nameRegex.acls.0.id
}
```

## Argument Reference

The following arguments are supported:

* `acl_name` - (Optional, ForceNew) The name of the acl.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Acl IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Acl name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `active`, `configuring`, `deleting`, `init`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Acl names.
* `acls` - A list of Ga Acls. Each element contains the following attributes:
	* `acl_entries` - The entries of the Acl.
		* `entry` - The IP entry that you want to add to the ACL.
		* `entry_description` - The description of the IP entry.
	* `acl_id` - The  ID of the Acl.
	* `acl_name` - The name of the acl.
	* `address_ip_version` - The address ip version.
	* `id` - The ID of the Acl. Its value is same as `acl_id`.
	* `status` - The status of the resource.