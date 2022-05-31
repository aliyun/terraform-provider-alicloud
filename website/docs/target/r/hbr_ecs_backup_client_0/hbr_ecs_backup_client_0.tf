data "alicloud_instances" "default" {
  name_regex = "ecs_instance_name"
  status     = "Running"
}

resource "alicloud_hbr_ecs_backup_client" "example" {
  instance_id        = data.alicloud_instances.default.instances.0.id
  use_https          = false
  data_network_type  = "PUBLIC"
  max_cpu_core       = 2
  max_worker         = 4
  data_proxy_setting = "USE_CONTROL_PROXY"
  proxy_host         = "192.168.11.101"
  proxy_port         = 80
  proxy_user         = "user"
  proxy_password     = "password"
}
