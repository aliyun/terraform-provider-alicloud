variable "name" {
  default = "tf-testAccSlbRuleBasic"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}_test"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name       = var.name
  vswitch_id = alicloud_vswitch.default.id
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 22
  frontend_port             = 22
  protocol                  = "http"
  bandwidth                 = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id

  servers {
    server_ids = alicloud_instance.default.*.id
    port       = 80
    weight     = 100
  }
}

resource "alicloud_slb_rule" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  frontend_port             = alicloud_slb_listener.default.frontend_port
  name                      = var.name
  domain                    = "*.aliyun.com"
  url                       = "/image"
  server_group_id           = alicloud_slb_server_group.default.id
  cookie                    = "23ffsa"
  cookie_timeout            = 100
  health_check_http_code    = "http_2xx"
  health_check_interval     = 10
  health_check_uri          = "/test"
  health_check_connect_port = 80
  health_check_timeout      = 30
  healthy_threshold         = 4
  unhealthy_threshold       = 4
  sticky_session            = "on"
  sticky_session_type       = "server"
  listener_sync             = "off"
  scheduler                 = "rr"
  health_check_domain       = "test"
  health_check              = "on"
}

