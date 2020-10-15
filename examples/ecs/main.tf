data "alicloud_instance_types" "instance_type" {
  instance_type_family = "ecs.n4"
  cpu_core_count       = "1"
  memory_size          = "2"
}

resource "alicloud_security_group" "group" {
  name        = var.short_name
  description = "New security group"
  vpc_id      = alicloud_vpc.vpc.id
}

resource "alicloud_security_group_rule" "allow_http_80" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = var.nic_type
  policy            = "accept"
  port_range        = "80/80"
  priority          = 1
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "allow_https_443" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = var.nic_type
  policy            = "accept"
  port_range        = "443/443"
  priority          = 1
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_disk" "disk" {
  availability_zone = alicloud_instance.instance[0].availability_zone
  category          = var.disk_category
  size              = var.disk_size
  count             = var.number
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
}

data "alicloud_zones" "zones_ds" {
  available_instance_type = data.alicloud_instance_types.instance_type.instance_types[0].id
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.zones_ds.zones[0].id
}

resource "alicloud_instance" "instance" {
  instance_name   = "${var.short_name}-${var.role}-${format(var.count_format, count.index + 1)}"
  host_name       = "${var.short_name}-${var.role}-${format(var.count_format, count.index + 1)}"
  image_id        = var.image_id
  instance_type   = data.alicloud_instance_types.instance_type.instance_types[0].id
  count           = var.number
  security_groups = alicloud_security_group.group.*.id
  vswitch_id      = alicloud_vswitch.vswitch.id

  internet_charge_type       = var.internet_charge_type
  internet_max_bandwidth_out = var.internet_max_bandwidth_out

  password = var.ecs_password

  instance_charge_type          = "PostPaid"
  system_disk_category          = "cloud_efficiency"
  system_disk_name              = "${var.short_name}-${var.role}-${format(var.count_format, count.index + 1)}-system-disk"
  system_disk_description       = "${var.short_name}-${var.role}-${format(var.count_format, count.index + 1)}-system-disk-description"
  security_enhancement_strategy = "Deactive"

  data_disks {
    name        = "disk1"
    size        = "20"
    category    = "cloud"
    description = "disk1"
  }
  data_disks {
    name        = "disk2"
    size        = "20"
    category    = "cloud"
    description = "disk2"
  }

  tags = {
    role = var.role
    dc   = var.datacenter
  }
}

resource "alicloud_disk_attachment" "instance-attachment" {
  count       = var.number
  disk_id     = alicloud_disk.disk.*.id[count.index]
  instance_id = alicloud_instance.instance.*.id[count.index]
}

