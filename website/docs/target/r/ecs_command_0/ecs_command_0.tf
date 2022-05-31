resource "alicloud_ecs_command" "example" {
  name            = "tf-testAcc"
  command_content = "bHMK"
  description     = "For Terraform Test"
  type            = "RunShellScript"
  working_dir     = "/root"
}

