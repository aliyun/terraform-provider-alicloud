variable "name" {
  default = "tf-testaccvirtualnode"
}

data "alicloud_eci_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_eci_zones.default.zones.0.zone_ids.1
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name   = var.name
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eci_virtual_node" "default" {
  security_group_id     = alicloud_security_group.default.id
  virtual_node_name     = var.name
  vswitch_id            = data.alicloud_vswitches.default.ids.1
  enable_public_network = false
  eip_instance_id       = alicloud_eip_address.default.id
  resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  kube_config           = "kube config"
  tags = {
    Created = "TF"
  }
  taints {
    effect = "NoSchedule"
    key    = "Tf1"
    value  = "Test1"
  }
}
