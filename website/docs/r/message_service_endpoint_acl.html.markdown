---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_endpoint_acl"
description: |-
  Provides a Alicloud Message Service Endpoint Acl resource.
---

# alicloud_message_service_endpoint_acl

Provides a Message Service Endpoint Acl resource.



For information about Message Service Endpoint Acl and how to use it, see [What is Endpoint Acl](https://www.alibabacloud.com/help/en/mns/developer-reference/api-mns-open-2022-01-19-authorizeendpointacl).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_message_service_endpoint" "default" {
  endpoint_enabled = true
  endpoint_type    = "public"
}

resource "alicloud_message_service_endpoint_acl" "default" {
  cidr          = "192.168.1.1/23"
  endpoint_type = alicloud_message_service_endpoint.default.id
  acl_strategy  = "allow"
}
```

## Argument Reference

The following arguments are supported:
* `acl_strategy` - (Required, ForceNew) The ACL policy. Valid value:
  - allow: indicates that the current endpoint allows access from the corresponding CIDR block. (Only allow is supported)
* `cidr` - (Required, ForceNew) The CIDR block.
-> **NOTE:** To ensure business stability, the system is configured by default with a CIDR (0.0.0.0/0) that allows access from all source addresses. If you need to remove this default configuration, you can do so by importing and deleting the CIDR using Terraform, or by manually deleting it in the console.
* `endpoint_type` - (Required, ForceNew) Access point type. Value:
  - public: indicates a public access point. (Currently only public is supported)

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Endpoint Acl. It formats as `<endpoint_type>:<acl_strategy>:<cidr>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Endpoint Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Endpoint Acl.

## Import

Message Service Endpoint Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_endpoint_acl.example <endpoint_type>:<acl_strategy>:<cidr>
```
