---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_security_policies"
sidebar_current: "docs-alicloud-datasource-nlb-security-policies"
description: |-
  Provides a list of Nlb Security Policies to the user.
---

# alicloud\_nlb\_security\_policies

This data source provides the Nlb Security Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_security_policies" "ids" {}
output "nlb_security_policy_id_1" {
  value = data.alicloud_nlb_security_policies.ids.policies.0.id
}

data "alicloud_nlb_security_policies" "nameRegex" {
  name_regex = "^my-SecurityPolicy"
}
output "nlb_security_policy_id_2" {
  value = data.alicloud_nlb_security_policies.nameRegex.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Security Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Security Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `security_policy_names` - (Optional, ForceNew) The names of the TLS security policies.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Configuring`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Security Policy names.
* `policies` - A list of Nlb Security Policies. Each element contains the following attributes:
	* `ciphers` - The supported cipher suites, which are determined by the TLS protocol version.
	* `status` - The status of the resource.
	* `tls_versions` - The TLS protocol versions that are supported.
	* `resource_group_id` - The ID of the resource group.
	* `security_policy_name` - The name of the TLS security policy.
	* `tags` - A mapping of tags to assign to the resource.
	* `id` - The id of the TLS security policy.