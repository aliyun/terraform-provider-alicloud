data "alicloud_actiontrail_history_delivery_jobs" "ids" {
  ids = ["example_id"]
}
output "actiontrail_history_delivery_job_id_1" {
  value = data.alicloud_actiontrail_history_delivery_jobs.ids.jobs.0.id
}

data "alicloud_actiontrail_history_delivery_jobs" "status" {
  ids    = ["example_id"]
  status = "2"
}
output "actiontrail_history_delivery_job_id_2" {
  value = data.alicloud_actiontrail_history_delivery_jobs.status.jobs.0.id
}

