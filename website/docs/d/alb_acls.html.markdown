---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_acls"
sidebar_current: "docs-alicloud-datasource-alb-acls"
description: |-
  Provides a list of Application Load Balancer (ALB) Acls to the user.
---

# alicloud\_alb\_acls

This data source provides the Application Load Balancer (ALB) Acls of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_acls" "ids" {}
output "alb_acl_id_1" {
  value = data.alicloud_alb_acls.ids.acls.0.id
}

data "alicloud_alb_acls" "nameRegex" {
  name_regex = "^my-Acl"
}
output "alb_acl_id_2" {
  value = data.alicloud_alb_acls.nameRegex.acls.0.id
}

```

## Argument Reference

The following arguments are supported:

* `acl_ids` - (Optional, ForceNew) The acl ids.
* `acl_name` - (Optional, ForceNew) The ACL Name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Acl IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Acl name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) Resource Group to Which the Number.
* `status` - (Optional, ForceNew) The state of the ACL. Valid values:`Provisioning`,`Available`and`Configuring`.  `Provisioning`: The ACL is being created. `Available`: The ACL is available. `Configuring`: The ACL is being configured.
* `tag` - (Optional, ForceNew) The tag.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Acl names.
* `acls` - A list of Alb Acls. Each element contains the following attributes:
	* `acl_entries` - ACL Entries.
		* `description` - Access Control Entries Note Description Length Is Limited to 1 to 256 Characters, Letters, digital, the Dash (-), a Forward Slash (/), Half a Period (.) and Underscores (_), Support Chinese Characters.
		* `entry` - The resource ID in terraform of Acl.
		* `status` - The status of the ACL entry. Valid values: `Adding` , `Available` and `Removing`. `Adding`: The entry is being added. `Available`: The entry is added and available. `Removing`: The entry is being removed.
	* `acl_id` - Access Control Policy ID.
	* `acl_name` - The ACL Name.
	* `address_ip_version` - Address Protocol Version.
	* `id` - The ID of the Acl.
	* `resource_group_id` - Resource Group to Which the Number.
	* `status` - The state of the ACL. Valid values:`Provisioning` , `Available` and `Configuring`. `Provisioning`: The ACL is being created. `Available`: The ACL is available. `Configuring`: The ACL is being configured.
