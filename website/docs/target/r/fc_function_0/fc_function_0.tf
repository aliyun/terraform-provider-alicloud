variable "name" {
  default = "alicloudfcfunctionconfig"
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "tf unit test"
}

resource "alicloud_log_store" "default" {
  project          = alicloud_log_project.default.name
  name             = var.name
  retention_period = "3000"
  shard_count      = 1
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "tf unit test"
  log_config {
    project  = alicloud_log_project.default.name
    logstore = alicloud_log_store.default.name
  }
  role       = alicloud_ram_role.default.arn
  depends_on = [alicloud_ram_role_policy_attachment.default]
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.id
  key    = "fc/hello.zip"
  source = "./hello.zip"
}

resource "alicloud_ram_role" "default" {
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

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_function" "foo" {
  service     = alicloud_fc_service.default.name
  name        = var.name
  description = "tf"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python2.7"
  handler     = "hello.handler"
  environment_variables = {
    prefix = "terraform"
  }
}
