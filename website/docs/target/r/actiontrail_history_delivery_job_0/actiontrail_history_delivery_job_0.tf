data "alicloud_regions" "example" {
  current = true
}

data "alicloud_account" "example" {}

resource "alicloud_log_project" "example" {
  name        = "example_value"
  description = "tf actiontrail test"
}


resource "alicloud_actiontrail_trail" "example" {
  trail_name      = "example_value"
  sls_project_arn = "acs:log:${data.alicloud_regions.example.regions.0.id}:${data.alicloud_account.example.id}:project/${alicloud_log_project.example.name}"
}

resource "alicloud_actiontrail_history_delivery_job" "example" {
  trail_name = alicloud_actiontrail_trail.example.id
}
