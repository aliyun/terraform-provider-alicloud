provider "archive" {}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "hello.py"
  output_path = "hello.zip"
}

data "alicloud_account" "current" {}

resource "alicloud_fc_service" "foo" {
  name            = "${var.service_name}"
  description     = "${var.service_description}"
  internet_access = "${var.service_internet_access}"
}

resource "alicloud_fc_function" "foo" {
  service     = "${alicloud_fc_service.foo.name}"
  name        = "${var.function_name}"
  description = "${var.function_description}"
  filename    = "${var.function_filename}"
  memory_size = "${var.function_memory_size}"
  runtime     = "${var.function_runtime}"
  handler     = "${var.function_handler}"
}

resource "alicloud_ram_role" "foo" {
  name = "AliyunApiGatewayAccessingFCRole1"

  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "apigateway.aliyuncs.com"
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

resource "alicloud_api_gateway_group" "apiGatewayGroup" {
  name        = "${var.apigateway_group_name}"
  description = "${var.apigateway_group_description}"
}

resource "alicloud_api_gateway_api" "apiGatewayApi" {
  name        = "terraformapi"
  group_id    = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
  description = "description"
  auth_type   = "APP"

  request_config = {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path1"
    mode     = "MAPPING"
  }

  service_type = "FunctionCompute"

  fc_service_config = {
    region        = "${var.fc_region}"
    function_name = "${alicloud_fc_function.foo.name}"
    service_name  = "${alicloud_fc_service.foo.name}"
    arn_role      = "${alicloud_ram_role.foo.arn}"
    timeout       = 10
  }

  request_parameters = [
    {
      name         = "aa"
      type         = "STRING"
      required     = "REQUIRED"
      in           = "QUERY"
      in_service   = "QUERY"
      name_service = "testparams"
    },
  ]

  stage_names = [
    "RELEASE",
    "PRE",
    "TEST",
  ]
}
