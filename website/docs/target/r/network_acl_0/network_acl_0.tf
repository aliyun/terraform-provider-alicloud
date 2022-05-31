data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "VpcConfig"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  vswitch_name = "vswitch"
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 4, 4)
  zone_id      = data.alicloud_zones.default.ids.0
}

resource "alicloud_network_acl" "default" {
  vpc_id           = alicloud_vpc.default.id
  network_acl_name = "network_acl"
  description      = "network_acl"
  ingress_acl_entries {
    description            = "tf-testacc"
    network_acl_entry_name = "tcp23"
    source_cidr_ip         = "196.168.2.0/21"
    policy                 = "accept"
    port                   = "22/80"
    protocol               = "tcp"
  }
  egress_acl_entries {
    description            = "tf-testacc"
    network_acl_entry_name = "tcp23"
    destination_cidr_ip    = "0.0.0.0/0"
    policy                 = "accept"
    port                   = "-1/-1"
    protocol               = "all"
  }
  resources {
    resource_id   = alicloud_vswitch.default.id
    resource_type = "VSwitch"
  }
}
