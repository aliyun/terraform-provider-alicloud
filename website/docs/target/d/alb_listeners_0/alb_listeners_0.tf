data "alicloud_alb_listeners" "ids" {
  ids = ["example_id"]
}
output "alb_listener_id_1" {
  value = data.alicloud_alb_listeners.ids.listeners.0.id
}

