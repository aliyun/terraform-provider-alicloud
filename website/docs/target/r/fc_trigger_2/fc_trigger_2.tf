variable "name" {
  default = "fctriggercdneventsconfig"
}

data "alicloud_account" "current" {
}

resource "alicloud_cdn_domain_new" "domain" {
  domain_name = "${var.name}.tf.com"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = 20
    port     = 80
    weight   = 10
  }
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
        "Version": "1",
        "Statement": [
            {
                "Action": "cdn:Describe*",
                "Resource": "*",
                "Effect": "Allow",
		        "Principal": {
                "Service":
                    ["log.aliyuncs.com"]
                }
            }
        ]
    }
    EOF
  description = "this is a test"
  force       = true
}

resource "alicloud_ram_policy" "foo" {
  name        = "${var.name}-trigger"
  document    = <<EOF
    {
        "Version": "1",
        "Statement": [
        {
            "Action": [
            "fc:InvokeFunction"
            ],
        "Resource": [
            "acs:fc:*:*:services/tf_cdnEvents/functions/*",
            "acs:fc:*:*:services/tf_cdnEvents.*/functions/*"
        ],
        "Effect": "Allow"
        }
        ]
    }
    EOF
  description = "this is a test"
  force       = true
}
resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name   = "${alicloud_ram_role.foo.name}"
  policy_name = "${alicloud_ram_policy.foo.name}"
  policy_type = "Custom"
}
resource "alicloud_fc_trigger" "default" {
  service    = "${alicloud_fc_service.foo.name}"
  function   = "${alicloud_fc_function.foo.name}"
  name       = "${var.name}"
  role       = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:cdn:*:${data.alicloud_account.current.id}"
  type       = "cdn_events"
  config     = <<EOF
      {"eventName":"LogFileCreated",
     "eventVersion":"1.0.0",
     "notes":"cdn events trigger",
     "filter":{
        "domain": ["${alicloud_cdn_domain_new.domain.domain_name}"]
        }
    }
EOF
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}
