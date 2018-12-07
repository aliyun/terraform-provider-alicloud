resource "alicloud_log_project" "example" {
  name        = "${var.project_name}"
  description = "created by terraform"
}

resource "alicloud_log_store" "example" {
  project               = "${alicloud_log_project.example.name}"
  name                  = "${var.logstore_name}"
  retention_period      = 3650
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
  append_meta           = true
}

resource "alicloud_log_store_index" "example" {
  project  = "${alicloud_log_project.example.name}"
  logstore = "${alicloud_log_store.example.name}"

  full_text {
    case_sensitive = true
    token          = " #$%^*\r\n\t"
  }

  field_search = [
    {
      name             = "terraform"
      enable_analytics = true
    },
  ]
}

resource "alicloud_log_logs" "example" {
    project = "${alicloud_log_project.example.name}"
    logstore = "${alicloud_log_store.example.name}"
    retry_seconds = 60
    source = "10.1.2.3"
    topic = "test_topic"
    logs = [
        {
            contents = {
                provider = "terraform"
                key1 = "value1"
                key2 = "value2"
                key3 = "value3"
            }
        }
    ]
    tags = {
        tag1 = "value1"
        tag2 = "value2"
    }
}

data "alicloud_log_logs" "example" {
    project = "${alicloud_log_project.example.name}"
    logstore = "${alicloud_log_store.example.name}"
    query = "* and terraform"
    output_file = "./logs.json"
}
