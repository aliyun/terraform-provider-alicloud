---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_certificate"
description: |-
  Provides a Alicloud ESA Certificate resource.
---

# alicloud_esa_certificate

Provides a ESA Certificate resource.



For information about ESA Certificate and how to use it, see [What is Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/SetCertificate).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_certificate" "default" {
  created_type = "free"
  domains      = "101.gositecdn.cn"
  site_id      = data.alicloud_esa_sites.default.sites.0.id
  type         = "lets_encrypt"
}
```

## Argument Reference

The following arguments are supported:
* `cas_id` - (Optional) Cloud certificate ID.
* `cert_id` - (Optional, ForceNew, Computed) The certificate Id.
* `cert_name` - (Optional, Computed) The certificate name.
* `certificate` - (Optional, Computed) Certificate content.
* `created_type` - (Required) The certificate type.
  - cas (Certificate Center Certificate)
  - upload (custom upload certificate)
  - free( Free certificate).

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `domains` - (Optional, ForceNew) A list of domain names. Multiple domain names are separated by commas.
* `private_key` - (Optional) The certificate private key.

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `region` - (Optional, Computed) Region. This parameter is required if the type is CAS.
For accounts on the Chinese site, this parameter value is: cn-hangzhou
For accounts on the international site, this parameter value is: ap-southeast-1
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites interface.
* `type` - (Optional, ForceNew, Computed) Certificate type. Possible values: lets_encrypt: Let's Encrypt certificate; 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<cert_id>`.
* `create_time` - Creation time.
* `status` - Certificate status.(within 30 days).- issued.- applying.- application failed.- canceled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate.
* `update` - (Defaults to 5 mins) Used when update the Certificate.

## Import

ESA Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_certificate.example <site_id>:<cert_id>
```