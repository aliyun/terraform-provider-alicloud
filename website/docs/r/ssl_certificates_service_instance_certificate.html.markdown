---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_instance_certificate"
description: |-
  Provides a Alicloud Certificate Management Service (Original SSL Certificate) Instance Certificate resource.
---

# alicloud_ssl_certificates_service_instance_certificate

Provides a Certificate Management Service (Original SSL Certificate) Instance Certificate resource.

A certificate issued from a certificate instance.

For information about Certificate Management Service (Original SSL Certificate) Instance Certificate and how to use it, see [What is Instance Certificate](https://next.api.alibabacloud.com/document/cas/2020-04-07/GetCertificateDetail).

-> **NOTE:** Available since v1.287.0.

~> **NOTE:** Certificates are issued by the certificate authority once an application passes review and cannot be created through an API, so this resource is read-only. Declaring it brings an already-issued certificate under management so that its attributes are tracked in state; destroying it removes it from state and leaves the certificate untouched. Use `alicloud_ssl_certificates_service_certificate_apply` to request a certificate, and `alicloud_ssl_certificates_service_instance_certificates` to look certificates up without managing them.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_instance_certificate&exampleId=a4a1a9bc-a012-a8b4-7889-b4b497a5b6d77cce5701&activeTab=example&spm=docs.r.ssl_certificates_service_instance_certificate.0.a4a1a9bca0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "certificate_id" {
  default = 12345678
}

resource "alicloud_ssl_certificates_service_instance_certificate" "default" {
  certificate_id = var.certificate_id
}

output "cert_identifier" {
  value = alicloud_ssl_certificates_service_instance_certificate.default.cert_identifier
}
```

### Deleting `alicloud_ssl_certificates_service_instance_certificate` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ssl_certificates_service_instance_certificate`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ssl_certificates_service_instance_certificate&spm=docs.r.ssl_certificates_service_instance_certificate.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `certificate_id` - (Required, ForceNew, Int) The ID of the certificate.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `instance_id` - The ID of the certificate instance that issued this certificate.
* `cert_identifier` - The global certificate identifier, formatted as the certificate ID plus `-` plus the site region ID. Alibaba Cloud services such as ALB, CDN and WAF reference a certificate by this value.
* `certificate_name` - The name of the certificate.
* `certificate_status` - The status of the certificate. Valid values: `issued`, `revoked`, `willExpire` and `expired`.
* `certificate_source` - The source of the certificate. Valid values: `BUY`, `UPLOAD` and `TEST`.
* `common_name` - The common name of the certificate.
* `domain` - All domain names covered by the certificate, separated by commas.
* `subject_alternative_names` - The subject alternative names of the certificate.
* `algorithm` - The encryption algorithm of the certificate.
* `key_size` - The key length of the certificate algorithm.
* `not_before` - The start of the certificate validity period.
* `not_after` - The end of the certificate validity period.
* `issuer` - The certificate authority that issued the certificate.
* `serial` - The serial number of the certificate.
* `finger_print` - The public key fingerprint of the certificate.
* `exist_private_key` - Indicates whether the private key of the certificate is held on the server side.
* `using_product_list` - The Alibaba Cloud services the certificate is currently deployed to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Instance Certificate.

## Import

Certificate Management Service (Original SSL Certificate) Instance Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_instance_certificate.example <certificate_id>
```
