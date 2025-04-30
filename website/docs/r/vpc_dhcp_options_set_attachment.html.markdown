---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_dhcp_options_set_attachment"
sidebar_current: "docs-alicloud-resource-vpc-dhcp-options-set-attachment"
description: |-
  Provides a Alicloud VPC Dhcp Options Set Attachment resource.
---

# alicloud_vpc_dhcp_options_set_attachment

Provides a VPC Dhcp Options Set Attachment resource.

For information about VPC Dhcp Options Set and how to use it, see [What is Dhcp Options Set](https://www.alibabacloud.com/help/doc-detail/174112.htm).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_dhcp_options_set_attachment&exampleId=1654693c-23b8-043a-cdbf-0c06993e7b4e0b3ef90e&activeTab=example&spm=docs.r.vpc_dhcp_options_set_attachment.0.1654693c23&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vpc_dhcp_options_set" "example" {
  dhcp_options_set_name        = var.name
  dhcp_options_set_description = var.name
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

* `dhcp_options_set_id` - (Required, ForceNew) The ID of the DHCP options set.
* `vpc_id` - (Required, ForceNew) The ID of the VPC network that is to be associated with the DHCP options set..
* `dry_run` - (Optional) Specifies whether to precheck this request only. Default values: `false`. Valid values:
  * `true` - Runs a precheck without associating the DHCP options set with the VPC network. The system checks whether your AccessKey pair is valid, whether the RAM user is authorized, and whether required parameters are specified. An error message is returned if the request fails the precheck. If the request passes the precheck, the `DryRunOperation` error code is returned.
  * `false` - Runs a precheck and returns a 2XX HTTP status code. After the request passes the precheck, the DHCP options set is associated with the VPC network. 

## Attributes Reference

The following attributes are exported:

* `id` - The Disk Attachment ID and it formats as `<vpc_id>:<dhcp_options_set_id>`.
* `status` -The status of the VPC network that is associated with the DHCP options set.  Valid values: `InUse` or `Pending`. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Dhcp Options Set.
* `delete` - (Defaults to 1 mins) Used when delete the Dhcp Options Set.

## Import

VPC Dhcp Options Set Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_dhcp_options_set_attachment.example <id>
```
