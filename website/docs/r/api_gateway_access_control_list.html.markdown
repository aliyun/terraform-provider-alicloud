---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_access_control_list"
description: |-
  Provides a Alicloud Api Gateway Access Control List resource.
---

# alicloud_api_gateway_access_control_list

Provides a Api Gateway Access Control List resource. Access control list.

For information about Api Gateway Access Control List and how to use it, see [What is Access Control List](https://www.alibabacloud.com/help/en/api-gateway/developer-reference/api-cloudapi-2016-07-14-createaccesscontrollist).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_access_control_list&exampleId=118d82d0-f617-a50e-927d-660972288cf5a6442c79&activeTab=example&spm=docs.r.api_gateway_access_control_list.0.118d82d0f6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_api_gateway_instance" "defaultxywS8c" {
  instance_name = var.name
  instance_spec = "api.s1.small"
  https_policy  = "HTTPS2_TLS1_0"
  zone_id       = "cn-hangzhou-MAZ6"
  payment_type  = "PayAsYouGo"
}

resource "alicloud_api_gateway_access_control_list" "default" {
  access_control_list_name = var.name
  address_ip_version       = "ipv4"
}
```

## Argument Reference

The following arguments are supported:
* `access_control_list_name` - (Required, ForceNew) Access control list name.
* `acl_entrys` - (Optional, Deprecated from v1.228.0) Information list of access control policies. You can add at most 50 IP addresses or CIDR blocks to an ACL in each call. If the IP address or CIDR block that you want to add to an ACL already exists, the IP address or CIDR block is not added. The entries that you add must be CIDR blocks. See [`acl_entrys`](#acl_entrys) below.
**NOTE:** Field 'acl_entrys' has been deprecated from provider version 1.228.0, and it will be removed in the future version. Please use the new resource 'alicloud_api_gateway_acl_entry_attachment'.
* `address_ip_version` - (Optional, ForceNew, Computed) The IP version. Valid values: ipv4 and ipv6.

### `acl_entrys`

The acl_entrys supports the following:
* `acl_entry_comment` - (Optional) The description of the ACL.
* `acl_entry_ip` - (Optional) The entries that you want to add to the ACL. You can add CIDR blocks. Separate multiple CIDR blocks with commas (,).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Control List.
* `delete` - (Defaults to 5 mins) Used when delete the Access Control List.
* `update` - (Defaults to 5 mins) Used when update the Access Control List.

## Import

Api Gateway Access Control List can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_access_control_list.example <id>
```