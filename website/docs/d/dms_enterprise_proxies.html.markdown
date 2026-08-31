---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_proxies"
sidebar_current: "docs-alicloud-datasource-dms-enterprise-proxies"
description: |-
  Provides a list of Dms Enterprise Proxies to the user.
---

# alicloud\_dms\_enterprise\_proxies

This data source provides the Dms Enterprise Proxies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.188.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dms_enterprise_proxies" "ids" {}
output "dms_enterprise_proxy_id_1" {
  value = data.alicloud_dms_enterprise_proxies.ids.proxies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Proxy IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tid` - (Optional, ForceNew) The ID of the tenant.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `proxies` - A list of Dms Enterprise Proxies. Each element contains the following attributes:
	* `creator_id` - The ID of the user who enabled the secure access proxy feature.
	* `creator_name` - The nickname of the user who enabled the secure access proxy feature.
	* `https_port` - The port that was used by HTTPS clients to connect to the database instance.
	* `id` - The ID of the Proxy.
	* `instance_id` - The ID of the database instance.
	* `private_enable` - Indicates whether the internal endpoint is enabled. Default value: true.
	* `private_host` - The internal endpoint.
	* `protocol_port` - Database protocol connection port number.
	* `protocol_type` - Database protocol type, for example, MYSQL.
	* `proxy_id` - The ID of the secure access proxy.
	* `public_enable` - Indicates whether the public endpoint is enabled.
	* `public_host` - The public endpoint. A public endpoint is returned no matter whether the public endpoint is enabled or disabled. **Note:** When the public network address is in the **true** state, the returned public network address is a valid address with DNS resolution capability. When the public address is in the **false** state, the returned Public address is an invalid address without DNS resolution.