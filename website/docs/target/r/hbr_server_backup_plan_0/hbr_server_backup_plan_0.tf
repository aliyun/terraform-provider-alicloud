data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-server-backup-plan"
  status     = "Running"
}

resource "alicloud_hbr_server_backup_plan" "example" {
  ecs_server_backup_plan_name = "server_backup_plan"
  instance_id                 = data.alicloud_instances.default.instances.0.id
  schedule                    = "I|1602673264|PT2H"
  retention                   = 1
  detail {
    app_consistent = true
    snapshot_group = true
  }
  disabled = false
}
