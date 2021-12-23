---
subcategory: "Anti-DDoS Pro"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_domain_resources"
sidebar_current: "docs-alicloud-datasource-ddoscoo-domain-resources"
description: |-
  Provides a list of Ddoscoo Domain Resources to the user.
---

# alicloud\_ddoscoo\_domain\_resources

This data source provides the Ddoscoo Domain Resources of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.123.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ddoscoo_domain_resources" "example" {
  ids = ["tftestacc1234.abc"]
}

output "first_ddoscoo_domain_resource_id" {
  value = data.alicloud_ddoscoo_domain_resources.example.resources.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Domain Resource IDs.
* `instance_ids` - (Optional, ForceNew) A ID list of Ddoscoo instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `query_domain_pattern` - (Optional, ForceNew) Match the pattern.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `resources` - A list of Ddoscoo Domain Resources. Each element contains the following attributes:
	* `black_list` - The IP addresses in the blacklist for the domain name.
	* `cc_enabled` - Whether frequency control guard (CC guard) is enabled. Values: `True`: Opened, `False`: Not enabled.
	* `cc_rule_enabled` - Whether custom frequency control guard (CC guard) is enabled. Values: `True`: Opened, `False`: Not enabled.
	* `cc_template` - The mode of the Frequency Control policy.
	* `cert_name` - The name of the certificate.
	* `id` - The ID of the Domain Resource.
	* `domain` - The domain name of the website that you want to add to the instance.
	* `http2_enable` - Whether Http2.0 is enabled.
	* `https_ext` - The advanced HTTPS settings.
	* `instance_ids` - A list ID of instance that you want to associate.
	* `policy_mode` - The type of backload algorithm.
	* `proxy_enabled` - Whether the website service forwarding rules have been turned on.
	* `proxy_types` - Protocol type and port number information.
		* `proxy_ports` - The forwarding port.
		* `proxy_type` - Protocol type.
	* `real_servers` - Server address information of the source station.
	* `rs_type` - Server address type.
	* `ssl_ciphers` - The type of the cipher suite.
	* `ssl_protocols` - The version of the TLS protocol.
	* `white_list` - The IP addresses in the whitelist for the domain name.
