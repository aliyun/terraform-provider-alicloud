---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_tls_cipher_policies"
sidebar_current: "docs-alicloud-datasource-slb-tls-cipher-policies"
description: |-
  Provides a list of Slb Tls Cipher Policies to the user.
---

# alicloud\_slb\_tls\_cipher\_policies

This data source provides the Slb Tls Cipher Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_slb_tls_cipher_policies" "ids" {
  ids = ["example_value-1", "example_value-2"]
}
output "slb_tls_cipher_policy_id_1" {
  value = data.alicloud_slb_tls_cipher_policies.ids.policies.0.id
}

data "alicloud_slb_tls_cipher_policies" "nameRegex" {
  name_regex = "^My-TlsCipherPolicy"
}
output "slb_tls_cipher_policy_id_2" {
  value = data.alicloud_slb_tls_cipher_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Tls Cipher Policy IDs.
* `include_listener` - (Optional, ForceNew) The include listener.
* `tls_cipher_policy_name` - (Optional, ForceNew) TLS policy name. Length is from 2 to 128, or in both the English and Chinese characters must be with an uppercase/lowercase letter or a Chinese character and the beginning, may contain numbers, in dot `.`, underscore `_` or dash `-`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Tls Cipher Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) TLS policy instance state. Valid values: `configuring`, `normal`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Tls Cipher Policy names.
* `policies` - A list of Slb Tls Cipher Policies. Each element contains the following attributes:
	* `id` - The ID of the Tls Cipher Policy.
	* `tls_versions` - The version of TLS protocol. 
	* `ciphers` - The encryption algorithms supported. It depends on the value of `tls_versions`.
	* `create_time` - The creation time timestamp.
	* `relate_listeners` - Array of Relate Listeners.
		* `load_balancer_id` - The ID of SLB instance.
		* `port` - Listening port. Valid value: 1 to 65535.
		* `protocol` - Snooping protocols. Valid values: `TCP`, `UDP`, `HTTP`, or `HTTPS`.
	* `status` - TLS policy instance state.
	* `tls_cipher_policy_id` - The ID of TLS cipher policy.
	* `tls_cipher_policy_name` - TLS policy name. Length is from 2 to 128, or in both the English and Chinese characters must be with an uppercase/lowercase letter or a Chinese character and the beginning, may contain numbers, in dot `.`, underscore `_` or dash `-`.
