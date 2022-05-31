variable "name" {
  default = "tf-testaccalicloudfcservice"
}

resource "alicloud_fc_custom_domain" "default" {
  domain_name = "terraform.functioncompute.com"
  protocol    = "HTTP"
  route_config {
    path          = "/login/*"
    service_name  = alicloud_fc_service.default.name
    function_name = alicloud_fc_function.default.name
    qualifier     = "v1"
    methods       = ["GET", "POST"]
  }
  cert_config {
    cert_name   = "your certificate name"
    private_key = "your private key"
    certificate = "your certificate data"
  }
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = "${var.name}-description"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "fc/hello.zip"
  content = <<EOF
		# -*- coding: utf-8 -*-
	def handler(event, context):
		print "hello world"
		return 'hello world'
	EOF
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = var.name
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = 512
  runtime     = "python2.7"
  handler     = "hello.handler"
}
