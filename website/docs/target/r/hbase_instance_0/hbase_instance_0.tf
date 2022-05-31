resource "alicloud_hbase_instance" "default" {
  name                   = "tf_testAccHBase_vpc"
  zone_id                = "cn-shenzhen-b"
  vswitch_id             = "vsw-123456"
  engine                 = "hbaseue"
  engine_version         = "2.0"
  master_instance_type   = "hbase.sn1.large"
  core_instance_type     = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type         = "cloud_efficiency"
  core_disk_size         = 400
  pay_type               = "PostPaid"
  cold_storage_size      = 0
}
