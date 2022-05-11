---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_listener_additional_certificate_attachment"
sidebar_current: "docs-alicloud-resource-alb-listener-additional-certificate-attachment"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Listener Additional Certificate Attachment resource.
---

# alicloud\_alb\_listener\_additional\_certificate\_attachment"

Provides a Application Load Balancer (ALB) Listener Additional Certificate Attachment resource.

For information about Application Load Balancer (ALB) Listener Additional Certificate Attachment and how to use it, see [What is Listener Additional Certificate Attachment](https://www.alibabacloud.com/help/en/doc-detail/302356.html).

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc-alb"
}

data "alicloud_alb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}


data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}

data "alicloud_resource_manager_resource_groups" "default" {}


resource "alicloud_alb_load_balancer" "default_3" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Standard"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default_1.ids[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default_2.ids[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
}

resource "alicloud_alb_server_group" "default_4" {
  protocol          = "HTTP"
  vpc_id            = data.alicloud_vpcs.default.vpcs.0.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
  tags = {
    Created = "TF"
  }
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  count            = 2
  certificate_name = join("", [var.name, count.index])
  cert             = file("${path}/test.crt")
  key              = file("${path}/test.key")
}

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default_3.id
  listener_protocol    = "HTTPS"
  listener_port        = 8081
  listener_description = var.name
  default_actions {
    type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default_4.id
      }
    }
  }
  certificates {
    certificate_id = join("", [alicloud_ssl_certificates_service_certificate.default.0.id, "-cn-hangzhou"])
  }
}

resource "alicloud_alb_listener_additional_certificate_attachment" "default" {
  certificate_id = join("", [alicloud_ssl_certificates_service_certificate.default.1.id, "-cn-hangzhou"])
  listener_id    = alicloud_alb_listener.default.id
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, ForceNew) The Certificate ID.
* `listener_id` - (Required, ForceNew) The ID of the ALB listener.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Listener Additional Certificate Attachment. The value formats as `<listener_id>:<certificate_id>`.
* `status` - The status of the certificate.
* `certificate_type` - The type of the certificate.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Listener.
* `delete` - (Defaults to 2 mins) Used when delete the Listener.


## Import

Application Load Balancer (ALB) Listener Additional Certificate Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_alb_listener_additional_certificate_attachment.example <listener_id>:<certificate_id>
```