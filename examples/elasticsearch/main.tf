resource "alicloud_elasticsearch_instance" "instance" {
  instance_charge_type     = var.instance_charge_type
  data_node_amount         = var.data_node_amount
  data_node_spec           = var.data_node_spec
  data_node_disk_size      = var.data_node_disk_size
  data_node_disk_type      = var.data_node_disk_type
  data_node_disk_encrypted = var.data_node_disk_encrypted
  vswitch_id               = var.vswitch_id
  password                 = var.password
  version                  = var.es_version
  description              = var.description
  zone_count               = var.zone_count
  master_node_spec         = var.master_node_spec
  client_node_amount       = var.client_node_amount
  client_node_spec         = var.client_node_spec
  protocol                 = var.protocol
  setting_config           = var.setting_config
  warm_node_spec            = var.warm_node_spec
  warm_node_amount          = var.warm_node_amount
  warm_node_disk_size       = var.warm_node_disk_size
  warm_node_disk_type       = var.warm_node_disk_type
  kibana_node_spec          = var.kibana_node_spec,
  enable_kibana_private_network = var.enable_kibana_private_network
  enable_kibana_public_network  = var.enable_kibana_public_network
  kibana_private_security_group_id = var.kibana_private_security_group_id
}

