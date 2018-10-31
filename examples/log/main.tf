resource "alicloud_log_project" "example" {
  name       = "${var.project_name}"
  description = "created by terraform"
}

resource "alicloud_log_store" "example" {
  project = "${alicloud_log_project.example.name}"
  name       = "${var.logstore_name}"
  retention_period = 3650
  shard_count = 3
  auto_split = true
  max_split_shard_count = 60
  append_meta = true
}

resource "alicloud_log_store_index" "example" {
  project = "${alicloud_log_project.example.name}"
  logstore = "${alicloud_log_store.example.name}"
  full_text {
    case_sensitive = true
    token = " #$%^*\r\n\t"
  }
  field_search = [
    {
      name = "terraform"
      enable_analytics = true
    }
  ]
}
