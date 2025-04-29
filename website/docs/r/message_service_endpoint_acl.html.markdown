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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_endpoint_acl&exampleId=34e18dbb-b48d-014a-3ff7-e69fe36253931d386189&activeTab=example&spm=docs.r.message_service_endpoint_acl.0.34e18dbbb4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Endpoint Acl.
* `delete` - (Defaults to 5 mins) Used when delete the Endpoint Acl.

## Import

Message Service Endpoint Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_endpoint_acl.example <endpoint_type>:<acl_strategy>:<cidr>
```
