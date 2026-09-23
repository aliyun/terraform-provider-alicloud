variable "instance_charge_type" {
  default = "PostPaid"
}

variable "description" {
  default = "tf-testAccESResource"
}

variable "period" {
  default = "2"
}

variable "data_node_spec" {
  default = "elasticsearch.sn2ne.large"
}

variable "data_node_amount" {
  default = "4"
}

variable "data_node_disk_size" {
  default = "20"
}

variable "data_node_disk_type" {
  default = "cloud_ssd"
}

variable "data_node_disk_encrypted" {
  default = false
}

variable "es_version" {
  default = "7.16_with_X-Pack"
}

variable "password" {
  default = "MTest12345"
}

variable "private_whitelist" {
  type    = list(string)
  default = ["10.1.0.0/16", "10.0.0.0/16"]
}

variable "kibana_whitelist" {
  type    = list(string)
  default = ["10.1.0.0/16", "10.0.0.0/16", "127.0.0.1"]
}

variable "master_node_spec" {
  default = "elasticsearch.sn2ne.xlarge"
}

variable "master_node_disk_type" {
  default = "cloud_ssd"
}

variable "zone_count" {
  default = "1"
}

variable "client_node_amount" {
  default = "2"
}

variable "client_node_spec" {
  default = "elasticsearch.sn2ne.large"
}

variable "protocol" {
  default = "HTTPS"
}

variable "setting_config" {
  type = map(string)
  default = {
    "action.auto_create_index": "+.*,-*",
    "action.destructive_requires_name": "true",
    "xpack.security.audit.enabled": "true"
  }
}

variable "warm_node_amount" {
  default = "3"
}

variable "warm_node_spec" {
  default = "elasticsearch.sn2ne.large"
}

variable "warm_node_disk_size" {
  default = "500"
}

variable "warm_node_disk_type" {
  default = "cloud_efficiency"
}

variable "warm_node_disk_encrypted" {
  default = false
}

variable "kibana_node_spec" {
  default = "elasticsearch.sn1ne.large"
}

variable "enable_kibana_public_network" {
  default = true
}