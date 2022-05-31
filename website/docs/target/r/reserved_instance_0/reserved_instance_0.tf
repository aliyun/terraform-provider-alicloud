resource "alicloud_reserved_instance" "default" {
  instance_type   = "ecs.g6.large"
  instance_amount = "1"
  period_unit     = "Year"
  offering_type   = "All Upfront"
  name            = name
  description     = "ReservedInstance"
  zone_id         = "cn-hangzhou-h"
  scope           = "Zone"
  period          = "1"
}
