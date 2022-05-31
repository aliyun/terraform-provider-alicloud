data "alicloud_actiontrail_trails" "default" {
  name_regex = "tf-testacc-actiontrail"
}

output "trail_name" {
  value = data.alicloud_actiontrail_trails.default.trails.0.id
}
