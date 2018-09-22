// Output VPC
output "vpc_id" {
  description = "The ID of the VPC."
  value = "${alicloud_cs_kubernetes.k8s.0.vpc_id}"
}

output "vswitch_ids" {
  description = "List ID of the VSwitches."
  value = ["${alicloud_cs_kubernetes.k8s.*.vswitch_ids}"]
}

output "nat_gateway_id" {
  value = "${alicloud_cs_kubernetes.k8s.0.nat_gateway_id}"
}

// Output kubernetes resource
output "cluster_id" {
  description = "ID of the kunernetes cluster."
  value = ["${alicloud_cs_kubernetes.k8s.*.id}"]
}

output "security_group_id" {
  description = "ID of the Security Group used to deploy kubernetes cluster."
  value = "${alicloud_cs_kubernetes.k8s.0.security_group_id}"
}

output "worker_nodes" {
  description = "List worker nodes of cluster."
  value = ["${alicloud_cs_kubernetes.k8s.*.worker_nodes}"]
}

output "master_nodes" {
  description = "List master nodes of cluster."
  value = ["${alicloud_cs_kubernetes.k8s.*.master_nodes}"]
}