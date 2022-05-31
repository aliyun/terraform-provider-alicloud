resource "alicloud_brain_industrial_pid_loop" "example" {
  pid_loop_configuration = "YourLoopConfiguration"
  pid_loop_dcs_type      = "standard"
  pid_loop_is_crucial    = true
  pid_loop_name          = "tf-testAcc"
  pid_loop_type          = "0"
  pid_project_id         = "856c6b8f-ca63-40a4-xxxx-xxxx"
}

