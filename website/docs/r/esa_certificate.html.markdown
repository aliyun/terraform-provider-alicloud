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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_certificate&exampleId=854a5788-03f5-0064-32c8-713a80e71dcbc607a971&activeTab=example&spm=docs.r.esa_certificate.0.854a578803&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `cert_name` - (Optional) The certificate name.
* `certificate` - (Optional) The certificate type. Valid values:

  - cas: a certificate purchased by using Certificate Management Service.
  - upload: a custom certificate that you upload.
* `created_type` - (Required) The certificate type.
  - cas (Certificate Center Certificate)
  - upload (custom upload certificate)
  - free( Free certificate).
* `domains` - (Optional, ForceNew) A list of domain names. Multiple domain names are separated by commas.
* `private_key` - (Optional) The certificate private key.
* `region` - (Optional) Geographical information.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites interface.
* `type` - (Optional, ForceNew) Certificate type. Possible values: lets_encrypt: Let's Encrypt certificate; 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<cert_id>`.
* `create_time` - Creation time.
* `status` - Certificate status.(within 30 days).- issued.- applying.- application failed.- canceled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate.
* `update` - (Defaults to 5 mins) Used when update the Certificate.

## Import

ESA Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_certificate.example <site_id>:<cert_id>
```