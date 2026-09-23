resource "alicloud_vpc" "main" {
  vpc_name   = var.long_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "main" {
  vpc_id     = alicloud_vpc.main.id
  count      = length(split(",", var.availability_zones))
  cidr_block = var.cidr_blocks["az${count.index}"]
  zone_id    = split(",", var.availability_zones)[count.index]

  depends_on = [alicloud_vpc.main]
}

resource "alicloud_slb_load_balancer" "instance" {
  load_balancer_name                 = var.name
  vswitch_id           = alicloud_vswitch.main[0].id
  internet_charge_type = var.internet_charge_type
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id = alicloud_slb_load_balancer.instance.id
  backend_port     = "2111"
  frontend_port    = "21"
  protocol         = "tcp"
  bandwidth        = "5"
}

