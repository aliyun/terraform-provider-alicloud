---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_endpoint"
description: |-
  Provides a Alicloud Message Service Endpoint resource.
---

# alicloud_message_service_endpoint

Provides a Message Service Endpoint resource.


For information about Message Service Endpoint and how to use it, see [What is Endpoint](https://www.alibabacloud.com/help/en/mns/developer-reference/api-mns-open-2022-01-19-enableendpoint).

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
```

### Deleting `alicloud_message_service_endpoint` or removing it from your configuration

Terraform cannot destroy resource `alicloud_message_service_endpoint`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `endpoint_enabled` - (Required, Bool) Specifies whether the endpoint is enabled. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `endpoint_type` - (Required, ForceNew) Access point type. Value:
  - public: indicates a public access point. (Currently only public is supported)

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Endpoint.

## Import

Message Service Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_endpoint.example <id>
```
