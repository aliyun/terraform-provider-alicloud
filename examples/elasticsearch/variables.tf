variable "instance_charge_type" {
  default = "PostPaid"
}

variable "description" {
  default = "testABC"
}

variable "period" {
  default = "2"
}

variable "data_node_spec" {
  default = "elasticsearch.sn2ne.large"
}

variable "data_node_amount" {
  default = "3"
}

variable "data_node_disk" {
  default = "40"
}

variable "data_node_disk_type" {
  default = "cloud_ssd"
}

variable "es_version" {
  default = "6.3_with_X-Pack"
}

variable "vswitch_id" {
  default = "switch id"
}

variable "es_admin_password" {
  default = "MTest1234"
}

variable "private_whitelist" {
    type = "list"
    default = [ "10.1.0.0/16", "10.0.0.0/16" ]
}

variable "kibana_whitelist" {
    type = "list"
    default = [ "10.1.0.0/16","10.0.0.0/16", "127.0.0.1" ]
}

variable "master_node_spec" {
    default = "elasticsearch.sn2ne.xlarge"
}
