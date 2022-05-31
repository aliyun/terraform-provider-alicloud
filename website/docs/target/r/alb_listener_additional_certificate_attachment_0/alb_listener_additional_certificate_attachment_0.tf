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
