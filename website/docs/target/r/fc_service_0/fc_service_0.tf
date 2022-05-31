variable "name" {
  default = "tf-testaccalicloudfcservice"
}

resource "alicloud_log_project" "foo" {
  name = var.name
}

resource "alicloud_log_store" "foo" {
  project = alicloud_log_project.foo.name
  name    = var.name
}

resource "alicloud_ram_role" "role" {
  name        = var.name
  document    = <<EOF
  {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "fc.aliyuncs.com"
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

resource "alicloud_ram_role_policy_attachment" "attac" {
  role_name   = alicloud_ram_role.role.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "foo" {
  name        = var.name
  description = "tf unit test"
  role        = alicloud_ram_role.role.arn
  depends_on  = [alicloud_ram_role_policy_attachment.attac]
}
