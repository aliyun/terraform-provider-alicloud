---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_dhcp_options_set"
sidebar_current: "docs-alicloud-resource-vpc-dhcp-options-set"
description: |-
  Provides a Alicloud VPC Dhcp Options Set resource.
---

# alicloud\_vpc\_dhcp\_options\_set

Provides a VPC Dhcp Options Set resource.

For information about VPC Dhcp Options Set and how to use it, see [What is Dhcp Options Set](https://www.alibabacloud.com/help/doc-detail/174112.htm).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_dhcp_options_set" "example" {
  dhcp_options_set_name        = "example_value"
  dhcp_options_set_description = "example_value"
  domain_name                  = "example.com"
  domain_name_servers          = "100.100.2.136"
}

```

## Argument Reference

The following arguments are supported:

* `dhcp_options_set_description` - (Optional) The description of the DHCP options set. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `dhcp_options_set_name` - (Optional) The name of the DHCP options set. The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character.
* `domain_name` - (Optional) The root domain, for example, example.com. After a DHCP options set is associated with a Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS instances in the VPC network.
* `domain_name_servers` - (Optional) The DNS server IP addresses. Up to four DNS server IP addresses can be specified. IP addresses must be separated with commas (,).Before you specify any DNS server IP address, all ECS instances in the associated VPC network use the IP addresses of the Alibaba Cloud DNS servers, which are `100.100.2.136` and `100.100.2.138`.
* `associate_vpcs` - (Optional, Deprecated) AssociateVpcs. Number of VPCs that can be associated with each DHCP options set is 10. Field `associate_vpcs` has been deprecated from provider version 1.153.0. It will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc.
  * `vpc_id` - (Optional) The ID of the VPC network that is associated with the DHCP options set.
* `dry_run` - (Optional) Specifies whether to precheck this request only. Valid values: `true` or `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dhcp Options Set.
* `status` - The status of the DHCP options set. Valid values: `Available`, `InUse` or `Pending`. `Available`: The DHCP options set is available for use. `InUse`: The DHCP options set is in use. `Pending`: The DHCP options set is being configured.
* `owner_id` - The ID of the account to which the DHCP options set belongs.
* `associate_vpcs` - AssociateVpcs.
  * `vpc_id` - The ID of the VPC network that is associated with the DHCP options set.
  * `associate_status` - The status of the VPC network that is associated with the DHCP options set. Valid values:`InUse` or `Pending`. `InUse`: The VPC network is in use. `Pending`: The VPC network is being configured.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Dhcp Options Set.
* `delete` - (Defaults to 2 mins) Used when delete the Dhcp Options Set.
* `update` - (Defaults to 2 mins) Used when update the Dhcp Options Set.

## Import

VPC Dhcp Options Set can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_dhcp_options_set.example <id>
```
