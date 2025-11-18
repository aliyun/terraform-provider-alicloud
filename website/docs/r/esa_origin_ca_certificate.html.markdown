---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_origin_ca_certificate"
description: |-
  Provides a Alicloud ESA Origin Ca Certificate resource.
---

# alicloud_esa_origin_ca_certificate

Provides a ESA Origin Ca Certificate resource.



For information about ESA Origin Ca Certificate and how to use it, see [What is Origin Ca Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/UploadOriginCaCertificate).

-> **NOTE:** Available since v1.256.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "example" {
  site_name   = "bcd.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}



resource "alicloud_esa_origin_ca_certificate" "default" {
  site_id     = alicloud_esa_site.example.id
  certificate = "-----BEGIN CERTIFICATE-----\nMIIDRTCCAi2gAwIBAgIUHRPTIPKP2zN9on/NCzBe0BV68UUwDQYJKoZIhvcNAQEF\nBQAwMzEPMA0GA1UEAwwGU1NMZXllMRMwEQYDVQQKDApTU0xleWUgSW5jMQswCQYD\nVQQGEwJDTjAeFw0yNTA3MzAwODQzMDBaFw0yNTEyMzEwODQwMDBaMGQxCzAJBgNV\nBAYTAkNOMQ8wDQYDVQQIDAbljJfkuqwxEDAOBgNVBAcMB0JlaWppbmcxGzAZBgNV\nBAoMEuenkeaKgOaciemZkOWFrOWPuDEVMBMGA1UEAwwMZ29zaXRlY2RuLmNuMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtu2oW3t2bj9LsFnXj1C2EmaR\nJYJwNgHsTBKl3DxeL2+Ext0qN2Z+UgTqYM1c1HOdwN9x13pnAVe4PmiLAkxpp/4u\n5gKsH1+6p3aXFUk0NvEoLXfESoQpyvoB0o/8oryxNs3+iUfvAk+a7IKAr99a1P9F\nTkpyE6t+dgSLYhHc49ZRdYImmZcYQLmpygYOwWBdv6hlQUFi/tvX16fRZ0GgyUOK\n7xsTWG6qUhPJyLRtj9zn+0khgh5DJhfJQ4KTWZMX63UPiIx7sPu9sR+TPWqJsEuq\nVipxouMys+NNMjDtn55+PE6/sDbkvULHeFUglGMZ9qHcl3ej31zmkhu6bmvNcQID\nAQABoyAwHjALBgNVHREEBDACggAwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0B\nAQUFAAOCAQEAF6J9TdaDYQ96EaWvb2ttQ6jNrDe4k3t1cdfhPEWMJzxZFxoDBYZ2\nAl9vB2JICEsGDkCwpqYz2UXJsGnq2rHjUxouYo1568K/loownWjwdCgdLGbQpnXY\nQeqPSTRLT71ikH+RqCpoYxcN63i3j9oYWm9KoD5F4arcqlLrEUZ1TqW5csGSY1h6\n2HmGPsINl9KCxwUS+76dxsdHIqLFx0qdnD6S5vmd0sin33jdYhj9ltp0KvhEgMvS\nXMuzECVRvI4MZxebf7gkV3EByqV6XvazBSxuMhplygpAaLra11yV1M/m9wzVwlnS\nS2GNvRkNym9WnH0IQ0kn9hS8hj52Nh12JQ==\n-----END CERTIFICATE-----"
  name        = "example"
}
```

## Argument Reference

The following arguments are supported:
* `certificate` - (Required, ForceNew) Certificate content.
* `name` - (Optional, ForceNew) The certificate name.
* `site_id` - (Required, ForceNew) Site Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<origin_ca_certificate_id>`.
* `create_time` - Creation time.
* `origin_ca_certificate_id` - OriginCaCertificate Id
* `status` - Certificate status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Origin Ca Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Origin Ca Certificate.

## Import

ESA Origin Ca Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_origin_ca_certificate.example <site_id>:<origin_ca_certificate_id>
```