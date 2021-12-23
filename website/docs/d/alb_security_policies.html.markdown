---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_security_policies"
sidebar_current: "docs-alicloud-datasource-alb-security-policies"
description: |-
  Provides a list of Alb Security Policies to the user.
---

# alicloud\_alb\_security\_policies

This data source provides the Alb Security Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_security_policies" "ids" {}
output "alb_security_policy_id_1" {
  value = data.alicloud_alb_security_policies.ids.policies.0.id
}

data "alicloud_alb_security_policies" "nameRegex" {
  name_regex = "^my-SecurityPolicy"
}
output "alb_security_policy_id_2" {
  value = data.alicloud_alb_security_policies.nameRegex.policies.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Security Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Security Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `security_policy_ids` - (Optional, ForceNew) The security policy ids.
* `security_policy_name` - (Optional, ForceNew) The name of the resource.
* `status` - (Optional, Computed) The status of the resource. Valid values : `Available`, `Configuring`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Security Policy names.
* `policies` - A list of Alb Security Policies. Each element contains the following attributes:
	* `id` - The ID of the Security Policy.
	* `resource_group_id` - The ID of the resource group.
	* `security_policy_id` - The first ID of the resource.
	* `security_policy_name` - The name of the resource. The name must be 2 to 128 characters in length and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-).
	* `status` - The status of the resource.
	* `tls_versions` - The TLS protocol versions that are supported. Valid values: TLSv1.0, TLSv1.1, TLSv1.2 and TLSv1.3.
	* `ciphers` - The supported cipher suites, which are determined by the TLS protocol version.
