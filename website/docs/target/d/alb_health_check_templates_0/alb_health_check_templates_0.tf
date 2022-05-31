data "alicloud_alb_health_check_templates" "ids" {
  ids = ["example_id"]
}
output "alb_health_check_template_id_1" {
  value = data.alicloud_alb_health_check_templates.ids.templates.0.id
}

data "alicloud_alb_health_check_templates" "nameRegex" {
  name_regex = "^my-HealthCheckTemplate"
}
output "alb_health_check_template_id_2" {
  value = data.alicloud_alb_health_check_templates.nameRegex.templates.0.id
}

