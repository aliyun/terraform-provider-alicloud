output "app_id" {
  value       = alicloud_edas_slb_attachment.default.app_id
  description = "The ID of the applicaton to which you want to bind an SLB instance."
}

output "slb_id" {
  value       = alicloud_edas_slb_attachment.default.slb_id
  description = "The ID of the SLB instance that is going to be bound."
}

output "slb_ip" {
  value       = alicloud_edas_slb_attachment.default.slb_ip
  description = "The IP address that is allocated to the bound SLB instance."
}

output "type" {
  value       = alicloud_edas_slb_attachment.default.type
  description = "The type of the bound SLB instance."
}

output "listener_port" {
  value       = alicloud_edas_slb_attachment.default.listener_port
  description = "The listening port for the bound SLB instance."
}

output "vserver_group_id" {
  value       = alicloud_edas_slb_attachment.default.vserver_group_id
  description = ""
}

output "vswitch_id" {
  value       = alicloud_edas_slb_attachment.default.vswitch_id
  description = "The ID of the virtual server (VServer) group associated with the intranet SLB instance."
}

output "slb_status" {
  value       = alicloud_edas_slb_attachment.default.slb_status
  description = "Running Status of SLB instance. Inactive：The instance is stopped, and listener will not monitor and foward traffic. Active：The instance is running. After the instance is created, the default state is active. Locked：The instance is locked, the instance has been owed or locked by Alibaba Cloud. Expired: The instance has expired."
}