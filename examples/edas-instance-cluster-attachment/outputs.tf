output "cluster_id" {
  value       = alicloud_edas_instance_cluster_attachment.default.cluster_id
  description = "The ID of the cluster that you want to create the application. The default cluster will be used if you do not specify this parameter."
}

output "instance_ids" {
  value       = alicloud_edas_instance_cluster_attachment.default.instance_ids
  description = "The ID of instance. e.g. instanceId1, instanceId2."
}

output "status_map" {
  value       = alicloud_edas_instance_cluster_attachment.default.status_map
  description = "ECS instance's import status. Valid values: 1 means running; 0 means being converted; -1 means failed to be converted; -2 means Offline. "
}


output "ecu_map" {
  value       = alicloud_edas_instance_cluster_attachment.default.ecu_map
  description = "The ecu map generated after ECS is imported into the cluster. Key: EcsId. Value: EcuId."
}

output "cluster_member_ids" {
  value       = alicloud_edas_instance_cluster_attachment.default.cluster_member_ids
  description = "The list of cluster_member_id. Key: ECSid. Value: cluster_member_id."
}