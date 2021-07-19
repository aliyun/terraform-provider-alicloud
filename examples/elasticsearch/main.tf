resource "alicloud_elasticsearch_instance" "instance" {
  instance_charge_type             = var.instance_charge_type
  data_node_amount                 = var.data_node_amount
  data_node_spec                   = var.data_node_spec
  data_node_disk_size              = var.data_node_disk_size
  data_node_disk_type              = var.data_node_disk_type
  data_node_disk_encrypted         = var.data_node_disk_encrypted
  data_node_disk_performance_level = var.data_node_disk_performance_level
  vswitch_id                       = var.vswitch_id
  password                         = var.password
  version                          = var.es_version
  description                      = var.description
  zone_count                       = var.zone_count
  master_node_spec                 = var.master_node_spec
  master_node_disk_type            = var.master_node_disk_type
  client_node_amount               = var.client_node_amount
  client_node_spec                 = var.client_node_spec
  warm_data_node_amount            = var.warm_data_node_amount
  warm_data_node_spec              = var.warm_data_node_spec
  warm_data_node_disk_size         = var.warm_data_node_disk_size
  warm_data_node_disk_encrypted    = var.warm_data_node_disk_encrypted
  kibana_node_spec                 = var.kibana_node_spec
  instance_category                = var.instance_category
  protocol                         = var.protocol
  setting_config                   = var.setting_config
}

