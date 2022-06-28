---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_simple_office_sites"
sidebar_current: "docs-alicloud-datasource-ecd-simple-office-sites"
description: |-
  Provides a list of Ecd Simple Office Sites to the user.
---

# alicloud\_ecd\_simple\_office\_sites

This data source provides the Ecd Simple Office Sites of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_simple_office_sites" "default" {
  ids    = ["example_id"]
  status = "REGISTERED"
}
output "desktop_access_type" {
  value = data.alicloud_ecd_simple_office_sites.default.sites.0.desktop_access_type
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Simple Office Site IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Simple Office Site name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Workspace State. Valid values: `REGISTERED`,`REGISTERING`. 

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Simple Office Site names.
* `sites` - A list of Ecd Simple Office Sites. Each element contains the following attributes:
	* `bandwidth` - (Deprecated from 1.142.0) The Internet Bandwidth Peak. It has been deprecated from version 1.142.0 and can be found in the new datasource alicloud_ecd_network_packages.
	* `cen_id` - Cloud Enterprise Network Instance Id.
	* `cidr_block` - Workspace Corresponds to the Security Office Network of IPv4 Segment.
	* `create_time` - Workspace Creation Time.
	* `custom_security_group_id` - Security Group ID.
	* `desktop_access_type` - Connect to the Cloud Desktop Allows the Use of the Access Mode of. Possible Values: the Internet: Only Allows the Client to Public Cloud Desktop. Virtual Private Cloud (VPC): Only Allows in the Virtual Private Cloud (VPC) in the Client to Connect to the Cloud Desktop. Any: Not by Way of Limitation. Use Client to Connect to the Cloud Desktop When It Is Possible to Choose the Connection.
	* `desktop_vpc_endpoint` - The Desktop Vpc Endpoint.
	* `dns_address` - Enterprise Ad Corresponding DNS Address.
	* `dns_user_name` - Easy-to-Use DNS Name.
	* `domain_name` - Enterprise of Ad Domain Name.
	* `domain_password` - Domain of the User Who Will Administer This Target Application Password.
	* `domain_user_name` - The Domain Administrator's Username.
	* `enable_admin_access` - Whether to Use Cloud Desktop User Empowerment of Local Administrator Permissions.
	* `enable_cross_desktop_access` - Enable Cross-Desktop Access.
	* `enable_internet_access` - (Deprecated from 1.142.0) Whether the Open Internet Access Function.	
	* `file_system_ids` - NAS File System ID.
	* `id` - The ID of the Simple Office Site.
	* `mfa_enabled` - Whether to Enable Multi-Factor Authentication MFA.
	* `network_package_id` - Internet Access ID.
	* `office_site_id` - The Workspace ID.
	* `office_site_type` - Workspace Account System Type. Possible Values: Simple: Convenient Account. AD_CONNECTOR: Enterprise Ad Account.
	* `simple_office_site_name` - The simple office site name.
	* `sso_enabled` - Whether to Enable Single Sign-on (SSO) for User-Based SSO.
	* `sso_status` - Whether to Enable Single Sign-on (SSO) for User-Based SSO.
	* `status` - Workspace State. Possible Values: Registering: Registered in the Registered: Registered.
	* `sub_dns_address` - AD Subdomain of the DNS Address.
	* `sub_domain_name` - AD Domain DNS Name.
	* `trust_password` - AD Trust Password.
	* `users` - AD User Name Array.
	* `vpc_id` - Security Office VPC ID.
	* `vswitch_ids` - The vswitch ids.
