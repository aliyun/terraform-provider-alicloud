---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_system_security_policies"
sidebar_current: "docs-alicloud-datasource-alb-system-security-policies"
description: |-
  Provides a list of ALB System Security Policies to the user.
---

# alicloud\_alb\_system\_security\_policies

This data source provides the ALB System Security Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.183.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_system_security_policies" "defaults" {
  ids = ["tls_cipher_policy_1_0"]
}

output "alb_system_security_policy_id_1" {
  value = data.alicloud_alb_system_security_policies.defaults.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of System Security Policy IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of System Security Policy IDs.
* `policies` - A list of ALB Security Policies. Each element contains the following attributes:
	* `id` - The ID of the Security Policy.
	* `security_policy_id` - The first ID of the resource.
	* `tls_versions` - The TLS protocol versions are supported. Valid values: TLSv1.0, TLSv1.1, TLSv1.2 and TLSv1.3.
	* `ciphers` - The supported cipher suites, which are determined by the TLS protocol version.
