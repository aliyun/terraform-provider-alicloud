variable name {
  default = "example-name"
}

data alicloud_cloud_sso_directories default {}

data alicloud_resource_manager_resource_directories default {}

resource alicloud_cloud_sso_directory default {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

resource alicloud_cloud_sso_access_configuration default {
  access_configuration_name = var.name
  directory_id              = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id[""])[0]
}

resource alicloud_cloud_sso_access_configuration_provisioning default {
  directory_id            = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id[""])[0]
  access_configuration_id = alicloud_cloud_sso_access_configuration.default.access_configuration_id
  target_type             = "RD-Account"
  target_id               = data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id
}
