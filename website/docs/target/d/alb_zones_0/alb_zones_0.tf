data "alicloud_alb_zones" "example" {}

output "first_alb_zones_id" {
  value = data.alicloud_alb_zones.example.zones.0.zone_id
}
