resource "alicloud_vpc" "this" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_security_group" "default" {
  name   = var.security_group_name
  vpc_id = var.vpc_id == "" ? alicloud_vpc.this.0.id : var.vpc_id
}

resource "alicloud_security_group_rule" "http-in" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "80/80"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "ssh-in" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
}