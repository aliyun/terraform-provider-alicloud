output "app_id" {
  value       = alicloud_edas_application_deployment.default.app_id
  description = "The ID of the application that you want to deploy."
}

output "group_id" {
  value       = alicloud_edas_application_deployment.default.group_id
  description = "The ID of the instance group where the application is going to be deployed. "
}


output "package_version" {
  value       = alicloud_edas_application_deployment.default.package_version
  description = "The version of the application that you want to deploy. It must be unique for every application. The length cannot exceed 64 characters. We recommended you to use a timestamp. "
}

output "war_url" {
  value       = alicloud_edas_application_deployment.default.war_url
  description = "The address to store the uploaded web application (WAR) package for application deployment. This parameter is required when the deployType parameter is set as url. We recommend you to set this parameter to the address of an Object Storage Service (OSS) system."
}

output "last_package_version" {
  value       = alicloud_edas_application_deployment.default.last_package_version
  description = "Version of the last deployment package."
}