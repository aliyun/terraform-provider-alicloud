variable "region" {
  default = "cn-hangzhou"
}

variable "account" {
  default = "12345"
}

provider "alicloud" {
  account_id = var.account
  region     = var.region
}

resource "alicloud_fc_trigger" "foo" {
  service    = "my-fc-service"
  function   = "hello-world"
  name       = "hello-trigger"
  role       = alicloud_ram_role.foo.arn
  source_arn = "acs:log:${var.region}:${var.account}:project/${alicloud_log_project.foo.name}"
  type       = "log"
  config     = <<EOF
    {
        "sourceConfig": {
            "project": "project-for-fc",
            "logstore": "project-for-fc"
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
            "project": "project-for-fc",
            "logstore": "project-for-fc"
        },
        "enable": true
    }
  
EOF


  depends_on = [alicloud_ram_role_policy_attachment.foo]
}

resource "alicloud_ram_role" "foo" {
  name     = "${var.name}-trigger"
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
