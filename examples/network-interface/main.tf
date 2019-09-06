resource "alicloud_vpc" "vpc" {
  name       = "tf-testAcc-vpc"
  cidr_block = var.vpc_cidr
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
  name              = "tf-testAcc-vswitch"
  cidr_block        = var.vswitch_cidr
  availability_zone = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.vpc.id
}

resource "alicloud_security_group" "sg" {
  name   = "tf-testAcc-sg"
  vpc_id = alicloud_vpc.vpc.id
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  most_recent = var.most_recent
  owners      = var.image_owners
  name_regex  = var.name_regex
}

resource "alicloud_instance" "instance" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  security_groups   = [alicloud_security_group.sg.id]

  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  system_disk_category       = var.system_disk_category
  image_id                   = data.alicloud_images.default.images[0].id
  instance_name              = "tf-testAcc-i"
  vswitch_id                 = alicloud_vswitch.vswitch.id
  internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface" "eni" {
  name            = "tf-testAcc-eni"
  vswitch_id      = alicloud_vswitch.vswitch.id
  security_groups = [alicloud_security_group.sg.id]
}

resource "alicloud_network_interface_attachment" "at" {
  instance_id          = alicloud_instance.instance.id
  network_interface_id = alicloud_network_interface.eni.id
}

data "alicloud_network_interfaces" "enis" {
  ids = [alicloud_network_interface.eni.id]
}

