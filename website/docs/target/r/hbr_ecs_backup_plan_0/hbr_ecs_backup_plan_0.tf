variable "name" {
  default = "valut-name"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-backup-plan"
  status     = "Running"
}

resource "alicloud_hbr_ecs_backup_plan" "example" {
  ecs_backup_plan_name = "example_value"
  instance_id          = data.alicloud_instances.default.instances.0.id
  vault_id             = alicloud_hbr_vault.default.id
  retention            = "1"
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  speed_limit          = "0:24:5120"
  path                 = ["/home", "/var"]
  exclude              = <<EOF
  ["/home/exclude"]
  EOF
  include              = <<EOF
  ["/home/include"]
  EOF
}
