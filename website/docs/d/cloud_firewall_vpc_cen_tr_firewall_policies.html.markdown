---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewall_policies"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-vpc-cen-tr-firewall-policies"
description: |-
  Provides a list of Cloud Firewall Vpc Cen Tr Firewall Policy owned by an Alibaba Cloud account.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewall_policies

This data source provides Cloud Firewall Vpc Cen Tr Firewall Policy available to the user. [What is Vpc Cen Tr Firewall Policy](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/DescribeTrFirewallV2RoutePolicyList)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_cloud_firewall_vpc_cen_tr_firewall_policies" "default" {
  firewall_id = "xxx"
  ids         = ["xxx:xxx"]
}

output "first_policy_id" {
  value = data.alicloud_cloud_firewall_vpc_cen_tr_firewall_policies.default.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `firewall_id` - (Optional) The ID of the VPC firewall instance.
* `lang` - (Optional) The language type of the received message. Valid values: `zh`, `en`.
* `ids` - (Optional) A list of policy IDs.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `policies` - A list of Cloud Firewall Vpc Cen Tr Firewall Policies. Each element contains the following attributes:
  * `id` - The ID of the policy, formatted as `<firewall_id>:<tr_firewall_route_policy_id>`.
  * `firewall_id` - The ID of the VPC firewall instance.
  * `tr_firewall_route_policy_id` - The ID of the firewall routing policy.
  * `policy_name` - The name of the traffic redirection template.
  * `policy_description` - The description of the traffic redirection.
  * `policy_type` - The drainage scenario type. Valid values: `fullmesh`, `one_to_one`, `end_to_end`.
  * `status` - The policy state. Valid values: `creating`, `deleting`, `opening`, `opened`, `closing`, `closed`.
  * `should_recover` - Whether to restore the drainage configuration. Valid values: `true`, `false`.
  * `src_candidate_list` - The list of primary traffic redirection instances.
    * `candidate_id` - The ID of the traffic redirection instance.
    * `candidate_type` - The type of the traffic redirection instance.
  * `dest_candidate_list` - The list of secondary traffic redirection instances.
    * `candidate_id` - The ID of the traffic redirection instance.
    * `candidate_type` - The type of the traffic redirection instance.
