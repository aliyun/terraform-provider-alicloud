provider "archive" {}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "index.js"
  output_path = "index.zip"
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

resource "alicloud_ram_policy" "policy" {
  name = "${var.ram_policy_name}"

  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
  EOF

  description = "this is a policy test"
  force       = true
}

resource "alicloud_ram_role" "role" {
  name = "${var.ram_role_name}"

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

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name   = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
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

  request_config {
    protocol = "HTTP"
    method   = "POST"
    path     = "/test/path2"
    mode     = "MAPPING"
    body_format = "STREAM"
  }

  service_type = "FunctionCompute"

  fc_service_config {
    region        = "${var.fc_region}"
    function_name = "${alicloud_fc_function.foo.name}"
    service_name  = "${alicloud_fc_service.foo.name}"
    arn_role      = "${alicloud_ram_role.role.arn}"
    timeout       = 10
  }

  request_parameters {
      name         = "aa"
      type         = "STRING"
      required     = "REQUIRED"
      in           = "QUERY"
      in_service   = "QUERY"
      name_service = "testparams"
    }

  stage_names = [
    "RELEASE",
    "PRE",
    "TEST",
  ]
}

resource "alicloud_api_gateway_app" "apiGatewayApp" {
  name        = "${var.apigateway_app_name_test}"
  description = "${var.apigateway_app_description_test}"
}

resource "alicloud_api_gateway_app_attachment" "foo" {
  api_id     = "${alicloud_api_gateway_api.apiGatewayApi.api_id}"
  group_id   = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
  stage_name = "RELEASE"
  app_id     = "${alicloud_api_gateway_app.apiGatewayApp.id}"
}
