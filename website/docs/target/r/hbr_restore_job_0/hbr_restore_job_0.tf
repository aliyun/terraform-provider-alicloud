data "alicloud_hbr_ecs_backup_plans" "default" {
  name_regex = "plan-tf-used-dont-delete"
}
data "alicloud_hbr_oss_backup_plans" "default" {
  name_regex = "plan-tf-used-dont-delete"
}
data "alicloud_hbr_nas_backup_plans" "default" {
  name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_snapshots" "ecs_snapshots" {
  source_type = "ECS_FILE"
  vault_id    = data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id
  instance_id = data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id
}

data "alicloud_hbr_snapshots" "oss_snapshots" {
  source_type = "OSS"
  vault_id    = data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
  bucket      = data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
}

data "alicloud_hbr_snapshots" "nas_snapshots" {
  source_type    = "NAS"
  vault_id       = data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
  file_system_id = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
  create_time    = data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
}

resource "alicloud_hbr_restore_job" "nasJob" {
  snapshot_hash         = data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_hash
  vault_id              = data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
  source_type           = "NAS"
  restore_type          = "NAS"
  snapshot_id           = data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_id
  target_file_system_id = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
  target_create_time    = data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
  target_path           = "/"
  options               = <<EOF
    {"includes":[], "excludes":[]}
  EOF
}

resource "alicloud_hbr_restore_job" "ossJob" {
  snapshot_hash = data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_hash
  vault_id      = data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
  source_type   = "OSS"
  restore_type  = "OSS"
  snapshot_id   = data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_id
  target_bucket = data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
  target_prefix = ""
  options       = <<EOF
    {"includes":[], "excludes":[]}
  EOF
}

resource "alicloud_hbr_restore_job" "ecsJob" {
  snapshot_hash      = data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_hash
  vault_id           = data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id
  source_type        = "ECS_FILE"
  restore_type       = "ECS_FILE"
  snapshot_id        = data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_id
  target_instance_id = data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id
  target_path        = "/"
}
