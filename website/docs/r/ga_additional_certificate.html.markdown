---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_additional_certificate"
sidebar_current: "docs-alicloud-resource-ga-additional-certificate"
description: |-
  Provides a Alicloud Global Accelerator (GA) Additional Certificate resource.
---

# alicloud\_ga\_additional\_certificate

Provides a Global Accelerator (GA) Additional Certificate resource.

For information about Global Accelerator (GA) Additional Certificate and how to use it, see [What is Additional Certificate](https://www.alibabacloud.com/help/en/doc-detail/302356.html).

-> **NOTE:** Available in v1.150.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc-ga"
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth       = 20
  type            = "Basic"
  bandwidth_type  = "Basic"
  duration        = 1
  ratio           = 30
  auto_pay        = true
  auto_use_coupon = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  count            = 2
  certificate_name = var.name
  cert             = file("${path}/test.crt")
  key              = file("${path}/test.key")
}

resource "alicloud_ga_listener" "default" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.default]
  accelerator_id = alicloud_ga_accelerator.default.id
  name           = var.name
  protocol       = "HTTPS"
  port_ranges {
    from_port = 8080
    to_port   = 8080
  }
  certificates {
    id = join("-", [alicloud_ssl_certificates_service_certificate.default.0.id, "cn-hangzhou"])
  }
}

resource "alicloud_ga_additional_certificate" "default" {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.1.id, "cn-hangzhou"])
  domain         = "test"
  accelerator_id = alicloud_ga_accelerator.default.id
  listener_id    = alicloud_ga_listener.default.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the GA instance.
* `certificate_id` - (Required, ForceNew) The Certificate ID.
* `domain` - (Required, ForceNew) The domain name specified by the certificate. **NOTE:** You can associate each domain name with only one additional certificate.
* `listener_id` - (Required, ForceNew) The ID of the listener. **NOTE:** Only HTTPS listeners support this parameter.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used to wait accelerator and listener to be active after creating the Ga additional certificate.
* `delete` - (Defaults to 1 mins) Used to wait accelerator and listener to be active after deleting the Ga additional certificate


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Additional Certificate. The value formats as `<accelerator_id>:<listener_id>:<domain>`.

## Import

Global Accelerator (GA) Additional Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_ga_additional_certificate.example <accelerator_id>:<listener_id>:<domain>
```