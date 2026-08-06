---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_certificate_apply"
description: |-
  Provides a Alicloud Certificate Management Service (Original SSL Certificate) Certificate Apply resource.
---

# alicloud_ssl_certificates_service_certificate_apply

Provides a Certificate Management Service (Original SSL Certificate) Certificate Apply resource.

Submits the certificate application recorded on a certificate instance and exposes the domain ownership validation records that have to be published before the certificate authority can issue the certificate.

For information about Certificate Management Service (Original SSL Certificate) Certificate Apply and how to use it, see [What is Certificate Apply](https://next.api.alibabacloud.com/document/cas/2020-04-07/ApplyCertificate).

-> **NOTE:** Available since v1.287.0.

~> **NOTE:** Creating this resource returns as soon as the application has been submitted; it does **not** wait for the certificate to be issued. Publish the records from `domain_validation_list` on your authoritative DNS, then use `alicloud_ssl_certificates_service_certificate_validation` to wait for issuance.

-> **NOTE:** A certificate instance can have at most one application in progress at a time, which is why this resource is identified by the instance ID.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_certificate_apply&exampleId=45819ed0-3d53-5417-be25-263d061cc0ca2e1a8ac0&activeTab=example&spm=docs.r.ssl_certificates_service_certificate_apply.0.45819ed03d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ssl_certificates_service_instance" "default" {
  product_type      = "cas"
  period            = 12
  pricing_cycle     = 2
  instance_name     = var.name
  domain            = "example.com"
  validation_method = "DNS"

  parameter {
    code  = "fullSpec"
    value = "ws.dv.f"
  }
  parameter {
    code  = "fullDomainCount"
    value = "1"
  }
}

resource "alicloud_ssl_certificates_service_certificate_apply" "default" {
  instance_id       = alicloud_ssl_certificates_service_instance.default.id
  domain            = alicloud_ssl_certificates_service_instance.default.domain
  validation_method = alicloud_ssl_certificates_service_instance.default.validation_method
}

output "domain_validation_list" {
  value = alicloud_ssl_certificates_service_certificate_apply.default.domain_validation_list
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ssl_certificates_service_certificate_apply&spm=docs.r.ssl_certificates_service_certificate_apply.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the certificate instance to submit the application on.
* `domain` - (Optional, ForceNew, Computed) The domain names the certificate is requested for, separated by commas.
* `csr` - (Optional, ForceNew, Computed) The content of the certificate signing request.
* `validation_method` - (Optional, ForceNew, Computed) The domain ownership validation method. Valid values: `DNS` and `HTTP`.
* `key_algorithm` - (Optional, ForceNew, Computed) The key algorithm the certificate is issued with, such as `RSA_2048`, `RSA_4096`, `ECC_256` or `SM2`.
* `generate_csr_method` - (Optional, ForceNew, Computed) How the certificate signing request is produced. Valid values: `online` and `upload`.

~> **NOTE:** The five application settings above are stored on the certificate instance and are written there by the instance resource. `ApplyCertificate` accepts none of them — it submits whatever the instance holds at that moment. They are accepted here so that changing one of them replaces this resource and resubmits the application, since a different configuration means a different certificate is being requested. **Reference the corresponding attribute of `alicloud_ssl_certificates_service_instance` rather than hardcoding a value**; a configured value that disagrees with the instance is rejected when the resource is created.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `certificate_status` - The status of the certificate this application produces.
* `pending_result` - The reason the application is pending, if the review has not passed.
* `certificate_id` - The ID of the certificate produced by this application. Populated only after the certificate has been issued.
* `cert_identifier` - The global certificate identifier of the certificate produced by this application. Populated only after the certificate has been issued.
* `domain_validation_list` - The domain ownership validation records that have to be published before the certificate authority can validate the domains. It contains the following attributes:
  * `domain` - The domain name to be validated.
  * `root_domain` - The root domain of the domain name to be validated. Use this value to locate the DNS zone when adding the validation record.
  * `validation_type` - The type of the validation record. Valid values: `TXT`, `CNAME` and `FILE`.
  * `validation_key` - The host record of the validation record.
  * `validation_value` - The value of the validation record.
  * `cname` - The record value used when the domain is validated through a CNAME record.
  * `cname_key` - The host record used when the domain is validated through a CNAME record.

-> **NOTE:** Reference `certificate_id` and `cert_identifier` from `alicloud_ssl_certificates_service_certificate_validation` rather than from this resource. The values here are only populated once the certificate happens to have been issued by the time of a refresh, whereas the validation resource guarantees issuance has completed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate Apply.
* `delete` - (Defaults to 20 mins) Used when delete the Certificate Apply. Withdrawing the application is processed asynchronously and the instance can stay in the pending state for several minutes afterwards; it has to leave that state before the instance itself can be refunded and removed.

## Import

Certificate Management Service (Original SSL Certificate) Certificate Apply can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_certificate_apply.example <instance_id>
```
