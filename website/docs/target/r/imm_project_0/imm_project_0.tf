resource "alicloud_ram_role" "role" {
  name        = "example_role_name"
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "imm.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description = "this is a role test."
  force       = true
}
resource "alicloud_imm_project" "example" {
  project      = "example_name"
  service_role = alicloud_ram_role.role.name
}
