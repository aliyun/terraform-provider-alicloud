resource "alicloud_governance_account" "default" {
  account_name_prefix = var.account_name_prefix
  folder_id           = var.folder_id
  baseline_id         = var.baseline_id
  display_name        = var.display_name
  default_domain_name = var.default_domain_name
}