---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_ipa_domains"
sidebar_current: "docs-alicloud-datasource-dcdn-ipa-domains"
description: |-
  Provides a list of Dcdn Ipa Domains to the user.
---

# alicloud\_dcdn\_ipa\_domains

This data source provides the Dcdn Ipa Domains of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.158.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dcdn_ipa_domains" "ids" {
  domain_name = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "dcdn_ipa_domain_id_1" {
  value = data.alicloud_dcdn_ipa_domains.ids.domains.0.id
}

data "alicloud_dcdn_ipa_domains" "status" {
  status = "online"
}
output "dcdn_ipa_domain_id_2" {
  value = data.alicloud_dcdn_ipa_domains.status.domains.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `domain_name` - (Optional, ForceNew) The name of the Domain.
* `ids` - (Optional, ForceNew, Computed)  A list of Ipa Domain IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the accelerated domain name. Valid values: `check_failed`, `checking`, `configure_failed`, `configuring`, `offline`, `online`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Dcdn Ipa Domain names.
* `domains` - A list of Dcdn Ipa Domains. Each element contains the following attributes:
	* `cert_name` - CertName.
	* `cname` - The CNAME assigned to the domain name.
	* `create_time` - The time when the accelerated domain name was created.
	* `description` - The description.
	* `domain_name` - The accelerated domain names.
	* `id` - The ID of the Ipa Domain.
	* `resource_group_id` - The ID of the resource group.
	* `scope` - The accelerated region.
	* `sources` - The information about the origin server.
		* `priority` - The priority of the origin server if multiple origin servers are specified.
		* `type` - The type of the origin server.
		* `weight` - The weight of the origin server if multiple origin servers are specified.
		* `content` - The address of the origin server.
		* `port` - The custom port.
	* `ssl_protocol` - Indicates whether the Security Socket Layer (SSL) certificate is enabled.
	* `ssl_pub` - Indicates the public key of the certificate if the HTTPS protocol is enabled.
	* `status` - The status of the accelerated domain name.