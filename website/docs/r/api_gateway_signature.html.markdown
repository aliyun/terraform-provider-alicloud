---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_signature"
sidebar_current: "docs-alicloud-resource-api-gateway-signature"
description: |-
  Provides an Api Gateway Signature resource.
---

# alicloud_api_gateway_signature

Provides an API Gateway signature resource. A signature is used to verify the identity of the API caller when the API is called through the signature algorithm. For information about API Gateway Signature and how to use it, see [CreateSignature](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-createsignature)

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_signature" "example" {
  signature_name   = "tf_example"
  signature_key    = "example_key"
  signature_secret = "example_secret"
}
```

## Argument Reference

The following arguments are supported:

* `signature_name` - (Required) The name of the signature. It must be 4 to 15 characters in length and can contain Chinese characters, English letters, numbers, and underscores (_). It must start with a letter or underscore.
* `signature_key` - (Required, Sensitive) The key of the signature. It is used together with the signature secret to verify the identity of the API caller.
* `signature_secret` - (Required, Sensitive) The secret of the signature. It is used together with the signature key to verify the identity of the API caller.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the signature allocated by the system.
* `created_time` - The creation time of the signature.
* `modified_time` - The last modification time of the signature.

## Import

Api gateway signature can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_signature.example "abc123"
```
