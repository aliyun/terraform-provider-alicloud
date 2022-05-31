variable "name" {
  default = "tftest"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name           = var.name
  topic_quota    = 50
  disk_type      = 1
  disk_size      = 500
  deploy_type    = 5
  io_max         = 20
  vswitch_id     = data.alicloud_vswitches.default.ids.0
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_instance_allowed_ip_attachment" "default" {
  allowed_ip   = "114.237.9.78/32"
  allowed_type = "vpc"
  instance_id  = alicloud_alikafka_instance.default.id
  port_range   = "9092/9092"
}
