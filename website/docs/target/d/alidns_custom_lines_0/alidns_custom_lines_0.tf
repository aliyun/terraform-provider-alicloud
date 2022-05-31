data "alicloud_alidns_custom_lines" "ids" {
  enable_details = true
  domain_name    = "your_domain_name"
}
output "alidns_custom_line_id_1" {
  value = data.alicloud_alidns_custom_lines.ids.lines.0.id
}
