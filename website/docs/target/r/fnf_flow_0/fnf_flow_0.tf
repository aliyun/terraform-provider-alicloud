resource "alicloud_ram_role" "default" {
  name     = "tf-testacc-fnfflow"
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

resource "alicloud_fnf_flow" "example" {
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: pass
      name: helloworld
  EOF
  role_arn    = alicloud_ram_role.default.arn
  description = "Test for terraform fnf_flow."
  name        = "tf-testacc-flow"
  type        = "FDL"
}
