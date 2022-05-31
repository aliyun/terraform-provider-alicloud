variable "name" {
  default = "tf-test112358"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_hbr_oss_backup_plan" "default" {
  oss_backup_plan_name = var.name
  prefix               = "/"
  bucket               = alicloud_oss_bucket.default.bucket
  vault_id             = alicloud_hbr_vault.default.id
  schedule             = "I|1602673264|PT2H"
  backup_type          = "COMPLETE"
  retention            = "2"
}
