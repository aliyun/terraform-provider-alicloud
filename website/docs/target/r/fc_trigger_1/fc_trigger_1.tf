variable "name" {
  default = "fctriggermnstopic"
}
data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
}
resource "alicloud_log_project" "foo" {
  name        = "${var.name}"
  description = "tf unit test"
}
resource "alicloud_log_store" "bar" {
  project          = "${alicloud_log_project.foo.name}"
  name             = "${var.name}-source"
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_log_store" "foo" {
  project          = "${alicloud_log_project.foo.name}"
  name             = "${var.name}"
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_mns_topic" "foo" {
  name = "${var.name}"
}
resource "alicloud_fc_service" "foo" {
  name            = "${var.name}"
  internet_access = false
}
resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key    = "fc/hello.zip"
  source = "./hello.zip"
}
resource "alicloud_fc_function" "foo" {
  service     = "${alicloud_fc_service.foo.name}"
  name        = "${var.name}"
  oss_bucket  = "${alicloud_oss_bucket.foo.id}"
  oss_key     = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = 512
  runtime     = "python2.7"
  handler     = "hello.handler"
}
resource "alicloud_ram_role" "foo" {
  name        = "${var.name}-trigger"
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "mns.aliyuncs.com"
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
  role_name   = "${alicloud_ram_role.foo.name}"
  policy_name = "AliyunMNSNotificationRolePolicy"
  policy_type = "System"
}
resource "alicloud_fc_trigger" "foo" {
  service    = "${alicloud_fc_service.foo.name}"
  function   = "${alicloud_fc_function.foo.name}"
  name       = "${var.name}"
  role       = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:mns:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:/topics/${alicloud_mns_topic.foo.name}"
  type       = "mns_topic"
  config_mns = <<EOF
  {
    "filterTag":"testTag",
    "notifyContentFormat":"STREAM",
    "notifyStrategy":"BACKOFF_RETRY"
  }
  EOF
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}
