---
subcategory: "SSL Certificates"
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

## Argument Reference

The following arguments are supported:
* `algorithm` - (Optional, ForceNew, Computed) The key algorithm type of the root CA certificate. The key algorithm is expressed in the format `_`. Valid values:
  - `RSA_1024`: Corresponds to the signature algorithm Sha256WithRSA.
  - `RSA_2048`: Corresponds to the signature algorithm Sha256WithRSA.
  - `RSA_4096`: Corresponds to the signature algorithm Sha256WithRSA.
  - `ECC_256`: Corresponds to the signature algorithm Sha256WithECDSA.
  - `ECC_384`: Corresponds to the signature algorithm Sha256WithECDSA.
  - `ECC_512`: Corresponds to the signature algorithm Sha256WithECDSA.
  - `SM2_256`: Corresponds to the signature algorithm SM3WithSM2.

The encryption algorithm of the root CA certificate must match the **certificate algorithm** of the private root CA you purchased. For example, if you selected `RSA` as the **certificate algorithm** when purchasing the private root CA, the key algorithm of the root CA certificate must be `RSA_1024`, `RSA_2048`, or `RSA_4096`.

-> **NOTE:** If `certificate_type` is set to `SUB_ROOT`, `algorithm` is required.

* `alias_name` - (Optional, Available since v1.266.0) A custom alias for the certificate, used to define a user-friendly name.

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `common_name` - (Required, ForceNew) The common name or short name of the organization. Chinese characters, English letters, and other characters are supported.
* `country_code` - (Optional, ForceNew) The two-letter uppercase alphabetic code representing the country or region where the organization is located. For example, `CN` represents China and `US` represents the United States.
For country codes, see the **International Codes** section in [Managing Company Information](https://help.aliyun.com/document_detail/198289.html).
* `crl_day` - (Optional, ForceNew, Int, Available since v1.269.0) The interval (in days) for updating the Certificate Revocation List (CRL).
* `enable_crl` - (Optional, Available since v1.269.0) Specifies whether to enable CRL.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `extended_key_usages` - (Optional, List, Available since v1.269.0) Extended attributes of the certificate, used to define extended key usages.  

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `locality` - (Required, ForceNew) The name of the city where the organization is located.  
* `organization` - (Required, ForceNew) The name of the organization associated with the CA certificate.
* `organization_unit` - (Required, ForceNew) The name of the department or branch within the organization
* `parent_identifier` - (Optional, ForceNew, Available since v1.269.0) Parent node identifier.  
* `path_len_constraint` - (Optional, Int, Available since v1.269.0) The maximum depth of subordinate CA levels allowed under this CA.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `resource_group_id` - (Optional, Computed, Available since v1.266.0) A resource property field representing the resource group.
* `state` - (Required, ForceNew) The name of the province, municipality directly under the central government, or autonomous region where the organization is located
* `tags` - (Optional, Map, Available since v1.266.0) Tags
* `years` - (Required, ForceNew, Int) The validity period of the root CA certificate, in years.

-> **NOTE:**  We recommend setting it to 5–10 years.

* `certificate_type` - (Optional, ForceNew, Available since v1.269.0) The type of the CA certificate. Default value: `ROOT`. Valid values:
  - `ROOT`: A root CA certificate.
  - `SUB_ROOT`: A subordinate CA certificate.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `status` - The current CA status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Pca Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Pca Certificate.
* `update` - (Defaults to 5 mins) Used when update the Pca Certificate.

## Import

SSL Certificates Pca Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_pca_certificate.example <identifier>
```