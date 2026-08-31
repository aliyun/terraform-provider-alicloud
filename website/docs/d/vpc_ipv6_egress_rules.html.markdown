---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_egress_rules"
sidebar_current: "docs-alicloud-datasource-vpc-ipv6-egress-rules"
description: |-
  Provides a list of Vpc Ipv6 Egress Rules to the user.
---

# alicloud\_vpc\_ipv6\_egress\_rules

This data source provides the Vpc Ipv6 Egress Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_ipv6_egress_rules" "ids" {
  ipv6_gateway_id = "example_value"
  ids             = ["example_value-1", "example_value-2"]
}
output "vpc_ipv6_egress_rule_id_1" {
  value = data.alicloud_vpc_ipv6_egress_rules.ids.rules.0.id
}

data "alicloud_vpc_ipv6_egress_rules" "nameRegex" {
  ipv6_gateway_id = "example_value"
  name_regex      = "^my-Ipv6EgressRule"
}
output "vpc_ipv6_egress_rule_id_2" {
  value = data.alicloud_vpc_ipv6_egress_rules.nameRegex.rules.0.id
}

data "alicloud_vpc_ipv6_egress_rules" "status" {
  ipv6_gateway_id = "example_value"
  status          = "Available"
}
output "vpc_ipv6_egress_rule_id_3" {
  value = data.alicloud_vpc_ipv6_egress_rules.status.rules.0.id
}

data "alicloud_vpc_ipv6_egress_rules" "ipv6EgressRuleName" {
  ipv6_gateway_id       = "example_value"
  ipv6_egress_rule_name = "example_value"
}
output "vpc_ipv6_egress_rule_id_4" {
  value = data.alicloud_vpc_ipv6_egress_rules.ipv6EgressRuleName.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ipv6 Egress Rule IDs.
* `instance_id` - (Optional, ForceNew) The ID of the instance that is associated with the IPv6 address to which the egress-only rule is applied.
* `ipv6_egress_rule_name` - (Optional, ForceNew) The name of the resource.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ipv6 Egress Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Deleting`, `Pending`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipv6 Egress Rule names.
* `rules` - A list of Vpc Ipv6 Egress Rules. Each element contains the following attributes:
	* `description` - The description of the egress-only rule.
	* `id` - The ID of the Ipv6 Egress Rule. The value formats as `<ipv6_gateway_id>:<ipv6_egress_rule_id>`.
	* `instance_id` - The ID of the instance to which the egress-only rule is applied.
	* `instance_type` - The type of the instance to which the egress-only rule is applied.
	* `ipv6_egress_rule_id` - The first ID of the resource.
	* `ipv6_egress_rule_name` - The name of the resource.
	* `ipv6_gateway_id` - The ID of the IPv6 gateway.
	* `status` - The status of the resource. Valid values: `Available`, `Pending` and `Deleting`.
