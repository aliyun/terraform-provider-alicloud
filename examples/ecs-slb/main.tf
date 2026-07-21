data "alicloud_zones" "zone" {
  available_instance_type = var.ecs_type
}

resource "alicloud_security_group" "group" {
  name        = var.short_name
  description = "New security group"
}

resource "alicloud_security_group_rule" "http-in" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "internet"
  policy            = "accept"
  port_range        = "80/80"
  priority          = 1
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "https-in" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "internet"
  policy            = "accept"
  port_range        = "443/443"
  priority          = 1
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "ssh-in" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "internet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.group.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "instance" {
  instance_name              = "${var.short_name}-${var.role}-${format(var.count_format, count.index + 1)}"
  host_name                  = "${var.short_name}-${var.role}-${format(var.count_format, count.index + 1)}"
  image_id                   = var.image_id
  instance_type              = var.ecs_type
  count                      = var.number
  security_groups            = alicloud_security_group.group.*.id
  internet_charge_type       = var.internet_charge_type
  internet_max_bandwidth_out = var.internet_max_bandwidth_out
  password                   = var.ecs_password
  availability_zone          = data.alicloud_zones.zone.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"

  tags = {
    role = var.role
    dc   = var.datacenter
  }
}

resource "alicloud_slb_load_balancer" "instance" {
  load_balancer_name                 = var.slb_name
  internet_charge_type = var.slb_internet_charge_type
  internet             = var.internet
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id = alicloud_slb_load_balancer.instance.id
  backend_port     = 2111
  frontend_port    = 21
  protocol         = "tcp"
  bandwidth        = 5
}

resource "alicloud_slb_attachment" "default" {
  load_balancer_id = alicloud_slb_load_balancer.instance.id
  instance_ids     = alicloud_instance.instance.*.id
}

