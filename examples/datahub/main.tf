resource "alicloud_datahub_project" "example" {
  name    = var.project_name
  comment = "Datahub project: a terraform example."
}

resource "alicloud_datahub_topic" "example" {
  name         = var.topic_name
  project_name = alicloud_datahub_project.example.name
  shard_count  = 3
  life_cycle   = 7
  record_type  = "BLOB"
  comment      = "Datahub blob topic: a terraform example."
}

resource "alicloud_datahub_subscription" "example" {
  project_name = alicloud_datahub_project.example.name
  topic_name   = alicloud_datahub_topic.example.name
  comment      = "Datahub subscription: a terraform example."
}

