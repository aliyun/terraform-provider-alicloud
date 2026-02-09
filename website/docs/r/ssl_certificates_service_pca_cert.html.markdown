---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_pca_cert"
description: |-
  Provides a Alicloud SSL Certificates Pca Cert resource.
---

# alicloud_ssl_certificates_service_pca_cert

Provides a SSL Certificates Pca Cert resource.



For information about SSL Certificates Pca Cert and how to use it, see [What is Pca Cert](https://next.api.alibabacloud.com/document/cas/2020-06-30/CreateClientCertificate).

-> **NOTE:** Available since v1.270.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_pca_cert&exampleId=3f88b177-bf06-c1c3-13e3-fb4ecab5f5b9a411eebd&activeTab=example&spm=docs.r.ssl_certificates_service_pca_cert.0.3f88b177bf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "root" {
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = "cbc.certqa.cn"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "sub" {
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.root.id
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = "cbc.certqa.cn"
  algorithm         = "RSA_2048"
  certificate_type  = "SUB_ROOT"
  enable_crl        = true
}

resource "alicloud_ssl_certificates_service_pca_cert" "default" {
  immediately       = "0"
  organization      = "terraform"
  years             = "1"
  upload_flag       = "0"
  locality          = "terraform"
  months            = "1"
  custom_identifier = "181"
  algorithm         = "RSA_2048"
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.sub.id
  san_value         = "somebody@example.com"
  enable_crl        = "1"
  organization_unit = "aliyun"
  state             = "Beijing"
  before_time       = "1767948807"
  days              = "1"
  san_type          = "1"
  after_time        = "1768035207"
  country_code      = "cn"
  common_name       = "exampleTerraform"
  alias_name        = "AliasName"
  status            = "ISSUE"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ssl_certificates_service_pca_cert&spm=docs.r.ssl_certificates_service_pca_cert.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `after_time` - (Optional, Int) The service expiration time of the client certificate, specified as a Unix timestamp in seconds.  

-> **NOTE:**  The `before_time` and `after_time` parameters must either both be empty or both be specified.  


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `algorithm` - (Optional, ForceNew, Computed) The key algorithm of the client certificate. The key algorithm is specified in the format `_`. Valid values:
  - `RSA_1024`: corresponds to the signature algorithm Sha256WithRSA.
  - `RSA_2048`: corresponds to the signature algorithm Sha256WithRSA.
  - `RSA_4096`: corresponds to the signature algorithm Sha256WithRSA.
  - `ECC_256`: corresponds to the signature algorithm Sha256WithECDSA.
  - `ECC_384`: corresponds to the signature algorithm Sha256WithECDSA.
  - `ECC_512`: corresponds to the signature algorithm Sha256WithECDSA.
  - `SM2_256`: corresponds to the signature algorithm SM3WithSM2.

The encryption algorithm of the client certificate must match that of the subordinate CA certificate, but the key length may differ. For example, if the key algorithm of the subordinate CA certificate is RSA_2048, the key algorithm of the client certificate must be one of RSA_1024, RSA_2048, or RSA_4096.

-> **NOTE:** You can call [DescribeCACertificate](https://help.aliyun.com/document_detail/465954.html) to query the key algorithm of the subordinate CA certificate.

* `alias_name` - (Optional) The name assigned to the issued certificate.  
* `before_time` - (Optional, Int) The issuance time of the client certificate, in timestamp format. By default, it is set to the time when you call this API. Unit: seconds.

-> **NOTE:**  The `before_time` and `after_time` parameters must either both be empty or both be specified.

* `common_name` - (Optional, ForceNew) Name of the certificate subject. For a Client Authentication (ClientAuth) certificate, the subject is typically an individual, company, organization, or application. We recommend using the common name of the subject—for example, Zhang San, Alibaba, Alibaba Cloud KMS, or Tmall Genie.  
* `country_code` - (Optional, ForceNew) Country code of the organization associated with the subordinate CA certificate that issued this certificate.  
For the meanings of different country codes, see the **International Codes** section in [Manage Company Information](https://help.aliyun.com/document_detail/198289.html).  
* `custom_identifier` - (Optional, ForceNew) A user-defined unique identifier.
* `days` - (Optional, ForceNew, Computed, Int) Validity period of the client certificate, in days.  

The `days`, `before_time`, and `after_time` parameters cannot all be empty. Additionally, `before_time` and `after_time` must either both be set or both remain unset. The specific rules are as follows:  
  - If you set the `days` parameter, you may optionally also set `before_time` and `after_time`.  
  - If you do not set the `days` parameter, you must set both `before_time` and `after_time`.  

-> **NOTE:** - If you set `days`, `before_time`, and `after_time` simultaneously, the validity period of the client certificate is determined by the value of `days`.  
  - The validity period of the client certificate cannot exceed that of the issuing subordinate CA certificate. You can call [DescribeCACertificate](https://help.aliyun.com/document_detail/465954.html) to check the validity period of the subordinate CA certificate.  
* `enable_crl` - (Optional, Int) Whether to include the CRL URL. Valid values:
  - `0`: No.
  - `1`: Yes.  
  
-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `immediately` - (Optional, Int) Specifies whether to return the digital certificate immediately. Valid values:  
  - `0`: Do not return the certificate. This is the default value.  
  - `1`: Return the certificate.  
  - `2`: Return the certificate and its certificate chain.  

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `locality` - (Optional, ForceNew) The name of the city where the organization associated with the certificate is located. Chinese characters, English letters, and other characters are supported.
By default, this value is the same as the city name of the organization associated with the issuing subordinate CA certificate.
* `months` - (Optional, Int) The duration for which the certificate is purchased, in months.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `organization` - (Optional, ForceNew) Name of the organization associated with the subordinate CA certificate that issued this certificate.  
* `organization_unit` - (Optional, ForceNew) Department name. Default: Aliyun CDN.
* `parent_identifier` - (Required, ForceNew) The unique identifier of the subordinate CA certificate that issued this certificate.
* `resource_group_id` - (Optional, Computed) The resource group ID. You can obtain this ID by calling the [ListResources](https://help.aliyun.com/document_detail/2716559.html) operation.  
* `san_type` - (Optional) The Subject Alternative Name (SAN) type supported by the client certificate. Valid values:
  - `1`: Email address.
  - `6`: Uniform Resource Identifier (URI).

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `san_value` - (Optional) Specific extension information for the client certificate. You can enter multiple extensions. If you need to specify multiple extensions, separate them with commas (,).  

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `state` - (Optional, ForceNew) The name of the province, municipality, or autonomous region where the certificate's organization is located. Chinese and English characters are supported. By default, this value is the same as the province, municipality, or autonomous region of the organization associated with the subordinate CA certificate that issued this certificate.
The name of the state or province where the certificate's organization is located. Chinese and English characters are supported. By default, this value is the same as the state or province of the organization associated with the subordinate CA certificate that issued this certificate.
* `status` - (Optional, Computed) The status of the certificate. Valid values:
  - `REVOKE`: indicates that the certificate has been revoked.
-> **NOTE:** If you want to destroy `alicloud_ssl_certificates_service_pca_cert`, `status` must be set to `REVOKE`
* `tags` - (Optional, Map) Information about the queried instances and their associated tags.  
* `upload_flag` - (Optional, Int) Indicates whether the certificate has been uploaded to the SSL certificate management platform.
* `years` - (Optional, Int) The duration for which the certificate is purchased, in years.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Pca Cert.
* `delete` - (Defaults to 5 mins) Used when delete the Pca Cert.
* `update` - (Defaults to 5 mins) Used when update the Pca Cert.

## Import

SSL Certificates Pca Cert can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_pca_cert.example <id>
```
