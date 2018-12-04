data "alicloud_zones" "default" {
  "available_disk_category"     = "cloud_efficiency"
  "available_resource_creation" = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "image" {
  name_regex  = "^ubuntu_14.*_64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}

resource "alicloud_vpc" "main" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id            = "${alicloud_vpc.main.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_security_group" "group" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  image_id                   = "${data.alicloud_images.image.images.0.id}"
  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  count                      = "2"
  security_groups            = ["${alicloud_security_group.group.*.id}"]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb" "instance" {
  name       = "${var.name}"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  name             = "${var.name}"

  servers = [
    {
      server_ids = ["${alicloud_instance.instance.0.id}", "${alicloud_instance.instance.1.id}"]
      port       = 100
      weight     = 10
    },
    {
      server_ids = ["${alicloud_instance.instance.*.id}"]
      port       = 80
      weight     = 100
    },
  ]
}

resource "alicloud_slb_listener" "tcp" {
  load_balancer_id          = "${alicloud_slb.instance.id}"
  backend_port              = "22"
  frontend_port             = "22"
  protocol                  = "tcp"
  bandwidth                 = "10"
  health_check_type         = "tcp"
  persistence_timeout       = 3600
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx"
  health_check_connect_port = 20
  health_check_uri          = "/console"
  established_timeout       = 600
  server_group_id           = "${alicloud_slb_server_group.group.id}"
}
