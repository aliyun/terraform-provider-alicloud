provider "archive" {
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "hello.py"
  output_path = "hello.zip"
}

data "alicloud_account" "current" {
}

resource "alicloud_fc_service" "foo" {
  name            = var.service_name
  description     = var.service_description
  internet_access = var.service_internet_access
}

resource "alicloud_fc_function" "foo" {
  service     = alicloud_fc_service.foo.name
  name        = var.function_name
  description = var.function_description
  filename    = var.function_filename
  memory_size = var.function_memory_size
  runtime     = var.function_runtime
  handler     = var.function_handler
}

resource "alicloud_fc_trigger" "foo" {
  service    = alicloud_fc_service.foo.name
  function   = alicloud_fc_function.foo.name
  name       = var.trigger_name
  role       = alicloud_ram_role.foo.arn
  source_arn = "acs:log:${var.region}:${data.alicloud_account.current.id}:project/${alicloud_log_project.foo.name}"
  type       = "log"

  config = <<EOF
    {
        "sourceConfig": {
            "project": alicloud_log_project.foo.name,
            "logstore": alicloud_log_store.source_store.name
        },
        "jobConfig": {
            "maxRetryTime": 3,
            "triggerInterval": 60
        },
        "functionParameter": {
            "a": "b",
            "c": "d"
        },
        "logConfig": {
            "project": alicloud_log_project.foo.name,
            "logstore": alicloud_log_store.fc_store.name
        },
        "enable": true
    }
  
EOF


  depends_on = [alicloud_ram_role_policy_attachment.foo]
}

resource "alicloud_ram_role" "foo" {
  name = "role-for-fc"

  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "log.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  
EOF


  description = "this is a test"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name   = alicloud_ram_role.foo.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_log_project" "foo" {
  name        = "project-for-fc-hello"
  description = "created by terraform"
}

resource "alicloud_log_store" "source_store" {
  project = alicloud_log_project.foo.name
  name    = "store-for-source-store"
}

resource "alicloud_log_store" "fc_store" {
  project = alicloud_log_project.foo.name
  name    = "store-for-fc-store"
}

