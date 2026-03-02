---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_pca_certificate"
description: |-
  Provides a Alicloud SSL Certificates Pca Certificate resource.
---

# alicloud_ssl_certificates_service_pca_certificate

Provides a SSL Certificates Pca Certificate resource.



For information about SSL Certificates Pca Certificate and how to use it, see [What is Pca Certificate](https://next.api.alibabacloud.com/document/cas/2020-06-30/CreateRootCACertificate).

-> **NOTE:** Available since v1.257.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_pca_certificate&exampleId=2878a18a-0923-ad7b-372b-28a565840241968dd61e&activeTab=example&spm=docs.r.ssl_certificates_service_pca_certificate.0.2878a18a09&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "default" {
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  country_code      = "cn"
  common_name       = "cbc.certqa.cn"
  algorithm         = "RSA_2048"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ssl_certificates_service_pca_certificate&spm=docs.r.ssl_certificates_service_pca_certificate.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `algorithm` - (Optional, ForceNew) The key algorithm type of the CA certificate. The key algorithm is in the <encryption algorithm>_<key length> format. Valid values:
  - `RSA_1024`: The corresponding signature algorithm is Sha256WithRSA.
  - `RSA_2048`: The corresponding signature algorithm is Sha256WithRSA.
  - `RSA_4096`: The corresponding signature algorithm is Sha256WithRSA.
  - `ECC_256`: The signature algorithm is Sha256WithECDSA.
  - `ECC_384`: The corresponding signature algorithm is Sha256WithECDSA.
  - `ECC_512`: The signature algorithm is Sha256WithECDSA.
  - `SM2_256`: The corresponding signature algorithm is SM3WithSM2.
-> **NOTE:** If `certificate_type` is set to `SUB_ROOT`, `algorithm` is required.
* `alias_name` - (Optional, Available since v1.266.0) A custom alias for the certificate, used to define a user-friendly name.

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `common_name` - (Required, ForceNew) The common name or abbreviation of the organization. Support the use of Chinese, English characters.
* `country_code` - (Optional, ForceNew) The code of the country or region in which the organization is located, using a two-digit capital abbreviation. For example, `CN` represents China and `US` represents the United States.
* `crl_day` - (Optional, ForceNew, Int, Available since v1.269.0) The validity period for the CRL, in days. Valid values: `1` to `365`. **Note:** `crl_day` takes effect only if `certificate_type` is set to `SUB_ROOT`.
* `enable_crl` - (Optional, Available since v1.269.0) This setting turns the Certificate Revocation List (CRL) service on or off. Valid values:
  - `true`: Yes.
  - `false`: No.
-> **NOTE:** `enable_crl` takes effect only if `certificate_type` is set to `SUB_ROOT`.
* `extended_key_usages` - (Optional, List, Available since v1.269.0) The extended key usages. **Note:** `extended_key_usages` takes effect only if `certificate_type` is set to `SUB_ROOT`.
* `locality` - (Required, ForceNew) Name of the city where the organization is located. Support the use of Chinese, English characters.
* `organization` - (Required, ForceNew) The name of the organization (corresponding to your enterprise or company) associated with the CA certificate. Support the use of Chinese, English characters.
* `organization_unit` - (Required, ForceNew) The name of the department or branch under the organization. Support the use of Chinese, English characters.
* `parent_identifier` - (Optional, ForceNew, Available since v1.269.0) The unique identifier of the root CA certificate.
-> **NOTE:** If `certificate_type` is set to `SUB_ROOT`, `parent_identifier` is required.
* `path_len_constraint` - (Optional, Int, Available since v1.269.0) The certificate path length. Default value: `0`. **Note:** `path_len_constraint` takes effect only if `certificate_type` is set to `SUB_ROOT`.
* `resource_group_id` - (Optional, Available since v1.266.0) A resource property field representing the resource group.
* `state` - (Required, ForceNew)  The name of the province, municipality, or autonomous region in which the organization is located. Support the use of Chinese, English characters. 
* `tags` - (Optional, Map, Available since v1.266.0) The tag of the resource.
* `years` - (Required, ForceNew, Int) The validity period of the CA certificate, in years.
-> **NOTE:**  It is recommended to set to `5` to `10` years.
* `certificate_type` - (Optional, ForceNew, Available since v1.269.0) The type of the CA certificate. Default value: `ROOT`. Valid values:
  - `ROOT`: A root CA certificate.
  - `SUB_ROOT`: A subordinate CA certificate.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the CA certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Pca Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Pca Certificate.
* `update` - (Defaults to 5 mins) Used when update the Pca Certificate.

## Import

SSL Certificates Pca Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_pca_certificate.example <id>
```
