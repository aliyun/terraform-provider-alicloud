---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-api-gateway-acl-entry-attachment"
description: |-
  Provides an Alicloud Api Gateway ACL Entry Attachment Resource.
---

# alicloud_api_gateway_acl_entry_attachment

Provides an ACL entry attachment resource for attaching ACL entry to an API Gateway ACL.

-> **NOTE:** Available since v1.228.0

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_api_gateway_access_control_list" "default" {
  access_control_list_name = var.name
  address_ip_version       = "ipv4"
}

resource "alicloud_api_gateway_acl_entry_attachment" "default" {
  acl_id  = alicloud_api_gateway_access_control_list.default.id
  entry   = "128.0.0.1/32"
  comment = "test comment"
}
```

## Argument Reference

The following arguments are supported:
* `acl_id` - (Required, ForceNew) The ID of the ACL that the entry will be attached to.
* `entry` - (Required, ForceNew) The CIDR block of the entry to attach. 
* `comment` - (Optional, ForceNew) The comment for the entry.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. The value formats as `<acl_id>:<entry>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the ACL Entry Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the ACL Entry Attachment.
* `update` - (Defaults to 5 mins) Used when update the ACL Entry Attachment.

## Import

Api Gateway Acl Entry Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_acl_entry_attachment.example <acl_id>:<entry>
```