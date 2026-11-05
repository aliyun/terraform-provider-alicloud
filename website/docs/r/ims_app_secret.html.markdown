---
subcategory: "IMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ims_app_secret"
sidebar_current: "docs-alicloud-resource-ims-app-secret"
description: |-
  Provides an IMS App Secret resource.
---

# alicloud_ims_app_secret

Provides an IMS App Secret resource. An application secret is a credential of an IMS application, usable as an OAuth client secret.

For information about IMS App Secret and how to use it, see [CreateAppSecret](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ims-2019-08-15-createappsecret).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The IMS application that the secret belongs to is managed out-of-band (there is no `alicloud_ims_application` resource). You must provide an existing `app_id`.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/terraform-provider-alicloud"
    }
  }
}

// Source: aliyun/terraform-provider-alicloud
// An IMS application is managed out-of-band; provide its AppId here.
resource "alicloud_ims_app_secret" "example" {
  app_id = "472457090344041****"
}

output "app_secret_value" {
  value     = alicloud_ims_app_secret.example.app_secret_value
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew, String) The ID of the IMS application that the secret belongs to. Changing this forces replacement of the secret. The application is managed out-of-band; this resource only creates and destroys a single secret for it.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource, in the format `<app_id>:<app_secret_id>`.
* `app_secret_id` - (Computed, String) The ID of the application secret, server-generated.
* `app_secret_value` - (Computed, Sensitive, String) The value of the application secret, usable as an OAuth client secret. Only returned by CreateAppSecret and GetAppSecret; it is not part of ListAppSecretIds output.
* `create_date` - (Computed, String) The time the secret was created, in RFC 3339 format (UTC).

## Import

IMS App Secret can be imported using the composite id `<app_id>:<app_secret_id>`, e.g.

```shell
$ terraform import alicloud_ims_app_secret.example 472457090344041****:2efd5004-005c-4f05-83c6-5b1dd176****
```
