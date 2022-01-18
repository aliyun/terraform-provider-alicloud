---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_dhcp_options_set_attachment"
sidebar_current: "docs-alicloud-resource-vpc-dhcp-options-set-attachment"
description: |-
  Provides a Alicloud VPC Dhcp Options Set Attachment resource.
---

# alicloud\_vpc\_dhcp\_options\_set\_attachment

Provides a VPC Dhcp Options Set Attachment resource.

For information about VPC Dhcp Options Set and how to use it, see [What is Dhcp Options Set](https://www.alibabacloud.com/help/doc-detail/174112.htm).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_vpc" "example" {
  vpc_name   = "test"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc_dhcp_options_set" "example" {
  dhcp_options_set_name        = "example_value"
  dhcp_options_set_description = "example_value"
  domain_name                  = "example.com"
  domain_name_servers          = "100.100.2.136"
}
resource "alicloud_vpc_dhcp_options_set_attachment" "example" {
  vpc_id              = alicloud_vpc.example.id
  dhcp_options_set_id = alicloud_vpc_dhcp_options_set.example.id
}

```

## Argument Reference

The following arguments are supported:

* `dhcp_options_set_id` - (Required) The ID of the DHCP options set.
* `vpc_id` - (Required) The ID of the VPC network that is to be associated with the DHCP options set..
* `dry_run` - (Optional) Specifies whether to precheck this request only. Default values: `false`. Valid values:
  * `true` - Runs a precheck without associating the DHCP options set with the VPC network. The system checks whether your AccessKey pair is valid, whether the RAM user is authorized, and whether required parameters are specified. An error message is returned if the request fails the precheck. If the request passes the precheck, the `DryRunOperation` error code is returned.
  * `false` - Runs a precheck and returns a 2XX HTTP status code. After the request passes the precheck, the DHCP options set is associated with the VPC network. 

## Attributes Reference

The following attributes are exported:

* `id` - The Disk Attachment ID and it formats as `<vpc_id>:<dhcp_options_set_id>`.
* `status` -The status of the VPC network that is associated with the DHCP options set.  Valid values: `InUse` or `Pending`. 

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Dhcp Options Set.
* `delete` - (Defaults to 1 mins) Used when delete the Dhcp Options Set.

## Import

VPC Dhcp Options Set Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_dhcp_options_set_attachment.example <id>
```
