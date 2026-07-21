---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_ca_certificate"
sidebar_current: "docs-alicloud-resource-slb-ca-certificate"
description: |-
  Provides a Load Banlancer CA Certificate resource.
---

# alicloud_slb_ca_certificate

-> **NOTE:** Available since v1.21.0.

A Load Balancer CA Certificate is used by the listener of the protocol https.

For information about slb and how to use it, see [What is Server Load Balancer](https://www.alibabacloud.com/help/doc-detail/27539.htm).

For information about CA Certificate and how to use it, see [Configure CA Certificate](https://www.alibabacloud.com/help/doc-detail/85968.htm).


## Example Usage

* using CA certificate content

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_ca_certificate&exampleId=ecd8639f-0e20-6099-1e2c-21fa2aeea4c507872863&activeTab=example&spm=docs.r.slb_ca_certificate.0.ecd8639f0e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_slb_ca_certificate" "foo" {
  ca_certificate_name = "tf-testAccSlbCACertificate"
  ca_certificate      = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

* using CA certificate file

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_ca_certificate&exampleId=1b9fd772-bb42-b2aa-fb8e-348f5bc67501d884a63d&activeTab=example&spm=docs.r.slb_ca_certificate.1.1b9fd772bb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_slb_ca_certificate" "foo-file" {
  ca_certificate_name = "tf-testAccSlbCACertificate"
  ca_certificate      = file("${path.module}/ca_certificate.pem")
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_slb_ca_certificate&spm=docs.r.slb_ca_certificate.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `ca_certificate_name` - (Optional, Available since v1.123.1) Name of the CA Certificate.
* `ca_certificate` - (Required, ForceNew) the content of the CA certificate.
* `resource_group_id` - (Optional, ForceNew, Available since v1.58.0) The Id of resource group which the slb_ca certificate belongs.
* `tags` - (Optional, Available since v1.66.0) A mapping of tags to assign to the resource.
* `name` - (Deprecated) Field `name` has been deprecated from provider version 1.123.1. New field `ca_certificate_name` instead

## Attributes Reference

The following attributes are exported:

* `id` - The Id of CA Certificate .

## Timeouts

-> **NOTE:** Available since v1.123.1

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `delete` - (Defaults to 5 mins) Used when delete the SLB CA Certificate.

## Import

Server Load balancer CA Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_ca_certificate.example abc123456
```
