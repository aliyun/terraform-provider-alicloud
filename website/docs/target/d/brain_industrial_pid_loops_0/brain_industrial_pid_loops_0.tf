data "alicloud_brain_industrial_pid_loops" "example" {
  pid_project_id = "856c6b8f-ca63-40a4-xxxx-xxxx"
  ids            = ["742a3d4e-d8b0-47c8-xxxx-xxxx"]
  name_regex     = "tf-testACC"
}

output "first_brain_industrial_pid_loop_id" {
  value = data.alicloud_brain_industrial_pid_loops.example.loops.0.id
}
