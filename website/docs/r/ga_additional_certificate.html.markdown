---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_additional_certificate"
sidebar_current: "docs-alicloud-resource-ga-additional-certificate"
description: |-
  Provides a Alicloud Global Accelerator (GA) Additional Certificate resource.
---

# alicloud_ga_additional_certificate

Provides a Global Accelerator (GA) Additional Certificate resource.

For information about Global Accelerator (GA) Additional Certificate and how to use it, see [What is Additional Certificate](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-associateadditionalcertificateswithlistener).

-> **NOTE:** Available since v1.150.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_additional_certificate&exampleId=1ad68600-f4be-90f7-92bd-c8f14d5a3648e7cc5837&activeTab=example&spm=docs.r.ga_additional_certificate.0.1ad68600f4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region = var.region
}

variable "name" {
  default = "tf-example"
}


data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth              = 100
  type                   = "Basic"
  bandwidth_type         = "Basic"
  payment_type           = "PayAsYouGo"
  billing_type           = "PayBy95"
  ratio                  = 30
  bandwidth_package_name = var.name
  auto_pay               = true
  auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  count            = 2
  certificate_name = join("-", [var.name, count.index])
  cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7zCCAtegAwIBAgIRAKi2/Fx1cUTyhV839x42ockwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjMwODA5MDQ1NDU3WhcNMjYwODA4MDQ1NDU3WjAsMQswCQYDVQQGEwJDTjEd
MBsGA1UEAxMUYWxpY2xvdWQtcHJvdmlkZXIuY24wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQDdkot9e0pMCTPAtA29Sz5sF+aPT/l9+3sOnQeJ1kKLNkqK
iQgwADexoAqlmTaZM03gh/GnkqPw9gxN/fJHWdVzxE03Fs8bKgMdS6cf0v/xArrQ
zm6N4vmsbuE8SX2eu303PAsyBMqPByTODZ5i+5LkZcrxMFQsbA3xnBouzS5e+T+a
7YTyyVv5WDy871/sdRAYTfnUttdnqkKGeMKgQgRlJ2pDk5/k2iwmQmSh/wbk465+
1U5w2npPYGPvGAkzl7RRc4/VckqlV8P0cmgguqIRyllJwFEnvcpqpOHTxBOBq9iZ
4b/h7ynrfB/GbAw574eSEl0gzLBW60bT9YedbTeXAgMBAAGjgdkwgdYwDgYDVR0P
AQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME
GDAWgBQogSYF0TQaP8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYB
BQUHMAGGFWh0dHA6Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDov
L2NhLm15c3NsLmNvbS9teXNzbHRlc3Ryc2EuY3J0MB8GA1UdEQQYMBaCFGFsaWNs
b3VkLXByb3ZpZGVyLmNuMA0GCSqGSIb3DQEBCwUAA4IBAQCwUBeznv6cAjcTLCDb
SSvgkM9HFcbWnuGS8Nf5P4YfmSs52VuHZyjzwphjAU6B/danI/nMdZe52PXyvjVV
02Y8ld/tMpqPV5SpaOadLtdg6TGBNJieOAt9doM8WNEgq/JycAL9ivIOjChUetZf
ZEV7HDIgiHSpqAPWMZYL71MS/p5zYkyOnPqmGyLNdi1neotwVCQopQXRNC2iLlVV
yQONfXH5iijqr1iTWkB0ESK/xBt1PB655PlTjzFQUOovE1SyoQS8K3u7TP6+BqtD
G9TYNTNZvxl5I/iU/KdWVip+qJbxRA8Skc8gHkkzeIEStw3l5cjnrp9h7EhnhkOh
ltGN
-----END CERTIFICATE-----
EOF
  key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA3ZKLfXtKTAkzwLQNvUs+bBfmj0/5fft7Dp0HidZCizZKiokI
