---
subcategory: "SSL Certificates"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_certificate"
description: |-
  Provides a Alicloud SSL Certificates Certificate resource.
---

# alicloud_ssl_certificates_service_certificate

Provides a SSL Certificates Certificate resource.



For information about SSL Certificates Certificate and how to use it, see [What is Certificate](https://www.alibabacloud.com/help/product/28533.html).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_certificate&exampleId=484f1229-91fa-a3e8-ded0-68c2f1583a7561832250&activeTab=example&spm=docs.r.ssl_certificates_service_certificate.0.484f122991&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = "terraform-example-${random_integer.default.result}"
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID1jCCAr6gAwIBAgIQQ7/8/QOOTbywxdgSX9aMqDANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNTA5MjIwNTU3NDVaFw0zMDA5MjEwNTU3NDVaMCAxCzAJBgNVBAYTAkNOMREw
DwYDVQQDEwgxNjg4LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AMEl04gKBqJxV+8KideZb7S4mPysehPzr/cXu4i1RXT7UFtNVZuqc4IdIzOja2SU
6uNn8mY6Pfc5FNybg98bYx0ADbub55TUaw2Pz1CFEbiMvLpzMkp4EZadvmJWZk8t
dNb+ClKqdXUWhxApS3Lz+wjCNYQnlODk4KmxmM8/U/CyQS7lgWS/1G72UFB09Skg
sfvWdoHLrFfIlbVkp9XVELCtOkjj8Nn/rPOhc31NbstrwV4Whl6jngGAkaEtImJ7
//sL+sPPsutefCgfZPrC+Zwru2En1BuIo5KW02NYLdjXbABH8xjkUobqRoro7eY3
VySBr7adD6QmNv5hWohOuykCAwEAAaOBzTCByjAOBgNVHQ8BAf8EBAMCBaAwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiBJgXRNBo/
wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYVaHR0cDov
L29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlzc2wuY29t
L215c3NsdGVzdHJzYS5jcnQwEwYDVR0RBAwwCoIIMTY4OC5jb20wDQYJKoZIhvcN
AQELBQADggEBAHa0ATVeHtPPw1+a6kajlW6OQUjhiJg+Sk9fVA1eJ2Hzl1yDDw3K
yAyl1gkxGI6BwWdX/C8IE6PuPYcG2CmJGoFoEAAIbAE76AKABvHoA8I6wyDruxFz
06bNM8104TxAHTxe2zaHgBQnYIRk07uA8gxjZKFp1//eYbxj8HiP0Q9zXqYjF79G
Le4PDw7Q6U22CP+cT9Sz5ZEoJCzmUtx3uQWhLzNxvyISrXeSqAFJzjtL0KKSR1cr
8he6FoeU37oKdmrnweLeBe+no3OMChETa2JN4VAzXj/nPpQcyB7nXDfLUHe01+BB
ZBXKFLD2H38e97mFl/7mgNP5Nc1sycI5Sp4=
-----END CERTIFICATE-----
EOF
  key              = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEowIBAAKCAQEAwSXTiAoGonFX7wqJ15lvtLiY/Kx6E/Ov9xe7iLVFdPtQW01V
m6pzgh0jM6NrZJTq42fyZjo99zkU3JuD3xtjHQANu5vnlNRrDY/PUIURuIy8unMy
SngRlp2+YlZmTy101v4KUqp1dRaHEClLcvP7CMI1hCeU4OTgqbGYzz9T8LJBLuWB
ZL/UbvZQUHT1KSCx+9Z2gcusV8iVtWSn1dUQsK06SOPw2f+s86FzfU1uy2vBXhaG
XqOeAYCRoS0iYnv/+wv6w8+y6158KB9k+sL5nCu7YSfUG4ijkpbTY1gt2NdsAEfz
GORShupGiujt5jdXJIGvtp0PpCY2/mFaiE67KQIDAQABAoIBAAKF9CZTUd8zvDKE
azo/Ur0Zf5omxgOBC/vzj0DLyXKr89KgMdhHmPG1YBKFIIU0XYCHXkclR05LAcbu
BdeCJpXS5zBbwDdAB9P/XHXQqeNvfJRc++ZgJ4QAXzkuqBssXK87ALcwFeUShxot
cphiWpW0inlwVkVn3WLUzfUV0+ARljn8VOf+aAmfCiQMl4gsBpvD3dxF84aihS+1
blqar5dE1GCJWHW67R1uSaAqHf7nwbBkZY8nTWF8n4+ELAAtlOgQKZlrQ+JxB3Ar
rWzgMj4M6F1/man1y/XPR56px9Xv3DwBZHuLufsqPr10q/nI9VIIQHe49sFgnN4+
48Q7wIECgYEAwxlrgBJI8gua4mJZxJRT8gBv2Mb1Kk1k7HVX11I+yF4eXr+cm+24
Cq7MjqmBXSnqvdQkwGFZ+C3cTKXJBPONWGF8NgiXaHSKjPEoFuHLdKBpgZMAax/L
aZBQRw6g12nz3XUCK0DE0wGgPkoDxc65s4NEWS+ua43LZ4TUOzWwwWECgYEA/XB1
ARNHyARy+P3iTeebh3t7qJoNoptLWHMlKjSjIZ1VZ4+9ilKsi5ZKVkPaLIjo8MGv
Ank3vzSrFSYhId0XfmSqoWySWc0eBkc6NERvopxuIV1WwRKf/18lLhxiEjHIcgds
G2KmfeiXdCKSgGlWvJmLITY4gJpOYMjpEDxipskCgYAdxnljmGbNmfvPZRcyKzkM
jAiF2wd7p0gp1lbLo9+1ELgt2ax7F7Ko3riVZUU7BLSwt/nL6o+iks02XW7qdIkz
3dzpGjKRXIfwrrVhmKBGclzny5mav8V5nO7DiXX+qkrvl3X3R/FCCtN77ivZOo2Y
2gXKXr6N55wNdnY1eyI4wQKBgQDXjZo2O+vFVuNimqyrjd1eMcxO7hfCwUooBGcL
qpFEucg1uK+Awig24LCBBly9nARjIJh1Bhw/58/KwQ9U+fJNcdkeSnV/I1HyDQqY
AczhBSM2BWkP9YNXc9jvivxudSECuwVblV/9nqGSCQWJag53gjAvIyqTVqpq7vYq
9PEC4QKBgGY2pj0ZNqGkq16jD3iS+DDBpX+TPnoHzu5GZCM/1GLZ6xXbpNWtZQt4
/m+6koRWeGvNAULnp8RSnhBzm+ZglpbwYcvsqRNDqIPGhJ2JruVA/bY3S0ebkRlD
xDn0dJVMvNyRR83ZpjTQhxoq5l56TN5xk1vdJ9nZdwJMmXiz2TrA
-----END PRIVATE KEY-----
EOF
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ssl_certificates_service_certificate&spm=docs.r.ssl_certificates_service_certificate.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `cert` - (Optional, ForceNew) The content of a non-SM certificate in PEM format.
* `key` - (Optional, ForceNew) The private key of a non-SM certificate in PEM format.
* `encrypt_cert` - (Optional, ForceNew, Available since v1.260.1) The content of an SM encryption certificate in PEM format.
* `encrypt_private_key` - (Optional, ForceNew, Available since v1.260.1) The private key of an SM encryption certificate in PEM format.
* `sign_cert` - (Optional, ForceNew, Available since v1.260.1) The content of an SM signing certificate in PEM format.
* `sign_private_key` - (Optional, ForceNew, Available since v1.260.1) The private key of an SM signing certificate in PEM format.
* `certificate_name` - (Optional) A custom name for the certificate. The name can be up to 64 characters long and can contain any character type, such as letters, numbers, and underscores. **NOTE:** From version 1.260.1, `certificate_name` can be modified.
* `resource_group_id` - (Optional, Available since v1.260.1) The ID of the resource group.
* `tags` - (Optional, Map, Available since v1.260.1) The tag of the resource.
* `name` - (Optional, Deprecated since v1.129.0) Field `name` has been deprecated from provider version 1.129.0 and it will be removed in the future version. Please use the new attribute `certificate_name` instead.
* `lang` - (Deprecated since v1.260.1) Field `lang` has been deprecated from provider version 1.260.1 and it will be removed in the future version.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

-> **NOTE:** Available since 1.260.1.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate.
* `update` - (Defaults to 5 mins) Used when update the Certificate.

## Import

SSL Certificates Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_certificate.example <id>
```
