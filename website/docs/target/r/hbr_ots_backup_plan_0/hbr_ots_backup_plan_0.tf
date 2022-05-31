variable "name" {
  default = "testAcc"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  vault_type = "OTS_BACKUP"
}

resource "alicloud_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "basic" {
  instance_name = alicloud_ots_instance.foo.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}

resource "alicloud_hbr_ots_backup_plan" "example" {
  ots_backup_plan_name = var.name
  vault_id             = alicloud_hbr_vault.default.id
  backup_type          = "COMPLETE"
  schedule             = "I|1602673264|PT2H"
  retention            = "2"
  instance_name        = alicloud_ots_instance.foo.name
  ots_detail {
    table_names = [alicloud_ots_table.basic.table_name]
  }
}
