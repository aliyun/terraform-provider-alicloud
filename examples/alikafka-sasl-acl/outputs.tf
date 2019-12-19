output "instance_id" {
  description = "The ID used to identify alikafka.instance resource."
  value       = "${alicloud_alikafka_sasl_acl.default.instance_id}"
}

output "username" {
  description = "Username of ALIKAFKA Sasl Acl."
  value       = "${alicloud_alikafka_sasl_acl.default.username}"
}

output "acl_resource_type" {
  description = "Resource type of ALIKAFKA Sasl Acl."
  value       = "${alicloud_alikafka_sasl_acl.default.acl_resource_type}"
}

output "acl_resource_name" {
  description = "Resource name of ALIKAFKA Sasl Acl."
  value       = "${alicloud_alikafka_sasl_acl.default.acl_resource_name}"
}

output "acl_resource_pattern_type" {
  description = "Resource pattern type of ALIKAFKA Sasl Acl."
  value       = "${alicloud_alikafka_sasl_acl.default.acl_resource_pattern_type}"
}

output "acl_operation_type" {
  description = "Operation type of ALIKAFKA Sasl Acl."
  value       = "${alicloud_alikafka_sasl_acl.default.acl_operation_type}"
}

output "host" {
  description = "The host of ALIKAFKA Sasl Acl."
  value       = "${alicloud_alikafka_sasl_acl.default.host}"
}