provider "alicloud" {
  region = var.region
}

data "alicloud_instance_types" "c2g4" {
  cpu_core_count = 4
  memory_size    = 8
}

data "alicloud_images" "centos" {
  most_recent = true
  name_regex  = "^centos_7\\w.*"
}

data "alicloud_zones" "default" {
  available_instance_type = data.alicloud_instance_types.c2g4.instance_types[0].id
  available_disk_category = var.disk_category
}

resource "alicloud_vpc" "default" {
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "vsw" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = var.vswitch_cidr
  zone_id    = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_security_group" "sg" {
  name   = "sg"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_security_group_rule" "in-all" {
  type              = "ingress"
  ip_protocol       = "all"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = 1
  security_group_id = alicloud_security_group.sg.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "en-all" {
  type              = "egress"
  ip_protocol       = "all"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = 1
  security_group_id = alicloud_security_group.sg.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "webserver" {
  security_groups = [alicloud_security_group.sg.id]
  vswitch_id      = alicloud_vswitch.vsw.id

  # series II
  instance_charge_type       = "PostPaid"
  instance_type              = data.alicloud_instance_types.c2g4.instance_types[0].id
  internet_max_bandwidth_out = 0

  system_disk_category = var.disk_category
  image_id             = data.alicloud_images.centos.images[0].id

  instance_name = "tf_lnmp"
  password      = var.ecs_password

  user_data = data.template_file.shell.rendered
}

data "template_file" "shell" {
  template = file("userdata.sh")

  vars = {
    db_name     = var.db_name
    db_user     = var.db_user
    db_pwd      = var.db_password
    db_root_pwd = var.db_root_password
  }
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = alicloud_vpc.default.id
  specification = "Small"
}

resource "alicloud_eip_address" "default" {
  count     = 2
  bandwidth = 10
}

resource "alicloud_eip_association" "default" {
  count         = 2
  allocation_id = alicloud_eip_address.default.*.id[count.index]
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_forward_entry" "dnat" {
  forward_table_id = alicloud_nat_gateway.default.forward_table_ids
  external_ip      = alicloud_eip_address.default[1].ip_address
  external_port    = "any"
  ip_protocol      = "any"
  internal_ip      = alicloud_instance.webserver.private_ip
  internal_port    = "any"
}

resource "alicloud_snat_entry" "snat" {
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vsw.id
  snat_ip           = alicloud_eip_address.default[0].ip_address
}

