---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_origin_client_certificate"
description: |-
  Provides a Alicloud ESA Origin Client Certificate resource.
---

# alicloud_esa_origin_client_certificate

Provides a ESA Origin Client Certificate resource.



For information about ESA Origin Client Certificate and how to use it, see [What is Origin Client Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/UploadOriginClientCertificate).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_OriginClientCertificate_example" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_origin_client_certificate" "default" {
  site_id     = alicloud_esa_site.resource_Site_OriginClientCertificate_example.id
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC+7dgpkQifIqzOU6KNkFRjZtMZOoN7/ihNf/BrYcPhLQSkcPOf\nUsTP/qvH0u965GnYFiAoK3uWGQo9aCBuoawRFKNBa9ZpJVyVbamBWTBQ/Fxsforq\n9jJNR7OWA3fxvDxgwyEkv0qsnh1zaOkjyUlxFYwDiFxZ1/RHAj/SABCubQIDAQAB\nAoGADiobBUprN1MdOtldj98LQ6yXMKH0qzg5yTYaofzIyWXLmF+A02sSitO77sEp\nXxae+5b4n8JKEuKcrd2RumNoHmN47iLQ0M2eodjUQ96kzm5Esq6nln62/NF5KLuK\nJDw63nTsg6K0O+gQZv4SYjZAL3cswSmeQmvmcoNgArfcaoECQQDgYy6S91ZIUsLx\n6BB3tW+x7APYnvKysYbcKUEP8AutZSo4hdMfPQkOD0LwP5dWsrNippDWjNDiPZmt\nVKuZDoDdAkEA2dPxy1eQeJsRYTZmTWIuh3UY9xlL3G9skcSOM4LbFidroHWW9UDJ\nJDSSEMH2+/4quYTdPr28cj7RCjqL0brC0QJABXDCL1QJ5oUDLwRWaeCfTawQR89K\nySRexbXGWxGR5uleBbLQ9J/xOUMLd3HDRJnemZS6TElrwyCFOlukMXjVjQJBALr5\nQC0opmu/vzVQepOl2QaQrrM7VXCLfAfLTbxNcD0d7TY4eTFfQMgBD/euZpB65LWF\npFs8hcsSvGApTObjhmECQEydB1zzjU6kH171XlXCtRFnbORu2IB7rMsDP2CBPHyR\ntYBjBNVHIUGcmrMVFX4LeMuvvmUyzwfgLmLchHxbDP8=\n-----END RSA PRIVATE KEY-----\n"
  hostnames = [
    "www.example1.com",
    "www.example2.com",
    "www.example3.com"
  ]
  origin_client_certificate_name = "exampleCertificate"
  certificate                    = "-----BEGIN CERTIFICATE-----\nMIICQTCCAaoCCQCFfdyqahygLzANBgkqhkiG9w0BAQUFADBlMQswCQYDVQQGEwJj\nbjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwI\nYWxpY2xvdWQxEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwHhcNMjAw\nODA2MTAwMDAyWhcNMzAwODA0MTAwMDAyWjBlMQswCQYDVQQGEwJjbjEQMA4GA1UE\nCAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamluZzERMA8GA1UECgwIYWxpY2xvdWQx\nEDAOBgNVBAsMB2FsaWJhYmExDTALBgNVBAMMBHRlc3QwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAL7t2CmRCJ8irM5Too2QVGNm0xk6g3v+KE1/8Gthw+EtBKRw\n859SxM/+q8fS73rkadgWICgre5YZCj1oIG6hrBEUo0Fr1mklXJVtqYFZMFD8XGx+\niur2Mk1Hs5YDd/G8PGDDISS/SqyeHXNo6SPJSXEVjAOIXFnX9EcCP9IAEK5tAgMB\nAAEwDQYJKoZIhvcNAQEFBQADgYEAavYdM9s5jLFP9/ZPCrsRuRsjSJpe5y9VZL+1\n+Ebbw16V0xMYaqODyFH1meLRW/A4xUs15Ny2vLYOW15Mriif7Sixty3HUedBFa4l\ny6/gQ+mBEeZYzMaTTFgyzEZDMsfZxwV9GKfhOzAmK3jZ2LDpHIhnlJN4WwVf0lME\npCPDN7g=\n-----END CERTIFICATE-----\n"
}
```

## Argument Reference

The following arguments are supported:
* `certificate` - (Required, ForceNew) The certificate content.
* `hostnames` - (Optional, List) The domain names to associate.
* `origin_client_certificate_name` - (Optional, ForceNew) The certificate name.
* `private_key` - (Required, ForceNew) The private key of the certificate.
* `site_id` - (Required, ForceNew) Site ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<origin_client_certificate_id>`.
* `create_time` - The time when the certificate was created.
* `origin_client_certificate_id` - The certificate ID.
* `status` - The status of the certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Origin Client Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Origin Client Certificate.
* `update` - (Defaults to 5 mins) Used when update the Origin Client Certificate.

## Import

ESA Origin Client Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_origin_client_certificate.example <site_id>:<origin_client_certificate_id>
```