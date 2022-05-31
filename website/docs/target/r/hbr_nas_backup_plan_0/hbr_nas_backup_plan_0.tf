variable "name" {
  default = "tf-testAccHBRNas"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = var.name
  encrypt_type  = "1"
}

data "alicloud_nas_file_systems" "default" {
  protocol_type     = "NFS"
  description_regex = alicloud_nas_file_system.default.description
}

resource "alicloud_hbr_nas_backup_plan" "default" {
  depends_on           = ["alicloud_nas_file_system.default"]
  nas_backup_plan_name = var.name
  file_system_id       = alicloud_nas_file_system.default.id
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  vault_id             = alicloud_hbr_vault.default.id
  create_time          = data.alicloud_nas_file_systems.default.systems.0.create_time
  retention            = "2"
  path                 = ["/"]
}
