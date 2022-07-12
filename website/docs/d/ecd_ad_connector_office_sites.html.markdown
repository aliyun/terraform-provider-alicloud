---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_ad_connector_office_sites"
sidebar_current: "docs-alicloud-datasource-ecd-ad-connector-office-sites"
description: |-
  Provides a list of Ecd Ad Connector Office Sites to the user.
---

# alicloud\_ecd\_ad\_connector\_office\_sites

This data source provides the Ecd Ad Connector Office Sites of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.176.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_ad_connector_office_sites" "ids" {}
output "ecd_ad_connector_office_site_id_1" {
  value = data.alicloud_ecd_ad_connector_office_sites.ids.sites.0.id
}

data "alicloud_ecd_ad_connector_office_sites" "nameRegex" {
  name_regex = "^my-AdConnectorOfficeSite"
}
output "ecd_ad_connector_office_site_id_2" {
  value = data.alicloud_ecd_ad_connector_office_sites.nameRegex.sites.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ad Connector Office Site IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ad Connector Office Site name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The workspace status. Valid values:
  - `REGISTERING`: The workspace is being registered.
  - `REGISTERED`: The workspace is registered.
  - `DEREGISTERING`: The workspace is being deregistered.
  - `DEREGISTERED`: The workspace is deregistered.
  - `ERROR`: The configurations of the workspace are invalid.
  - `NEEDCONFIGTRUST`: The trust relationship needs to be configured.
  - `NEEDCONFIGUSER`: Users need to be configured.
  - `CONFIGTRUSTING`: The trust relationship is being configured.
  - `CONFIGTRUSTFAILED`: The trust relationship fails to be configured.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ad Connector Office Site names.
* `sites` - A list of Ecd Ad Connector Office Sites. Each element contains the following attributes:
	* `ad_connector_office_site_name` - The Name of the ad connector office site.
	* `ad_connectors` - AD Connector Collection of Information.
		* `ad_connector_address` - AD Connector across Zones, Its Connection Addresses.
		* `connector_status` - AD Connector of the State. Possible Values: Creating: in the Creation of. Connecting: Connection. Requires the User to Your Own Ad Configured on the Domain to Which. Running: Run. Expired: If You Are out-of-Date. CONNECT_ERROR: Connection Error.
		* `network_interface_id` - AD Connector Mount of the Card ID.
		* `vswitch_id` - AD Connector in the Network Corresponding to the ID of the VSwitch in.
	* `bandwidth` - The Internet Bandwidth Peak. Possible Values: 0~200. If This Field Is Set to 0, Indicates That There Is No Open Internet Access.
	* `cen_id` - Cloud Enterprise Network Instance Id.
	* `cidr_block` - Workspace Corresponds to the Security Office Network of IPv4 Segment.
	* `create_time` - Workspace Creation Time.
	* `custom_security_group_id` - Security Group ID.
	* `desktop_access_type` - The method that is used to connect the client to cloud desktops.
	* `desktop_vpc_endpoint` - The endpoint that is used to connect to cloud desktops over a VPC.
	* `dns_address` - Enterprise Ad Corresponding DNS Address.
	* `dns_user_name` - The Easy-to-Use DNS Name.
	* `domain_name` - Enterprise of Ad Domain Name.
	* `domain_user_name` - The Domain Administrator's Username.
	* `enable_admin_access` - Whether to Use Cloud Desktop User Empowerment of Local Administrator Permissions.
	* `enable_cross_desktop_access` - Indicates whether the desktop communication feature is enabled for cloud desktops in the same workspace. After the feature is enabled, the cloud desktops in the same workspace can access each other.
	* `enable_internet_access` - Whether the Open Internet Access Function.
	* `file_system_ids` - NAS File System ID.
	* `id` - The ID of the Ad Connector Office Site.
	* `logs` - Registered Log Information.
		* `level` - Log Level. Possible Values: Info: Information Error: Error Warn: Warning.
		* `message` - The Log Details.
		* `step` - Log Information Corresponding to the Step.
		* `time_stamp` - Log Print Time.
	* `mfa_enabled` - Whether to Enable Multi-Factor Authentication MFA.
	* `network_package_id` - The ID of the Internet Access.
	* `office_site_id` - The ID of the Workspace.
	* `office_site_type` - Workspace Account System Type. Possible Values: Simple: Convenient Account. AD_CONNECTOR: Enterprise Ad Account.
	* `sso_enabled` - Whether to Enable Single Sign-on (SSO) for User-Based SSO.
	* `status` - The workspace status.
	* `sub_domain_dns_address` - Sub-Domain DNS Address.
	* `sub_domain_name` - The AD Domain DNS Name.
	* `trust_password` - The AD Trust Password.
	* `vpc_id` - Security Office VPC ID.
	* `vswitch_ids` - The vswitch ids.