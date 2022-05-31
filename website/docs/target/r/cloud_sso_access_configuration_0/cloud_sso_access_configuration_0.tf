variable "name" {
  default = "example-value"
}
data "alicloud_cloud_sso_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}
resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = var.name
  directory_id              = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
  permission_policies {
    permission_policy_document = <<EOF
		{
        "Statement":[
        {
        "Action":"ecs:Get*",
        "Effect":"Allow",
        "Resource":[
            "*"
        ]
        }
        ],
			"Version": "1"
		}
	  EOF
    permission_policy_type     = "Inline"
    permission_policy_name     = var.name
  }
}
