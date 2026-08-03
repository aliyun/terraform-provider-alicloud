---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_instance"
description: |-
  Provides a Alicloud Certificate Management Service (Original SSL Certificate) Instance resource.
---

# alicloud_ssl_certificates_service_instance

Provides a Certificate Management Service (Original SSL Certificate) Instance resource.



For information about Certificate Management Service (Original SSL Certificate) Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

-> **NOTE:** Available since v1.287.0.

~> **NOTE:** This resource manages the certificate **instance** — the purchased subscription slot that a certificate can later be issued into — and the certificate application details recorded on it. It does **not** submit the certificate application and does **not** issue a certificate. A newly created instance stays in the `inactive` state. Submitting the application against this instance is done with `alicloud_ssl_certificates_service_certificate_apply`; publishing the domain validation records it reports, and waiting for the certificate authority to issue the certificate, are separate steps again.

-> **NOTE:** Creating this resource places a real prepaid subscription order through `BssOpenApi::CreateInstance` and is billed on creation. The available specifications and their prices differ per account; call [DescribePricingModule](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/DescribePricingModule) with `ProductCode` and `ProductType` set to `cas` to list what your account can order. Destroying the resource refunds the order when the instance is still protected, then deletes it.

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
  product_type  = "cas"
  period        = 12
  pricing_cycle = 2
  instance_name = var.name

  parameter {
    code  = "fullSpec"
    value = "ws.dv.f"
  }
  parameter {
    code  = "fullDomainCount"
    value = "1"
  }
}
```

## Argument Reference

The following arguments are supported:
* `auto_reissue` - (Optional, Computed) Specifies whether to enable managed renewal.
enable: Enabled
disable: Disabled

~> **NOTE:** When managed renewal is enabled, the certificate authority issues a **new certificate** for the instance instead of extending the existing one. The renewed certificate has a different certificate ID and certificate identifier, so any value read from this instance or referenced by downstream resources changes after a renewal. Make sure the resources that consume the certificate take their value from the certificate attributes rather than from a hardcoded ID.

* `company_id` - (Optional, Computed) The company information ID. Required for OV and EV certificates — without it the certificate application cannot be submitted for this instance.

~> **NOTE:** The address recorded on the instance — `city` and `province` — comes from the company referenced by `company_id`, and from a server-side default when no company is referenced. Certificate authorities reject non-standard addresses, so the address is taken from the company record rather than accepted per instance; set it on `alicloud_ssl_certificates_service_company` instead. Both attributes are read-only here.

* `contact_id_list` - (Optional, Set) The list of contact IDs. At least one contact ID is required, otherwise the certificate application cannot be submitted for this instance.
* `country_code` - (Optional, Computed) The code of the country or region where the certificate organization is located. For example, CN indicates China and US indicates the United States. This field is required when a Certificate Signing Request (CSR) is generated for a DV certificate. The default value is CN.
* `csr` - (Optional, Computed) The CSR content. You can generate a CSR by using OpenSSL or Keytool. For more information, see [How to create a CSR file](https://help.aliyun.com/document_detail/42218.html).
* `domain` - (Optional, Computed) The domain names to be bound to the certificate. Specific requirements are as follows:
  - You can specify a single domain name or a wildcard domain name (such as `*.aliyundoc.com`).
  - You can specify multiple domain names. Separate multiple domain names with commas (,). Whether free domain names are provided is determined based on the first domain name.

-> **NOTE:** When you bind multiple domain names to a certificate, this parameter is required. This parameter and `csr` cannot both be empty.

~> **NOTE:** Which of `domain` and `csr` is authoritative depends on `generate_csr_method`:

| `generate_csr_method` | Which value wins |
| --- | --- |
| `online` | The system generates the CSR and `domain` is used as you set it. Leave `csr` unset; the generated value is read back into state. |
| `upload` | The domain is taken from the `CN` field of the CSR you upload, so leave `domain` unset and let it be read back from the certificate. Setting it to a value that disagrees with the CSR is rejected by the API. |

* `generate_csr_method` - (Optional, Computed) The method used to generate the Certificate Signing Request (CSR). Default value: online.
online: The system generates the CSR. In this case, the Csr parameter is ignored.
upload: You upload the CSR. In this case, the Csr parameter is required.
* `instance_name` - (Optional, Computed) The name of the instance. When a certificate is issued, this value is used as the default name of the certificate.
* `key_algorithm` - (Optional, Computed) The certificate algorithm. Default value: RSA_2048.
  - `RSA_2048`
  - `RSA_3072`
  - `RSA_4096`
  - `ECC_256`
  - `SM2`
* `parameter` - (Optional, ForceNew, List) The list of modules. See [`parameter`](#parameter) below.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `period` - (Optional, ForceNew, Int) The subscription period. Unit: months. For products billed on a yearly basis, enter an integer multiple of 12.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `pricing_cycle` - (Optional, ForceNew, Int) The unit of the subscription period.
  - PricingCycle=1: The unit of the subscription period is year.
  - PricingCycle=2: The unit of the subscription period is month.
  - PricingCycle=3: The unit of the subscription period is day.

Default value: PricingCycle=2

This parameter applies only to specific product types (ProductType is ddos_originpre_public_cn, ddosDip, ddoscoo, ddos_originpre_public_intl, ddosDip_intl, or ddoscoo_intl).

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `product_type` - (Optional, ForceNew, Computed) The product type. Valid values: `cas` (China site) and `cas_intl` (international site). Defaults to `cas_intl` for international accounts and `cas` otherwise.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The list of tags.
* `validation_method` - (Optional, Computed) The verification method for the certificate application.
DNS: DNS verification, using TXT or CNAME records.
HTTP: File verification.

### `parameter`

The parameter supports the following:
* `code` - (Required, ForceNew) The module code. Common values are `fullSpec` (certificate brand and grade for a single-domain certificate), `wildcardSpec` (the same for a wildcard certificate), `fullDomainCount` (number of single domain names) and `wildcardDomainCount` (number of wildcard domain names).
* `value` - (Required, ForceNew) The value of the module identified by `code`. For the domain-count modules this is an integer from 1 to 150. For `fullSpec` and `wildcardSpec` this is a specification code such as `ws.dv.f`.

Specification codes accepted by `fullSpec` follow a `brand.grade.type` pattern, for example:

| Value | Brand | Grade |
| --- | --- | --- |
| `ws.dv.f` | WoSign | DV |
| `rap.dv.f` | Rapid | DV |
| `vt.dv.f` / `vt.ov.f` | vTrus | DV / OV |
| `cf.ov.f` / `cf.ev.f` | CFCA | OV / EV |
| `gs.dv.f` / `gs.ov.f` | GlobalSign | DV / OV |
| `geo.dv.f` / `geo.ov.f` / `geo.ev.f` | GeoTrust | DV / OV / EV |
| `ss.ov.f` / `ss.ev.f` | DigiCert | OV / EV |

-> **NOTE:** The list above is a snapshot and is not exhaustive. Specifications, prices and availability differ per account and change over time. Call [DescribePricingModule](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/DescribePricingModule) with `ProductCode` and `ProductType` set to `cas` to get the values your account can actually order.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `city` - The city of the address recorded on the instance.
* `province` - The province or region of the address recorded on the instance.
* `average_waiting_time` - Average waiting time for issuing a certificate of this specification, in seconds.
* `brand` - The certificate brand.
* `certificate_type` - The type of the certificate.
* `full_domain_count` - The number of single domain names included in the instance.
* `instance_end_time` - Instance expiration time, a UNIX timestamp in seconds.
* `instance_start_time` - Instance start time, a UNIX timestamp in seconds.
* `instance_type` - The instance type.
* `order_end_time` - Order end time, a UNIX timestamp in seconds.
* `order_start_time` - Order start time, a UNIX timestamp in seconds.
* `spec` - The specification of the purchased instance.
* `status` - The instance status.
* `upgrade_status` - The upgrade status of the instance.
* `wildcard_domain_count` - The number of wildcard domain names included in the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Certificate Management Service (Original SSL Certificate) Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_instance.example <instance_id>
```