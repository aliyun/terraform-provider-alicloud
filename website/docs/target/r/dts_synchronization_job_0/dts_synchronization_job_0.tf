resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "PolarDB"
  source_endpoint_region           = "cn-hangzhou"
  destination_endpoint_engine_name = "ADB30"
  destination_endpoint_region      = "cn-hangzhou"
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_synchronization_job" "default" {
  dts_instance_id                    = alicloud_dts_synchronization_instance.default.id
  dts_job_name                       = "tf-testAccCase1"
  source_endpoint_instance_type      = "PolarDB"
  source_endpoint_instance_id        = "pc-xxxxxxxx"
  source_endpoint_engine_name        = "PolarDB"
  source_endpoint_region             = "cn-hangzhou"
  source_endpoint_database_name      = "tf-testacc"
  source_endpoint_user_name          = "root"
  source_endpoint_password           = "password"
  destination_endpoint_instance_type = "ads"
  destination_endpoint_instance_id   = "am-xxxxxxxx"
  destination_endpoint_engine_name   = "ADB30"
  destination_endpoint_region        = "cn-hangzhou"
  destination_endpoint_database_name = "tf-testacc"
  destination_endpoint_user_name     = "root"
  destination_endpoint_password      = "password"
  db_list                            = "{\"tf-testacc\":{\"name\":\"tf-test\",\"all\":true,\"state\":\"normal\"}}"
  structure_initialization           = "true"
  data_initialization                = "true"
  data_synchronization               = "true"
  status                             = "Synchronizing"
}
