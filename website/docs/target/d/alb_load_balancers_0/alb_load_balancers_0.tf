data "alicloud_alb_load_balancers" "ids" {}
output "alb_load_balancer_id_1" {
  value = data.alicloud_alb_load_balancers.ids.balancers.0.id
}

data "alicloud_alb_load_balancers" "nameRegex" {
  name_regex = "^my-LoadBalancer"
}
output "alb_load_balancer_id_2" {
  value = data.alicloud_alb_load_balancers.nameRegex.balancers.0.id
}

