---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_instance_acl_attachment"
sidebar_current: "docs-alicloud-resource-api-gateway-instance-acl-attachment"
description: |-
  Provides an Alicloud Api Gateway Instance ACL Attachment Resource.
---

# alicloud_api_gateway_instance_acl_attachment

Provides an Instance ACL attachment resource for attaching an ACL to a specific API Gateway instance.

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_api_gateway_instance" "default" {
  instance_name = var.name
  instance_spec = "api.s1.small"
  https_policy  = "HTTPS2_TLS1_0"
  zone_id       = "cn-hangzhou-MAZ6"
  payment_type  = "PayAsYouGo"
  instance_type = "normal"
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

resource "alicloud_api_gateway_instance_acl_attachment" "default" {
  instance_id = alicloud_api_gateway_instance.default.id
  acl_id      = alicloud_api_gateway_access_control_list.default.id
  acl_type    = "white"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the API Gateway instance that the ACL will be attached to.
* `acl_id` - (Required, ForceNew) The ID of the ACL to attach.
* `acl_type` - (Required, ForceNew) The type of the ACL. Valid values: `white`, `black`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<insntance_id>:<acl_id>:<acl_type>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance ACL Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Instance ACL Attachment.
* `update` - (Defaults to 5 mins) Used when update the Instance ACL Attachment.

## Import

Api Gateway Instance Acl Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_instance_acl_attachment.example <instance_id>:<acl_id>:<acl_type>
```
