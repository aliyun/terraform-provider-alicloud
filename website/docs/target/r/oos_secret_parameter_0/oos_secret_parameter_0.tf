data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_kms_keys" "default" {
  status = "Enabled"
}

resource "alicloud_kms_key" "default" {
  count                  = length(data.alicloud_kms_keys.default.ids) > 0 ? 0 : 1
  description            = var.name
  status                 = "Enabled"
  pending_window_in_days = 7
}

resource "alicloud_oos_secret_parameter" "example" {
  secret_parameter_name = "example_value"
  value                 = "example_value"
  type                  = "Secret"
  key_id                = length(data.alicloud_kms_keys.default.ids) > 0 ? data.alicloud_kms_keys.default.ids.0 : concat(alicloud_kms_key.default.*.id, [""])[0]
  description           = "example_value"
  tags = {
    Created = "TF"
    For     = "OosSecretParameter"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

