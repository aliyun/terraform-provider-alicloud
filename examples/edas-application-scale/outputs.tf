output "app_id" {
  value       = alicloud_edas_application_scale.default.app_id
  description = "The ID of the application that you want to deploy."
}

output "deploy_group" {
  value       = alicloud_edas_application_scale.default.deploy_group
  description = "The ID of the instance group to which you want to add ECS instances to scale out the application."
}

output "ecu_info" {
  value       = alicloud_edas_application_scale.default.ecu_info
  description = "The ID of the Elastic Compute Unit (ECU) where you want to deploy the application."
}

output "force_status" {
  value       = alicloud_edas_application_scale.default.force_status
  description = "This parameter specifies whether to forcibly remove an ECS instance where the application is deployed. It is set as true only after the ECS instance expires. In normal cases, this parameter do not need to be specified."
}

output "ecc_info" {
  value       = alicloud_edas_application_scale.default.ecc_info
  description = "The ID of the Elastic Compute Container (ECC) is corresponding to the ECS instance that you want to remove for the application."
}