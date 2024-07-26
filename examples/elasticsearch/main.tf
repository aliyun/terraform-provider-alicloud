provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  # If not set, cn-beijing will be used.
  region = var.region
}

data "alicloud_elasticsearch_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.description
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.description
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_elasticsearch_zones.default.zones[0].id
}

resource "alicloud_elasticsearch_instance" "instance" {
  instance_charge_type     = var.instance_charge_type
  data_node_amount         = var.data_node_amount
  data_node_spec           = var.data_node_spec
  data_node_disk_size      = var.data_node_disk_size
  data_node_disk_type      = var.data_node_disk_type
  data_node_disk_encrypted = var.data_node_disk_encrypted
  vswitch_id               = alicloud_vswitch.default.id
  password                 = var.password
  version                  = var.es_version
  description              = var.description
  zone_count               = var.zone_count
  master_node_spec         = var.master_node_spec
  master_node_disk_type    = var.master_node_disk_type
  client_node_amount       = var.client_node_amount
  client_node_spec         = var.client_node_spec
  protocol                 = var.protocol
  setting_config           = var.setting_config
  warm_node_spec            = var.warm_node_spec
  warm_node_amount          = var.warm_node_amount
  warm_node_disk_size       = var.warm_node_disk_size
  warm_node_disk_type       = var.warm_node_disk_type
  kibana_node_spec          = var.kibana_node_spec
  enable_kibana_public_network  = var.enable_kibana_public_network
}

