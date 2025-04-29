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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_endpoint&exampleId=d96a3eaa-04e2-197f-4043-68de4784dcbe83e7663b&activeTab=example&spm=docs.r.message_service_endpoint.0.d96a3eaa04&intl_lang=EN_US" target="_blank">
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Endpoint.

## Import

Message Service Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_endpoint.example <id>
```
