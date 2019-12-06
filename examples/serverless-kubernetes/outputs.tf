// Output VPC
output "vpc_id" {
  description = "The ID of the VPC."
  value       = alicloud_cs_serverless_kubernetes.serverless.vpc_id
}

output "vswitch_id" {
  description = "The ID of the VSwitch."
  value       = alicloud_cs_serverless_kubernetes.serverless.vswitch_id
}

// Output kubernetes resource
output "cluster_id" {
  description = "ID of the kunernetes cluster."
  value       = [alicloud_cs_serverless_kubernetes.serverless.id]
}

output "deletion_protection" {
  description = "ID of the Security Group used to deploy kubernetes cluster."
  value       = alicloud_cs_serverless_kubernetes.serverless.deletion_protection
}

output "endpoint_public_access_enabled" {
  description = "Whether enable public access or not"
  value       = alicloud_cs_serverless_kubernetes.serverless.endpoint_public_access_enabled
}



