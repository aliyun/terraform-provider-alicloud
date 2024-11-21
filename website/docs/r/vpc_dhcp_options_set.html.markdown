---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_dhcp_options_set"
description: |-
  Provides a Alicloud VPC Dhcp Options Set resource.
---

# alicloud_vpc_dhcp_options_set

Provides a VPC Dhcp Options Set resource. DHCP option set.

For information about VPC Dhcp Options Set and how to use it, see [What is Dhcp Options Set](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/dhcp-options-sets-overview).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_dhcp_options_set&exampleId=1647740c-87ae-8291-0b18-cb346be3cbd09fc7e4b2&activeTab=example&spm=docs.r.vpc_dhcp_options_set.0.1647740c87&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

variable "domain" {
  default = "terraform-example.com"
}

resource "alicloud_vpc_dhcp_options_set" "example" {
  dhcp_options_set_name        = var.name
  dhcp_options_set_description = var.name
  domain_name                  = var.domain
  domain_name_servers          = "100.100.2.136"
}
```

## Argument Reference

The following arguments are supported:
* `associate_vpcs` - (Optional, Computed, Deprecated since v1.153.0) Field 'associate_vpcs' has been deprecated from provider version 1.153.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc. See [`associate_vpcs`](#associate_vpcs) below.
* `dhcp_options_set_description` - (Optional) The description can be blank or contain 1 to 256 characters. It must start with a letter or Chinese character but cannot start with http:// or https://.
* `dhcp_options_set_name` - (Optional) The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character.
* `domain_name` - (Optional) The root domain, for example, example.com. After a DHCP options set is associated with a Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS instances in the VPC network.
* `domain_name_servers` - (Optional) The DNS server IP addresses. Up to four DNS server IP addresses can be specified. IP addresses must be separated with commas (,).Before you specify any DNS server IP address, all ECS instances in the associated VPC network use the IP addresses of the Alibaba Cloud DNS servers, which are 100.100.2.136 and 100.100.2.138.
* `dry_run` - (Optional) Whether to PreCheck only this request, value:
  - **true**: sends a check request and does not delete the DHCP option set. Check items include whether required parameters are filled in, request format, and restrictions. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - **false** (default): Sends a normal request and directly deletes the DHCP option set after checking.
* `ipv6_lease_time` - (Optional, Computed, Available since v1.207.0) The lease time of the IPv6 DHCP option set.When the lease time is set to hours: Unit: h. Value range: 24h ~ 1176h,87600h ~ 175200h. Default value: 87600h.When the lease time is set to day: Unit: d. Value range: 1d ~ 49d,3650d ~ 7300d. Default value: 3650d.
* `lease_time` - (Optional, Computed, Available since v1.207.0) The lease time of the IPv4 DHCP option set.When the lease time is set to hours: Unit: h. Value range: 24h ~ 1176h,87600h ~ 175200h. Default value: 87600h.When the lease time is set to day: Unit: d. Value range: 1d ~ 49d,3650d ~ 7300d. Default value: 3650d.
* `resource_group_id` - (Optional, Computed, Available since v1.207.0) The ID of the resource group.
* `tags` - (Optional, Map, Available since v1.207.0) Tags of the current resource.

### `associate_vpcs`

The associate_vpcs supports the following:
* `vpc_id` - (Required) The ID of the VPC network that is associated with the DHCP options set.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `associate_vpcs` - Field 'associate_vpcs' has been deprecated from provider version 1.153.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc.
  * `associate_status` - The status of the VPC associated with the DHCP option set.
* `owner_id` - The ID of the account to which the DHCP options set belongs.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dhcp Options Set.
* `delete` - (Defaults to 5 mins) Used when delete the Dhcp Options Set.
* `update` - (Defaults to 5 mins) Used when update the Dhcp Options Set.

## Import

VPC Dhcp Options Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_dhcp_options_set.example <id>
```