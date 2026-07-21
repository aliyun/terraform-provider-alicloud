# Stack Group Information
output "stack_group_name" {
  description = "The name of the ROS stack group"
  value       = alicloud_ros_stack_group.example.stack_group_name
}

output "stack_group_id" {
  description = "The ID of the ROS stack group"
  value       = alicloud_ros_stack_group.example.stack_group_id
}

output "stack_group_status" {
  description = "The status of the ROS stack group"
  value       = alicloud_ros_stack_group.example.status
}

# Stack Instances Information
output "stack_instances_count" {
  description = "Total number of stack instances created"
  value       = length(alicloud_ros_stack_instances.self_managed.stack_instances)
}

output "stack_instances_details" {
  description = "Detailed information about all stack instances"
  value       = alicloud_ros_stack_instances.self_managed.stack_instances
}

output "stack_instance_ids" {
  description = "List of stack instance IDs (last_operation_id) for tracking operations"
  value       = [for inst in alicloud_ros_stack_instances.self_managed.stack_instances : inst.last_operation_id]
}

output "stack_instance_regions" {
  description = "List of regions where stack instances were deployed"
  value       = distinct([for inst in alicloud_ros_stack_instances.self_managed.stack_instances : inst.region_id])
}

output "stack_instance_accounts" {
  description = "List of accounts where stack instances were deployed"
  value       = distinct([for inst in alicloud_ros_stack_instances.self_managed.stack_instances : inst.account_id])
}

# Deployment Summary
output "deployment_summary" {
  description = "Summary of the deployment including region and account distribution"
  value = {
    total_instances        = length(alicloud_ros_stack_instances.self_managed.stack_instances)
    regions_deployed       = distinct([for inst in alicloud_ros_stack_instances.self_managed.stack_instances : inst.region_id])
    accounts_deployed      = distinct([for inst in alicloud_ros_stack_instances.self_managed.stack_instances : inst.account_id])
    instance_statuses      = { for inst in alicloud_ros_stack_instances.self_managed.stack_instances : "${inst.account_id}/${inst.region_id}" => inst.status }
    instance_drift_status  = { for inst in alicloud_ros_stack_instances.self_managed.stack_instances : "${inst.account_id}/${inst.region_id}" => inst.stack_drift_status }
    last_operation_ids     = { for inst in alicloud_ros_stack_instances.self_managed.stack_instances : "${inst.account_id}/${inst.region_id}" => inst.last_operation_id }
  }
}
