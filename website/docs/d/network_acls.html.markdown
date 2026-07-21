---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_network_acls"
sidebar_current: "docs-alicloud-datasource-network-acls"
description: |-
  Provides a list of Network Acls to the user.
---

# alicloud\_network\_acls

This data source provides the Network Acls of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_network_acls" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_network_acl_id" {
  value = data.alicloud_network_acls.example.acls.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Network Acl ID.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Network Acl name.
* `network_acl_name` - (Optional, ForceNew) The name of the network ACL.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The state of the network ACL. Valid values: `Available` and `Modifying`.
* `vpc_id` - (Optional, ForceNew) The ID of the associated VPC.
* `resource_type` - (Optional, ForceNew) The type of the associated resource. Valid values `VSwitch`. `resource_type` and `resource_id` need to be specified at the same time to take effect.
* `resource_id` - (Optional, ForceNew) The ID of the associated resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Network Acl names.
* `acls` - A list of Network Acls. Each element contains the following attributes:
	* `description` - Description of network ACL information.
	* `egress_acl_entries` - Output direction rule information.
		* `description` - Give the description information of the direction rule.
		* `destination_cidr_ip` - The destination address segment.
		* `network_acl_entry_name` - The name of the entry for the direction rule.
		* `policy` - The  authorization policy.
		* `port` - Destination port range.
		* `protocol` - Transport  layer protocol.
	* `id` - The ID of the Network Acl.
	* `ingress_acl_entries` - Entry direction rule information.
		* `source_cidr_ip` - The source address field.
		* `description` - Description of the entry direction rule.
		* `network_acl_entry_name` - The name of the entry direction rule entry.
		* `policy` - The authorization policy.
		* `port` - Source port range.
		* `protocol` - Transport layer protocol.
	* `network_acl_id` - The first ID of the resource.
	* `network_acl_name` - The name of the network ACL.
	* `resources` - The associated resource.
		* `resource_id` - The ID of the associated resource.
		* `resource_type` - The type of the associated resource.
		* `status` - The state of the associated resource.
	* `status` - The state of the network ACL.
	* `vpc_id` - The ID of the associated VPC.
