---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_dhcp_options_sets"
sidebar_current: "docs-alicloud-datasource-vpc-dhcp-options-sets"
description: |-
  Provides a list of Vpc Dhcp Options Sets to the user.
---

# alicloud\_vpc\_dhcp\_options\_sets

This data source provides the Vpc Dhcp Options Sets of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_dhcp_options_sets" "ids" {
  ids = ["example_value"]
}
output "vpc_dhcp_options_set_id_1" {
  value = data.alicloud_vpc_dhcp_options_sets.ids.sets.0.id
}

data "alicloud_vpc_dhcp_options_sets" "nameRegex" {
  name_regex = "^my-DhcpOptionsSet"
}
output "vpc_dhcp_options_set_id_2" {
  value = data.alicloud_vpc_dhcp_options_sets.nameRegex.sets.0.id
}

data "alicloud_vpc_dhcp_options_sets" "dhcpOptionsSetName" {
  dhcp_options_set_name = "my-DhcpOptionsSet"
}
output "vpc_dhcp_options_set_id_3" {
  value = data.alicloud_vpc_dhcp_options_sets.dhcpOptionsSetName.sets.0.id
}

data "alicloud_vpc_dhcp_options_sets" "domainName" {
  ids         = ["example_value"]
  domain_name = "example.com"
}
output "vpc_dhcp_options_set_id_4" {
  value = data.alicloud_vpc_dhcp_options_sets.domainName.sets.0.id
}

data "alicloud_vpc_dhcp_options_sets" "status" {
  ids    = ["example_value"]
  status = "Available"
}
output "vpc_dhcp_options_set_id_5" {
  value = data.alicloud_vpc_dhcp_options_sets.status.sets.0.id
}

```

## Argument Reference

The following arguments are supported:

* `dhcp_options_set_name` - (Optional, ForceNew) The name of the DHCP options set.The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character.
* `domain_name` - (Optional, ForceNew) The root domain, for example, example.com. After a DHCP options set is associated with a Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS instances in the VPC network.
* `ids` - (Optional, ForceNew, Computed)  A list of Dhcp Options Set IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Dhcp Options Set name.
* `status` - (Optional, ForceNew) The status of the DHCP options set. Valid values: `Available`, `InUse` or `Pending`. `Available`: The DHCP options set is available for use. `InUse`: The DHCP options set is in use. `Pending`: The DHCP options set is being configured.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Dhcp Options Set names.
* `sets` - A list of Vpc Dhcp Options Sets. Each element contains the following attributes:
    * `associate_vpcs` - AssociateVpcs.
        * `associate_status` - The status of the VPC network that is associated with the DHCP options set. Valid values:`InUse` or `Pending`. `InUse`: The VPC network is in use. `Pending`: The VPC network is being configured.
        * `vpc_id` - The ID of the VPC network that is associated with the DHCP options set.
    * `dhcp_options_set_description` - The description of the DHCP options set. The description must be 2 to 256
      characters in length and cannot start with `http://` or `https://`.
    * `dhcp_options_set_id` - The resource ID in terraform of Dhcp Options Set.
    * `dhcp_options_set_name` - The root domain, for example, example.com. After a DHCP options set is associated with a
      Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the
      ECS instances in the VPC network.
    * `domain_name` - The root domain, for example, example.com. After a DHCP options set is associated with a Virtual
      Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS
      instances in the VPC network.
    * `domain_name_servers` - The DNS server IP addresses. Up to four DNS server IP addresses can be specified. IP
      addresses must be separated with commas (,).
    * `id` - The resource ID in terraform of Dhcp Options Set.
    * `owner_id` - The ID of the account to which the DHCP options set belongs.
    * `status` - The status of the DHCP options set. Valid values: `Available`, `InUse` or `Pending`. `Available`: The DHCP options set is available for use. `InUse`: The DHCP options set is in use. `Pending`: The DHCP options set is being configured.
