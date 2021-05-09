output "instance_id" {
  description = "The ID used to identify alikafka.instance resource."
  value       = alicloud_alikafka_sasl_user.default.instance_id
}

output "username" {
  description = "Name of ALIKAFKA Sasl User."
  value       = alicloud_alikafka_sasl_user.default.username
}

output "password" {
  description = "Password of ALIKAFKA Sasl User."
  value       = alicloud_alikafka_sasl_user.default.password
}