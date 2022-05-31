
variable "name" {
  default = "tf-testAccCassandrBackupPlan"
}

data "alicloud_cassandra_zones" "example" {
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  vpc_id       = alicloud_cassandra_cluster.example.id
  zone_id      = data.alicloud_cassandra_zones.example.zones[length(data.alicloud_cassandra_zones.example.ids) + (-1)].id
  cidr_block   = cidrsubnet(alicloud_vpc.example.vpcs.0.cidr_block, 8, 4)
}

resource "alicloud_cassandra_cluster" "example" {
  cluster_name        = var.name
  data_center_name    = var.name
  auto_renew          = "false"
  instance_type       = "cassandra.c.large"
  major_version       = "3.11"
  node_count          = "2"
  pay_type            = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.example[0].id
  disk_size           = "160"
  disk_type           = "cloud_ssd"
  maintain_start_time = "18:00Z"
  maintain_end_time   = "20:00Z"
  ip_white            = "127.0.0.1"
}

resource "alicloud_cassandra_backup_plan" "example" {
  cluster_id     = alicloud_cassandra_cluster.example.id
  data_center_id = alicloud_cassandra_cluster.example.zone_id
  backup_time    = "00:10Z"
  active         = false

}