MAA3saAKpZk2mTNN4Ifxp5Kj8PYMTf3yR1nVc8RNNxbPGyoDHUunH9L/8QK60M5u
jeL5rG7hPEl9nrt9NzwLMgTKjwckzg2eYvuS5GXK8TBULGwN8ZwaLs0uXvk/mu2E
8slb+Vg8vO9f7HUQGE351LbXZ6pChnjCoEIEZSdqQ5Of5NosJkJkof8G5OOuftVO
cNp6T2Bj7xgJM5e0UXOP1XJKpVfD9HJoILqiEcpZScBRJ73KaqTh08QTgavYmeG/
4e8p63wfxmwMOe+HkhJdIMywVutG0/WHnW03lwIDAQABAoIBAQCe5rHS09B8pzzO
PlJ8JrIlox5eOOScTPX7jPITD+25GL5si8mrYvyODlCUYkSdqgV3uQa9PpUEAfDh
HfXa5boGxAj8MQdmW8LQB6lbUV7r4SFJDkKKzvRvjTVKnwnQBHXQXudIf9ckq+Lh
QzMLmY/G7JmWTyqOkQ+O7nx4g/11bcU7uQrQdvWPfc0+IiT1TYQdyLQ/Chlj3RF/
iwF8ZL2sfKF+Z5O49+Q6cXvUcQOvqtkIXbQijayyVNBMJwDB7aOZRA7JBNj9/ib6
N0iTo81dJVz/nnpbWRaFTVinIsDF1heDfQ1qDx06T/Mpi6pjoWjRUcyIHEbZJTel
0nXDJD1BAoGBAPZB/PN8MP+o9gkf2jnoU9LzctDJrQwD1J2XElq4RomimPIMqDQP
5TRAJThf0O0X4Mv2n9EzV457OpJL+fz9htRWEYogWl9bkbzZ1AoX4K/acuGeawTT
YEhPjJ2ZETsBsCeDkDDuHHzYwRQv+EfoXH36z9PBDxG1ZDb7kWwAILXdAoGBAOZW
jXG7m4I7cxUtXGtjwydh4K7nwH/5QoH2m928HM2AT48eQCl3CMQ089+qeJGgfHQv
GyVOO/FGhcFsFi10FMQ7IlwWgZODg64qnrNhi4zbV1M2wKem1T2dlEpkd82EFdnS
GYRIEkFORMxEDyzx3Th2TajpWC8YKKG3Tnm0bQ4DAoGBAIZTEEswHvoVi78GZN7Z
X3/d028X0xCOtlcPpK9ffPpuesbtKILdeMS7iJHrkecB81jOOfa+7q+FgDl0v/PD
xtvj5sVVSHZjWGeO2h53T9QccDWpV+7V7dsDqUv9xmxNS20CUpCeEWP4R7lfQSrY
EDuXp+11jWa3buae6n/iwfTxAoGABEYW2cVhXUk9GWd+D4AKXvCx+ozSRY2abk7l
FXgoEKgQ0db92ccboohY/g1rr0gLBxzYpBiPhCqK0MvwnWdJ+1odiRfhz5rhFpoz
16A3tqVbOXAKoxG1Yy9JURgMIQQSY7hCQPIVZKDPJfsdTPgv4pxPVJL/z9/i4R1F
l3yBiYECgYEA0+vpzL24nHZYdwgBF4qbmYhv8baRi07/BNgV1+d6vESuO/MwwoE/
2UZ9Drf5yoX2Bvi5/vVMbyc7cSluO7icPBkl0D8F7E3x0v5mzwPxtpR8BTRoJKOL
/rMdLscMz2VQsL5DJd/9OZg60fHRaRtWtV0afXzL5zUxnfDLot24IG4=
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  name           = var.name
  protocol       = "HTTPS"
  port_ranges {
    from_port = 8080
    to_port   = 8080
  }
  certificates {
    id = join("-", [alicloud_ssl_certificates_service_certificate.default.1.id, var.region])
  }
}

locals {
  domain = "alicloud-provider.cn"
}

resource "alicloud_ga_additional_certificate" "default" {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.1.id, var.region])
  domain         = local.domain
  accelerator_id = alicloud_ga_listener.default.accelerator_id
  listener_id    = alicloud_ga_listener.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_additional_certificate&spm=docs.r.ga_additional_certificate.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `listener_id` - (Required, ForceNew) The ID of the listener. **NOTE:** Only HTTPS listeners support this parameter.
* `domain` - (Required, ForceNew) The domain name specified by the certificate. **NOTE:** You can associate each domain name with only one additional certificate.
* `certificate_id` - (Required) The Certificate ID. **NOTE:** From version 1.209.1, `certificate_id` can be modified.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Additional Certificate. The value format as `<accelerator_id>:<listener_id>:<domain>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used to wait accelerator and listener to be active after creating the Ga additional certificate.
* `update` - (Defaults to 3 mins) Used to wait accelerator and listener to be active after updating the Ga additional certificate.
* `delete` - (Defaults to 3 mins) Used to wait accelerator and listener to be active after deleting the Ga additional certificate.

## Import

Global Accelerator (GA) Additional Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_additional_certificate.example <accelerator_id>:<listener_id>:<domain>
```
