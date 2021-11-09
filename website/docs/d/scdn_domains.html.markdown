---
subcategory: "SCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_scdn_domains"
sidebar_current: "docs-alicloud-datasource-scdn-domains"
description: |-
  Provides a list of Scdn Domains to the user.
---

# alicloud\_scdn\_domains

This data source provides the Scdn Domains of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_scdn_domains" "nameRegex" {
  name_regex = "^my-Domain"
}
output "scdn_domain_id" {
  value = data.alicloud_scdn_domains.nameRegex.domains.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Domain IDs. Its element value is same as Domain Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Domain name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The Resource Group ID.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: "check_failed", "checking", "configure_failed", "configuring", "offline", "online".

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Domain names.
* `domains` - A list of Scdn Domains. Each element contains the following attributes:
	* `cert_infos` - Certificate Information.
		* `cert_name` - If You Enable HTTPS Here Certificate Name.
		* `ssl_pub` - If You Enable HTTPS Here Key.
		* `cert_type` - Certificate Type. Value Range: Upload: Certificate. CAS: Certificate Authority Certificate. Free: Free Certificate.
		* `ssl_protocol` - Whether to Enable SSL Certificate. Valid Values: on, off.
	* `cname` - In Order to Link the CDN Domain Name to Generate a CNAME Domain Name, in the Domain Name Resolution Service Provider at the Acceleration Domain Name CNAME Resolution to the Domain.
	* `create_time` - Creation Time.
	* `description` - Review the Reason for the Failure Is Displayed.
	* `domain_name` - Your Domain Name.
	* `gmt_modified` - Last Modified Date.
	* `id` - The ID of the Domain. Its value is same as Queue Name.
	* `resource_group_id` - The Resource Group ID.
	* `sources` - the Origin Server Information.
		* `type` - the Origin Server Type. Valid Values: Ipaddr: IP Source Station Domain: the Domain Name, See Extra Domain Quota OSS: OSS Bucket as a Source Station.
		* `content` - The Back-to-Source Address.
		* `enabled` - State.
		* `port` - Port.
		* `priority` - Priority.
	* `status` - The status of the resource.
