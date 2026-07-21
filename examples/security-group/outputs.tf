output "http_rule_id" {
  value = alicloud_security_group_rule.http-in.id
}

output "ssh_rule_id" {
  value = alicloud_security_group_rule.ssh-in.id
}

