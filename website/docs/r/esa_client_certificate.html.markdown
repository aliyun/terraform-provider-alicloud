---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_client_certificate"
description: |-
  Provides a Alicloud ESA Client Certificate resource.
---

# alicloud_esa_client_certificate

Provides a ESA Client Certificate resource.



For information about ESA Client Certificate and how to use it, see [What is Client Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateClientCertificate).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_client_certificate&exampleId=6005b8b7-3b33-6cbf-7750-be15a86377d0bb7bed12&activeTab=example&spm=docs.r.esa_client_certificate.0.6005b8b73b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_client_certificate" "default" {
  site_id       = data.alicloud_esa_sites.default.sites.0.id
  pkey_type     = "RSA"
  validity_days = "365"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_client_certificate&spm=docs.r.esa_client_certificate.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `csr` - (Optional) Certificate signing request content.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `pkey_type` - (Optional) The private key algorithm type.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `site_id` - (Required, ForceNew) Site Id
* `validity_days` - (Required) Certificate validity period.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `status` - (Optional) The certificate status. Valid values: `revoked`, `active`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<client_cert_id>`.
* `client_cert_id` - ClientCertificate Id
* `create_time` - The time when the certificate was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Client Certificate.
* `update` - (Defaults to 5 mins) Used when update the Client Certificate.

## Import

ESA Client Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_client_certificate.example <site_id>:<client_cert_id>
```