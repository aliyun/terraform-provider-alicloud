output "app_id" {
  value       = alicloud_edas_deploy_group.default.app_id
  description = "The ID of the application that you want to deploy."
}

output "group_name" {
  value       = alicloud_edas_deploy_group.default.group_name
  description = "The name of the instance group that you want to create. The length cannot exceed 64 characters."
}

output "group_type" {
  value       = alicloud_edas_deploy_group.default.group_type
  description = "The type of the instance group that you want to create. Valid values: 0: Default group. 1: Phased release is disabled for traffic management. 2: Phased release is enabled for traffic management."
}