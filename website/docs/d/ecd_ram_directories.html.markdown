---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_ram_directories"
sidebar_current: "docs-alicloud-datasource-ecd-ram-directories"
description: |-
  Provides a list of Ecd Ram Directories to the user.
---

# alicloud\_ecd\_ram\_directories

This data source provides the Ecd Ram Directories of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.174.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_ram_directories" "ids" {
  ids = ["example_id"]
}
output "ecd_ram_directory_id_1" {
  value = data.alicloud_ecd_ram_directories.ids.directories.0.id
}

data "alicloud_ecd_ram_directories" "nameRegex" {
  name_regex = "^my-RamDirectory"
}
output "ecd_ram_directory_id_2" {
  value = data.alicloud_ecd_ram_directories.nameRegex.directories.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ram Directory IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ram Directory name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of directory. Valid values: `REGISTERING`, `REGISTERED`, `DEREGISTERING`, `NEEDCONFIGTRUST`, `CONFIGTRUSTFAILED`, `DEREGISTERED`, `ERROR`, `CONFIGTRUSTING`, `NEEDCONFIGUSER`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ram Directory names.
* `directories` - A list of Ecd Ram Directories. Each element contains the following attributes:
	* `ad_connectors` - The AD connectors.
		* `ad_connector_address` - The address of AD connector.
		* `connector_status` - The status of connector.
		* `network_interface_id` - The ID of the network interface.
		* `vswitch_id` - The ID of VSwitch.
	* `create_time` - The CreateTime of resource.
	* `custom_security_group_id` - The id of the custom security group.
	* `desktop_access_type` - The desktop access type.
	* `desktop_vpc_endpoint` - The desktop vpc endpoint.
	* `directory_type` - The directory type.
	* `dns_address` - The address of DNSAddress.
	* `dns_user_name` - The username of DNS.
	* `domain_name` - The name of the domain.
	* `domain_password` - The domain password.
	* `domain_user_name` - The username of the domain.
	* `enable_admin_access` - Whether to enable admin access.
	* `enable_cross_desktop_access` - Whether to enable cross desktop access.
	* `enable_internet_access` - Whether enable internet access.
	* `file_system_ids` - The ids of filesystem.
	* `id` - The ID of the Ram Directory.
	* `logs` - The register log information.
		* `level` - The level of log.
		* `message` - The message of log.
		* `step` - The step of log.
		* `time_stamp` - The time stamp of log.
	* `mfa_enabled` - Whether to enable MFA.
	* `ram_directory_id` - The ID of ram directory.
	* `ram_directory_name` - The name of directory.
	* `sso_enabled` - Whether to enable SSO.
	* `status` - The status of directory.
	* `sub_dns_address` - The address of sub DNS.
	* `sub_domain_name` - The Name of the sub-domain.
	* `trust_password` - The trust password.
	* `vpc_id` - The ID of the vpc.
	* `vswitch_ids` - List of VSwitch IDs in the directory.