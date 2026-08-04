---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_certificate_validation"
description: |-
  Provides a Alicloud Certificate Management Service (Original SSL Certificate) Certificate Validation resource.
---

# alicloud_ssl_certificates_service_certificate_validation

Provides a Certificate Management Service (Original SSL Certificate) Certificate Validation resource.

Waits until the certificate application submitted on a certificate instance has been issued, and exposes the resulting certificate for downstream services to reference.

Certificate validation is the step between applying and issuance: the certificate authority verifies that the applicant controls the domain — for a DV certificate typically by checking a DNS record published under it, carrying the values that `alicloud_ssl_certificates_service_certificate_apply` reports in `domain_validation_list`. Once that verification passes, the authority issues the certificate a few minutes later. This resource represents the wait for that outcome.

For information about Certificate Management Service (Original SSL Certificate) and domain validation, see [SSL Certificates Service](https://www.alibabacloud.com/help/product/28533.html).

-> **NOTE:** Available since v1.287.0.

~> **NOTE:** This resource creates nothing in the cloud. It only waits: after the application has been submitted and the domain validation records have been published, it polls the instance until the certificate is issued. Destroying it calls no API and leaves the certificate untouched. It exists so that downstream resources can depend on a certificate that has actually been issued, rather than on one that has merely been requested.

## Example Usage

Basic Usage

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

resource "alicloud_alidns_record" "default" {
  for_each = {
    for v in alicloud_ssl_certificates_service_certificate_apply.default.domain_validation_list :
    v.domain => v
  }
  domain_name = each.value.root_domain
  rr          = each.value.validation_key
  type        = each.value.validation_type
  value       = each.value.validation_value
  ttl         = 600
}

resource "alicloud_ssl_certificates_service_certificate_validation" "default" {
  instance_id           = alicloud_ssl_certificates_service_instance.default.id
  validation_record_ids = [for r in alicloud_alidns_record.default : r.id]

  timeouts {
    create = "120m"
  }
}
```

Downstream services should reference the certificate exposed by this resource, so that the reference is only resolved once the certificate actually exists:

```terraform
resource "alicloud_alb_listener" "default" {
  certificates {
    certificate_id = alicloud_ssl_certificates_service_certificate_validation.default.cert_identifier
  }
}
```

### Deleting `alicloud_ssl_certificates_service_certificate_validation` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ssl_certificates_service_certificate_validation`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the certificate instance whose application is being waited on.
* `validation_record_ids` - (Optional, ForceNew, List) The IDs of the DNS records carrying the domain ownership validation information. The values themselves are never read; they exist so that Terraform creates the validation records before it starts waiting for issuance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `certificate_id` - The ID of the issued certificate.
* `cert_identifier` - The global certificate identifier, formatted as the certificate ID plus `-` plus the site region ID. Alibaba Cloud services such as ALB, CDN and WAF reference a certificate by this value.
* `certificate_status` - The status of the issued certificate.

-> **NOTE:** Enabling managed renewal on the certificate instance causes a **new certificate** to be issued when the certificate is renewed, which changes `certificate_id` and `cert_identifier`. Resources referencing these attributes pick up the new certificate on the next apply.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 75 mins) Used when waiting for the certificate to be issued. A DV certificate is usually issued within minutes of the validation records going live. OV and EV certificates additionally go through a manual review that can take hours and exposes no progress API, so raise this timeout for them.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate Validation.

## Import

Certificate Management Service (Original SSL Certificate) Certificate Validation can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_certificate_validation.example <instance_id>
```
