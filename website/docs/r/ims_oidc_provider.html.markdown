---
subcategory: "IMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ims_oidc_provider"
description: |-
  Provides a Alicloud IMS Oidc Provider resource.
---

# alicloud_ims_oidc_provider

Provides a IMS Oidc Provider resource.

OpenID Connect Provider.

For information about IMS Oidc Provider and how to use it, see [What is Oidc Provider](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-createoidcprovider).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "oidc_provider_name" {
  default = "amp-resource-example-oidc-provider"
}


resource "alicloud_ims_oidc_provider" "default" {
  description         = var.oidc_provider_name
  issuer_url          = "https://oauth.aliyun.com"
  fingerprints        = ["902ef2deeb3c5b13ea4c3d5193629309e231ae55"]
  issuance_limit_time = "12"
  oidc_provider_name  = var.name
  client_ids          = ["123", "456"]
}
```

## Argument Reference

The following arguments are supported:
* `client_ids` - (Optional, Set) Client ID. 
* `description` - (Optional) Description of OIDC identity provider.
* `fingerprints` - (Optional, Set) The authentication fingerprint of the HTTPS CA certificate.
* `issuance_limit_time` - (Optional, Computed, Int) The earliest time when an external IdP is allowed to issue an ID Token. If the iat field in the ID Token is greater than the current time, the request is rejected.
Unit: hours. Value range: 1~168.
* `issuer_url` - (Required, ForceNew) The issuer URL of the OIDC identity provider.
* `oidc_provider_name` - (Required, ForceNew) The name of the OIDC identity provider.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `arn` - ARN of OIDC identity provider.
* `create_time` - Creation Time (UTC time).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Oidc Provider.
* `delete` - (Defaults to 5 mins) Used when delete the Oidc Provider.
* `update` - (Defaults to 5 mins) Used when update the Oidc Provider.

## Import

IMS Oidc Provider can be imported using the id, e.g.

```shell
$ terraform import alicloud_ims_oidc_provider.example <id>
```