---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_tls_cipher_policies"
sidebar_current: "docs-alicloud-datasource-slb-tls-cipher-policies"
description: |-
  Provides a list of Slb Tls Cipher Policies to the user.
---

# alicloud_slb_tls_cipher_policies

This data source provides the Slb Tls Cipher Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.135.0.

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

* `ids` - (Optional) A list of Tls Cipher Policy IDs.
* `include_listener` - (Optional) The include listener.
* `tls_cipher_policy_name` - (Optional) TLS policy name. Length is from 2 to 128, or in both the English and Chinese characters must be with an uppercase/lowercase letter or a Chinese character and the beginning, may contain numbers, in dot `.`, underscore `_` or dash `-`.
* `name_regex` - (Optional) A regex string to filter results by Tls Cipher Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional) TLS policy instance state. Valid values: `configuring`, `normal`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Tls Cipher Policy names.
* `policies` - A list of Slb Tls Cipher Policies. Each element contains the following attributes:

### `policies`

* `ciphers` - The encryption algorithms supported. It depends on the value of `tls_versions`.
* `create_time` - The creation time timestamp.
* `id` - The ID of the Tls Cipher Policy.
* `relate_listeners` - Array of Relate Listeners.
* `status` - TLS policy instance state.
* `tls_cipher_policy_id` - The ID of TLS cipher policy.
* `tls_cipher_policy_name` - TLS policy name. Length is from 2 to 128, or in both the English and Chinese characters must be with an uppercase/lowercase letter or a Chinese character and the beginning, may contain numbers, in dot `.`, underscore `_` or dash `-`.
* `tls_versions` - The version of TLS protocol.

### `policies-relate_listeners`

* `load_balancer_id` - The ID of SLB instance.
* `port` - Listening port. Valid value: 1 to 65535.
* `protocol` - Snooping protocols. Valid values: `TCP`, `UDP`, `HTTP`, or `HTTPS`.
