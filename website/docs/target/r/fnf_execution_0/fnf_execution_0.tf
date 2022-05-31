variable "name" {
  default = "tf-testacc-fnfflow"
}

resource "alicloud_ram_role" "default" {
  name     = var.name
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "fnf.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}

resource "alicloud_fnf_flow" "default" {
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: wait
      name: custom_wait
      duration: $.wait
  EOF
  role_arn    = alicloud_ram_role.default.arn
  description = "Test for terraform fnf_flow."
  name        = var.name
  type        = "FDL"
}

resource "alicloud_fnf_execution" "default" {
  execution_name = var.name
  flow_name      = alicloud_fnf_flow.default.name
  input          = "{\"wait\": 600}"
}
